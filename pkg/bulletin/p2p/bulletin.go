package p2p

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/ipfs/boxo/util"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p/core/event"
	"github.com/libp2p/go-libp2p/core/network"
	libp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"google.golang.org/protobuf/proto"

	eventbus "github.com/sourcenetwork/eventbus-go"
	rpc "github.com/sourcenetwork/go-libp2p-pubsub-rpc"

	"github.com/sourcenetwork/orbis-go/config"
	gossipbulletinv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/gossipbulletin/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin/memmap"
	"github.com/sourcenetwork/orbis-go/pkg/host"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/util/glob"
)

var log = logging.Logger("orbis/bulletin/p2p")

const (
	ProtocolID = "/orbis-bulletion/1.0.0"
	name       = "p2pbb"

	readMessageType     = "read"
	postMessageType     = "post"
	responseMessageType = "response"
	queryMessageType    = "query"
)

const (
	// wait a random amount of time from this interval
	// before dialing peers or reconnecting to help prevent DoS
	dialRandomizerIntervalMilliseconds = 3000

	// repeatedly try to reconnect for a few minutes
	// ie. 5 * 20 = 100s
	reconnectAttempts = 20
	reconnectInterval = 5 * time.Second

	// then move into exponential backoff mode for ~1day
	// ie. 3**10 = 16hrs
	reconnectBackOffAttempts    = 10
	reconnectBackOffBaseSeconds = 3

	readTimeout     = 10 * time.Second
	netQueryTimeout = 10 * time.Second

	queryResponseBuffer = 100
)

var _ bulletin.Bulletin = (*Bulletin)(nil)

type Message = gossipbulletinv1alpha1.Message

type Bulletin struct {
	h   *host.Host
	mem *memmap.Bulletin
	ctx context.Context

	bus eventbus.Bus

	topics map[string]*rpc.Topic

	reonnecting     sync.Map
	persistentPeers map[peer.ID]peer.AddrInfo
}

func New(ctx context.Context, host *host.Host, cfg config.Bulletin) (*Bulletin, error) {
	bus := eventbus.NewBus()
	bb := &Bulletin{
		h:               host,
		ctx:             ctx,
		topics:          make(map[string]*rpc.Topic),
		bus:             bus,
		mem:             memmap.New(memmap.WithBus(bus)),
		persistentPeers: make(map[peer.ID]peer.AddrInfo),
	}

	host.SetStreamHandler(ProtocolID, bb.HandleStream)

	err := host.Discover(ctx, cfg.P2P.Rendezvous)
	if err != nil {
		return nil, fmt.Errorf("discover: %w", err)
	}

	// parse persistent peers
	for _, pstr := range strings.Split(cfg.P2P.PersistentPeers, ",") {
		if pstr == "" {
			continue
		}

		pma, err := ma.NewMultiaddr(strings.TrimSpace(pstr))
		if err != nil {
			return nil, fmt.Errorf("parse persistent peer: %w", err)
		}

		paddr, err := peer.AddrInfoFromP2pAddr(pma)
		if err != nil {
			return nil, fmt.Errorf("convert multiaddr to peer addr: %w", err)
		}

		bb.persistentPeers[paddr.ID] = *paddr
	}

	go bb.maintainPeers(ctx)

	return bb, nil
}

func (bb *Bulletin) Name() string {
	return name
}

func (bb *Bulletin) Init(ctx context.Context) error {
	return nil
}

// Register a namespace for this bulletin
func (bb *Bulletin) Register(ctx context.Context, namespace string) error {
	if namespace == "" {
		return bulletin.ErrEmptyNamespace
	}

	if _, exists := bb.topics[namespace]; exists {
		return bulletin.ErrDuplicateTopic
	}

	topic, err := rpc.NewTopic(ctx, bb.h.PubSub(), bb.h.ID(), namespace, true)
	if err != nil {
		return fmt.Errorf("create new topic: %w", err)
	}

	bb.topics[namespace] = topic
	topic.SetMessageHandler(bb.topicMessageHandler)

	return nil
}

func (bb *Bulletin) findTopicForMessageID(id string) (string, *rpc.Topic) {
	for name, topic := range bb.topics {
		if strings.HasPrefix(id, name) {
			return name, topic
		}
	}

	return "", nil
}

func (bb *Bulletin) findTopicForNamespace(namespace string) (string, *rpc.Topic) {
	for name, topic := range bb.topics {
		if namespace == name {
			return name, topic
		}
	}

	return "", nil
}

func (bb *Bulletin) Post(ctx context.Context, namespace, id string, msg *transport.Message) (bulletin.Response, error) {
	resp, err := bb.mem.Post(ctx, namespace, id, msg)
	if err != nil {
		return bulletin.Response{}, fmt.Errorf("post to local store: %w", err)
	}

	// gossip
	name, topic := bb.findTopicForNamespace(namespace)
	if topic == nil {
		return bulletin.Response{}, bulletin.ErrTopicNotFound
	}
	log.Debugf("Publising post on topic: %s", name)

	payload, err := proto.Marshal(msg)
	if err != nil {
		return bulletin.Response{}, fmt.Errorf("marshal post message payload: %w", err)
	}

	bbMessage := &Message{
		Type:      postMessageType,
		Namespace: namespace,
		Id:        id,
		Payload:   payload,
	}
	msgbuf, err := proto.Marshal(bbMessage)
	if err != nil {
		return bulletin.Response{}, fmt.Errorf("marshal post message: %w", err)
	}

	_, err = topic.Publish(ctx, msgbuf, rpc.WithIgnoreResponse(true))
	if err != nil {
		return bulletin.Response{}, fmt.Errorf("publish post on: %w", err)
	}

	return resp, nil
}

func (bb *Bulletin) Peers(topic string) []peer.ID {
	return bb.h.PubSub().ListPeers(topic)
}

func (bb *Bulletin) Read(ctx context.Context, namespace, id string) (bulletin.Response, error) {
	// check if the read key is in our local store, otherwise ask the network
	resp, err := bb.mem.Read(ctx, namespace, id)
	if errors.Is(err, bulletin.ErrMessageNotFound) {
		log.Debugf("not found locally, fetching from pubsub")

		name, topic := bb.findTopicForNamespace(namespace)
		if topic == nil {
			return bulletin.Response{}, bulletin.ErrTopicNotFound
		}
		log.Debugf("publishing read request on topic: %s", name)

		buf, err := proto.Marshal(&Message{
			Type:      readMessageType,
			Namespace: namespace,
			Id:        id,
		})
		if err != nil {
			return bulletin.Response{}, fmt.Errorf("marshal read message: %w", err)
		}

		// check or set timeout on context
		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithDeadline(ctx, time.Now().Add(readTimeout))
			defer cancel()
		}

		respCh, err := topic.Publish(ctx, buf, rpc.WithIgnoreResponse(false))
		if err != nil {
			return bulletin.Response{}, fmt.Errorf("publish read request: %w", err)
		}

		select {
		case r := <-respCh:
			if r.Err != nil {
				return bulletin.Response{}, fmt.Errorf("read request response: %w", r.Err)
			}

			msg := new(Message)
			err := proto.Unmarshal(r.Data, msg)
			if err != nil {
				return bulletin.Response{}, fmt.Errorf("unmarshal response: %w", err)
			}
			if msg.Type != responseMessageType {
				return bulletin.Response{}, bulletin.ErrBadResponseType
			}

			tMsg := new(transport.Message)
			err = proto.Unmarshal(msg.Payload, tMsg)
			if err != nil {
				return bulletin.Response{}, fmt.Errorf("unmarshal message payload: %w", err)
			}

			return bulletin.Response{
				Data: tMsg,
			}, nil
		case <-ctx.Done():
			return bulletin.Response{}, bulletin.ErrReadTimeout
		}

	} else if err != nil {
		return bulletin.Response{}, fmt.Errorf("read from local store: %w", err)
	}

	return resp, nil
}

type queryConfig struct {
	timeout time.Time
}

type QueryOption func(*queryConfig) error

func WithTimeFilter(t time.Time) QueryOption {
	return func(q *queryConfig) error {
		q.timeout = t
		return nil
	}
}

// Query
// TODO? Options for enable/disable net query?
func (bb *Bulletin) Query(ctx context.Context, namespace, query string) (<-chan bulletin.QueryResponse, error) {
	// q := &queryConfig{}

	if query == "" {
		return nil, fmt.Errorf("query can't be empty")
	}

	// dedicate response channel so we can merge
	respCh := make(chan bulletin.QueryResponse, queryResponseBuffer)

	// local query
	localRespCh, err := bb.mem.Query(ctx, namespace, query)
	if err != nil {
		return nil, fmt.Errorf("query local store: %w", err)
	}

	// forward
	for resp := range localRespCh {
		respCh <- resp
	}

	// net query
	bbMessage := &Message{
		Type:      queryMessageType,
		Namespace: namespace,
		Payload:   []byte(query),
	}
	msgbuf, err := proto.Marshal(bbMessage)
	if err != nil {
		return nil, fmt.Errorf("marshal bulletin message: %s", err)
	}

	var wg sync.WaitGroup
	for name, topic := range bb.topics {
		// is topic related to query?
		// ex
		// topic name: "/ring/123"
		// query: *, /ring/123/dkg/0*
		//
		// the topic matches the glob pattern
		if !glob.Glob(namespace, name) {
			continue
		}

		wg.Add(1)
		go func(topic *rpc.Topic) {
			ctx, cancel := context.WithTimeout(ctx, netQueryTimeout)
			defer cancel()
			p2pRespCh, err := topic.Publish(ctx, msgbuf, rpc.WithMultiResponse(true))
			if err != nil {
				log.Errorf("Failed to publish net query request: %s", err)
			}

			// consume p2pRespCh, read into local store
			// if we already have it, ignore
			log.Infof("Waiting for responses on query topic")
			for resp := range p2pRespCh {
				log.Infof("Got response on query topic")
				if resp.Err != nil {
					log.Errorf("Net query request event: %s", resp.Err)
					continue
				}

				bbMessage := new(Message)
				err := proto.Unmarshal(resp.Data, bbMessage)
				if err != nil {
					log.Errorf("Unmarshal query message: %s", err)
					continue
				}

				if bbMessage.Type != responseMessageType {
					continue
				}

				if bb.mem.Has(ctx, bbMessage.Namespace, bbMessage.Id) {
					continue
				}

				tMsg := new(transport.Message)
				err = proto.Unmarshal(bbMessage.Payload, tMsg)
				if err != nil {
					log.Errorf("Unmarshal query message payload: %s", err)
					continue
				}

				// copy into our local bulletin
				identifer := bbMessage.Namespace + bbMessage.Id
				localResp, err := bb.mem.PostByString(ctx, identifer, tMsg, false)
				if err != nil {
					log.Errorf("Post query message: %s", err)
					continue
				}

				respCh <- bulletin.QueryResponse{
					Resp: localResp,
				}
			}
			wg.Done()
		}(topic)
	}

	// wait until all our outstanding net queries are completed
	// before closing the response channel
	go func() {
		wg.Wait()
		close(respCh)
	}()

	return respCh, nil

}

func (bb *Bulletin) Verify(context.Context, bulletin.Proof, string, bulletin.Message) bool {
	return true
}

// Events
func (bb *Bulletin) Events() eventbus.Bus {
	return bb.mem.Events()
}

func (bb *Bulletin) HandleStream(stream libp2pnetwork.Stream) {
	log.Infof("Received stream: %s", stream.Conn().RemotePeer())
}

func (bb *Bulletin) topicMessageHandler(from peer.ID, topic string, msg []byte) ([]byte, error) {
	log.Debugf("Handling topic %s message from %s", topic, from)
	bbMessage := new(Message)
	err := proto.Unmarshal(msg, bbMessage)
	if err != nil {
		return nil, fmt.Errorf("unmarshal topic message: %w", err)
	}

	var messageResponse []byte

	switch bbMessage.Type {
	case postMessageType:
		log.Debugf("Handling topic message as post request")
		// store posted message, no response necessary
		tMsg := new(transport.Message)
		err = proto.Unmarshal(bbMessage.Payload, tMsg)
		if err != nil {
			return nil, fmt.Errorf("unmarshal post message payload: %w", err)
		}

		identifier := bbMessage.Namespace + bbMessage.Id
		_, err = bb.mem.PostByString(bb.ctx, identifier, tMsg, true)
		if err != nil {
			return nil, fmt.Errorf("post message to local store: %w", err)
		}

	case readMessageType:
		log.Debug("Handling topic message as read request")
		identifier := bbMessage.Namespace + bbMessage.Id
		resp, err := bb.mem.ReadByString(bb.ctx, identifier)
		if err != nil {
			return nil, fmt.Errorf("read message from local store: %w", err)
		}

		tBuf, err := proto.Marshal(resp.Data)
		if err != nil {
			return nil, fmt.Errorf("marshal read message response: %w", err)
		}

		buf, err := proto.Marshal(&Message{
			Type:    responseMessageType,
			Payload: tBuf,
		})
		if err != nil {
			return nil, fmt.Errorf("marshal read message: %w", err)
		}

		messageResponse = buf
	case queryMessageType:
		log.Debug("handling topic message as query request")
		respCh, err := bb.mem.Query(bb.ctx, bbMessage.Namespace, string(bbMessage.Payload))
		if err != nil {
			return nil, fmt.Errorf("query local store: %w", err)
		}

		t := bb.topics[topic]
		for resp := range respCh {
			if resp.Err != nil {
				return nil, fmt.Errorf("query response: %w", resp.Err)
			}

			// original message CID for identifiying responses
			msgID := cid.NewCidV1(cid.Raw, util.Hash(msg))

			tBuf, err := proto.Marshal(resp.Resp.Data)
			if err != nil {
				return nil, fmt.Errorf("unmarshal query response data: %w", err)
			}

			buf, err := proto.Marshal(&Message{
				Type:    responseMessageType,
				Payload: tBuf,
				Id:      resp.Resp.ID,
			})
			if err != nil {
				return nil, fmt.Errorf("marshal query response: %w", err)
			}

			// manually publish response instead of returning via messageReponse var
			// so that we can have multiple (streamed) responses
			t.PublishResponse(from, msgID, buf, nil)
		}

	default:
		log.Warnf("Received unknown message type '%s' on topic %s from %s", bbMessage.Type, topic, from)
		return nil, nil // ignore for now
	}

	return messageResponse, nil
}

func (b *Bulletin) maintainPeers(ctx context.Context) {
	go func() {
		for _, p := range b.persistentPeers {
			go b.reconnectToPeer(ctx, p)
		}
	}()

	subCh, err := b.h.EventBus().Subscribe(new(event.EvtPeerConnectednessChanged))
	if err != nil {
		log.Fatalf("Error subscribing to peer connectedness changes: %s", err)
	}
	defer subCh.Close()

	for {
		select {
		case ev, ok := <-subCh.Out():
			if !ok {
				return
			}

			evt := ev.(event.EvtPeerConnectednessChanged)
			if evt.Connectedness != network.NotConnected {
				continue
			}

			if _, ok := b.persistentPeers[evt.Peer]; !ok {
				continue
			}

			paddr := b.persistentPeers[evt.Peer]
			go b.reconnectToPeer(ctx, paddr)

		case <-ctx.Done():
			return
		}
	}
}

func (b *Bulletin) reconnectToPeer(ctx context.Context, paddr peer.AddrInfo) {
	if _, ok := b.reonnecting.Load(paddr.ID.String()); ok {
		log.Debug("duplicate peer maintainence goroutine", paddr.ID)
		return
	}

	b.reonnecting.Store(paddr.ID.String(), struct{}{})
	defer b.reonnecting.Delete(paddr.ID.String())

	start := time.Now()
	log.Infof("Reconnecting to peer %s", paddr)
	for i := 0; i < reconnectAttempts; i++ {
		select {
		case <-ctx.Done():
			log.Debug("peer maintainence goroutine context finished", paddr.ID)
			return
		default:
			// noop fallthrough
		}

		err := b.h.Connect(ctx, paddr)
		if err == nil {
			log.Infof("reconnected to peer %s during regular backoff", paddr.ID)
			return //success
		}

		log.Infof("Error reconnecting to peer %s: %s, Retrying %d/%d attemps", paddr, err, i, reconnectAttempts)
		randomSleep(reconnectInterval)
	}

	log.Errorf("Failed to reconnect to peer %s. Beginning exponential backoff. Elapsed %s", paddr, time.Since(start))
	for i := 0; i < reconnectBackOffAttempts; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			// noop fallthrough
		}

		// sleep an exponentially increasing amount
		sleepIntervalSeconds := math.Pow(reconnectBackOffBaseSeconds, float64(i))
		randomSleep(time.Duration(sleepIntervalSeconds) * time.Second)

		err := b.h.Connect(ctx, paddr)
		if err == nil {
			log.Infof("reconnected to peer %s during exponential backoff", paddr.ID)
			return //success
		}

		log.Infof("Error reconnecting to peer %s: %s, Retrying %d/%d attemps", paddr, err, i, reconnectAttempts)
	}
	log.Errorf("Failed to reconnect to peer %s. Giving up. Elapsed %s", paddr, time.Since(start))
}

func (bb *Bulletin) Host() *host.Host {
	return bb.h
}

func randomSleep(interval time.Duration) {
	r := time.Duration(rand.Int63n(dialRandomizerIntervalMilliseconds)) * time.Millisecond
	time.Sleep(r + interval)
}

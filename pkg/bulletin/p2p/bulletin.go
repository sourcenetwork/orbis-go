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

	"github.com/ipfs/go-cid"
	util "github.com/ipfs/go-ipfs-util"
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

const ()

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
	host.Discover(ctx, cfg.Rendezvous)

	// parse persistent peers
	for _, pstr := range strings.Split(cfg.PersistentPeers, ",") {
		if pstr == "" {
			continue
		}
		pma, err := ma.NewMultiaddr(strings.TrimSpace(pstr))
		if err != nil {
			return nil, fmt.Errorf("failed to parse persistent peer: %w", err)
		}
		paddr, err := peer.AddrInfoFromP2pAddr(pma)
		if err != nil {
			return nil, fmt.Errorf("failed to convert multiaddr to peer addr: %w", err)
		}
		bb.persistentPeers[paddr.ID] = *paddr
	}

	go bb.maintainPeers(ctx)

	return bb, nil
}

func (bb *Bulletin) Name() string {
	return name
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
		return err
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

func (bb *Bulletin) Post(ctx context.Context, id string, msg *transport.Message) (bulletin.Response, error) {
	resp, err := bb.mem.Post(ctx, id, msg)
	if err != nil {
		return bulletin.Response{}, err
	}

	// gossip
	name, topic := bb.findTopicForMessageID(id)
	if topic == nil {
		return bulletin.Response{}, bulletin.ErrTopicNotFound
	}
	log.Debug("bulletin post: publising post request on topic:", name)

	buf, err := proto.Marshal(msg)
	if err != nil {
		return bulletin.Response{}, err
	}

	bbMessage := &Message{
		Type:    postMessageType,
		Id:      id,
		Payload: buf,
	}
	msgbuf, err := proto.Marshal(bbMessage)
	if err != nil {
		return bulletin.Response{}, err
	}
	if _, err := topic.Publish(ctx, msgbuf, rpc.WithIgnoreResponse(true)); err != nil {
		return bulletin.Response{}, err
	}
	return resp, nil
}

func (bb *Bulletin) Read(ctx context.Context, id string) (bulletin.Response, error) {
	// check if the read key is in our local store, otherwise ask the network
	log.Debug("bulletin read:", id)
	resp, err := bb.mem.Read(ctx, id)
	if errors.Is(err, bulletin.ErrMessageNotFound) {
		log.Debug("bulletin read: not found locally, fetching from pubsub")

		name, topic := bb.findTopicForMessageID(id)
		if topic == nil {
			return bulletin.Response{}, bulletin.ErrTopicNotFound
		}
		log.Debug("bulletin read: publishing read request on topic:", name)

		buf, err := proto.Marshal(&Message{
			Type: readMessageType,
			Id:   id,
		})
		if err != nil {
			return bulletin.Response{}, err
		}

		// check or set timeout on context
		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithDeadline(ctx, time.Now().Add(readTimeout))
			defer cancel()
		}

		respCh, err := topic.Publish(ctx, buf, rpc.WithIgnoreResponse(false))
		if err != nil {
			return bulletin.Response{}, nil
		}

		select {
		case r := <-respCh:
			if r.Err != nil {
				return bulletin.Response{}, err
			}

			msg := new(Message)
			if err := proto.Unmarshal(r.Data, msg); err != nil {
				return bulletin.Response{}, err
			}
			if msg.Type != responseMessageType {
				return bulletin.Response{}, bulletin.ErrBadResponseType
			}

			tMsg := new(transport.Message)
			proto.Unmarshal(msg.Payload, tMsg)

			return bulletin.Response{
				Data: tMsg,
			}, nil
		case <-ctx.Done():
			return bulletin.Response{}, bulletin.ErrReadTimeout
		}

	} else if err != nil {
		return bulletin.Response{}, err
	}

	return resp, nil
}

// Query
// TODO? Options for enable/disable net query?
func (bb *Bulletin) Query(ctx context.Context, query string) (<-chan bulletin.QueryResponse, error) {
	if query == "" {
		return nil, fmt.Errorf("query can't be empty")
	}

	// dedicate response channel so we can merge
	respCh := make(chan bulletin.QueryResponse, queryResponseBuffer)

	// local query
	localRespCh, err := bb.mem.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	// forward
	for resp := range localRespCh {
		respCh <- resp
	}

	// net query
	bbMessage := &Message{
		Type:    queryMessageType,
		Payload: []byte(query),
	}
	msgbuf, err := proto.Marshal(bbMessage)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for name, topic := range bb.topics {
		// is topic related to query?
		// ex
		// topic name: "/ring/123"
		// query: *, /ring/123/dkg/0*
		//
		// either the topic matches the glob pattern, or its a prefix of the glob pattern
		if !glob.Glob(query, name) && !strings.HasPrefix(query, name) {
			continue
		}

		wg.Add(1)
		go func(topic *rpc.Topic) {
			ctx, cancel := context.WithTimeout(ctx, netQueryTimeout)
			defer cancel()
			p2pRespCh, err := topic.Publish(ctx, msgbuf, rpc.WithMultiResponse(true))
			if err != nil {
				log.Error("failed to publish net query request", err)
			}

			// consume p2pRespCh, read into local store
			// if we already have it, ignore
			fmt.Println("waiting for responses on query topic")
			for resp := range p2pRespCh {
				fmt.Println("got response on query topic")
				if resp.Err != nil {
					log.Error("error on net query request event ", resp.Err)
					continue
				}

				bbMessage := new(Message)
				err := proto.Unmarshal(resp.Data, bbMessage)
				if err != nil {
					log.Error("p2p bulletin net query: proto unmarshal bulletin message:", err)
					continue
				}

				if bbMessage.Type != responseMessageType {
					continue
				}

				if bb.mem.Has(ctx, bbMessage.Id) {
					continue
				}

				tMsg := new(transport.Message)
				err = proto.Unmarshal(bbMessage.Payload, tMsg)
				if err != nil {
					log.Error("p2p bulletin net query: proto unmarshal transport message:", err)
					continue
				}

				// copy into our local bulletin
				localResp, err := bb.mem.PostByString(ctx, bbMessage.Id, tMsg, false)
				if err != nil {
					log.Error("p2p bulletin net query: post:", err)
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
	log.Infof("Received stream: %s", stream.Conn().RemotePeer().String())
}

func (bb *Bulletin) topicMessageHandler(from peer.ID, topic string, msg []byte) ([]byte, error) {
	log.Debugf("handling topic %s message from %s", topic, from)
	bbMessage := new(Message)
	err := proto.Unmarshal(msg, bbMessage)
	if err != nil {
		return nil, err
	}

	var messageResponse []byte

	switch bbMessage.Type {
	case postMessageType:
		log.Debug("handling topic message as post request")
		// store posted message, no response necessary
		tMsg := new(transport.Message)
		err = proto.Unmarshal(bbMessage.Payload, tMsg)
		if err != nil {
			return nil, err
		}

		_, err = bb.mem.PostByString(bb.ctx, bbMessage.Id, tMsg, true)
		if err != nil {
			return nil, err
		}

	case readMessageType:
		log.Debug("handling topic message as read request")
		resp, err := bb.mem.ReadByString(bb.ctx, bbMessage.Id)
		if err != nil {
			return nil, err
		}

		tBuf, err := proto.Marshal(resp.Data)
		if err != nil {
			return nil, err
		}

		buf, err := proto.Marshal(&Message{
			Type:    responseMessageType,
			Payload: tBuf,
		})
		if err != nil {
			return nil, err
		}
		messageResponse = buf
	case queryMessageType:
		log.Debug("handling topic message as query request")
		respCh, err := bb.mem.Query(bb.ctx, string(bbMessage.Payload))
		if err != nil {
			return nil, err
		}

		t := bb.topics[topic]
		for resp := range respCh {
			if resp.Err != nil {
				return nil, fmt.Errorf("local query response error: %w", resp.Err)
			}

			// original message CID for identifiying responses
			msgID := cid.NewCidV1(cid.Raw, util.Hash(msg))

			tBuf, err := proto.Marshal(resp.Resp.Data)
			if err != nil {
				return nil, fmt.Errorf("local query response: proto unmarshal: %w", err)
			}
			buf, err := proto.Marshal(&Message{
				Type:    responseMessageType,
				Payload: tBuf,
				Id:      resp.Resp.ID,
			})
			if err != nil {
				return nil, err
			}

			// manually publish response instead of returning via messageReponse var
			// so that we can have multiple (streamed) responses
			t.PublishResponse(from, msgID, buf, nil)
		}

	default:
		log.Warn("received unknown message type '%s' on topic %s from %s", bbMessage.Type, topic, from)
		return nil, nil // ignore for now
	}

	return messageResponse, nil
}

func (b *Bulletin) maintainPeers(ctx context.Context) {
	go func() {
		for _, p := range b.persistentPeers {
			b.h.Connect(ctx, p)
		}
	}()

	subCh, err := b.h.EventBus().Subscribe(new(event.EvtPeerConnectednessChanged))
	if err != nil {
		panic(err)
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

			go b.reconnectToPeer(ctx, evt.Peer)

		case <-ctx.Done():
			return
		}
	}
}

func (b *Bulletin) reconnectToPeer(ctx context.Context, pid peer.ID) {
	if _, ok := b.reonnecting.Load(pid.String()); ok {
		return
	}

	b.reonnecting.Store(pid.String(), struct{}{})
	defer b.reonnecting.Delete(pid.String())

	paddr := b.persistentPeers[pid]

	start := time.Now()
	log.Info("Reconnecting to peer %s", paddr)
	for i := 0; i < reconnectAttempts; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			// noop fallthrough
		}

		err := b.h.Connect(ctx, paddr)
		if err == nil {
			return //success
		}

		log.Info("Error reconnecting to peer. Trying again", "tries", i, "err", err, "addr", paddr)
		randomSleep(reconnectInterval)
	}

	log.Error("Failed to reconnect to peer. Beginning exponential backoff", "addr", paddr, "elapsed", time.Since(start))
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
			return //success
		}

		log.Info("Error reconnecting to peer. Trying again", "tries", i, "err", err, "addr", paddr)
	}
	log.Error("Failed to reconnect to peer. Giving up", "addr", paddr, "elapsed", time.Since(start))
}

func (bb *Bulletin) Host() *host.Host {
	return bb.h
}

func randomSleep(interval time.Duration) {
	r := time.Duration(rand.Int63n(dialRandomizerIntervalMilliseconds)) * time.Millisecond
	time.Sleep(r + interval)
}

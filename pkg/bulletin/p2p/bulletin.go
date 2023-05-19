package p2p

import (
	"context"
	"errors"
	"time"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/peer"
	libp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	rpc "github.com/textileio/go-libp2p-pubsub-rpc"
	"google.golang.org/protobuf/proto"

	"github.com/sourcenetwork/orbis-go/config"
	gossipbulletinv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/gossipbulletin/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin/memmap"
	"github.com/sourcenetwork/orbis-go/pkg/host"
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

var _ bulletin.Bulletin = (*Bulletin)(nil)

var (
	readTimeout = time.Second
)

type Message = gossipbulletinv1alpha1.Message

type Bulletin struct {
	h   *host.Host
	mem memmap.Bulletin
	ctx context.Context

	topics map[string]*rpc.Topic
}

func New(ctx context.Context, host *host.Host, cfg config.Bulletin) (*Bulletin, error) {
	bb := &Bulletin{
		h:      host,
		ctx:    ctx,
		topics: make(map[string]*rpc.Topic),
	}

	host.SetStreamHandler(ProtocolID, bb.HandleStream)
	host.Discover(ctx, cfg.Rendezvous)

	return bb, nil
}

func (bb *Bulletin) Name() string {
	return name
}

func (bb *Bulletin) Register(ctx context.Context, namespace string) error {
	if namespace == "" {
		return bulletin.ErrEmptyNamespace
	}

	if _, exists := bb.topics[namespace]; exists {
		return bulletin.ErrDuplicateTopic
	}
	topic, err := rpc.NewTopic(ctx, bb.h.PubSub(), peer.ID(bb.h.ID()), namespace, true)
	if err != nil {
		return err
	}

	bb.topics[namespace] = topic
	topic.SetMessageHandler(bb.topicMessageHandler)
	return nil
}

func (bb *Bulletin) Post(ctx context.Context, id bulletin.ID, msg bulletin.Message) (bulletin.Response, error) {
	resp, err := bb.mem.Post(ctx, id, msg)
	if err != nil {
		return bulletin.Response{}, err
	}

	// gossip
	topic, exists := bb.topics[id.ServiceNamespace()]
	if !exists {
		return bulletin.Response{}, bulletin.ErrTopicNotFound
	}

	bbMessage := &Message{
		Type:    postMessageType,
		Id:      id.String(),
		Payload: msg,
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

func (bb *Bulletin) Read(ctx context.Context, key bulletin.ID) (bulletin.Response, error) {
	// check if the read key is in our local store, otherwise ask the network
	resp, err := bb.mem.Read(ctx, key)
	if errors.Is(err, bulletin.ErrMessageNotFound) {

		topic, exists := bb.topics[key.ServiceNamespace()]
		if !exists {
			return bulletin.Response{}, bulletin.ErrTopicNotFound
		}

		buf, err := proto.Marshal(&Message{
			Type: readMessageType,
		})
		if err != nil {
			return bulletin.Response{}, err
		}

		// check or set timeout on context
		if _, ok := ctx.Deadline(); !ok {
			ctx, _ = context.WithDeadline(ctx, time.Now().Add(readTimeout))
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

			var msg *Message
			if err := proto.Unmarshal(r.Data, msg); err != nil {
				return bulletin.Response{}, err
			}
			if msg.Type != responseMessageType {
				return bulletin.Response{}, bulletin.ErrBadResponseType
			}

			return bulletin.Response{
				Data: msg.Payload,
			}, nil
		case <-ctx.Done():
			return bulletin.Response{}, bulletin.ErrReadTimeout
		}

	} else if err != nil {
		return bulletin.Response{}, err
	}

	return resp, nil
}

func (bb *Bulletin) Query(ctx context.Context, query string) ([]bulletin.Response, error) {
	panic("implement me")
}

func (bb *Bulletin) Verify(context.Context, bulletin.Proof, string, bulletin.Message) bool {
	return true
}

// EventBus
// Events() eventbus.Bus

func (bb *Bulletin) HandleStream(stream libp2pnetwork.Stream) {
	log.Infof("Received stream: %s", stream.Conn().RemotePeer().Pretty())
}

func (bb *Bulletin) topicMessageHandler(from peer.ID, topic string, msg []byte) ([]byte, error) {
	var bbMessage *Message
	err := proto.Unmarshal(msg, bbMessage)
	if err != nil {
		return nil, err
	}

	switch bbMessage.Type {
	case postMessageType:
		// store posted message, no response necessary
		_, err := bb.mem.PostByString(bb.ctx, bbMessage.Id, bbMessage.Payload)
		if err != nil {
			return nil, err
		}

	case readMessageType:
		resp, err := bb.mem.ReadByString(bb.ctx, bbMessage.Id)
		if err != nil {
			return nil, err
		}

		buf, err := proto.Marshal(&Message{
			Type:    responseMessageType,
			Payload: resp.Data,
		})
		if err != nil {
			return nil, err
		}
		return buf, nil
	default:
		log.Warn("received unknown message type '%s' on topic %s from %s", bbMessage.Type, topic, from)
		return nil, nil // ignore for now
	}

	return nil, nil // unreachable due to default case in switch
}

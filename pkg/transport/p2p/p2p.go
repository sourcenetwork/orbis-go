package p2p

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/pkg/host"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
	"google.golang.org/protobuf/proto"
)

var log = logging.Logger("orbis/transport/p2p")

var (
	_ transport.Transport = (*Transport)(nil)
)

const (
	ProtocolID = "/orbis-transport/1.0.0"
	name       = "p2p"
)

type Transport struct {
	h *host.Host
}

func New(ctx context.Context, host *host.Host, cfg config.Transport) (*Transport, error) {
	return &Transport{h: host}, nil
}

func (t *Transport) Name() string {
	return name
}

func (t *Transport) Send(ctx context.Context, node transport.Node, msg *transport.Message) error {

	// todo: telemetry
	// todo: verify msg is of type p2p.message
	// todo sign message

	peerID := peer.ID(node.ID())
	// todo protocol formatting
	protocolID := protocol.ConvertFromStrings([]string{msg.GetType()})
	stream, err := t.h.NewStream(ctx, peerID, protocolID...)
	if err != nil {
		return fmt.Errorf("open stream: %w", err)
	}
	defer stream.Close()

	buf, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	_, err = stream.Write(buf)
	if err != nil {
		return fmt.Errorf("write stream: %w", err)
	}

	return nil
}

func (t *Transport) Gossip(ctx context.Context, topic string, msg *transport.Message) error {
	panic("not implemented") // TODO: Implement
}

func (t *Transport) Connect(ctx context.Context, node transport.Node) error {

	id, err := peer.Decode(node.ID())
	if err != nil {
		return fmt.Errorf("decode peer id: %w", err)
	}

	pi := peer.AddrInfo{
		ID:    id,
		Addrs: []ma.Multiaddr{node.Address()},
	}

	err = t.h.Connect(ctx, pi)
	if err != nil {
		return fmt.Errorf("connect to peer: %w", err)
	}

	return err
}

func (t *Transport) Host() transport.Host {
	return &Host{t.h}
}

func (t *Transport) NewMessage(rid types.RingID, id string, gossip bool, payload []byte, msgType string) (*transport.Message, error) {

	pubkeyBytes, err := t.Host().PublicKey().Raw()
	if err != nil {
		return nil, fmt.Errorf("get raw public key: %w", err)
	}

	// todo: Signature (should be done on send)
	// replay? nonce?
	return &transport.Message{
		Timestamp:  time.Now().Unix(),
		Id:         id,
		RingId:     string(rid),
		NodeId:     t.Host().ID(),
		NodePubKey: pubkeyBytes,
		Type:       msgType,
		Payload:    payload,
		Gossip:     gossip,
	}, nil
}

func (t *Transport) AddHandler(pid protocol.ID, handler transport.Handler) {
	streamHandler := streamHandlerFrom(handler)
	t.h.SetStreamHandler(pid, streamHandler)
}

func (t *Transport) RemoveHandler(pid protocol.ID) {
	t.h.RemoveStreamHandler(pid)
}

func streamHandlerFrom(handler transport.Handler) func(network.Stream) {
	return func(stream network.Stream) {

		log.Infof("new stream from %s", stream.Conn().RemotePeer())

		buf, err := io.ReadAll(stream)
		if err != nil {
			if err != io.EOF {
				log.Errorf("read stream: %s", err)
			}
			err = stream.Reset()
			if err != nil {
				log.Errorf("reset stream: %s", err)
			}
			return
		}

		err = stream.Close()
		if err != nil {
			log.Errorf("close stream: %s", err)
			return
		}

		data := &transport.Message{}
		err = proto.Unmarshal(buf, data)
		if err != nil {
			log.Errorf("unmarshal data: %s", err)
			return
		}

		err = handler(data)
		if err != nil {
			log.Errorf("handle data: %s", err)
			return
		}
	}
}

package p2p

import (
	"context"
	"io"
	"time"

	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"

	libp2phost "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
	"google.golang.org/protobuf/proto"
)

type Transport struct {
	h libp2phost.Host
}

func NewTransport(h libp2phost.Host) transport.Transport {
	return &Transport{
		h: h,
	}
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
		return err // todo: wrap
	}
	defer stream.Close()

	buf, err := proto.Marshal(msg)
	if err != nil {
		return err // todo: wrap
	}

	_, err = stream.Write(buf)
	return err
}

func (t *Transport) Gossip(ctx context.Context, topic string, msg *transport.Message) error {
	panic("not implemented") // TODO: Implement
}

func (t *Transport) Connect(ctx context.Context, node transport.Node) error {
	pi := peer.AddrInfo{
		ID:    peer.ID(node.ID()),
		Addrs: []ma.Multiaddr{node.Address()},
	}
	return t.h.Connect(ctx, pi)
}

func (t *Transport) Host() transport.Host {
	return t.host()
}

func (t *Transport) NewMessage(ringID types.RingID, id string, gossip bool, payload []byte, msgType string) (*transport.Message, error) {
	h := t.host()
	pubkeyBytes, err := h.PublicKey().Raw()
	if err != nil {
		return nil, err // todo: wrap
	}
	// todo: Signature (should be done on send)
	// replay? nonce?
	return &transport.Message{
		Timestamp:  time.Now().Unix(),
		Id:         id,
		RingId:     string(ringID),
		NodeId:     h.ID(),
		NodePubKey: pubkeyBytes,
		Type:       msgType,
		Payload:    payload,
		Gossip:     gossip,
	}, nil
}

func (t *Transport) host() *host {
	return &host{t.h}
}

func (t *Transport) node() node {
	return node{
		id:        t.h.ID(),
		publicKey: t.h.Peerstore().PubKey(t.h.ID()),
		address:   t.h.Addrs()[0],
	}
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
		// parse message via protobuf
		// send to handler
		data := &transport.Message{}
		buf, err := io.ReadAll(stream)
		if err != nil {
			stream.Reset()
			// todo: log err
			return
		}
		stream.Close()

		err = proto.Unmarshal(buf, data)
		if err != nil {
			// todo: log err
			return
		}

		handler(data)
	}
}

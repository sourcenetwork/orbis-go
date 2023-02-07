package p2p

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	libp2pHost "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

type p2pTransport struct {
	h libp2pHost.Host
}

func NewTransport(h libp2pHost.Host) transport.Transport {
	return &p2pTransport{
		h: h,
	}
}

func (pt *p2pTransport) Send(ctx context.Context, node transport.Node, msg transport.Message) error {
	// todo: telemetry

	// todo: verify msg is of type p2p.message

	peerID := peer.ID(node.ID())
	protocolID := protocol.ConvertFromStrings([]string{msg.Type()})
	stream, err := pt.h.NewStream(ctx, peerID, protocolID...)
	if err != nil {
		return err // todo: wrap
	}
	defer stream.Close()

	buf, err := msg.Marshal()
	if err != nil {
		return err // todo: wrap
	}

	_, err = stream.Write(buf)
	return err
}

func (pt *p2pTransport) Gossip(ctx context.Context, topic string, msg transport.Message) error {
	panic("not implemented") // TODO: Implement
}

func (pt *p2pTransport) Connect(ctx context.Context, node transport.Node) error {
	pi := peer.AddrInfo{
		ID:    peer.ID(node.ID()),
		Addrs: []ma.Multiaddr{node.Address()},
	}
	return pt.h.Connect(ctx, pi)
}

func (pt *p2pTransport) Host() transport.Host {
	return pt.host()
}

func (pt *p2pTransport) NewMessage(id string, gossip bool, payload []byte, msgType string) (transport.Message, error) {
	h := pt.host()
	signBuf, err := h.Sign(payload)
	if err != nil {
		return transport.Message{}, err // todo: wrap
	}

	return transport.Message{
		Timestamp: uint64(time.Now().Unix()),
		ID:        id,
		Node:      pt.node(),
		Type:      msgType,
		Payload:   payload,
		Signature: signBuf,
		Gossip:    gossip,
	}, nil
}

func (pt *p2pTransport) host() *host {
	return &host{pt.h}
}

func (pt *p2pTransport) node() node {
	return node{
		id:        pt.h.ID(),
		publicKey: pt.h.Peerstore().PubKey(pt.h.ID()),
		address:   pt.h.Addrs()[0],
	}
}

func (pt *p2pTransport) AddHandler(pid protocol.ID, handler transport.Handler) {
	h := streamHandlerFrom(handler)

}

func (pt *p2pTransport) RemoveHandler(pid protocol.ID) {
	panic("not implemented") // TODO: Implement
}

func streamHandlerFrom(handler transport.Handler) func(network.Stream) {

}

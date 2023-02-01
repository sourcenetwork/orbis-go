package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	libp2pHost "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/protocol"

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

	peerID := peer.IDFromString(node.ID())
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
	panic("not implemented") // TODO: Implement
}

func (pt *p2pTransport) Host() transport.Host {
	return hostFromLibP2P(pt.h)
}

func (pt *p2pTransport) NewMessage(id string, gossip bool, payload []byte, msgType string) (transport.Message, error) {
	return pt.newMessage(id, gossip, payload, msgType)
}

func (pt *p2pTransport) newMessage(id string, gossip bool, payload []byte, msgType string) (message, error) {
	return message{}, nil
}

func (pt *p2pTransport) AddHandler(pid protocol.ID, handler transport.Handler) {
	panic("not implemented") // TODO: Implement
}

func (pt *p2pTransport) RemoveHandler(pid protocol.ID) {
	panic("not implemented") // TODO: Implement
}

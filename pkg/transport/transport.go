package transport

import (
	"context"

	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/config"
	transportv1alpha "github.com/sourcenetwork/orbis-go/gen/proto/orbis/transport/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/types"

	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
)

type Message = transportv1alpha.Message

type Handler func(*Message) error

type Transport interface {
	Send(ctx context.Context, node Node, msg *Message) error
	Gossip(ctx context.Context, topic string, msg *Message) error
	Connect(ctx context.Context, node Node) error
	Host() Host
	NewMessage(rid types.RingID, id string, gossip bool, payload []byte, msgType string) (*Message, error)
	AddHandler(pid protocol.ID, handler Handler)
	RemoveHandler(pid protocol.ID)
}

type Node interface {
	ID() string
	PublicKey() crypto.PublicKey
	Address() ma.Multiaddr
}

type Host interface {
	Node
	Sign(msg []byte) ([]byte, error)
}

type Factory interface {
	New(ctx context.Context, inj *do.Injector, cfg config.Transport) (Transport, error)
}

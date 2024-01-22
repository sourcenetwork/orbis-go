package transport

import (
	"context"

	transportv1alpha "github.com/sourcenetwork/orbis-go/gen/proto/orbis/transport/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/types"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
)

type Message = transportv1alpha.Message

type Handler func(*Message) error

type Node interface {
	ID() peer.ID
	PublicKey() crypto.PublicKey
	Address() ma.Multiaddr
}

type Transport interface {
	Name() string
	PublicKey() crypto.PublicKey
	PrivateKey() (crypto.PrivateKey, error)
	NewMessage(rid types.RingID, id string, gossip bool, payload []byte, msgType string, target Node) (*Message, error)
	ID() peer.ID
	Address() ma.Multiaddr
	Send(ctx context.Context, node Node, msg *Message) error
	AddHandler(pid protocol.ID, handler Handler)
	RemoveHandler(pid protocol.ID)
	// P2P
	PubSub() *pubsub.PubSub
	Network() network.Network
}

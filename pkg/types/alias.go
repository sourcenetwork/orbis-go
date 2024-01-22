package types

import (
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"

	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
)

type Ring struct {
	*ringv1alpha1.Ring
}

func RingFromManifest(manifest []byte) (*Ring, RingID, error) {
	return nil, "", nil
}

type Secret struct {
	*ringv1alpha1.Secret
}

type Node struct {
	id        peer.ID
	address   ma.Multiaddr
	publicKey crypto.PublicKey
}

func NewNode(id peer.ID, addr ma.Multiaddr, pk crypto.PublicKey) *Node {
	return &Node{
		id:        id,
		address:   addr,
		publicKey: pk,
	}
}

func (n *Node) ID() peer.ID {
	return n.id
}

func (n *Node) Address() ma.Multiaddr {
	return n.address
}

func (n *Node) PublicKey() crypto.PublicKey {
	return n.publicKey
}

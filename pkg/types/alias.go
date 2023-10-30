package types

import (
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
	index     int // -1 means invalid index
	id        string
	address   ma.Multiaddr
	publicKey crypto.PublicKey
}

func NewNode(idx int, id string, addr ma.Multiaddr, pk crypto.PublicKey) *Node {
	return &Node{
		index:     idx,
		id:        id,
		address:   addr,
		publicKey: pk,
	}
}

func (n *Node) ID() string {
	return n.id
}

func (n *Node) Index() int {
	return n.index
}

func (n *Node) Address() ma.Multiaddr {
	return n.address
}

func (n *Node) PublicKey() crypto.PublicKey {
	return n.publicKey
}

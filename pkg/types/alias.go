package types

import (
	"crypto"

	ma "github.com/multiformats/go-multiaddr"

	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
	secretv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/secret/v1alpha1"
)

type Ring struct {
	ringv1alpha1.Ring
}

func RingFromManifest(manifest []byte) (*Ring, RingID, error) {
	return nil, "", nil
}

type Secret struct {
	secretv1alpha1.Secret
}

type Node struct {
	index int // -1 means invalid index
	ringv1alpha1.Node
}

func (n *Node) Index() int {
	return n.index
}

func (n *Node) Address() (ma.Multiaddr, error) {
	return ma.NewMultiaddr(n.Node.Address)
}

func (n *Node) PublicKey() crypto.PublicKey {
	return nil
}

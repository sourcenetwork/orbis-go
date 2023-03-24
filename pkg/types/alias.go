package types

import (
	"crypto"

	ma "github.com/multiformats/go-multiaddr"

	repov1alpha "github.com/sourcenetwork/orbis-go/gen/proto/orbis/repo/v1alpha"
)

type Ring struct {
	repov1alpha.Ring
}

func RingFromManifest(manifest []byte) (*Ring, RingID, error) {
	return nil, "", nil
}

type Secret struct {
	repov1alpha.Secret
}

type Node struct {
	index int // -1 means invalid index
	repov1alpha.Node
}

func (n *Node) Index() int {
	return n.index
}

func (n *Node) Address() (ma.Multiaddr, error) {
	return ma.NewMultiaddr(n.Node.Address)
}

func (n *Node) PublicKey() crypto.PublicKey {

}

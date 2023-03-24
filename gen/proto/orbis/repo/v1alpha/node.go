package repov1alpha

import (
	ma "github.com/multiformats/go-multiaddr"
)

func (n *Node) Mutliaddr() (ma.Multiaddr, error) {
	return ma.NewMultiaddr(n.Address)
}

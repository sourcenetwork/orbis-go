package p2p

import (
	"github.com/sourcenetwork/orbis-go/pkg/crypto"

	libp2phost "github.com/libp2p/go-libp2p/core/host"
	ma "github.com/multiformats/go-multiaddr"
)

type Host struct {
	host libp2phost.Host
}

func (h *Host) ID() string {
	return h.host.ID().String()
}

func (h *Host) PublicKey() crypto.PublicKey {
	libp2pPubKey := h.host.Peerstore().PubKey(h.host.ID())
	pubkey, _ := crypto.PublicKeyFromLibP2P(libp2pPubKey)
	return pubkey
}

func (h *Host) Address() ma.Multiaddr {
	return h.host.Addrs()[0]
}

func (h *Host) Sign(data []byte) ([]byte, error) {
	key := h.host.Peerstore().PrivKey(h.host.ID())
	res, err := key.Sign(data)
	return res, err
}

func (h *Host) Peers() ([]Node, error) {
	var nodes []Node
	s := h.host.Network().Peerstore()
	for _, p := range h.host.Network().Peers() {
		a := s.PeerInfo(p)
		n := Node{
			id:      p.String(),
			address: a.Addrs[0],
		}
		nodes = append(nodes, n)

	}

	return nodes, nil
}

func (h *Host) Close() error {
	return h.host.Close()
}

type Node struct {
	id        string
	publicKey crypto.PublicKey
	address   ma.Multiaddr
}

func NewNode(id string, publicKey crypto.PublicKey, address ma.Multiaddr) *Node {
	return &Node{
		id:        id,
		publicKey: publicKey,
		address:   address,
	}
}

func (n Node) ID() string {
	return n.id
}

func (n Node) PublicKey() crypto.PublicKey {
	return n.publicKey
}

func (n Node) Address() ma.Multiaddr {
	return n.address
}

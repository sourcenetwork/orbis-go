package p2p

import (
	"github.com/libp2p/go-libp2p-core/peer"
	lcrypto "github.com/libp2p/go-libp2p/core/crypto"
	libp2p "github.com/libp2p/go-libp2p/core/host"
	ma "github.com/multiformats/go-multiaddr"

	ocrypto "github.com/sourcenetwork/orbis-go/pkg/crypto"
)

type host struct {
	node
	host libp2p.Host // libp2p
}

func (h *host) ID() string {
	return h.host.ID().String()
}

func (h *host) PublicKey() ocrypto.PublicKey {
	libp2pPubKey := h.host.Peerstore().PubKey(h.host.ID())
	return ocrypto.PublicKeyFromLibP2P(libp2pPubKey)
}

func (h *host) Address() ma.Multiaddr {
	return h.host.Addrs()[0]
}

func (h *host) Sign(data []byte) ([]byte, error) {
	key := h.host.Peerstore().PrivKey(h.host.ID())
	res, err := key.Sign(data)
	return res, err
}

type node struct {
	id        peer.ID
	publicKey lcrypto.PubKey
	address   ma.Multiaddr
}

func (n *node) ID() string {
	return n.id.String()
}

func (n *node) PublicKey() ocrypto.PublicKey {
	return ocrypto.PublicKeyFromLibP2P(n.publicKey)
}

func (n *node) Address() ma.Multiaddr {
	return n.address
}

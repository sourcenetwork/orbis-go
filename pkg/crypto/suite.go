package crypto

import (
	ic "github.com/libp2p/go-libp2p/core/crypto"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
)

var suite = edwards25519.NewBlakeSHA256Ed25519()

// PublicKey
type PublicKey struct {
	key kyber.Point
}

type PrivateKey struct {
	key kyber.Scalar
}

func PublicKeyFromLibP2P(pubKey ic.PubKey) (PublicKey, error) {
	var pubkey PublicKey
	buf, err := pubKey.Raw()
	if err != nil {
		return pubkey, err
	}

	pubkey.key = suite.Point()
	err = pubkey.key.UnmarshalBinary(buf)

	return pubkey, err
}

// PriShare
type PriShare struct{}

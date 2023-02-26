package crypto

import (
	ic "github.com/libp2p/go-libp2p/core/crypto"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/share"
)

var suite = edwards25519.NewBlakeSHA256Ed25519()

// PublicKey
type PublicKey kyber.Point

type PrivateKey kyber.Scalar

func PublicKeyFromLibP2P(pubKey ic.PubKey) (PublicKey, error) {
	var pubkey PublicKey
	buf, err := pubKey.Raw()
	if err != nil {
		return pubkey, err
	}

	pubkey = suite.Point()
	err = pubkey.UnmarshalBinary(buf)

	return pubkey, err
}

// PriShare
type PriShare share.PriShare

// PubPoly
type PubPoly share.PubPoly

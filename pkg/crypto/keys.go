package crypto

import (
	"crypto/sha512"
	"fmt"

	ic "github.com/libp2p/go-libp2p/core/crypto"
	icpb "github.com/libp2p/go-libp2p/core/crypto/pb"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/suites"
)

var (
	// suite = edwards25519.NewBlakeSHA256Ed25519()
	ErrBadKeyType = fmt.Errorf("bad key type")
)

// PublicKey
type PublicKey interface {
	ic.PubKey
	Point() kyber.Point
}

type pubKey struct {
	ic.PubKey
	suite suites.Suite
}

// PublicKeyFromLibP2P creates a PublicKey from a given
// LibP2P based PubKey
func PublicKeyFromLibP2P(pubkey ic.PubKey) (PublicKey, error) {
	return publicKeyFromLibP2P(pubkey)
}

func PublicKeyFromProto(pk *icpb.PublicKey) (PublicKey, error) {
	icpk, err := ic.PublicKeyFromProto(pk)
	if err != nil {
		return nil, err
	}
	return publicKeyFromLibP2P(icpk)
}

func PublicKeyToProto(pk PublicKey) (*icpb.PublicKey, error) {
	return ic.PublicKeyToProto(pk)
}

func PublicKeyFromPoint(point kyber.Point) (PublicKey, error) {
	panic("todo")
}

func publicKeyFromLibP2P(pubkey ic.PubKey) (*pubKey, error) {
	suite, err := SuiteForType(pubkey.Type())
	if err != nil {
		return nil, err
	}

	return &pubKey{
		PubKey: pubkey,
		suite:  suite,
	}, nil

}

func (p *pubKey) Point() kyber.Point {
	buf, _ := p.PubKey.Raw()
	point := p.suite.Point()
	point.UnmarshalBinary(buf)
	return point
}

type libp2pPrivKey interface {
	ic.Key
	Sign([]byte) ([]byte, error)
}

type PrivateKey interface {
	libp2pPrivKey
	Scalar() kyber.Scalar
	GetPublic() PublicKey
}

type privKey struct {
	ic.PrivKey
	suite suites.Suite
}

func PrivateKeyFromLibP2P(privkey ic.PrivKey) (PrivateKey, error) {
	suite, err := SuiteForType(privkey.Type())
	if err != nil {
		return nil, err
	}

	return &privKey{
		PrivKey: privkey,
		suite:   suite,
	}, nil
}

// Scalar returns a numeric elliptic curve scalar
// representation of the private key.
//
// WARNING: THIS ONLY WORDS WITH Edwards25519 CURVES RIGHT NOW.
func (p *privKey) Scalar() kyber.Scalar {
	// There is a discrepency between LibP2P private keys
	// and "raw" EC scalars. LibP2P private keys is an
	// (x, y) pair, where x is the given "seed" and y is
	// the cooresponding publickey. Where y is computed as
	//
	// h := sha512.Hash(x)
	// s := scalar().SetWithClamp(h)
	// y := point().ScalarBaseMul(x)
	//
	// So to make sure future conversions of this scalar
	// to a public key, like in the DKG setup, we need to
	// convert this key to a scalar using the Hash and Clamp
	// method.
	//
	// To understand clamping, see here:
	// https://neilmadden.blog/2020/05/28/whats-the-curve25519-clamping-all-about/

	buf, err := p.PrivKey.Raw()
	if err != nil {
		panic(err)
	}

	// hash seed and clamp bytes
	digest := sha512.Sum512(buf[:32])
	digest[0] &= 0xf8
	digest[31] &= 0x7f
	digest[31] |= 0x40
	return p.suite.Scalar().SetBytes(digest[:32])
}

func (p *privKey) GetPublic() PublicKey {
	return &pubKey{
		PubKey: p.PrivKey.GetPublic(),
		suite:  p.suite,
	}
}

// PriShare
type PriShare struct {
	*share.PriShare
}

// PubPoly
type PubPoly struct {
	*share.PubPoly
}

func setScalarWithClamp(s kyber.Scalar, buf []byte) {

}

package crypto

import (
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

type PrivateKey interface {
	ic.PrivKey
	Scalar() kyber.Scalar
}

type privKey struct {
	ic.PrivKey
	suite suites.Suite
}

func (p *privKey) Scalar() kyber.Scalar {
	buf, _ := p.PrivKey.Raw()
	scalar := p.suite.Scalar()
	scalar.UnmarshalBinary(buf)
	return scalar
}

func (p *privKey) GetPublic() PublicKey {
	return &pubKey{
		PubKey: p.PrivKey.GetPublic(),
		suite:  p.suite,
	}
}

// PriShare
type PriShare struct {
	share.PriShare
}

// PubPoly
type PubPoly struct {
	*share.PubPoly
}

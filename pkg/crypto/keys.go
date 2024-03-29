package crypto

import (
	gocrypto "crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/sha512"
	"fmt"
	"io"
	"strings"

	ic "github.com/libp2p/go-libp2p/core/crypto"
	icpb "github.com/libp2p/go-libp2p/core/crypto/pb"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/suites"
)

type KeyType = icpb.KeyType

func KeyTypeFromString(keyType string) (KeyType, error) {
	switch strings.ToLower(keyType) {
	case "ed25519":
		return Ed25519, nil
	case "ecdsa":
		return ECDSA, nil
	case "secp256k1":
		return Secp256k1, nil
	default:
		return 0, ErrBadKeyType
	}
}

var (
	Ed25519   = icpb.KeyType_Ed25519
	ECDSA     = icpb.KeyType_ECDSA
	Secp256k1 = icpb.KeyType_Secp256k1
)

// PublicKey
type PublicKey interface {
	ic.PubKey
	Point() kyber.Point
	Std() (gocrypto.PublicKey, error)
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

func PublicKeyFromStdPublicKey(pubkey gocrypto.PublicKey) (PublicKey, error) {
	var icpk ic.PubKey
	var err error
	switch pkt := pubkey.(type) {
	case ed25519.PublicKey:
		icpk, err = ic.UnmarshalEd25519PublicKey(pkt)
	case ecdsa.PublicKey:
		icpk, err = ic.ECDSAPublicKeyFromPubKey(pkt)
	default:
		return nil, fmt.Errorf("unknown key type")
	}

	if err != nil {
		return nil, err
	}
	return publicKeyFromLibP2P(icpk)
}

func PublicKeyFromPoint(suite suites.Suite, point kyber.Point) (PublicKey, error) {

	buf, err := point.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal point: %w", err)
	}

	var pk ic.PubKey

	switch strings.ToLower(suite.String()) {
	case "ed25519":
		pk, err = ic.UnmarshalEd25519PublicKey(buf)
	case "secp256k1":
		pk, err = ic.UnmarshalSecp256k1PublicKey(buf)
	case "ecdsa":
		pk, err = ic.UnmarshalECDSAPublicKey(buf)
	case "rsa":
		pk, err = ic.UnmarshalRsaPublicKey(buf)
	default:
		return nil, ErrBadKeyType
	}

	if err != nil {
		return nil, fmt.Errorf("unmarshal public key: %w", err)
	}

	return PublicKeyFromLibP2P(pk)
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

func PublicKeyToProto(pk PublicKey) (*icpb.PublicKey, error) {
	return ic.PublicKeyToProto(pk)
}

func (p *pubKey) Point() kyber.Point {
	buf, _ := p.PubKey.Raw()
	point := p.suite.Point()
	point.UnmarshalBinary(buf)
	return point
}

func (p *pubKey) Std() (gocrypto.PublicKey, error) {
	return ic.PubKeyToStdKey(p.PubKey)
}

type libp2pPrivKey interface {
	ic.Key
	Sign([]byte) ([]byte, error)
}

type PrivateKey interface {
	libp2pPrivKey
	Scalar() kyber.Scalar
	GetPublic() PublicKey
	// Std() gocrypto.PrivateKey
}

type privKey struct {
	ic.PrivKey
	suite suites.Suite
}

func GenerateKeyPair(ste suites.Suite, src io.Reader) (PrivateKey, PublicKey, error) {
	keyType, err := KeyTypeFromString(ste.String())
	if err != nil {
		return nil, nil, err
	}
	sk, pk, err := ic.GenerateKeyPairWithReader(int(keyType), 0, src)
	if err != nil {
		return nil, nil, err
	}

	return &privKey{
			PrivKey: sk,
			suite:   ste,
		}, &pubKey{
			PubKey: pk,
			suite:  ste,
		}, nil
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

// DistKeyShare
type DistKeyShare struct {
	// Coefficients of the public polynomial holding the public key
	Commits []kyber.Point

	// PriShare of the distributed secret
	PriShare *share.PriShare
}

// PubPoly
type PubPoly struct {
	*share.PubPoly
}

func setScalarWithClamp(s kyber.Scalar, buf []byte) {

}

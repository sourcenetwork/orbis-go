package crypto

import (
	gocrypto "crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	scrypto "github.com/TBD54566975/ssi-sdk/crypto"
	"github.com/TBD54566975/ssi-sdk/did/key"
	cometcrypto "github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/lestrrat-go/jwx/v2/jwk"
	ic "github.com/libp2p/go-libp2p/core/crypto"
	icpb "github.com/libp2p/go-libp2p/core/crypto/pb"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/suites"
	"golang.org/x/crypto/ripemd160"
)

type KeyType int32

func (kt KeyType) String() string {
	if kt == 4 {
		return "octet"
	}
	return icpb.KeyType(kt).String()
}

func KeyPairTypeFromString(keyType string) (KeyType, error) {
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
	Ed25519   = KeyType(icpb.KeyType_Ed25519)
	ECDSA     = KeyType(icpb.KeyType_ECDSA)
	Secp256k1 = KeyType(icpb.KeyType_Secp256k1)
	// Generic octet sequence ([]byte) for symetric keys, icpb.KeyType uses 0-3
	Octet = KeyType(4)
)

type Key interface {
	// Equals checks whether two keys are the same
	Equals(Key) bool
	// Raw returns the raw bytes of the key
	Raw() ([]byte, error)
	// Type returns the key type
	Type() KeyType
}

func IsAsymmetric(key Key) bool {
	switch key.(type) {
	case PublicKey, PrivateKey:
		return true
	}
	return false
}

func IsPublic(key Key) bool {
	switch key.(type) {
	case PublicKey:
		return true
	}
	return false
}

func IsPrivate(key Key) bool {
	switch key.(type) {
	case PrivateKey:
		return true
	}
	return false
}

func GetPublic(key Key) (PublicKey, error) {
	if !IsAsymmetric(key) {
		return nil, fmt.Errorf("key must be asymmetric")
	}

	switch kt := key.(type) {
	case PublicKey:
		return kt, nil
	case PrivateKey:
		return kt.GetPublic(), nil
	}
	return nil, fmt.Errorf("unknown key type")
}

func GetPrivate(key Key) (PrivateKey, error) {
	if !IsAsymmetric(key) {
		return nil, fmt.Errorf("key must by asymmetric")
	}

	priv, ok := key.(PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key isn't private")
	}

	return priv, nil
}

// PublicKey
type PublicKey interface {
	Key
	Verify(data []byte, sig []byte) (bool, error)
	Point() kyber.Point
	Std() (gocrypto.PublicKey, error)
	String() string
	DID() (string, error)
}

var _ PublicKey = (*pubKeyLibP2P)(nil)

type pubKeyLibP2P struct {
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
	case *ecdsa.PublicKey:
		icpk, err = ic.ECDSAPublicKeyFromPubKey(*pkt)
	case secp256k1.PublicKey:
		sppk := ic.Secp256k1PublicKey(pkt)
		icpk = &sppk
	case *secp256k1.PublicKey:
		icpk = (*ic.Secp256k1PublicKey)(pkt)
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

func PublicKeyFromBytes(keyType KeyType, buf []byte) (PublicKey, error) {
	var pk ic.PubKey
	var err error
	switch keyType.String() {
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
		return nil, fmt.Errorf("public key from bytes: %w", err)
	}

	return PublicKeyFromLibP2P(pk)
}

func publicKeyFromLibP2P(pubkey ic.PubKey) (*pubKeyLibP2P, error) {
	suite, err := SuiteForType(KeyType(pubkey.Type()))
	if err != nil {
		return nil, err
	}

	return &pubKeyLibP2P{
		PubKey: pubkey,
		suite:  suite,
	}, nil

}

func ToLibP2PPublicKey(pk PublicKey) ic.PubKey {
	return pk.(*pubKeyLibP2P).PubKey
}

func PublicKeyToProto(pk PublicKey) (*icpb.PublicKey, error) {
	return ic.PublicKeyToProto(pk.(*pubKeyLibP2P).PubKey)
}

func (p *pubKeyLibP2P) Point() kyber.Point {
	buf, _ := p.PubKey.Raw()
	point := p.suite.Point()
	point.UnmarshalBinary(buf)
	return point
}

func (p *pubKeyLibP2P) Std() (gocrypto.PublicKey, error) {
	// our version of "standard" secp256k1 keys uses the
	// dcred type directly, as its more common among our
	// dependencies (like go-jose)
	switch pk := p.PubKey.(type) {
	case *ic.Secp256k1PublicKey:
		return (*secp256k1.PublicKey)(pk), nil
	}
	gokey, err := ic.PubKeyToStdKey(p.PubKey)
	if err != nil {
		return nil, err
	}

	// convert ed25519 keys to non pointers
	switch kt := gokey.(type) {
	case *ed25519.PublicKey:
		return *kt, nil
	}
	return gokey, nil
}

func (p *pubKeyLibP2P) String() string {
	buf, _ := p.Raw()

	enc := b64.StdEncoding.EncodeToString(buf)
	return enc
}

func (p *pubKeyLibP2P) DID() (string, error) {
	didKeyType, err := cryptoKeyTypeToDID(p.Type())
	if err != nil {
		return "", err
	}
	keyBuf, err := p.Raw()
	if err != nil {
		return "", err
	}
	didKey, err := key.CreateDIDKey(didKeyType, keyBuf)
	if err != nil {
		return "", fmt.Errorf("creating did: %w", err)
	}

	return didKey.String(), nil
}

func (p *pubKeyLibP2P) MarshalJWK() ([]byte, error) {
	goPubKey, err := ic.PubKeyToStdKey(p.PubKey)
	if err != nil {
		return nil, err
	}
	jPubKey, err := jwk.FromRaw(goPubKey)
	if err != nil {
		return nil, err
	}

	return json.Marshal(jPubKey)
}

func (p *pubKeyLibP2P) Type() KeyType {
	// we can safely ignore this error
	// since its already validated
	return KeyType(p.PubKey.Type())
}

func (p *pubKeyLibP2P) Equals(k1 Key) bool {
	if p == k1 {
		return true
	}

	return basicEquals(p, k1)
}

func cryptoKeyTypeToDID(kt KeyType) (scrypto.KeyType, error) {
	switch kt {
	case Ed25519:
		return scrypto.Ed25519, nil
	case Secp256k1:
		return scrypto.SECP256k1, nil
	}

	return "", fmt.Errorf("invalid key type")
}

func basicEquals(k1, k2 Key) bool {
	if k1.Type() != k2.Type() {
		return false
	}

	a, err := k1.Raw()
	if err != nil {
		return false
	}
	b, err := k2.Raw()
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(a, b) == 1
}

var _ PrivateKey = (*privKeyLibP2P)(nil)

type PrivateKey interface {
	Key
	Sign([]byte) ([]byte, error)
	Scalar() kyber.Scalar
	GetPublic() PublicKey
	Std() (gocrypto.PrivateKey, error)
}

type privKeyLibP2P struct {
	ic.PrivKey
	suite suites.Suite
}

func GenerateKeyPair(keyType KeyType, src io.Reader) (PrivateKey, PublicKey, error) {
	suite, err := SuiteForType(keyType)
	if err != nil {
		return nil, nil, err
	}
	sk, pk, err := ic.GenerateKeyPairWithReader(int(keyType), 0, src)
	if err != nil {
		return nil, nil, err
	}

	return &privKeyLibP2P{
			PrivKey: sk,
			suite:   suite,
		}, &pubKeyLibP2P{
			PubKey: pk,
			suite:  suite,
		}, nil
}

func PrivateKeyFromBytes(keyType KeyType, buf []byte) (PrivateKey, error) {
	var pk ic.PrivKey
	var err error
	switch keyType {
	case Ed25519:
		pk, err = ic.UnmarshalEd25519PrivateKey(buf)
	case Secp256k1:
		pk, err = ic.UnmarshalSecp256k1PrivateKey(buf)
	case ECDSA:
		pk, err = ic.UnmarshalECDSAPrivateKey(buf)
	// case "rsa":
	// 	pk, err = ic.UnmarshalRsaPrivateKey(buf)
	default:
		return nil, ErrBadKeyType
	}

	if err != nil {
		return nil, fmt.Errorf("public key from bytes: %w", err)
	}

	return PrivateKeyFromLibP2P(pk)
}

func PrivateKeyFromLibP2P(privkey ic.PrivKey) (PrivateKey, error) {
	suite, err := SuiteForType(KeyType(privkey.Type()))
	if err != nil {
		return nil, err
	}

	return &privKeyLibP2P{
		PrivKey: privkey,
		suite:   suite,
	}, nil
}

func ToLibP2PPrivateKey(pk PrivateKey) ic.PrivKey {
	return pk.(*privKeyLibP2P).PrivKey
}

func (p *privKeyLibP2P) Std() (gocrypto.PrivateKey, error) {
	// our version of "standard" secp256k1 keys uses the
	// dcred type directly, as its more common among our
	// dependencies (like go-jose)
	switch pk := p.PrivKey.(type) {
	case *ic.Secp256k1PrivateKey:
		return (*secp256k1.PrivateKey)(pk), nil
	}

	gokey, err := ic.PrivKeyToStdKey(p.PrivKey)
	if err != nil {
		return nil, err
	}

	// convert ed25519 keys to non pointers
	switch kt := gokey.(type) {
	case *ed25519.PrivateKey:
		return *kt, nil
	}
	return gokey, nil
}

// Scalar returns a numeric elliptic curve scalar
// representation of the private key.
//
// WARNING: THIS ONLY WORKS WITH ED25519 & SECP256K1 CURVES RIGHT NOW.
func (p *privKeyLibP2P) Scalar() kyber.Scalar {
	switch p.Type() {
	case Ed25519:
		return p.ed25519Scalar()
	case Secp256k1:
		return p.secp256k1Scalar()
	default:
		panic("only ed25519 and secp256k1 private key scalar conversion supported")
	}
}

func (p *privKeyLibP2P) secp256k1Scalar() kyber.Scalar {
	buf, err := p.Raw()
	if err != nil {
		panic(err) // todo
	}

	return p.suite.Scalar().SetBytes(buf)
}

func (p *privKeyLibP2P) ed25519Scalar() kyber.Scalar {
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

func (p *privKeyLibP2P) GetPublic() PublicKey {
	return &pubKeyLibP2P{
		PubKey: p.PrivKey.GetPublic(),
		suite:  p.suite,
	}
}

func (p *privKeyLibP2P) MarshalJWK() ([]byte, error) {
	goPubKey, err := ic.PrivKeyToStdKey(p.PrivKey)
	if err != nil {
		return nil, err
	}
	jPubKey, err := jwk.FromRaw(goPubKey)
	if err != nil {
		return nil, err
	}

	return json.Marshal(jPubKey)
}

func (p *privKeyLibP2P) Type() KeyType {
	return KeyType(p.PrivKey.Type())
}

func (p *privKeyLibP2P) Equals(k1 Key) bool {
	if p == k1 {
		return true
	}

	return basicEquals(p, k1)
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

func PublicKeyToBech32(pubkey PublicKey) (string, error) {
	buf, err := pubkey.Raw()
	if err != nil {
		return "", fmt.Errorf("getting public key bytes: %w", err)
	}
	sum := tmhash.SumTruncated(buf)
	return bech32.ConvertAndEncode("source", sum)
}

func PublicKeyBytesToBech32(pubkey []byte) (string, error) {
	// convert to address
	sha := sha256.Sum256(pubkey)
	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha[:]) // does not error
	buf := cometcrypto.Address(hasherRIPEMD160.Sum(nil))
	addr := sdk.AccAddress(buf)
	// convert to bech32 (stringified address)
	return addr.String(), nil
}

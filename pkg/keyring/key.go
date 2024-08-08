package keyring

/*

import (
	"github.com/sourcenetwork/orbis-go/pkg/crypto"

	scrypto "github.com/TBD54566975/ssi-sdk/crypto"
	"github.com/TBD54566975/ssi-sdk/did/key"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type Key struct {
	jwk.Key
}

// NewKey creates a new key record.
// It accepts:
// - ecdsa.PublicKey
// - ecdsa.PrivateKey
// - ed25519.PublicKey
// - ed25519.PrivateKey
// - []byte (symetric Key)
func NewKey(k any) (Key, error)

func (k Key) DID() (string, error) {
	if !k.isDIDKeyType() {
		return "", ErrInvalidDIDKeyType
	}
	crv, ok := k.Get("crv")
	if !ok {
		return "", ErrInvalidDIDKeyType
	}
	crvStr, ok := crv.(string)
	if !ok {
		return "", ErrInvalidDIDKeyType
	}

	didKeyType, err := jwkToDidKeyType(crvStr)
	if err != nil {
		return "", err
	}

	pk, _ := k.PublicKey()
	pkbuf, err := pk.Raw()
	if err != nil {
		return "", err
	}
	didKey, err := key.CreateDIDKey(didKeyType, pkbuf)
	if err != nil {
		return "", err
	}
	return didKey.String(), nil
}

func (k Key) PrivateKey() (crypto.PrivateKey, error)

func (k Key) PublicKey() (crypto.PrivateKey, error)

func (k Key) isDIDKeyType() bool {
	return k.Algorithm() == jwa.EC
}

func (k Key) isSignerType() bool {
	return k.Algorithm() == jwa.EC // yes this is duplicate of isDIDKeyType :)
}

var (
	jwkToDidKeyTypeMap = map[jwa.EllipticCurveAlgorithm]scrypto.KeyType{
		jwa.Ed25519:   scrypto.Ed25519,
		jwa.X25519:    scrypto.X25519,
		jwa.Secp256k1: scrypto.SECP256k1,
		jwa.P256:      scrypto.P256,
		jwa.P384:      scrypto.P384,
		jwa.P521:      scrypto.P521,
	}
)

func jwkToDidKeyType(kty string) (scrypto.KeyType, error) {
	sckt, ok := jwkToDidKeyTypeMap[jwa.EllipticCurveAlgorithm(kty)]
	if !ok {
		return "", ErrInvalidDIDKeyType
	}
	return sckt, nil
}

*/

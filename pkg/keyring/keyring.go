package keyring

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/keyring/keystore"
)

type Info struct {
	Name string
	Key  crypto.Key
}

type Keyring interface {
	Set(name string, key crypto.Key) error
	Get(name string) (crypto.Key, error)
	Delete(name string) error
	List() ([]Info, error)

	Sign(name string, msg []byte) ([]byte, crypto.PublicKey, error)
}

type keyring struct {
	store   keystore.Keystore
	backend string
}

func New(backend string, args ...any) (Keyring, error) {
	kstore, err := keystore.Open(backend, args...)
	if err != nil {
		return nil, err
	}

	return &keyring{
		store:   kstore,
		backend: backend,
	}, nil
}

func (k *keyring) Set(name string, key crypto.Key) error {
	jkey, err := keyToJWK(key)
	if err != nil {
		return fmt.Errorf("converting to JWK: %w", err)
	}
	data, err := json.Marshal(jkey)
	if err != nil {
		return fmt.Errorf("marshaling JWK: %w", err)
	}
	return k.store.Set(name, data)
}

func keyToJWK(key crypto.Key) (jwk.Key, error) {
	var gokey any
	var err error
	switch kt := key.(type) {
	case crypto.PrivateKey:
		// although these cases seem like they are the same, they aren't <3
		// the Std() method returns diff types depending, so it can't
		// be further abstracted/simplified
		gokey, err = kt.Std()
	case crypto.PublicKey:
		gokey, err = kt.Std()
	default:
		rawKey, err := kt.Raw()
		if err != nil {
			return nil, err
		}
		return jwk.FromRaw(rawKey)
	}
	if err != nil {
		return nil, err
	}

	// jwx/jwk package only understands go standard lib ecdsa public/private keys
	// not the decred/secp256k1 package types.
	//
	// TODO: Update `Std()` method on crypto.Pub/PrivKeys to match this
	switch rawkey := gokey.(type) {
	case *secp256k1.PublicKey:
		return jwk.FromRaw(rawkey.ToECDSA())
	case *secp256k1.PrivateKey:
		return jwk.FromRaw(rawkey.ToECDSA())
	default:
		return jwk.FromRaw(rawkey)
	}
}

func (k *keyring) Get(name string) (crypto.Key, error) {
	if name == "" {
		return nil, fmt.Errorf("missing key name")
	}
	data, err := k.store.Get(name)
	if err != nil {
		return nil, fmt.Errorf("get from store: %w", err)
	}

	jKey, err := jwk.ParseKey(data)
	if err != nil {
		return nil, fmt.Errorf("parse JWK: %w", err)
	}

	return jwkToKey(jKey)
}

// jwkToKey converts a jwk.JWK key to a crypto.Key key.
// However, there are some caveats when making this conversion
// due to a type mismatch in the orbis/crypto package (and libp2p).
// So we have to do more work when converting secp256k1 keys.
//
// Note: We don't yet support Symmetric keys in the orbis/crypto
// package so this is left as a TODO.
func jwkToKey(key jwk.Key) (crypto.Key, error) {
	asymmetric, ok := key.(jwk.AsymmetricKey)
	if !ok {
		// TODO HANDLE SYMMETRIC KEYS
		panic("todo symmetric keys in keyring")
	}

	crvI, ok := key.Get("crv")
	if !ok {
		return nil, fmt.Errorf("can't resolve JWK curve")
	}
	crv, ok := crvI.(jwa.EllipticCurveAlgorithm)
	if !ok {
		return nil, fmt.Errorf("bad JWK curve")
	}

	// we need to handle secp256k1 keys separetly
	// because of a type mismatch between dcrec/secp256k1 and
	// crypto/ecdsa
	//
	// Also, we only support secp256k1 and ed25519 keys!
	if asymmetric.IsPrivate() && crv == jwa.Secp256k1 {
		var ecKey ecdsa.PrivateKey
		key.Raw(&ecKey)
		// convert from crypto/ecdsa to dcrec/secp256k1
		var scalar secp256k1.ModNScalar
		scalar.SetByteSlice(ecKey.D.Bytes())
		gokey := secp256k1.NewPrivateKey(&scalar)
		return crypto.PrivateKeyFromBytes(crypto.Secp256k1, gokey.Serialize())
	} else if !asymmetric.IsPrivate() && crv == jwa.Secp256k1 {
		var ecKey ecdsa.PublicKey
		key.Raw(&ecKey)
		// convert from crypto/ecdsa to dcrec/secp256k1
		x := bigIntToFieldVal(ecKey.X)
		y := bigIntToFieldVal(ecKey.Y)
		gokey := secp256k1.NewPublicKey(x, y)
		return crypto.PublicKeyFromBytes(crypto.Secp256k1, gokey.SerializeUncompressed())
	} else if asymmetric.IsPrivate() && crv == jwa.Ed25519 {
		var ecKey ed25519.PrivateKey
		key.Raw(&ecKey)
		return crypto.PrivateKeyFromBytes(crypto.Ed25519, []byte(ecKey))
	} else if !asymmetric.IsPrivate() && crv == jwa.Ed25519 {
		var ecKey ed25519.PublicKey
		key.Raw(&ecKey)
		return crypto.PublicKeyFromBytes(crypto.Ed25519, []byte(ecKey))
	}

	return nil, fmt.Errorf("unsupported curve or algorithm")
}

// bigIntToFieldVal converts a big.Int to a secp256k1.FieldVal
func bigIntToFieldVal(b *big.Int) *secp256k1.FieldVal {
	fv := new(secp256k1.FieldVal)
	fv.SetByteSlice(b.Bytes())
	return fv
}

func (k *keyring) Delete(name string) error {
	return k.store.Delete(name)
}

func (k *keyring) List() ([]Info, error) {
	rawinfos, err := k.store.List()
	if err != nil {
		return nil, fmt.Errorf("keystore list: %w", err)
	}
	infos := make([]Info, len(rawinfos))
	for i, info := range rawinfos {
		rawKey, err := jwk.ParseKey(info.Data)
		if err != nil {
			return nil, fmt.Errorf("parse JWK: %w", err)
		}
		key, err := jwkToKey(rawKey)
		if err != nil {
			return nil, fmt.Errorf("converting JWK: %w", err)
		}
		infos[i] = Info{
			Name: info.Name,
			Key:  key,
		}
	}

	return infos, nil
}

func (k *keyring) Sign(name string, msg []byte) ([]byte, crypto.PublicKey, error) {
	key, err := k.Get(name)
	if err != nil {
		return nil, nil, err
	}

	// make sure the key we got supports signing
	privKey, ok := key.(crypto.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("signing only supported for private keys")
	}

	sig, err := privKey.Sign(msg)
	if err != nil {
		return nil, nil, err
	}

	return sig, privKey.GetPublic(), nil
}

type ctxKey string

const (
	keyringCtxKey ctxKey = "keyring"
)

// WithKeyring adds a keyring to a context
func WithKeyring(ctx context.Context, kr Keyring) context.Context {
	return context.WithValue(ctx, keyringCtxKey, kr)
}

// FromContext gets an existing keyring from a context
func FromContext(ctx context.Context) (Keyring, bool) {
	keyring, ok := ctx.Value(keyringCtxKey).(Keyring)
	return keyring, ok
}

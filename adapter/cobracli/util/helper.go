package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"

	ssicrypto "github.com/TBD54566975/ssi-sdk/crypto"
	"github.com/go-jose/go-jose/v3/jwt"
	ic "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/spf13/pflag"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
)

type result map[string]interface{}

func (r result) Output() error {

	j, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("marshal result to json: %w", err)
	}

	fmt.Println(string(j))

	return nil
}

func claimsFromFlag(flag *pflag.Flag) (jwt.Claims, error) {

	claims := flag.Value.String()
	c := jwt.Claims{}
	err := json.Unmarshal([]byte(claims), &c)
	if err != nil {
		return jwt.Claims{}, fmt.Errorf("unmarshal claims: %w", err)
	}
	return c, nil
}

func bytesFromFlag(flag *pflag.Flag) ([]byte, error) {

	return base64.StdEncoding.DecodeString(flag.Value.String())
}

func suiteFromFlag(flag *pflag.Flag) (suites.Suite, error) {

	keyType := flag.Value.String()
	ste, err := suites.Find(keyType)
	if err != nil {
		return nil, fmt.Errorf("find suite: %w", err)
	}

	return ste, nil
}

func pointFromFlag(ste suites.Suite, flag *pflag.Flag) (kyber.Point, error) {

	b64 := flag.Value.String()
	p := ste.Point()

	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return p, fmt.Errorf("decode: %w", err)
	}

	err = p.UnmarshalBinary(raw)
	if err != nil {
		return p, fmt.Errorf("unmarshal: %w", err)
	}

	return p, nil
}

func pointToB64(p kyber.Point) (string, error) {

	raw, err := p.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	b64 := base64.StdEncoding.EncodeToString(raw)

	return b64, nil
}

func skFromFlag(ste suites.Suite, flag *pflag.Flag) (kyber.Scalar, error) {

	b64 := flag.Value.String()

	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	icSk, err := ic.UnmarshalEd25519PrivateKey(raw)
	if err != nil {
		return nil, fmt.Errorf("unmarshal private key: %w", err)
	}

	sk, err := crypto.PrivateKeyFromLibP2P(icSk)
	if err != nil {
		return nil, fmt.Errorf("convert private key: %w", err)
	}

	return sk.Scalar(), nil
}

func kyberSuiteToSSIKeyType(ste suites.Suite) (ssicrypto.KeyType, error) {
	var keyType ssicrypto.KeyType

	switch strings.ToLower(ste.String()) {
	case "ed25519":
		keyType = ssicrypto.Ed25519
	case "secp256k1":
		keyType = ssicrypto.SECP256k1
	case "rsa":
		keyType = ssicrypto.RSA
	case "ecdsa":
		keyType = ssicrypto.SECP256k1ECDSA
	default:
		return keyType, fmt.Errorf("unsupported key type: %s", ste.String())
	}

	return keyType, nil
}

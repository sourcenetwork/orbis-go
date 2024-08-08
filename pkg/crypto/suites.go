package crypto

import (
	"strings"

	"github.com/sourcenetwork/orbis-go/pkg/crypto/suites/secp256k1"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/suites"
	"go.dedis.ch/protobuf"
)

// register protobuf custom reflect marshallers
func init() {
	ed25519 := edwards25519.NewBlakeSHA256Ed25519()
	spk1 := secp256k1.NewBlakeKeccackSecp256k1()

	protobuf.RegisterInterface(func() interface{} { return ed25519.Point() })
	protobuf.RegisterInterface(func() interface{} { return ed25519.Scalar() })

	protobuf.RegisterInterface(func() interface{} { return spk1.Point() })
	protobuf.RegisterInterface(func() interface{} { return spk1.Scalar() })
}

func SuiteForType(kt KeyType) (suites.Suite, error) {
	switch kt {
	case Secp256k1, ECDSA:
		return secp256k1.NewBlakeKeccackSecp256k1(), nil
	case Ed25519:
		// TODO
		// reader := rand.New(rand.NewSource(0))
		// r := random.New(reader)
		return edwards25519.NewBlakeSHA256Ed25519(), nil
	default:
		return nil, ErrBadKeyType
	}
}

func FindSuite(name string) (suites.Suite, error) {
	switch strings.ToLower(name) {
	case "ed25519":
		return edwards25519.NewBlakeSHA256Ed25519(), nil
	case "secp256k1":
		return secp256k1.NewBlakeKeccackSecp256k1(), nil
	default:
		return nil, ErrBadKeyType
	}
}

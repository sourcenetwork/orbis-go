package crypto

import (
	icpb "github.com/libp2p/go-libp2p/core/crypto/pb"
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

func SuiteForType(kt icpb.KeyType) (suites.Suite, error) {
	switch kt {
	case icpb.KeyType_Secp256k1:
		return secp256k1.NewBlakeKeccackSecp256k1(), nil
	case icpb.KeyType_Ed25519:
		return edwards25519.NewBlakeSHA256Ed25519(), nil
	default:
		return nil, ErrBadKeyType
	}
}

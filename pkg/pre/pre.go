package pre

import (
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

type ReencryptReply interface {
	Share() share.PubShare
	Challenge() kyber.Scalar
	Proof() kyber.Scalar
}

// PRE via Threshold MPC
type PRE interface {
	// Initialize the PRE system
	Init(ring types.RingID, n int32, t int32, nodes []types.Node) error
	// Name of the PRE implementation
	Name() string
	// Reencrypt using a nodes local private share
	Reencrypt(crypto.PublicKey, kyber.Point) (ReencryptReply, error)
	// Process incoming replies from other nodes
	// note: We can likely drop pss.Node
	Process(pss.Node, ReencryptReply) error
	// Recover the encrypted ReKey
	Recover() (kyber.Point, error)
}

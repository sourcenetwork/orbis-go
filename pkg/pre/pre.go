package pre

import (
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"

	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

type ReencryptReply interface {
	Share() share.PubShare
	Challenge() kyber.Scalar
	Proof() kyber.Scalar
}

// PRE via Threshold MPC
type PRE interface {
	Name() string
	// Reencrypt using a nodes local private share
	Reencrypt(crypto.PublicKey, kyber.Point) (ReencryptReply, error)
	// Process incoming replies from other nodes
	// note: We can likely drop pss.Node
	Process(pss.Node, ReencryptReply) error
	// Recover the encrypted ReKey
	Recover() (kyber.Point, error)
}

type Factory interface {
	New(types.RingID, int, int, transport.Transport, bulletin.Bulletin, []types.Node, dkg.DKG) (PRE, error)
}

type ProvideFactory func(*do.Injector) Factory

package pre

import (
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/suites"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

type ReencryptReply struct {
	Share     share.PubShare // nodes re-encrypted secret share
	Challenge kyber.Scalar   // random oracle challenge
	Proof     kyber.Scalar   // nizk proofi of re-encryption
}

// PRE via Threshold MPC
type PRE interface {

	// Initialize the PRE system
	Init(rid types.RingID, n int32, t int32) error

	// Name of the PRE implementation
	Name() string

	// Reencrypt using a nodes local private share
	Reencrypt(prishare crypto.DistKeyShare, scrt *types.Secret, rdrPk crypto.PublicKey) (ReencryptReply, error)

	// Verify incoming replies from other nodes
	Verify(rdrPk crypto.PublicKey, dkgCmt crypto.PubPoly, encCmt kyber.Point, reply ReencryptReply) error

	// Recover the encrypted ReKey
	Recover(ste suites.Suite, xncSki []*share.PubShare, t int, n int) (kyber.Point, error)
}

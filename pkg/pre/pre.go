package pre

import (
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"
)

type ReencryptReply interface {
	Share() share.PubShare
	Challenge() kyber.Scalar
	Proof() kyber.Scalar
}

// Threshold PRE
type Theshold interface {
	// Reencrypt using a nodes local private share
	Reencrypt(crypto.PublicKey, kyber.Point) (ReencryptReply, error)
	// Process incoming replies from other nodes
	// note: We can likely drop pss.Node
	Process(pss.Node, ReencryptReply) error
	// Recover the encrypted ReKey
	Recover() (kyber.Point, error)
}

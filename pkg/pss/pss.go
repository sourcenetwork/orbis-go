package pss

import (
	"context"

	"go.dedis.ch/kyber/v3/suites"

	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

type Message interface{}

type PSS interface {
	// Name of the PSS Algorithm
	Name() string
	// Cryptographic suite
	Suite() suites.Suite

	// Start the service
	Start()
	// Shutdown the service
	Shutdown() error
	// Process incoming messages relating to the
	// operations or maintenence of the PSS/DKG
	// algorithm
	ProcessMessage(context.Context, Message)

	// Aggregate public key of the PSS/DKG
	PublicKey() crypto.PublicKey
	// Public polynomial
	// might not need on interface
	PublicPoly() crypto.PubPoly

	// Private share of this node
	Share() crypto.PriShare

	// State of the PSS
	State() State

	Num() int
	Threshold() int
}

type Node interface {
	transport.Node
	Index() int
}

type Factory interface {
	New(types.RingID, int32, int32, transport.Transport, bulletin.Bulletin, []types.Node, dkg.DKG) (PSS, error)
	Name() string
}

package pss

import (
	"context"

	"go.dedis.ch/kyber/v3/suites"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

type Message interface{}

type PSS interface {
	// Initialze a new PSS
	Init(types.RingID, int32, int32, []types.Node) error
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
	Share() crypto.DistKeyShare

	// State of the PSS
	State() string

	Num() int
	Threshold() int
}

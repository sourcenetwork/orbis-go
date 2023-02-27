package pss

import (
	"context"

	"github.com/samber/do"
	"go.dedis.ch/kyber/v3/suites"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/pss/types"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

type Message interface{}

type Service interface {
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
	PublicPoly() crypto.PubPoly

	// Private share of this node
	Share() crypto.PriShare

	// State of the PSS
	State() types.State

	Num() int
	Threshold() int
}

type Node interface {
	transport.Node
	Index() int
}

type ProviderFn = func(*do.Injector) Service

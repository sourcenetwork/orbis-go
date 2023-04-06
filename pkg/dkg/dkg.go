package dkg

import (
	"context"

	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

const (
	ProtocolName = "dkg"
)

// enum?
type State struct{}

type Node = transport.Node

type DKG interface {
	Init(ctx context.Context, nodes []transport.Node, n int, threshold int) error
	Name() string

	PublicKey() crypto.PublicKey
	Share() crypto.PriShare

	State() State

	Start(context.Context) error
	Close(context.Context) error

	ProcessMessage(*transport.Message) error

	// hooks?
}

type Factory interface {
	New(db.Repository, transport.Transport, bulletin.Bulletin, crypto.PrivateKey) (DKG, error)
}

// ProvideFactory
// or FactoryProvider??
// type ProvideFactory func(*do.Injector) Factory

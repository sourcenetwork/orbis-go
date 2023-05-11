package dkg

import (
	"context"

	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

const (
	INITIALIZED State = iota // DKG group has initialized but not started the generation
	STARTED                  // Started the distributed key generation
	CERTIFIED                // Generated and cerified the shared key

	ProtocolName = "dkg"
)

// enum
type State uint8

type Node = transport.Node

type DKG interface {
	Init(ctx context.Context, pk crypto.PrivateKey, nodes []Node, n int, threshold int) error

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
	New(*do.Injector, []*db.RepoKey) (DKG, error)
	Name() string
	Repos() []string
}

package dkg

import (
	"context"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
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
	Init(ctx context.Context, pk crypto.PrivateKey, nodes []Node, n int32, threshold int32) error

	Name() string

	PublicKey() crypto.PublicKey
	Share() crypto.PriShare

	State() State

	Start(context.Context) error
	Close(context.Context) error

	ProcessMessage(*transport.Message) error

	// hooks?
}

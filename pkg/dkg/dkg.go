package dkg

import (
	"context"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

const (
	UNINITIALIZED State = iota // DKG group has not been initialized.
	INITIALIZED                // DKG group has initialized but not started the generation.
	STARTED                    // Started the distributed key generation.
	CERTIFIED                  // Generated and cerified the shared key.

	CUSTOM_STATE_MASK State = 0b10000000 // Mask to reserve usage of custom enums for implementations

	ProtocolName = "dkg"
)

// enum
type State uint8

func (s State) String() string {
	switch s {
	case UNINITIALIZED:
		return "UNINITIALIZED"
	case INITIALIZED:
		return "INITIALIZED"
	case STARTED:
		return "STARTED"
	case CERTIFIED:
		return "CERTIFIED"
	default:
		return "UNKNOWN"
	}
}

type Node = transport.Node

type DKG interface {
	Init(ctx context.Context, pk crypto.PrivateKey, rid types.RingID, nodes []Node, n int32, threshold int32, fromState bool) error

	Name() string

	PublicKey() crypto.PublicKey
	Share() crypto.PriShare

	State() State

	Start(context.Context) error
	Close(context.Context) error

	ProcessMessage(*transport.Message) error

	// hooks?
}

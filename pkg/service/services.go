package service

import (
	"context"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/crypto/proof"
	ptypes "github.com/sourcenetwork/orbis-go/pkg/pss/types"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

// type SSInfo interface {
// 	NodeInfo() NodeInfo
// }

// type NodeInfoConstraint any

// SecretRingService is a service that manages
// user created 'secrets'.
type SecretRingService[SSInfo interface {
	NodeInfo() NodeInfo
}, NodeInfo any] interface {
	// Secret Management Operations
	SecretsManagerService

	// DKG Operations
	DKGService[NodeInfo]
}

type SecretsManagerService interface {
	// Create(context.Context) types.Secret
	Store(context.Context, SID, types.Secret, proof.VerifiableEncryption) error
	Get(context.Context, SID) (types.Secret, error)
	GetShares(context.Context, SID) ([]types.PrivSecretShare, error)
	Delete(context.Context, SID) error
}

// DKGService
type DKGService[NodeInfo any] interface {
	PublicKey() (crypto.PublicKey, error)
	Refresh(context.Context, dkg.Config) (dkg.RefreshState, error)
	Threshold() int
	State() ptypes.State

	// Committee Operations
	CommitteeService[NodeInfo]
}

type CommitteeService[NodeInfo any] interface {
	Nodes() []NodeInfo
}

// SID is a Secret Identifier (SecretID)
type SID struct{}

// type NodeInfo struct {}

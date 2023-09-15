package service

import (
	"context"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/crypto/proof"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

// NOTE: This file will be broken down into repective service files.

// SecretRingService is a service that manages
// user created 'secrets'.
// type RingService interface {
// 	// Secret Management Operations
// 	SecretsManagerService

// 	// DKG + PSS Operations
// 	DistKeyService
// }

type SecretsManagerService interface {
	// Create(context.Context) types.Secret
	Store(context.Context, types.SecretID, *types.Secret, proof.VerifiableEncryption) error
	Get(context.Context, types.SecretID) (types.Secret, error)
	GetShares(context.Context, types.SecretID) ([]types.PrivSecretShare, error)
	Delete(context.Context, types.SecretID) error
}

// DistKeyService
type DistKeyService interface {
	PublicKey() (crypto.PublicKey, error)
	Refresh(context.Context, pss.Config) (pss.RefreshState, error)
	Threshold() int
	State() pss.State

	// Committee Operations
	CommitteeService
}

type CommitteeService interface {
	Nodes() []types.Node
}

// SID is a Secret Identifier (SecretID)
type SID struct{}

// type NodeInfo struct {}

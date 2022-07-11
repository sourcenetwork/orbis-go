package service

import (
	"context"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/proof"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

// SecretRingService is a service that manages
// user created 'secrets'.
type SecretRingService interface {
	// Secret Management Operations
	SecretsManagerService

	// DKG Operations
	DKGService

	// Network Operations
	GetNodes()
}

type SecretsManagerService interface {
	// Create(context.Context) types.Secret
	Store(context.Context, types.SID, types.Secret, proof.VerifiableEncryption) error
	Get(context.Context, types.SID) (types.Secret, error)
	GetShare(context.Context, types.SID) (types.SecretShare, error)
	GetShares(context.Context, types.SID) ([]types.SecretShare, error)
	Delete(context.Context, types.SID) error
}

// DKGService
type DKGService interface {
	PublicKey() (crypto.PublicKey, error)
	Refresh(context.Context) (dkg.RefreshState, error)
	Threshold() int
}

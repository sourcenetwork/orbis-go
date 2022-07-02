package service

import (
	"context"

	"github.com/sourcenetwork/orbis-go/pkg/types"
)

// SecretRingService is a service that manages
// user created 'secrets'.
type SecretRingService interface {
	// Secret Management Operations
	SecretManagementService

	// DKG Operations
	DKGService

	// Network Operations
	GetNodes()
}

type SecretManagementService interface {
	//
	Create(context.Context) types.Secret
	Store(context.Context, types.SID, types.Secret, types.Proof) error
	Get(context.Context, types.SID) (types.Secret, error)
	GetShare(context.Context, types.SID) (types.SecretShare, error)
	GetShares(context.Context, types.SID) ([]types.SecretShare, error)
	Delete(context.Context, types.SID) error
}

// DKGService
type DKGService interface {
	PublicKey() (types.PublicKey, error)
	Refresh(context.Context) (types.RefreshState, error)
	Threshold() int
}

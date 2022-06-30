package service

import (
	"context"

	"github.com/sourcenetwork/orbis-go/pkg/core"
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
	Create(context.Context) core.Secret
	Store(context.Context, core.SID, core.Secret, core.Proof) error
	Get(context.Context, core.SID) (core.Secret, error)
	GetShare(context.Context, core.SID) (core.SecretShare, error)
	GetShares(context.Context, core.SID) ([]core.SecretShare, error)
	Delete(context.Context, core.SID) error
}

// DKGService
type DKGService interface {
	PublicKey() (core.PublicKey, error)
	Refresh(context.Context) (core.RefreshState, error)
	Threshold() int
}

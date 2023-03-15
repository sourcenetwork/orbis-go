package dkg

import (
	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

type DKG interface{}

type Factory interface {
	New(types.RingID, int, int, transport.Transport, bulletin.Bulletin, []types.Node) (DKG, error)
}

// ProvideFactory
// or FactoryProvider??
type ProvideFactory func(*do.Injector) Factory

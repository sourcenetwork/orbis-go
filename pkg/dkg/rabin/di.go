package rabin

import (
	"fmt"

	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	odb "github.com/sourcenetwork/orbis-go/pkg/db"
	orbisdkg "github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

var (
	_       types.Factory[orbisdkg.DKG] = (*factory)(nil)
	Factory                             = factory{}
)

type factory struct{}

func (factory) New(inj *do.Injector, rkeys []odb.RepoKey, _ config.Config) (orbisdkg.DKG, error) {
	db, err := do.Invoke[*odb.DB](inj)
	if err != nil {
		return nil, fmt.Errorf("invoke db: %w", err)
	}

	t, err := do.Invoke[transport.Transport](inj)
	if err != nil {
		return nil, fmt.Errorf("invoke transport: %w", err)
	}

	b, err := do.Invoke[bulletin.Bulletin](inj)
	if err != nil {
		return nil, fmt.Errorf("invoke bulletin: %w", err)
	}

	return New(db, rkeys, t, b)
}

func (factory) Name() string {
	return name
}

func (factory) Repos() []string {
	return []string{"dkg"}
}

// /rabin/{deals,shares}

package elgamal

import (
	"fmt"

	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/pre"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

var _ types.Factory[pre.PRE] = (*factory)(nil)

var (
	Factory = factory{}
)

type factory struct{}

func (factory) New(inj *do.Injector, rkeys []db.RepoKey) (pre.PRE, error) {

	db, err := do.Invoke[*db.DB](inj)
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

	dkg, err := do.Invoke[dkg.DKG](inj)
	if err != nil {
		return nil, fmt.Errorf("invoke dkg: %w", err)
	}

	return New(db, rkeys, t, b, dkg)
}

func (factory) Name() string {
	return name
}

func (factory) Repos() []string {
	return []string{}
}

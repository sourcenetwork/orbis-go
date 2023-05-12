package rabin

import (
	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	orbisdkg "github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

var (
	Factory = factory{}
)

type factory struct{}

func (factory) New(inj *do.Injector, rkeys []*db.RepoKey) (orbisdkg.DKG, error) {
	db, err := do.Invoke[*db.DB](inj)
	if err != nil {
		return nil, err
	}
	t, err := do.Invoke[transport.Transport](inj)
	if err != nil {
		return nil, err
	}
	b, err := do.Invoke[bulletin.Bulletin](inj)
	if err != nil {
		return nil, err
	}
	return New(db, rkeys, t, b)
}

func (factory) Name() string {
	return name
}

func (factory) Repos() []string {
	return []string{"deals", "shares"}
}

// /orbis/dkg/rabin/{deals,shares}

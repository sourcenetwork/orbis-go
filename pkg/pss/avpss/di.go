package avpss

import (
	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

var (
	_       types.Factory[pss.PSS] = (*factory)(nil)
	Factory                        = factory{}
)

type factory struct{}

func (factory) New(inj *do.Injector, rkeys []db.RepoKey, _ config.Config) (pss.PSS, error) {
	db, err := do.Invoke[*db.DB](inj)
	if err != nil {
		return nil, err
	}

	tp, err := do.Invoke[transport.Transport](inj)
	if err != nil {
		return nil, err
	}

	bb, err := do.Invoke[bulletin.Bulletin](inj)
	if err != nil {
		return nil, err
	}

	dkg, err := do.Invoke[dkg.DKG](inj)
	if err != nil {
		return nil, err
	}

	return New(db, rkeys, tp, bb, dkg)
}

func (factory) Name() string {
	return name
}

func (factory) Repos() []string {
	return []string{}
}

package elgamal

import (
	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/pre"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

// func Provider(i *do.Injector) (pre.Factory, error) {
// 	return factory{}, nil
// }

var _ types.Factory[pre.PRE] = (*factory)(nil)

var (
	Factory = factory{}
)

type factory struct{}

func (factory) New(inj *do.Injector, rkeys []db.RepoKey) (pre.PRE, error) {
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
	dkg, err := do.Invoke[dkg.DKG](inj)
	if err != nil {
		return nil, err
	}

	return New(db, rkeys, t, b, dkg)
}

func (factory) Name() string {
	return name
}

func (factory) Repos() []db.Record {
	return []db.Record{}
}

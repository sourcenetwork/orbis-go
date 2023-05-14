package rabin

import (
	"github.com/samber/do"
	rabinv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/rabin/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	orbisdkg "github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

var (
	_       types.Factory[orbisdkg.DKG] = (*factory)(nil)
	Factory                             = factory{}
)

type factory struct{}

func (factory) New(inj *do.Injector, rkeys []db.RepoKey) (orbisdkg.DKG, error) {
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

func (factory) Repos() []db.Record {
	return []db.Record{&rabinv1alpha1.Deal{}, &rabinv1alpha1.Response{}}
}

// /rabin/{deals,shares}

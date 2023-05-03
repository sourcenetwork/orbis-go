package p2p

import (
	"context"

	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

var (
	Factory = factory{}
)

type factory struct{}

func (factory) New(ctx context.Context, inj *do.Injector, cfg config.Transport) (transport.Transport, error) {
	return New(ctx, inj, cfg)
}

func (factory) Name() string {
	return name
}

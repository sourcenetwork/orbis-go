package p2p

import (
	"context"

	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
)

const ProviderName = "p2pbb"

func Provider(i *do.Injector) (bulletin.Factory, error) {
	return factory{}, nil
}

type factory struct{}

func (factory) New(ctx context.Context, inj *do.Injector, cfg config.Bulletin) (bulletin.Bulletin, error) {
	return New(ctx, inj, cfg)
}

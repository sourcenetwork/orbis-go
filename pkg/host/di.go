package host

import (
	"context"

	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/config"
)

const name = "libp2p"

func Provider(i *do.Injector) (Factory, error) {
	return factory{}, nil
}

type Factory interface {
	Name() string
	New(ctx context.Context, cfg config.P2P) (*Host, error)
}

type factory struct{}

func (factory) New(ctx context.Context, cfg config.P2P) (*Host, error) {
	return New(ctx, cfg)
}

func (factory) Name() string {
	return name
}

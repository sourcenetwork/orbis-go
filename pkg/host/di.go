package host

import (
	"context"

	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/config"
)

const name = "libp2p"

type Factory struct{}

func (Factory) New(ctx context.Context, inj *do.Injector, cfg config.Host) (*Host, error) {
	return New(ctx, cfg)
}

func (Factory) Name() string {
	return name
}

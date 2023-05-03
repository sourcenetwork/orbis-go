package app

import (
	"context"
	"fmt"

	"github.com/samber/do"

	"github.com/sourcenetwork/orbis-go/pkg/host"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

// App implements App all services.
type App struct {
	p  *host.Host
	tp transport.Transport

	inj *do.Injector
}

func (a *App) Host() *host.Host {
	return a.p
}
func (a *App) Transport() transport.Transport {
	return a.tp
}

func (a *App) Injector() *do.Injector {
	return a.inj
}

func New(ctx context.Context, opts ...Option) (*App, error) {

	a := &App{
		inj: do.New(),
	}

	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, fmt.Errorf("apply orbis option: %w", err)
		}
	}

	return a, nil
}

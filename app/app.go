package app

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/orbis-go/infra/logger"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

// App implements Orbis all services.
type App struct {
	lg  logger.Logger
	bn  bulletin.Bulletin
	pss pss.PSS
	tp  transport.Transport
}

func (a *App) Logger() logger.Logger {
	return a.lg
}

func (a *App) Transport() transport.Transport {
	return a.tp
}

func New(ctx context.Context, opts ...Option) (*App, error) {

	app := &App{}

	for _, opt := range opts {
		err := opt(app)
		if err != nil {
			return nil, fmt.Errorf("apply app option: %w", err)
		}
	}

	return app, nil
}

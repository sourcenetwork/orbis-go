package app

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/orbis-go/infra/logger"
)

// App implements Orbis all services.
type App struct {
	l logger.Logger
}

func (a *App) Logger() logger.Logger {
	return a.l
}

func New(ctx context.Context, opts ...Option) (*App, error) {

	var cfg Config

	// Apply options to the config .
	for _, opt := range opts {
		err := opt(&cfg)
		if err != nil {
			return nil, fmt.Errorf("apply app option: %w", err)
		}
	}

	app := &App{}

	// Process config.
	app.l = cfg.logger

	return app, nil
}

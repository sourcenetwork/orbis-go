package main

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/orbis-go/app"
	"github.com/sourcenetwork/orbis-go/config"
)

func setupApp(ctx context.Context, cfg config.Config) (*app.App, error) {

	lg, err := setupLogger(ctx, cfg.Logger)
	if err != nil {
		return nil, err
	}

	t, err := setupTransport(ctx, lg, cfg.Transport)
	if err != nil {
		return nil, err
	}

	opts := []app.Option{
		app.DefaultOptions(),
		app.WithLogger(lg),
		app.WithTransport(t),
	}

	app, err := app.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("create app: %w", err)
	}

	return app, nil
}

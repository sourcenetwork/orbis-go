package main

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/orbis-go/app"
	"github.com/sourcenetwork/orbis-go/config"
)

func setupApp(ctx context.Context, cfg config.Config) (*app.App, error) {

	opts := configToAppOptions(cfg)

	app, err := app.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("create app: %w", err)
	}

	return app, nil
}

func configToAppOptions(cfg config.Config) []app.Option {

	opts := []app.Option{
		app.DefaultOptions(),
		app.WithLogger(cfg.Logger),
	}

	return opts
}

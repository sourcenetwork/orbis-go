package main

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/infra/logger"
	"github.com/sourcenetwork/orbis-go/infra/logger/zap"
)

func setupLogger(ctx context.Context, cfg config.Logger) (logger.Logger, error) {

	switch cfg.Logger {
	case "zap":
		return zap.New(cfg), nil
	default:
		return nil, fmt.Errorf("logger not supported: %s", cfg.Logger)
	}
}

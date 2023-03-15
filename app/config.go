package app

import (
	"fmt"

	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/infra/logger"
	"github.com/sourcenetwork/orbis-go/infra/logger/zap"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
)

type Config struct {
	logger   logger.Logger
	bulletin bulletin.Bulletin
	pss      pss.PSS
}

type Option func(c *Config) error

func DefaultOptions() Option {
	return func(c *Config) error {
		return nil
	}
}

func WithLogger(cfg config.Logger) Option {
	return func(c *Config) error {
		switch cfg.Logger {
		case "zap":
			c.logger = zap.New(cfg)
		default:
			return fmt.Errorf("logger not supported: %s", cfg.Logger)
		}
		return nil
	}
}

// WithBulletinService registers a BulletinBoard Service into the config
func WithBulletinService(b bulletin.Bulletin) Option {
	return func(c *Config) error {
		c.bulletin = b
		return nil
	}
}

// WithSharingService regisers a Proactive Sharing Service into the config
func WithSharingService(pss pss.PSS) Option {
	return func(c *Config) error {
		c.pss = pss
		return nil
	}
}

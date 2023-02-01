package orbis

import (
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
)

type Config struct{}

type Option func(cfg *Config) error

func DefaultOptions() Option {
	panic("todo")
}

// WithBulletinService registers a BulletinBoard Service into the config
func WithBulletinService(bulletin.ProviderFn) Option

// WithSharingService regisers a Proactive Sharing Service into the config
func WithSharingService(pss.ProviderFn) Option

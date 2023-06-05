package app

import (
	"fmt"

	"github.com/samber/do"

	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/host"
	"github.com/sourcenetwork/orbis-go/pkg/pre"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

type Option func(a *App) error

type Provider[T any] func(*do.Injector)

func DefaultOptions() Option {
	return func(o *App) error {
		return nil
	}
}

func WithHost(f *host.Host) Option {
	return func(a *App) error {
		do.ProvideValue(a.inj, f)
		return nil
	}
}

func WithDBData(path string) Option {
	return func(a *App) error {
		d, err := db.New(path)
		if err != nil {
			return fmt.Errorf("create db: %w", err)
		}
		do.ProvideValue(a.inj, d)
		a.db = d
		return nil
	}
}

func WithTransport(f transport.Transport) Option {
	return func(a *App) error {
		do.ProvideNamedValue(a.inj, f.Name(), f)
		return nil
	}
}

// WithBulletin registers BulletinBoard factory.
func WithBulletin(f bulletin.Bulletin) Option {
	return func(a *App) error {
		do.ProvideNamedValue(a.inj, f.Name(), f)
		return nil
	}
}

func WithDistKeyGenerator(f types.Factory[dkg.DKG]) Option {
	return func(a *App) error {
		if f.Name() == "" {
			return ErrFactoryEmptyName
		}
		do.ProvideNamedValue(a.inj, f.Name(), f)
		return a.setupRepoKeysForService(f.Name(), f.Repos())
	}
}

// WithProxyReencryption registers o Proxy-Reencryption Service.
func WithProxyReencryption(f types.Factory[pre.PRE]) Option {
	return func(a *App) error {
		do.ProvideNamedValue(a.inj, f.Name(), f)
		return a.setupRepoKeysForService(f.Name(), f.Repos())
	}
}

// WithProactiveSecretSharing registers o Proactive Secret Sharing Service.
func WithProactiveSecretSharing(f types.Factory[pss.PSS]) Option {
	return func(a *App) error {
		do.ProvideNamedValue(a.inj, f.Name(), f)
		return a.setupRepoKeysForService(f.Name(), f.Repos())
	}
}

package app

import (
	"github.com/samber/do"

	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/host"
	"github.com/sourcenetwork/orbis-go/pkg/pre"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

// type Factory[T any] interface {
// 	New(*do.Injector, []*db.RepoKey) (T, error)
// 	Name() string
// 	Repos() []string
// }

type Option func(a *App) error

type Provider[T any] func(*do.Injector)

func DefaultOptions() Option {
	return func(o *App) error {
		return nil
	}
}

func WithHost(f host.Factory) Option {
	return func(a *App) error {
		do.ProvideNamedValue(a.inj, f.Name(), f)
		return nil
	}
}

func WithTransport(f transport.Factory) Option {
	return func(a *App) error {
		do.ProvideNamedValue(a.inj, f.Name(), f)
		return nil
	}
}

// WithBulletin registers BulletinBoard factory.
func WithBulletin(f bulletin.Factory) Option {
	return func(a *App) error {
		do.ProvideNamedValue(a.inj, f.Name(), f)
		return nil
	}
}

func WithDistKeyGenerator(f types.Factory[dkg.DKG]) Option {
	return func(a *App) error {
		do.ProvideNamedValue(a.inj, f.Name(), f)
		return nil
	}
}

// WithProxyReencryption registers o Proxy-Reencryption Service.
func WithProxyReencryption(f pre.Factory) Option {
	return func(a *App) error {
		do.ProvideNamedValue(a.inj, f.Name(), f)
		return nil
	}
}

// WithProactiveSecretSharing registers o Proactive Secret Sharing Service.
func WithProactiveSecretSharing(f pss.Factory) Option {
	return func(a *App) error {
		do.ProvideNamedValue(a.inj, f.Name(), f)
		return nil
	}
}

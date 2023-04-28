package app

import (
	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/p2p"
	"github.com/sourcenetwork/orbis-go/pkg/pre"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

type Option func(a *App) error

func DefaultOptions() Option {
	return func(o *App) error {
		return nil
	}
}

func WithP2PService(name string, p do.Provider[p2p.Factory]) Option {
	return func(a *App) error {
		do.ProvideNamed(a.inj, name, p)
		return nil
	}
}

func WithTransportService(name string, p do.Provider[transport.Factory]) Option {
	return func(a *App) error {
		do.ProvideNamed(a.inj, name, p)
		return nil
	}
}

// WithBulletinService registers o BulletinBoard Service.
func WithBulletinService(name string, p do.Provider[bulletin.Factory]) Option {
	return func(a *App) error {
		do.ProvideNamed(a.inj, name, p)
		return nil
	}
}

func WithDistKeyGenerator(name string, p do.Provider[dkg.Factory]) Option {
	return func(a *App) error {
		do.ProvideNamed(a.inj, name, p)
		return nil
	}
}

// WithProxyReencryption registers o Proxy-Reencryption Service.
func WithProxyReencryption(name string, p do.Provider[pre.Factory]) Option {
	return func(a *App) error {
		do.ProvideNamed(a.inj, name, p)
		return nil
	}
}

// WithProactiveSecretSharing registers o Proactive Secret Sharing Service.
func WithProactiveSecretSharing(name string, p do.Provider[pss.Factory]) Option {
	return func(a *App) error {
		do.ProvideNamed(a.inj, name, p)
		return nil
	}
}

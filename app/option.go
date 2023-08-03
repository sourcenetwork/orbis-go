package app

import (
	"fmt"

	"github.com/samber/do"

	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	"github.com/sourcenetwork/orbis-go/pkg/host"
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

// func WithAuthz(az authz.Authz) Option {
// 	return func(a *App) error {
// 		do.ProvideNamedValue(a.inj, az.Name(), az)
// 		return nil
// 	}
// }

// func WithKeyResolver(r authn.KeyResolver) Option {
// 	return func(a *App) error {
// 		a.resolver = r
// 		return nil
// 	}
// }

func WithService[S any](s S) Option {
	return func(a *App) error {
		// quick runtime "implements" check for the Named interface
		// Since generic type constraint parameters cant use type
		// assertions. But as you can see, the original type
		// is used in the `ProvideNamedValue` call.
		if n, ok := interface{}(s).(types.Named); ok {
			do.ProvideNamedValue(a.inj, n.Name(), s)
		} else {
			do.ProvideValue(a.inj, s)
		}
		return nil
	}
}

// func WithStatelessFactory[T any](f types.StatelessFactory[T]) {}

func WithFactory[T any](f types.Factory[T]) Option {
	return func(a *App) error {
		if f.Name() == "" {
			return ErrFactoryEmptyName
		}
		do.ProvideNamedValue(a.inj, f.Name(), f)
		return a.setupRepoKeysForService(f.Name(), f.Repos())
	}
}

// func WithDistKeyGeneratorFactory(f types.Factory[dkg.DKG]) Option {
// 	return func(a *App) error {
// 		if f.Name() == "" {
// 			return ErrFactoryEmptyName
// 		}
// 		do.ProvideNamedValue(a.inj, f.Name(), f)
// 		return a.setupRepoKeysForService(f.Name(), f.Repos())
// 	}
// }

// // WithProxyReencryption registers o Proxy-Reencryption Service.
// func WithProxyReencryptionFactory(f types.Factory[pre.PRE]) Option {
// 	return func(a *App) error {
// 		do.ProvideNamedValue(a.inj, f.Name(), f)
// 		return a.setupRepoKeysForService(f.Name(), f.Repos())
// 	}
// }

// // WithProactiveSecretSharing registers o Proactive Secret Sharing Service.
// func WithProactiveSecretSharingFactory(f types.Factory[pss.PSS]) Option {
// 	return func(a *App) error {
// 		do.ProvideNamedValue(a.inj, f.Name(), f)
// 		return a.setupRepoKeysForService(f.Name(), f.Repos())
// 	}
// }

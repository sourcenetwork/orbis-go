package app

import (
	"fmt"

	"github.com/samber/do"

	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

type Option func(a *App) error

func DefaultOptions(cfg config.Config) Option {
	return func(a *App) error {
		a.config = cfg
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

// WithService registers a service into the dependency injection system
// so it can be used globally. Service values are shared, and not
// (re)initialized on each `Invoke` call, so whatever type used here must
// assume it will be shared by all ring instances. So it must be thread safe
func WithService[S any](s S) Option {
	return func(a *App) error {
		// quick runtime "implements" check for the Named interface
		// Since generic type constraint parameters cant use type
		// assertions. But as you can see, the original type
		// is used in the `ProvideNamedValue` call so we keep the
		// concrete type.
		if n, ok := any(s).(types.Named); ok {
			do.ProvideNamedValue(a.inj, n.Name(), s)
		} else {
			do.ProvideValue(a.inj, s)
		}
		return nil
	}
}

// WithStaticService is like `WithService` but only one of a given
// type can be injected, and it doesn't used a named injection.
func WithStaticService[S any](s S) Option {
	return func(a *App) error {
		do.ProvideValue(a.inj, s)
		return nil
	}
}

// func WithStatelessFactory[T any](f types.StatelessFactory[T]) {}

// WithFactory registers a `types.Factory[T]` into the dependency
// injection system. The Factory type itself should be immutable, so
// that it can be used statelessly. The newly created instances from
// factory are not shared, and are bound to a single ring object, so
// they do not need to be thread safe.
func WithFactory[T any](f types.Factory[T]) Option {
	return func(a *App) error {
		if f.Name() == "" {
			return ErrFactoryEmptyName
		}
		do.ProvideNamedValue(a.inj, f.Name(), f)
		return a.setupRepoKeysForService(f.Name(), f.Repos())
	}
}

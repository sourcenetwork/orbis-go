package types

import (
	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/pkg/db"
)

type Named interface {
	Name() string
}

type Factory[T any] interface {
	Named
	New(*do.Injector, []db.RepoKey) (T, error)
	Repos() []string
}

type StatelessFactory[T any] interface {
	Factory[T]
	// New(*do.Injector) (T, error)
}

type StatefulFactory[T any] interface {
	Factory[T]
	New(*do.Injector, []db.RepoKey) (T, error)
	Repos() []string
}

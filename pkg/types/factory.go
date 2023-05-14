package types

import (
	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/pkg/db"
)

type Factory[T any] interface {
	New(*do.Injector, []db.RepoKey) (T, error)
	Name() string
	Repos() []db.Record
}

package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-bond/bond"
	logging "github.com/ipfs/go-log"
)

// type DB interface {
// 	Rings() Repository[*types.Ring]
// 	Secrets() Repository[*types.Secret]
// }

var log = logging.Logger("orbis/db")

type RepoKey interface {
	Name() string
}

type repoKey struct {
	name string
}

func (rk *repoKey) Name() string {
	return rk.name
}

// NewRepoKey returns a new pointer to
// a RepoKey.
func NewRepoKey(name string) RepoKey {
	if name == "" {
		panic("empty key name not allowed")
	}
	return &repoKey{name}
}

type DB struct {
	bond  bond.DB
	repos map[RepoKey]any // map[tableKey]Repository
}

func GetRepo[R Record](db *DB, rkey RepoKey, pkFunc RepoPrimaryKeyFunc[R]) (Repository[R], error) {
	if rkey == nil {
		return nil, ErrRepoKeyInvalid
	}
	repo, ok := db.repos[rkey]
	if !ok {
		repo = newSimpleRepo[R](db.bond, pkFunc)
		db.repos[rkey] = repo
	}

	repoTyped, ok := repo.(Repository[R])
	if !ok {
		return nil, fmt.Errorf("repo type doesn't match")
	}
	return repoTyped, nil
}

func New(path string) (*DB, error) {
	opts := bond.DefaultOptions()

	dirname, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path = filepath.Join(dirname, ".orbis", path) //todo: Parameterize path
	bdb, err := bond.Open(path, opts)
	if err != nil {
		return nil, err
	}

	return &DB{
		bond:  bdb,
		repos: make(map[RepoKey]any),
	}, nil
}

// MountRepo will create and mount a typed repo under the given
// RepoKey
//
// The last parameter ...T is an *optional* generic type
// inference, if the call site can't be explicit
// func MountRepo[T Record](db *DB, key RepoKey, _ ...T) error {
// 	if key == nil {
// 		return ErrRepoKeyInvalid
// 	}
// 	log.Debugf("db.MountRepo(): mounting repo %s", key.Name())
// 	if _, exists := db.repos[key]; exists {
// 		return ErrDuplicateKey
// 	}
// 	repo := newSimpleRepo[T](db.bond)
// 	db.repos[key] = repo
// 	return nil
// }

// func thing() {
// 	// db, _ := New()
// 	// table, err := db.Get("users")
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	db.Repo[types.Deals](db, repoKey)
// }

// type ctxKey string

// var (
// 	dbCtxKey = ctxKey("db")
// )

// func DBFromContext(ctx context.Context) (*DB, bool) {
// 	db, ok := ctx.Value(dbCtxKey).(*DB)
// 	if !ok {
// 		return nil, false
// 	}
// 	return db, true
// }

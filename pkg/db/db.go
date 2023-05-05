package db

import (
	"fmt"

	"github.com/go-bond/bond"
)

// type DB interface {
// 	Rings() Repository[*types.Ring]
// 	Secrets() Repository[*types.Secret]
// }

type RepoKey struct {
	name string
}

type DB struct {
	bond  bond.DB
	repos map[*RepoKey]any // map[tableKey]Repository
}

func GetRepo[R Record](db DB, rkey *RepoKey) (Repository[R], error) {
	repo, ok := db.repos[rkey]
	if !ok {
		return nil, fmt.Errorf("no repo exists for %s", rkey.name)
	}

	repoTyped, ok := repo.(Repository[R])
	if !ok {
		return nil, fmt.Errorf("repo type doesn't match")
	}
	return repoTyped, nil
}

func New() (*DB, error) {
	opts := bond.DefaultOptions()
	opts.Serializer = protoSerializer{}
	bdb, err := bond.Open("~/.orbis/data", opts) //todo: Parameterize location
	if err != nil {
		return nil, err
	}

	return &DB{
		bond: bdb,
	}, nil
}

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

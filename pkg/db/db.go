package db

import (
	"context"

	"github.com/go-bond/bond"
	"github.com/gogo/protobuf/proto"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

const (
	ringRepoTableID   = 1
	ringRepoTableName = "ring_repo_table"
)

type Record interface {
	proto.Message
	GetId() string
}

type Repository struct {
	Rings   *simpleRepo[*types.Ring]
	Secrets *simpleRepo[*types.Secret]

	// map[string]simpleRepo[T]
	// secretRepos map[string]any
}

func New() (*Repository, error) {
	opts := bond.DefaultOptions()
	opts.Serializer = protoSerializer{}
	bdb, err := bond.Open("~/.orbis/data", opts)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Rings:   newSimpleRepo[*types.Ring](bdb),
		Secrets: newSimpleRepo[*types.Secret](bdb),
	}, nil
}

type simpleRepo[T Record] struct {
	table bond.Table[T]
}

func newSimpleRepo[T Record](bdb bond.DB) *simpleRepo[T] {
	rr := &simpleRepo[T]{}
	rr.table = bond.NewTable(bond.TableOptions[T]{
		DB:        bdb,
		TableID:   ringRepoTableID,
		TableName: ringRepoTableName,
		TablePrimaryKeyFunc: func(b bond.KeyBuilder, t T) []byte {
			return []byte(t.GetId())
		},
	})
	return rr
}

func (rr *simpleRepo[T]) Create(ctx context.Context, t T) error {
	return rr.table.Insert(ctx, []T{t})
}

func (rr *simpleRepo[T]) Get(ctx context.Context, t T) (T, error) {
	return rr.table.Get(t)
}

func (rr *simpleRepo[T]) GetAll(ctx context.Context) ([]T, error) {
	var ts *[]T
	if err := rr.table.Query().Execute(ctx, ts); err != nil {
		return nil, err
	}
	return *ts, nil
}

type Repo[T proto.Message] interface {
	Create(T) error
	Get(string) (T, error)
	GetAll() ([]T, error)
}

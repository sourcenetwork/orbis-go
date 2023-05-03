package db

import (
	"context"
	"fmt"
	"hash/fnv"

	"github.com/go-bond/bond"
)

const (
	ringRepoTableID = 1
)

type Repository[T Record] interface {
	Create(context.Context, T) error
	Get(context.Context, T) (T, error)
	GetAll(context.Context) ([]T, error)
}

type simpleRepo[T Record] struct {
	table bond.Table[T]
}

func newSimpleRepo[T Record](bdb bond.DB) *simpleRepo[T] {
	var t T
	name := getTableName(t)
	rr := &simpleRepo[T]{}
	rr.table = bond.NewTable(bond.TableOptions[T]{
		DB:        bdb,
		TableID:   bond.TableID(hash(name)),
		TableName: name,
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

func getTableName(r Record) string {
	return fmt.Sprintf("%T", r)
}

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

package db

import (
	"context"

	"github.com/go-bond/bond"
)

type Batch = bond.Batch

type FilterFunc[R any] func(r R) bool
type OrderLessFunc[R any] func(r, r2 R) bool

type Query[R any] interface {
	After(R) Query[R]
	Filter(FilterFunc[R]) Query[R]
	Limit(uint64) Query[R]
	Offset(uint64) Query[R]
	Order(OrderLessFunc[R]) Query[R]

	Execute(context.Context) ([]R, error)
}

type rawQuery[R any] struct {
	bondQuery bond.Query[R]
}

func (q rawQuery[R]) After(r R) Query[R] {
	return rawQuery[R]{q.bondQuery.After(r)}
}

func (q rawQuery[R]) Filter(filter FilterFunc[R]) Query[R] {
	return rawQuery[R]{q.bondQuery.Filter(bond.FilterFunc[R](filter))}
}

func (q rawQuery[R]) Limit(limit uint64) Query[R] {
	return rawQuery[R]{q.bondQuery.Limit(limit)}
}

func (q rawQuery[R]) Offset(offset uint64) Query[R] {
	return rawQuery[R]{q.bondQuery.Offset(offset)}
}

func (q rawQuery[R]) Order(less OrderLessFunc[R]) Query[R] {
	return rawQuery[R]{q.bondQuery.Order(bond.OrderLessFunc[R](less))}
}

func (q rawQuery[R]) Execute(ctx context.Context) ([]R, error) {
	records := new([]R)
	err := q.bondQuery.Execute(ctx, records)
	return *records, err
}

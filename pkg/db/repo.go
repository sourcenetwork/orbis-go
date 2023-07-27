package db

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/go-bond/bond"
	"google.golang.org/protobuf/proto"
)

var (
	_ Repository[Record] = (*simpleRepo[Record])(nil)
)

type Record interface {
	proto.Message
}

type RingIDGetter interface {
	GetRingID() string
}

type RepoPrimaryKeyFunc[T any] bond.TablePrimaryKeyFunc[T]

type Repository[T Record] interface {
	Create(context.Context, T) error
	CreateMany(context.Context, []T) error
	Save(context.Context, T) error
	Update(context.Context, T) error
	Get(context.Context, T) (T, error)
	GetAll(context.Context) ([]T, error)
	Query() Query[T]
	Exists(context.Context, T) bool
}

type simpleRepo[T Record] struct {
	table bond.Table[T]
}

type RepoOption func(*RepoOptions)

type RepoOptions struct {
	ringID string
}

func WithRingID(rid string) RepoOption {
	return func(opt *RepoOptions) {
		opt.ringID = rid
	}
}

// func defaultPrimaryKeyFunc[T Record](b bond.KeyBuilder, t T) []byte {
// 	// if possible, primary keys are:
// 	// /<ring_id>/<record_id>
// 	// otherwise:
// 	// /<record_id>
// 	if getter, ok := any(t).(RingIDGetter); ok {
// 		b = b.AddStringField(getter.GetRingID())
// 	}
// 	b = b.AddStringField(t.GetId())
// 	return b.Bytes()
// }

func newSimpleRepo[T Record](bdb bond.DB, pkFunc RepoPrimaryKeyFunc[T]) *simpleRepo[T] {
	if pkFunc == nil {
		panic("invalid primary key function, cant be nil")
	}
	var t T
	name := getTableName(t)
	rr := &simpleRepo[T]{}
	rr.table = bond.NewTable(bond.TableOptions[T]{
		DB:                  bdb,
		TableID:             bond.TableID(hash(name)),
		TableName:           name,
		TablePrimaryKeyFunc: bond.TablePrimaryKeyFunc[T](pkFunc),
		Serializer:          protoSerializer[T]{},
	})
	return rr
}

func (rr *simpleRepo[T]) Create(ctx context.Context, t T) error {
	err := rr.table.Insert(ctx, []T{t})
	if err != nil && strings.Contains(err.Error(), "already exists") {
		return errors.Join(ErrRecordAlreadyExists, err)
	} else if err != nil {
		return fmt.Errorf("repo create: %w", err)
	}
	return nil
}

func (rr *simpleRepo[T]) CreateMany(ctx context.Context, ts []T) error {
	err := rr.table.Insert(ctx, ts)
	if err != nil && strings.Contains(err.Error(), "already exists") {
		return errors.Join(ErrRecordAlreadyExists, err)
	} else if err != nil {
		return fmt.Errorf("repo create many: %w", err)
	}
	return nil
}

func (rr *simpleRepo[T]) Save(ctx context.Context, t T) error {
	log.Debugf("Saving repo entry for %T", t)
	if rr.table.Exist(t) {
		log.Debug("entry exists, updating")
		return rr.Update(ctx, t)
	}
	log.Debug("entry doesn't exist, creating")
	return rr.Create(ctx, t)
}

func (rr *simpleRepo[T]) Update(ctx context.Context, t T) error {
	return rr.table.Update(ctx, []T{t})
}

func (rr *simpleRepo[T]) Get(ctx context.Context, t T) (T, error) {
	var zeroT T
	ts, err := rr.table.Get(ctx, bond.NewSelectorPoint(t))
	if err != nil {
		return zeroT, err
	}
	return ts[0], nil
}

func (rr *simpleRepo[T]) GetAll(ctx context.Context) ([]T, error) {
	var ts []T
	if err := rr.table.Query().Execute(ctx, &ts); err != nil {
		return nil, err
	}
	return ts, nil
}

func (rr *simpleRepo[T]) Query() Query[T] {
	return rawQuery[T]{rr.table.Query()}
}

func (rr *simpleRepo[T]) Exists(ctx context.Context, t T) bool {
	return rr.table.Exist(t)
}

func getTableName(r Record) string {
	return string(r.ProtoReflect().Descriptor().Name())
}

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

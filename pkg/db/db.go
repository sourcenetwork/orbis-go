package db

import (
	"context"

	"github.com/go-bond/bond"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

const (
	ringRepoTableID   = 1
	ringRepoTableName = "ring_repo_table"
)

type Repository struct {
	Rings RingRepo
	// Secrets
}

type ringRepo struct {
	table bond.Table[*types.Ring]
}

func newRingRepo(bdb bond.DB) *ringRepo {
	rr := &ringRepo{}
	rr.table = bond.NewTable(bond.TableOptions[*types.Ring]{
		DB:        bdb,
		TableID:   ringRepoTableID,
		TableName: ringRepoTableName,
		TablePrimaryKeyFunc: func(b bond.KeyBuilder, r *types.Ring) []byte {
			return []byte(r.Id)
		},
	})
	return rr
}

func (rr *ringRepo) Create(ctx context.Context, r *types.Ring) error {
	return rr.table.Insert(ctx, []*types.Ring{r})
}

type RingRepo interface {
	Create(types.Ring) error
	Get(types.RingID) (types.Ring, error)
	GetAll() ([]types.Ring, error)
}

type SecretRepo interface{}

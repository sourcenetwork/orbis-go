package db

import (
	"github.com/go-bond/bond"
	"github.com/gogo/protobuf/proto"

	"github.com/sourcenetwork/orbis-go/pkg/types"
)

type Record interface {
	proto.Message
	GetId() string
}

type DB interface {
	Rings() Repository[*types.Ring]
	Secrets() Repository[*types.Secret]
}

type simpleDB struct {
	ringRepo    *simpleRepo[*types.Ring]
	secretsRepo *simpleRepo[*types.Secret]
}

func New() (DB, error) {
	opts := bond.DefaultOptions()
	opts.Serializer = protoSerializer{}
	bdb, err := bond.Open("~/.orbis/data", opts) //todo: Parameterize location
	if err != nil {
		return nil, err
	}

	return &simpleDB{
		ringRepo:    newSimpleRepo[*types.Ring](bdb),
		secretsRepo: newSimpleRepo[*types.Secret](bdb),
	}, nil
}

func (sdb *simpleDB) Rings() Repository[*types.Ring] {
	return sdb.ringRepo
}

func (sdb *simpleDB) Secrets() Repository[*types.Secret] {
	return sdb.secretsRepo
}

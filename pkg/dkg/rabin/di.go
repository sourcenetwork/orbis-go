package rabin

import (
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	orbisdkg "github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

// func Provider(i *do.Injector) (orbisdkg.Factory, error) {
// 	return factory{}, nil
// }

var (
	Factory = factory{}
)

type factory struct{}

func (factory) New(repo db.Repository, t transport.Transport, b bulletin.Bulletin, pk crypto.PrivateKey) (orbisdkg.DKG, error) {
	return New(repo, t, b, pk)
}

func (factory) Name() string {
	return name
}

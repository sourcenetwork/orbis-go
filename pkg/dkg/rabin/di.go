package rabin

import (
	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	orbisdkg "github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

const ProviderName = "rabin"

func Provider(i *do.Injector) (orbisdkg.Factory, error) {
	return factory{}, nil
}

type factory struct{}

func (factory) New(repo db.Repository, t transport.Transport, b bulletin.Bulletin, pk crypto.PrivateKey) (orbisdkg.DKG, error) {
	return New(repo, t, b, pk)
}

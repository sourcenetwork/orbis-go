package elgamal

import (
	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/pre"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

const ProviderName = "elgamal"

func Provider(i *do.Injector) (pre.Factory, error) {
	return factory{}, nil
}

type factory struct{}

func (factory) New(rid types.RingID, n int32, t int32, tp transport.Transport, bb bulletin.Bulletin, nodes []types.Node, dkg dkg.DKG) (pre.PRE, error) {
	return New(rid, n, t, tp, bb, nodes, dkg)
}

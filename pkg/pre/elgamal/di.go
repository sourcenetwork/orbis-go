package elgamal

import (
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/pre"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

// func Provider(i *do.Injector) (pre.Factory, error) {
// 	return factory{}, nil
// }

var (
	Factory = factory{}
)

type factory struct{}

func (factory) New(rid types.RingID, n int32, t int32, tp transport.Transport, bb bulletin.Bulletin, nodes []types.Node, dkg dkg.DKG) (pre.PRE, error) {
	return New(rid, n, t, tp, bb, nodes, dkg)
}

func (factory) Name() string {
	return name
}

package avpss

import (
	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

const ProviderName = "elgamal"

func Provider(i *do.Injector) (pss.Factory, error) {
	return factory{}, nil
}

type factory struct{}

func (factory) New(rid types.RingID, n int32, t int32, tp transport.Transport, bb bulletin.Bulletin, nodes []types.Node, d dkg.DKG) (pss.PSS, error) {
	return New(rid, n, t, tp, bb, nodes, d)
}

func New(types.RingID, int32, int32, transport.Transport, bulletin.Bulletin, []types.Node, dkg.DKG) (*AVPSS, error) {
	return &AVPSS{}, nil
}

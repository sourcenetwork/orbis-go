package avpss

import (
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

const ProviderName = "avpss"

var (
	Factory = factory{}
)

type factory struct{}

func (factory) New(rid types.RingID, n int32, t int32, tp transport.Transport, bb bulletin.Bulletin, nodes []types.Node, d dkg.DKG) (pss.PSS, error) {
	return New(rid, n, t, tp, bb, nodes, d)
}

func (factory) Name() string {
	return name
}

package avpss

import (
	"context"

	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
	"go.dedis.ch/kyber/v3/suites"
)

const name = "avpss"

type AVPSS struct {
}

func New(types.RingID, int32, int32, transport.Transport, bulletin.Bulletin, []types.Node, dkg.DKG) (*AVPSS, error) {
	return &AVPSS{}, nil
}

func (a *AVPSS) Name() string {
	return name
}

func (a *AVPSS) Suite() suites.Suite {
	return nil
}

func (a *AVPSS) Start() {

}

func (a *AVPSS) Shutdown() error {
	return nil
}

func (a *AVPSS) ProcessMessage(context.Context, pss.Message) {
}

func (a *AVPSS) PublicKey() crypto.PublicKey {
	return nil
}
func (a *AVPSS) PublicPoly() crypto.PubPoly {
	return crypto.PubPoly{}
}
func (a *AVPSS) Share() crypto.PriShare {
	return crypto.PriShare{}
}
func (a *AVPSS) State() pss.State {
	return pss.State{}
}
func (a *AVPSS) Num() int {
	return 0
}
func (a *AVPSS) Threshold() int {
	return 0
}

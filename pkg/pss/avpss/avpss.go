package avpss

import (
	"context"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"go.dedis.ch/kyber/v3/suites"
)

type AVPSS struct {
}

func (a *AVPSS) Name() string {
	return "avpss"
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

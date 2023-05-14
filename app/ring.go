package app

import (
	"context"
	"fmt"

	logging "github.com/ipfs/go-log"
	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/crypto/proof"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/pre"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

var log = logging.Logger("orbis/ring")

type Ring struct {
	ID types.RingID

	DKG dkg.DKG
	PSS pss.PSS
	PRE pre.PRE

	Transport transport.Transport
	Bulletin  bulletin.Bulletin
	DB        *db.DB
}

/*
ring1
manifest := {
	"N": 9,
	"T": 7,
	"curve": "Ed25519",

	"dkg": "rabin",
	"pss": "avpss",
	"pre": "elgamal",
	"bulletin": "sourcehub",
	"transport": "libp2p"
}

ring2
manifest := {
	"N": 9,
	"T": 7,
	"curve": "Ed25519",

	"dkg": "rabin",
	"pss": "avpss",
	"pre": "elgamal",
	"bulletin": "sourcehub",
	"transport": "libp2p"
}

ring3
manifest := {
	"N": 9,
	"T": 7,
	"curve": "Ed25519",

	"dkg": "rabin",
	"pss": "avpss",
	"pre": "elgamal",
	"bulletin": "cosmos:sourcehub",
	"transport": "libp2p"
}
*/

/*

nodeconfig {
	bulletin:
		cosmos:
			sourcehub:
				rpcURL: 'localhost:1234/rpc'
				chainID: sourcehub-1
}

*/

func (app *App) NewRing(ctx context.Context, ring *types.Ring) (*Ring, error) {
	rid := types.RingID(ring.Id)

	// rings get their own cloned dependency injector handler
	// since we actually create and initialize services which
	// are scoped to rings, compared to the factories
	// which are global.
	inj := app.inj.Clone()

	// factories
	dkgFactory, err := do.InvokeNamed[types.Factory[dkg.DKG]](inj, ring.Dkg)
	if err != nil {
		return nil, fmt.Errorf("invoke dkg factory: %w", err)
	}

	pssFactory, err := do.InvokeNamed[types.Factory[pss.PSS]](inj, ring.Pss)
	if err != nil {
		return nil, fmt.Errorf("invoke pss factory: %w", err)
	}

	preFactory, err := do.InvokeNamed[types.Factory[pre.PRE]](inj, ring.Pre)
	if err != nil {
		return nil, err
	}

	// get global services
	db, err := do.Invoke[*db.DB](inj)
	if err != nil {
		return nil, fmt.Errorf("invoke db: %w", err)
	}

	tp, err := do.InvokeNamed[transport.Transport](inj, ring.Transport)
	if err != nil {
		return nil, fmt.Errorf("invoke transport: %w", err)
	}
	do.ProvideValue(inj, tp) // register configured generic transport locally

	bb, err := do.InvokeNamed[bulletin.Bulletin](inj, ring.Bulletin)
	if err != nil {
		return nil, fmt.Errorf("invoke bulletin: %w", err)
	}
	do.ProvideValue(inj, bb) // // register configured generic bulletin locally

	// setup and register local services
	dkgSrv, err := dkgFactory.New(inj, nil) // @todo repokeys
	if err != nil {
		return nil, fmt.Errorf("create dkg service: %w", err)
	}

	// get the privkey for host
	h := app.Host()
	hpk := h.Peerstore().PrivKey(h.ID())
	cpk, err := crypto.PrivateKeyFromLibP2P(hpk)
	if err != nil {
		return nil, fmt.Errorf("converting libp2p private key: %w", err)
	}

	if err := dkgSrv.Init(ctx, cpk /* private key */, nil /* nodes */, ring.N, ring.T); err != nil {
		return nil, fmt.Errorf("initializing dkg: %w", err)
	}
	do.ProvideValue(inj, dkgSrv)

	preSrv, err := preFactory.New(inj, nil) // @todo repokeys
	if err != nil {
		return nil, fmt.Errorf("create pre service: %w", err)
	}
	if err := preSrv.Init(rid, ring.N, ring.T, []types.Node{}); err != nil {
		return nil, fmt.Errorf("initializing pre: %w", err)
	}
	do.ProvideValue(inj, preSrv)

	pssSrv, err := pssFactory.New(inj, nil) // @todo repokeys
	if err != nil {
		return nil, fmt.Errorf("create pss service: %w", err)
	}
	if err := pssSrv.Init(rid, ring.N, ring.T, []types.Node{}); err != nil {
		return nil, fmt.Errorf("create pss service: %w", err)
	}

	rs := &Ring{
		ID:        rid,
		DKG:       dkgSrv,
		PSS:       pssSrv,
		PRE:       preSrv,
		Transport: tp,
		Bulletin:  bb,
		DB:        db,
	}

	// called in ring.Join() - go rs.handleEvents()

	return rs, nil
}

func (r *Ring) Store(context.Context, types.SecretID, *types.Secret, proof.VerifiableEncryption) error {
	return nil
}

func (r *Ring) Get(context.Context, types.SecretID) (types.Secret, error) {

	return types.Secret{}, nil
}

func (r *Ring) GetShares(context.Context, types.SecretID) ([]types.PrivSecretShare, error) {
	return nil, nil
}

func (r *Ring) Delete(context.Context, types.SecretID) error {
	return nil
}

func (r *Ring) PublicKey() (crypto.PublicKey, error) {
	return nil, nil

}

func (r *Ring) Refresh(context.Context, pss.Config) (pss.RefreshState, error) {
	return pss.RefreshState{}, nil
}

func (r *Ring) Threshold() int {
	return 0
}

func (r *Ring) State() pss.State {
	return pss.State{}
}

func (r *Ring) Nodes() []pss.Node {
	return nil
}

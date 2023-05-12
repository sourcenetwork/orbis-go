package ring

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
	"bulletin": "sourcehub",
	"transport": "libp2p"
}
*/

type Factory[T any] interface {
	New(*do.Injector, []*db.RepoKey) (T, error)
	Name() string
	Repos() []string
}

func NewRing(ctx context.Context, inj *do.Injector, ring *types.Ring) (*Ring, error) {

	rid := types.RingID(ring.Id)

	// factories
	dkgFactory, err := do.InvokeNamed[types.Factory[dkg.DKG]](inj, ring.Dkg)
	if err != nil {
		return nil, fmt.Errorf("invoke dkg factory: %w", err)
	}

	pssFactory, err := do.InvokeNamed[pss.Factory](inj, ring.Pss)
	if err != nil {
		return nil, fmt.Errorf("invoke pss factory: %w", err)
	}

	preFactory, err := do.InvokeNamed[pre.Factory](inj, ring.Pre)
	if err != nil {
		return nil, err
	}

	// services
	p2p, err := do.InvokeNamed[transport.Transport](inj, ring.Transport)
	if err != nil {
		return nil, fmt.Errorf("invoke p2p factory: %w", err)
	}
	bb, err := do.InvokeNamed[bulletin.Bulletin](inj, ring.Bulletin)
	if err != nil {
		return nil, fmt.Errorf("invoke bulletin factory: %w", err)
	}

	dkgSrv, err := dkgFactory.New(inj, nil)
	if err != nil {
		return nil, fmt.Errorf("create dkg service: %w", err)
	}

	preSrv, err := preFactory.New(rid, ring.N, ring.T, p2p, bb, []types.Node{}, dkgSrv)
	if err != nil {
		return nil, fmt.Errorf("create pre service: %w", err)
	}

	pssSrv, err := pssFactory.New(rid, ring.N, ring.T, p2p, bb, []types.Node{}, dkgSrv)
	if err != nil {
		return nil, fmt.Errorf("create pss service: %w", err)
	}

	rs := &Ring{
		ID:        rid,
		DKG:       dkgSrv,
		PSS:       pssSrv,
		PRE:       preSrv,
		Transport: p2p,
		Bulletin:  bb,
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

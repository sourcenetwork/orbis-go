package orbis

import (
	"context"

	"github.com/samber/do"

	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/crypto/proof"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/pre"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"github.com/sourcenetwork/orbis-go/pkg/service"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

type ringService struct {
	ID types.RingID

	DKG dkg.DKG
	PSS pss.PSS
	PRE pre.PRE

	Transport transport.Transport
	Bulletin  bulletin.Bulletin

	repo db.Repository
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

// NewRing creates a new instance of a ring with all the depedant components
// and services wired together.
func (n *node) NewRing(ctx context.Context, manifest []byte, repo db.Repository) (service.RingService, error) {
	ring, rid, err := types.RingFromManifest(manifest)
	if err != nil {
		return nil, err
	}

	if err := repo.Rings.Create(ctx, ring); err != nil {
		return nil, err // todo wrap
	}

	// factories
	dkgFactory, err := do.InvokeNamed[dkg.Factory](n.injector, ring.Dkg)
	if err != nil {
		return nil, err
	}

	pssFactory, err := do.InvokeNamed[pss.Factory](n.injector, ring.Pss)
	if err != nil {
		return nil, err
	}

	preFactory, err := do.InvokeNamed[pre.Factory](n.injector, ring.Pre)
	if err != nil {
		return nil, err
	}

	// check err group

	// services
	p2p, err := do.InvokeNamed[transport.Transport](n.injector, ring.Transport)
	if err != nil {
		return nil, err
	}

	bb, err := do.InvokeNamed[bulletin.Bulletin](n.injector, ring.Bulletin)
	if err != nil {
		return nil, err
	}

	dkgSrv, err := dkgFactory.New(rid, ring.N, ring.T, p2p, bb, []types.Node{})
	if err != nil {
		return nil, err
	}

	preSrv, err := preFactory.New(rid, ring.N, ring.T, p2p, bb, []types.Node{}, dkgSrv)
	if err != nil {
		return nil, err
	}

	pssSrv, err := pssFactory.New(rid, ring.N, ring.T, p2p, bb, []types.Node{}, dkgSrv)
	if err != nil {
		return nil, err
	}

	rs := &ringService{
		ID:        rid,
		DKG:       dkgSrv,
		PSS:       pssSrv,
		PRE:       preSrv,
		Transport: p2p,
		Bulletin:  bb,
		repo:      repo,
	}

	// called in ring.Join() - go rs.handleEvents()

	return rs, nil
}

func (r *ringService) Store(context.Context, types.SecretID, *types.Secret, proof.VerifiableEncryption) error {
	return nil
}
func (r *ringService) Get(context.Context, types.SecretID) (types.Secret, error) {

	return types.Secret{}, nil
}
func (r *ringService) GetShares(context.Context, types.SecretID) ([]types.PrivSecretShare, error) {
	return nil, nil

}
func (r *ringService) Delete(context.Context, types.SecretID) error {
	return nil

}

func (r *ringService) PublicKey() (crypto.PublicKey, error) {
	return nil, nil

}
func (r *ringService) Refresh(context.Context, pss.Config) (pss.RefreshState, error) {
	return pss.RefreshState{}, nil

}
func (r *ringService) Threshold() int {
	return 0
}

func (r *ringService) State() pss.State {
	return pss.State{}
}

func (r *ringService) Nodes() []pss.Node {
	return nil
}

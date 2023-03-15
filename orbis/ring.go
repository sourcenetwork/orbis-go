package orbis

import (
	"context"

	"github.com/samber/do"

	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
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
	dkgFactory, err1 := do.InvokeNamed[dkg.Factory](n.injector, ring.Dkg)
	pssFactory, err2 := do.InvokeNamed[pss.Factory](n.injector, ring.Pss)
	preFactory, err2 := do.InvokeNamed[pre.Factory](n.injector, ring.Pre)

	// check err group

	// services
	p2p, err := do.InvokeNamed[transport.Transport](n.injector, ring.Transport)
	bb, err := do.InvokeNamed[bulletin.Bulletin](n.injector, ring.Bulletin)

	dkgSrv, err := dkgFactory.New(rid, ring.N, ring.T, p2p, bb, _)
	// if err
	preSrv, err := preFactory.New(rid, ring.N, ring.T, p2p, bb, _, dkgSrv)
	// if err
	pssSrv, err := pssFactory.New(rid, ring.N, ring.T, p2p, bb, _, dkgSrv)
	// if err

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

package app

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/samber/do"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"

	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/authn"
	"github.com/sourcenetwork/orbis-go/pkg/authz"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin/p2p"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/pre"
	"github.com/sourcenetwork/orbis-go/pkg/pre/elgamal"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	p2ptransport "github.com/sourcenetwork/orbis-go/pkg/transport/p2p"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

type Ring struct {
	ID       types.RingID
	manifest *ringv1alpha1.Ring

	DKG dkg.DKG
	PSS pss.PSS
	PRE pre.PRE

	Authz    authz.Authz
	Authn    authn.CredentialService
	Resolver authn.KeyResolver

	// collection of registered services
	// that require startup/shutdown and
	// expose hooks.
	services []service
	nodes    []types.Node

	Transport transport.Transport
	Bulletin  bulletin.Bulletin
	DB        *db.DB

	N int
	T int

	inj *do.Injector

	preReqMsg chan *transport.Message

	encScrts map[string][]byte            // preStoreMsgID
	encCmts  map[string][]byte            // preStoreMsgID
	xncCmts  map[string]chan kyber.Point  // preEncryptMsgID
	xncSki   map[string][]*share.PubShare // preEncryptMsgID
}
type State map[string]string

type service interface {
	Start(context.Context) error
	Close(context.Context) error
	State() string
	Name() string
}

func (app *App) GetRing(ctx context.Context, id string) (*Ring, error) {
	app.mu.Lock()
	defer app.mu.Unlock()

	r, ok := app.rings[types.RingID(id)]
	if !ok {
		return nil, fmt.Errorf("ring not found: %s", id)
	}

	return r, nil
}

func (app *App) ListRing(ctx context.Context) ([]*Ring, error) {
	app.mu.Lock()
	defer app.mu.Unlock()

	return app.listRing(ctx)
}

func (app *App) listRing(ctx context.Context) ([]*Ring, error) {
	var rings []*Ring
	for _, r := range app.rings {
		rings = append(rings, r)
	}

	return rings, nil
}

func (app *App) JoinRing(ctx context.Context, ring *ringv1alpha1.Ring) (*Ring, error) {
	app.mu.Lock()
	defer app.mu.Unlock()

	r, err := app.joinRing(ctx, ring, false /* fromState */)
	if err != nil {
		return nil, fmt.Errorf("join ring: %w", err)
	}

	err = app.ringRepo.Create(ctx, ring)
	if err != nil {
		return nil, fmt.Errorf("create ring: %w", err)
	}
	app.rings[r.ID] = r

	return r, nil
}

func (app *App) joinRing(ctx context.Context, ring *ringv1alpha1.Ring, fromState bool) (*Ring, error) {
	rid := types.RingID(ring.Id)

	if _, exists := app.rings[rid]; exists {
		return nil, fmt.Errorf("already joined ring %s", rid)
	}

	rs := &Ring{}

	// rings get their own cloned dependency injector handler,
	// since we actually create and initialize services which
	// are only scoped to rings, compared to the factories
	// which are global.
	inj := app.inj.Clone()

	// factories
	log.Info("pulling factory dependencies for ring")
	authnFactory, err := do.InvokeNamed[types.Factory[authn.CredentialService]](inj, ring.Authentication)
	if err != nil {
		return nil, fmt.Errorf("invoke authn credential service: %w", err)
	}

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
		return nil, fmt.Errorf("invoke pre factory: %w", err)
	}

	// get global services
	d, err := do.Invoke[*db.DB](inj)
	if err != nil {
		return nil, fmt.Errorf("invoke db: %w", err)
	}

	// register configured generic transport locally
	tp, err := do.InvokeNamed[transport.Transport](inj, ring.Transport)
	if err != nil {
		return nil, fmt.Errorf("invoke transport: %w", err)
	}
	do.ProvideValue(inj, tp)

	// register configured generic bulletin locally
	bb, err := do.InvokeNamed[bulletin.Bulletin](inj, ring.Bulletin)
	if err != nil {
		return nil, fmt.Errorf("invoke bulletin: %w", err)
	}
	do.ProvideValue(inj, bb)

	authz, err := do.InvokeNamed[authz.Authz](inj, ring.Authorization)
	if err != nil {
		return nil, fmt.Errorf("invoke authz: %w", err)
	}
	do.ProvideValue(inj, authz)

	// setup and register local services
	log.Info("Initializating local services for ring")
	authn, err := authnFactory.New(inj, []db.RepoKey{}) // empty repo keys
	if err != nil {
		return nil, fmt.Errorf("create authn: %w", err)
	}

	dkgRepoKeys := app.repoKeysForService(dkgFactory.Name())
	log.Debugf("dkg repo keys: %v", dkgRepoKeys)
	dkgSrv, err := dkgFactory.New(inj, dkgRepoKeys)
	if err != nil {
		return nil, fmt.Errorf("create dkg service: %w", err)
	}

	tpNodes, err := nodesFromIDs(ring.Nodes)
	if err != nil {
		return nil, fmt.Errorf("convert nodes from ring ids")
	}

	err = dkgSrv.Init(ctx, app.privateKey, rid, tpNodes, ring.N, ring.T, fromState)
	if err != nil {
		return nil, fmt.Errorf("initialize dkg: %w", err)
	}
	do.ProvideValue(inj, dkgSrv)

	err = rs.registerService(dkgSrv)
	if err != nil {
		return nil, fmt.Errorf("start dkg: %w", err)
	}

	preRepoKeys := app.repoKeysForService(preFactory.Name())
	preSrv, err := preFactory.New(inj, preRepoKeys)
	if err != nil {
		return nil, fmt.Errorf("create pre service: %w", err)
	}

	nodes := make([]types.Node, len(tpNodes))
	for i, n := range tpNodes {
		nodes[i] = *types.NewNode(i, n.ID(), n.Address(), n.PublicKey())
	}

	err = preSrv.Init(rid, ring.N, ring.T)
	if err != nil {
		return nil, fmt.Errorf("initialize pre: %w", err)
	}
	do.ProvideValue(inj, preSrv)

	pssRepoKeys := app.repoKeysForService(pssFactory.Name())
	pssSrv, err := pssFactory.New(inj, pssRepoKeys)
	if err != nil {
		return nil, fmt.Errorf("create pss service: %w", err)
	}

	err = pssSrv.Init(rid, ring.N, ring.T, nodes)
	if err != nil {
		return nil, fmt.Errorf("create pss service: %w", err)
	}

	rs = &Ring{
		ID:        rid,
		manifest:  ring,
		DKG:       dkgSrv,
		PSS:       pssSrv,
		PRE:       preSrv,
		Transport: tp,
		Bulletin:  bb,
		DB:        d,
		inj:       inj,
		N:         int(ring.N),
		T:         int(ring.T),

		nodes:     nodes,
		services:  rs.services, // this is dumb, but im being lazy, sorry.
		preReqMsg: make(chan *transport.Message, 10),
		encScrts:  make(map[string][]byte),
		encCmts:   make(map[string][]byte),
		xncCmts:   make(map[string]chan kyber.Point),
		xncSki:    make(map[string][]*share.PubShare),
		Authz:     authz,
		Authn:     authn,
	}

	go rs.preReencryptMessageHandler()

	tp.AddHandler(protocol.ID(elgamal.EncryptedSecretRequest), rs.preTransportMessageHandler)
	tp.AddHandler(protocol.ID(elgamal.EncryptedSecretReply), rs.preTransportMessageHandler)

	bbnamespace := fmt.Sprintf("/ring/%s/pre/store", string(rid))
	err = bb.Register(ctx, bbnamespace)
	if err != nil {
		return nil, fmt.Errorf("register bulletin: %w", err)
	}

	// TODO: this is a hack to wait for the bulletin to be registered
	time.Sleep(1 * time.Second)
	log.Infof("registered to topic %s with peers %v", bbnamespace, bb.(*p2p.Bulletin).Host().PubSub().ListPeers(bbnamespace))

	return rs, nil
}

func (a *App) ListRings() []*Ring {
	a.mu.Lock()
	defer a.mu.Unlock()

	var rings []*Ring
	for _, r := range a.rings {
		rings = append(rings, r)
	}

	return rings
}

func nodesFromIDs(nodes []*ringv1alpha1.Node) ([]transport.Node, error) {

	var tNodes []transport.Node
	for _, n := range nodes {

		id, err := peer.Decode(n.Id)
		if err != nil {
			return nil, fmt.Errorf("invalid peer id: %w", err)
		}

		pubKey, err := id.ExtractPublicKey()
		if err != nil {
			return nil, fmt.Errorf("extract publick key from id: %w", err)
		}

		key, err := crypto.PublicKeyFromLibP2P(pubKey)
		if err != nil {
			return nil, fmt.Errorf("invalid public key: %w", err)
		}

		addr, err := ma.NewMultiaddr(n.Address)
		if err != nil {
			return nil, fmt.Errorf("invalid address: %w", err)
		}

		node := p2ptransport.NewNode(n.Id, key, addr)
		tNodes = append(tNodes, node)
	}

	return tNodes, nil
}

func (r *Ring) Delete(context.Context, types.SecretID) error {
	// TODO: implement
	return nil
}

func (r *Ring) PublicKey() (crypto.PublicKey, error) {
	return r.DKG.PublicKey()
}

func (r *Ring) Refresh(context.Context, pss.Config) (pss.RefreshState, error) {
	return pss.RefreshState{}, nil
}

func (r *Ring) Nodes() []types.Node {
	return r.nodes
}
func (r *Ring) Threshold() int {
	return r.T
}

func (r *Ring) State() State {
	state := make(State)
	for _, s := range r.services {
		state[s.Name()] = s.State()
	}
	return state
}

func (r *Ring) Manifest() *ringv1alpha1.Ring {
	return r.manifest
}

func (r *Ring) Start(ctx context.Context) error {

	for _, srv := range r.services {
		err := srv.Start(ctx)
		if err != nil {
			return fmt.Errorf("start service %s: %w", srv.Name(), err)
		}
	}

	return nil
}

// LoadRings loads any existing rings into state
func (app *App) LoadRings(ctx context.Context) error {
	app.mu.Lock()
	defer app.mu.Unlock()
	log.Info("Loading rings from state")

	rings, err := app.ringRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	log.Infof("Found %d rings, rejoining", len(rings))
	for _, r := range rings {
		ring, err := app.joinRing(ctx, r, true /* fromState */)
		if err != nil {
			return fmt.Errorf("join ring: %w", err)
		}
		ring.manifest = r

		app.rings[ring.ID] = ring
	}
	log.Infof("Finished loading %d rings from state", len(rings))

	return nil
}

func (r *Ring) registerService(srv service) error {
	if srv == nil {
		return fmt.Errorf("service can't be nil")
	}
	r.services = append(r.services, srv)
	return nil
}

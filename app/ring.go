package app

import (
	"context"
	"errors"
	"fmt"

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
	manifest *ringv1alpha1.Manifest

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
	app *App

	preReqMsg chan *transport.Message

	xncCmts map[string]chan kyber.Point  // preEncryptMsgID
	xncSki  map[string][]*share.PubShare // preEncryptMsgID
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

func (app *App) JoinRing(ctx context.Context, manifest *ringv1alpha1.Manifest) (*Ring, error) {
	app.mu.Lock()
	defer app.mu.Unlock()

	r, err := app.joinRing(ctx, manifest, false /* fromState */)
	if err != nil {
		return nil, fmt.Errorf("join ring: %w", err)
	}

	ring := &ringv1alpha1.Ring{
		Id:       string(r.ID),
		Manifest: r.manifest,
	}
	err = app.ringRepo.Create(ctx, ring)
	if err != nil {
		return nil, fmt.Errorf("create ring: %w", err)
	}
	app.rings[r.ID] = r

	return r, nil
}

func (app *App) joinRing(ctx context.Context, manifest *ringv1alpha1.Manifest, fromState bool) (*Ring, error) {

	rid := ringID(manifest)
	log.Infof("Joining ring %s, nodes: %v", rid, manifest.Nodes)

	if _, exists := app.rings[rid]; exists {
		return nil, fmt.Errorf("already joined ring %s", rid)
	}

	rs := &Ring{app: app}

	// rings get their own cloned dependency injector handler,
	// since we actually create and initialize services which
	// are only scoped to rings, compared to the factories
	// which are global.
	inj := app.inj.Clone()

	// get global services
	d, err := do.Invoke[*db.DB](inj)
	if err != nil {
		return nil, fmt.Errorf("invoke db: %w", err)
	}

	// register configured generic transport locally
	tp, err := do.InvokeNamed[transport.Transport](inj, manifest.Transport)
	if err != nil {
		return nil, fmt.Errorf("invoke transport: %w", err)
	}
	do.ProvideValue(inj, tp)

	// register configured generic bulletin locally
	bb, err := do.InvokeNamed[bulletin.Bulletin](inj, manifest.Bulletin)
	if err != nil {
		return nil, fmt.Errorf("invoke bulletin: %w", err)
	}
	err = bb.Init(ctx)
	if err != nil {
		return nil, fmt.Errorf("init bulletin: %w", err)
	}
	do.ProvideValue(inj, bb)

	var (
		authnSrv authn.CredentialService
		authzSrv authz.Authz
	)

	authnSrv, err = do.InvokeNamed[authn.CredentialService](inj, manifest.Authentication)
	if err != nil && !errors.Is(err, do.ErrNotFound) { // ignore not found for now
		return nil, fmt.Errorf("invoke authn service: %w", err)
	} else if authnSrv != nil {
		do.ProvideValue(inj, authnSrv)
	}

	authzSrv, err = do.InvokeNamed[authz.Authz](inj, manifest.Authorization)
	if err != nil && !errors.Is(err, do.ErrNotFound) { // ignore not found for now
		return nil, fmt.Errorf("invoke authz service: %w", err)
	} else if authzSrv != nil {
		do.ProvideValue(inj, authzSrv)
	}

	// instanciating services from registered factories.

	// authn/authz can potentially use a global service instead of a factory
	// so we need to check if its already a configured service
	if authnSrv == nil {
		authnSrv, err = serviceFromFactory[authn.CredentialService](rs, inj, manifest.Authentication)
		if err != nil {
			return nil, fmt.Errorf("invoke authn credential service from factory: %w", err)
		}
	}
	if authzSrv == nil {
		authzSrv, err = serviceFromFactory[authz.Authz](rs, inj, manifest.Authorization)
		if err != nil {
			return nil, fmt.Errorf("invoke authz service from factory: %w", err)
		}
	}

	dkgSrv, err := serviceFromFactory[dkg.DKG](rs, inj, manifest.Dkg)
	if err != nil {
		return nil, fmt.Errorf("invoke dkg service from factory: %w", err)
	}

	pssSrv, err := serviceFromFactory[pss.PSS](rs, inj, manifest.Pss)
	if err != nil {
		return nil, fmt.Errorf("invoke pss service from factory: %w", err)
	}

	preSrv, err := serviceFromFactory[pre.PRE](rs, inj, manifest.Pre)
	if err != nil {
		return nil, fmt.Errorf("invoke pss service from factory: %w", err)
	}

	// setup and register local services
	log.Info("Initializating local services for ring")

	tpNodes, err := nodesFromIDs(manifest.Nodes)
	if err != nil {
		return nil, fmt.Errorf("convert nodes from ring ids")
	}

	err = dkgSrv.Init(ctx, app.privateKey, rid, tpNodes, manifest.N, manifest.T, fromState)
	if err != nil {
		return nil, fmt.Errorf("initialize dkg: %w", err)
	}

	// todo: Unify these nodes and the tpNodes for the DKG
	nodes := make([]types.Node, len(tpNodes))
	for i, n := range tpNodes {
		nodes[i] = *types.NewNode(i, n.ID(), n.Address(), n.PublicKey())
	}

	err = preSrv.Init(rid, manifest.N, manifest.T)
	if err != nil {
		return nil, fmt.Errorf("initialize pre: %w", err)
	}

	err = pssSrv.Init(rid, manifest.N, manifest.T, nodes)
	if err != nil {
		return nil, fmt.Errorf("create pss service: %w", err)
	}

	rs = &Ring{
		ID:        rid,
		manifest:  manifest,
		DKG:       dkgSrv,
		PSS:       pssSrv,
		PRE:       preSrv,
		Transport: tp,
		Bulletin:  bb,
		DB:        d,
		inj:       inj,
		N:         int(manifest.N),
		T:         int(manifest.T),

		nodes:     nodes,
		services:  rs.services, // this is dumb, but im being lazy, sorry.
		preReqMsg: make(chan *transport.Message, 10),
		xncCmts:   make(map[string]chan kyber.Point),
		xncSki:    make(map[string][]*share.PubShare),
		Authz:     authzSrv,
		Authn:     authnSrv,
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
	// time.Sleep(3 * time.Second)
	log.Infof("registered to namespace %s", bbnamespace)

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

func (r *Ring) Manifest() *ringv1alpha1.Manifest {
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
		ring, err := app.joinRing(ctx, r.Manifest, true /* fromState */)
		if err != nil {
			return fmt.Errorf("join ring: %w", err)
		}
		ring.manifest = r.Manifest

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

func serviceFromFactory[T any](ring *Ring, inj *do.Injector, name string) (T, error) {
	var zero T
	typeName := fmt.Sprintf("%T", zero)

	if name == "" {
		return zero, fmt.Errorf("missing name for service factory %s", typeName)
	}

	factory, err := do.InvokeNamed[types.Factory[T]](inj, name)
	if err != nil {
		return zero, fmt.Errorf("invoke %s(%s) factory: %w", typeName, name, err)
	}

	repoKeys := ring.app.repoKeysForService(factory.Name())
	srv, err := factory.New(inj, repoKeys, ring.app.config)
	if err != nil {
		return zero, fmt.Errorf("create %s(%s) service: %w", typeName, name, err)
	}

	do.ProvideValue(inj, srv)

	if s, ok := any(srv).(service); ok {
		err = ring.registerService(s)
		if err != nil {
			return zero, fmt.Errorf("register %s(%s) service: %w", typeName, name, err)
		}
	}

	return srv, nil
}

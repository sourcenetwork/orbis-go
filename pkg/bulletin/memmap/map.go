package memmap

import (
	"context"
	"sync"

	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
)

// type BaseBulletin = bulletin.Bulletin[string, []byte]

var _ bulletin.Bulletin = (*Bulletin)(nil)

// Bulletin is an in-memory testing bulletinboard
// implementation. It is *not* verifiable, doesn't use
// any BFT mechanics, nor connected to a network.
// For testing purposes only.
type Bulletin struct {
	mu       sync.RWMutex
	messages map[string][]byte
}

func (b *Bulletin) Name() string {
	return "memtest"
}

func (b *Bulletin) Register(ctx context.Context, namespace string) error {
	return nil // noop
}

// Post
func (b *Bulletin) Post(ctx context.Context, identifier bulletin.ID, msg bulletin.Message) (bulletin.Response, error) {
	idstr := identifier.String()
	return b.PostByString(ctx, idstr, msg)
}

func (b *Bulletin) PostByString(ctx context.Context, identifier string, msg bulletin.Message) (bulletin.Response, error) {
	if identifier == "" {
		return bulletin.Response{}, bulletin.ErrEmptyID
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	// check duplicate
	if _, exists := b.messages[identifier]; exists {
		return bulletin.Response{}, bulletin.ErrDuplicateMessage
	}
	b.messages[identifier] = msg
	return bulletin.Response{
		Data: msg,
	}, nil
}

// Read
func (b *Bulletin) Read(ctx context.Context, identifier bulletin.ID) (bulletin.Response, error) {
	idstr := identifier.String()
	return b.ReadByString(ctx, idstr)
}

func (b *Bulletin) ReadByString(ctx context.Context, identifier string) (bulletin.Response, error) {
	if identifier == "" {
		return bulletin.Response{}, bulletin.ErrEmptyID
	}

	b.mu.RLock()
	defer b.mu.RUnlock()
	msg, exists := b.messages[identifier]
	if !exists {
		return bulletin.Response{}, bulletin.ErrMessageNotFound
	}
	return bulletin.Response{
		Data: msg,
	}, nil
}

// Query
func (b *Bulletin) Query(ctx context.Context, query string) ([]bulletin.Response, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	resps := make([]bulletin.Response, 0)
	for id, msg := range b.messages {
		if glob(query, id) {
			resps = append(resps, bulletin.Response{
				Data: msg,
			})
		}
	}

	return resps, nil
}

func (b *Bulletin) Start() {}

func (b *Bulletin) Shutdown() {}

/*
p2ptp := p2p.NewTransport()
rabin := dkg.Service(rabin)

ring.New(rabin, avpss, cosmosBulletin, p2ptp)

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

func New(manifest, repo) (service.SecretRing, error) {
	ring, rid := types.RingFromManifest(manifest)
	repo.Ring.Create(ctx, rid, ring)

	// type safe factories for constructing named DKGs
	dkgFactory, err := do.InvokeNamed[dkg.Factory](manifest.dkg)
	pssFactory, err := do.InvokeNamed[pss.Factory](manifest.pss)
	preFactory, err := do.InvokeNamed[pre.Factory](manifest.pre)
	// services
	p2p, err := do.InvokeNamed[transport.Transport](manifest.transport)
	hub, err := do.InvokeNamed[bulletin.Bulletin](manifest.bulletin))
	dkgSrv, err := dkgFactory.New(rid, n, t, p2p, hub, nodes)
	preSrv, err := preFactory.New(rid, n, t, p2p, hub, nodes, dkgSrv)
	pssSrv, err := pssFactory.New(rid, n, t, p2p, hub, nodes, dkgSrv)

	rs := &RingService{
		ID: ringID,
		DKG: dkgSrv,
		PSS: pssSrv,
		PRE: preSrv,
		Transport: p2p,
		Bulletin: hub,
		Repo: repo,
	}

	go rs.handleEvents()

	return rs, nil
}

func (p *AVPSS) Reshare() error {
	share := ...
	p.repo.
}

type Events struct {
	Rounds eventbus.Channel[Rounds]
}

type AVPSS struct {
	events Events
}

func (p *AVPSS) RegisterTransport(t transport.Transport) error {
	p.transport = t
	t.AddHandler(avpss.ProtocolID(rid), p.transportHandler)
}

func (p *AVPSS) transportHandler(transport.Message)

func (r *RingService) handleEvents() {
	for {
		select {
		case <-ctx.Done():
			// close
		case pssmsg <-r.events.BulletinPSSMessageCh:
			// forward
			r.PSSService.ProcessMessage(r.ctx, pssmsg)
		case dkgmsg <-r.events.BulletinDKGMessageCh:
			//forward
			r.DKGService.ProcessMessage(r.ctx, dkgmsg)
		}
	}
}

*/

package memmap

import (
	"context"
	"sync"

	"github.com/golang/protobuf/proto"
	logging "github.com/ipfs/go-log"
	"github.com/sourcenetwork/eventbus-go"

	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/util/glob"
)

// type BaseBulletin = bulletin.Bulletin[string, []byte]

var log = logging.Logger("orbis/bulletin/map")

var _ bulletin.Bulletin = (*Bulletin)(nil)

type Option func(*Bulletin)

func WithBus(bus eventbus.Bus) Option {
	return func(b *Bulletin) {
		b.bus = bus
	}
}

// Bulletin is an in-memory testing bulletinboard
// implementation. It is *not* verifiable, doesn't use
// any BFT mechanics, nor connected to a network.
// For testing purposes only.
type Bulletin struct {
	mu       sync.RWMutex
	messages map[string][]byte

	bus eventbus.Bus
}

func New(opts ...Option) *Bulletin {
	b := &Bulletin{
		messages: make(map[string][]byte),
		bus:      eventbus.NewBus(),
	}

	for _, o := range opts {
		o(b)
	}

	return b
}

func (b *Bulletin) Name() string {
	return "memtest"
}

func (bb *Bulletin) Init(ctx context.Context) error {
	return nil
}

func (b *Bulletin) Register(ctx context.Context, namespace string) error {
	return nil // noop
}

// Post
func (b *Bulletin) Post(ctx context.Context, namespace, id string, msg *transport.Message) (bulletin.Response, error) {
	return b.PostByString(ctx, namespace+id, msg, true)
}

func (b *Bulletin) PostByString(ctx context.Context, identifier string, msg *transport.Message, emit bool) (bulletin.Response, error) {
	log.Debugf("handling post for namespace ID %s, emit=%v", identifier, emit)
	if identifier == "" {
		return bulletin.Response{}, bulletin.ErrEmptyID
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	// check duplicate
	if _, exists := b.messages[identifier]; exists {
		return bulletin.Response{}, bulletin.ErrDuplicateMessageF(identifier)
	}
	buf, err := proto.Marshal(msg)
	if err != nil {
		return bulletin.Response{}, err
	}
	b.messages[identifier] = buf

	if emit {
		log.Debugf("publishing post event locally for %s", identifier)
		evt := bulletin.Event{
			Message: msg,
			ID:      identifier,
		}
		err = eventbus.Publish(b.bus, evt) // publish the event locally
		if err != nil {
			log.Errorf("failed to publish event to channel: %w", err)
		}
	}

	return bulletin.Response{
		Data: msg,
		ID:   identifier,
	}, nil
}

// Read
func (b *Bulletin) Read(ctx context.Context, namespace, id string) (bulletin.Response, error) {
	return b.ReadByString(ctx, namespace+id)
}

func (b *Bulletin) ReadByString(ctx context.Context, identifier string) (bulletin.Response, error) {
	log.Debug("handling read for id %s", identifier)
	if identifier == "" {
		return bulletin.Response{}, bulletin.ErrEmptyID
	}

	b.mu.RLock()
	defer b.mu.RUnlock()
	msg, exists := b.messages[identifier]
	if !exists {
		return bulletin.Response{}, bulletin.ErrMessageNotFound
	}

	tMsg := new(transport.Message)
	err := proto.Unmarshal(msg, tMsg)
	if err != nil {
		return bulletin.Response{}, err
	}

	return bulletin.Response{
		Data: tMsg,
		ID:   identifier,
	}, nil
}

// Query
func (b *Bulletin) Query(ctx context.Context, namespace, query string) (<-chan bulletin.QueryResponse, error) {
	respCh := make(chan bulletin.QueryResponse, 0)

	query = namespace + query

	go func() {
		b.mu.RLock()
		defer func() {
			b.mu.RUnlock()
			close(respCh)
		}()

		for id, msg := range b.messages {
			if glob.Glob(query, id) {
				tMsg := new(transport.Message)
				err := proto.Unmarshal(msg, tMsg)
				if err != nil {
					respCh <- bulletin.QueryResponse{
						Err: err,
					}
					return // should we exit or continue loop?
				}

				respCh <- bulletin.QueryResponse{
					Resp: bulletin.Response{
						ID:   id,
						Data: tMsg,
					},
				}
			}
		}
	}()

	return respCh, nil
}

func (b *Bulletin) Has(ctx context.Context, namespace, id string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	identifier := namespace + id
	_, exists := b.messages[identifier]
	return exists
}

// Events
func (b *Bulletin) Events() eventbus.Bus {
	return b.bus
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

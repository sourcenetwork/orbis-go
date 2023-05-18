package app

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/samber/do"

	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	p2pbb "github.com/sourcenetwork/orbis-go/pkg/bulletin/p2p"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	"github.com/sourcenetwork/orbis-go/pkg/host"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	p2ptp "github.com/sourcenetwork/orbis-go/pkg/transport/p2p"
)

var (
	ErrDuplicateRepo    = fmt.Errorf("duplicate named repos")
	ErrFactoryEmptyName = fmt.Errorf("factory name can't be empty")
	ErrKeyMissing       = fmt.Errorf("key missing or nil")
)

// App implements App all services.
type App struct {
	host *host.Host
	tp   transport.Transport
	bb   bulletin.Bulletin
	db   *db.DB

	inj *do.Injector

	privateKey crypto.PrivateKey

	// namespaced key => repoParam
	// collected during app initialization
	repoParams map[string]repoParam

	// namespaced key => repo key
	// mounted repos after initialization
	repoKeys map[string]db.RepoKey

	// service name => []namespace keys
	// index for which keys are for which
	// service
	serviceRepos map[string][]string
}

type repoParam struct {
	key db.RepoKey
	typ db.Record
}

func (a *App) Host() *host.Host {
	return a.host
}
func (a *App) Transport() transport.Transport {
	return a.tp
}

func (a *App) Injector() *do.Injector {
	return a.inj
}

func New(ctx context.Context, host *host.Host, opts ...Option) (*App, error) {
	if host == nil {
		return nil, fmt.Errorf("host is nil")
	}

	// get the privkey for host
	hpk := host.Peerstore().PrivKey(host.ID())
	cpk, err := crypto.PrivateKeyFromLibP2P(hpk)
	if err != nil {
		return nil, fmt.Errorf("converting libp2p private key: %w", err)
	}

	inj := do.New()

	// register global services
	tp, err := p2ptp.New(ctx, host, config.Transport{})
	if err != nil {
		return nil, fmt.Errorf("creating transport: %w", err)
	}
	do.ProvideNamedValue[transport.Transport](inj, tp.Name(), tp)

	bb, err := p2pbb.New(ctx, host, config.Bulletin{})
	if err != nil {
		return nil, fmt.Errorf("creating bulletin: %w", err)
	}
	do.ProvideNamedValue[bulletin.Bulletin](inj, bb.Name(), bb)

	do.ProvideValue(inj, host)

	a := &App{
		host:         host,
		inj:          inj,
		tp:           tp,
		bb:           bb,
		privateKey:   cpk,
		repoParams:   make(map[string]repoParam),
		repoKeys:     make(map[string]db.RepoKey),
		serviceRepos: make(map[string][]string),
	}

	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, fmt.Errorf("apply orbis option: %w", err)
		}
	}

	return a, a.mountRepos()
}

func (a *App) mountRepos() error {
	for name, params := range a.repoParams {
		if params.key == nil {
			return ErrKeyMissing
		}
		if err := db.MountRepo(a.db, params.key, params.typ); err != nil {
			return err
		}
		a.repoKeys[name] = params.key
	}
	return nil
}

func (a *App) setupRepoKeysForService(namespace string, records []db.Record) error {
	repoKeys := keysForRepoTypes(records)
	serviceKeys := make([]string, len(records))
	for i, k := range repoKeys {
		name := namespaceKey(namespace, k)
		if _, exists := a.repoParams[name]; exists {
			return ErrDuplicateRepo
		}
		serviceKeys = append(serviceKeys, name)
		a.repoParams[name] = repoParam{
			key: k,
			typ: records[i],
		}
	}
	a.serviceRepos[namespace] = serviceKeys
	return nil
}
func (a *App) repoKeysForService(name string) []db.RepoKey {
	keys, exists := a.serviceRepos[name]
	if !exists {
		return nil
	}
	rkeys := make([]db.RepoKey, len(keys))
	for i, k := range keys {
		rkeys[i] = a.repoKeys[k]
	}
	return rkeys
}

func keysForRepoTypes(records []db.Record) []db.RepoKey {
	keys := make([]db.RepoKey, len(records))
	for i, r := range records {
		typeName := proto.MessageV2(r).ProtoReflect().Descriptor().Name()
		keys[i] = db.NewRepoKey(string(typeName))
	}
	return keys
}

func namespaceKey(namespace string, key db.RepoKey) string {
	return namespace + "/" + key.Name()
}

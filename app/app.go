package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/samber/do"

	logging "github.com/ipfs/go-log"
	"github.com/sourcenetwork/orbis-go/config"
	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

var log = logging.Logger("orbis/app")

var (
	ErrDuplicateRepo    = fmt.Errorf("duplicate named repos")
	ErrFactoryEmptyName = fmt.Errorf("factory name can't be empty")
	ErrKeyMissing       = fmt.Errorf("key missing or nil")
)

// App implements App all services.
type App struct {
	host transport.Transport
	db   *db.DB

	inj *do.Injector

	privateKey crypto.PrivateKey

	ringRepo db.Repository[*ringv1alpha1.Ring]

	rings map[types.RingID]*Ring

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

	config config.Config

	mu sync.Mutex
}

type repoParam struct {
	key db.RepoKey
	typ db.Record
}

func (a *App) Host() transport.Transport {
	return a.host
}

func (a *App) Injector() *do.Injector {
	return a.inj
}

func New(ctx context.Context, host transport.Transport, opts ...Option) (*App, error) {
	if host == nil {
		return nil, fmt.Errorf("host is nil")
	}

	cpk, err := host.PrivateKey()
	if err != nil {
		return nil, fmt.Errorf("convert libp2p private key: %w", err)
	}

	inj := do.New()

	a := &App{
		host:         host,
		inj:          inj,
		privateKey:   cpk,
		repoParams:   make(map[string]repoParam),
		repoKeys:     make(map[string]db.RepoKey),
		serviceRepos: make(map[string][]string),
		rings:        make(map[types.RingID]*Ring),
	}

	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, fmt.Errorf("apply orbis option: %w", err)
		}
	}

	a.ringRepo, err = db.GetRepo(a.db, db.NewRepoKey("ring"), ringPkFunc)
	if err != nil {
		return nil, fmt.Errorf("get ring repo: %w", err)
	}

	return a, nil
}

func (a *App) setupRepoKeysForService(namespace string, records []string) error {
	if len(records) == 0 {
		return nil
	}

	repoKeys := keysForRepoTypes(records)
	serviceKeys := make([]string, len(records))
	for i, k := range repoKeys {
		name := namespaceKey(namespace, k)
		if _, exists := a.repoParams[name]; exists {
			return ErrDuplicateRepo
		}
		serviceKeys[i] = name
		a.repoKeys[name] = k
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

func keysForRepoTypes(records []string) []db.RepoKey {
	keys := make([]db.RepoKey, len(records))
	for i, r := range records {
		keys[i] = db.NewRepoKey(r)
	}

	return keys
}

func namespaceKey(namespace string, key db.RepoKey) string {
	return namespace + "/" + key.Name()
}

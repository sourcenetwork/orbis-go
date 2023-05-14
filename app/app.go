package app

import (
	"context"
	"fmt"

	"github.com/samber/do"

	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	p2pbb "github.com/sourcenetwork/orbis-go/pkg/bulletin/p2p"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	"github.com/sourcenetwork/orbis-go/pkg/host"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	p2ptp "github.com/sourcenetwork/orbis-go/pkg/transport/p2p"
)

// App implements App all services.
type App struct {
	host *host.Host
	tp   transport.Transport
	bb   bulletin.Bulletin
	db   *db.DB

	inj *do.Injector

	repoKeys map[string]db.RepoKey
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

	db, err := db.New()
	if err != nil {
		return nil, fmt.Errorf("creating db: %w", err)
	}
	do.ProvideValue(inj, db)

	do.ProvideValue(inj, host)

	a := &App{
		host: host,
		inj:  inj,
		tp:   tp,
		bb:   bb,
		db:   db,
	}

	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, fmt.Errorf("apply orbis option: %w", err)
		}
	}

	return a, nil
}

package main

import (
	"context"
	"fmt"

	"github.com/TBD54566975/ssi-sdk/did/key"

	"github.com/sourcenetwork/orbis-go/app"
	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/pkg/authn"
	"github.com/sourcenetwork/orbis-go/pkg/authn/jws"
	"github.com/sourcenetwork/orbis-go/pkg/authz"
	"github.com/sourcenetwork/orbis-go/pkg/authz/acp"
	"github.com/sourcenetwork/orbis-go/pkg/authz/zanzi"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	p2pbb "github.com/sourcenetwork/orbis-go/pkg/bulletin/p2p"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin/sourcehub"
	"github.com/sourcenetwork/orbis-go/pkg/cosmos"
	"github.com/sourcenetwork/orbis-go/pkg/did"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/dkg/rabin"
	"github.com/sourcenetwork/orbis-go/pkg/host"
	"github.com/sourcenetwork/orbis-go/pkg/pre"
	"github.com/sourcenetwork/orbis-go/pkg/pre/elgamal"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
	"github.com/sourcenetwork/orbis-go/pkg/pss/avpss"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	p2ptp "github.com/sourcenetwork/orbis-go/pkg/transport/p2p"
)

func setupApp(ctx context.Context, cfg config.Config) (*app.App, error) {

	host, err := host.New(ctx, cfg.Host)
	if err != nil {
		return nil, fmt.Errorf("create host: %w", err)
	}

	var cosmosclient *cosmos.Client
	if cfg.Cosmos != (config.Cosmos{}) {
		cosmosclient, err = cosmos.New(ctx, cfg.Cosmos)
		if err != nil {
			return nil, fmt.Errorf("create cosmos client: %w", err)
		}
	}

	tp, err := p2ptp.New(ctx, host, cfg.Transport)
	if err != nil {
		return nil, fmt.Errorf("create transport: %w", err)
	}

	bb, err := p2pbb.New(ctx, host, cfg.Bulletin)
	if err != nil {
		return nil, fmt.Errorf("create p2p bulletin: %w", err)
	}

	hubbb, err := sourcehub.New(ctx, host, cosmosclient, cfg.Bulletin)
	if err != nil {
		return nil, fmt.Errorf("create sourcehub bulletin: %w", err)
	}

	// Services & Factory Options
	//
	// Services are global instances that are shared between all
	// consumers/callers (rings). They are instanciated once (like above).
	//
	// Factories are singletons that produce newly instanciated
	// objects for each new consumer/caller (rings)
	//
	// Options are called in order, `app.DefaultOptions` *should* be called first.
	opts := []app.Option{
		app.DefaultOptions(cfg),

		// shared global transport and bulletin.
		app.WithService[transport.Transport](tp),

		app.WithService[bulletin.Bulletin](bb),
		app.WithService[bulletin.Bulletin](hubbb),

		app.WithService[*cosmos.Client](cosmosclient),

		// Authentication and Authorization services
		app.WithService(authz.NewAllow(authz.ALLOW_ALL)),
		app.WithService(did.NewResolver(key.Resolver{})),
		app.WithFactory[authn.CredentialService](jws.SelfSignedFactory),
		app.WithFactory[authz.Authz](zanzi.Factory),
		app.WithFactory[authz.Authz](acp.Factory),

		// DKG, PRE, and PSS Factories
		app.WithFactory[dkg.DKG](rabin.Factory),
		app.WithFactory[pre.PRE](elgamal.Factory),
		app.WithFactory[pss.PSS](avpss.Factory),

		// TODO: Enable support the AVPSS, ECPSS, and CHURP based PSS systems.
		// Also enable basic VSS for testing (no networking/bulleting required).
		// app.WithProactiveSecretSharing(vss.Provider),

		// mount DB Tables
		app.WithDBData(cfg.DB.Path),
	}

	app, err := app.New(ctx, host, opts...)
	if err != nil {
		return nil, fmt.Errorf("create app: %w", err)
	}

	return app, nil
}

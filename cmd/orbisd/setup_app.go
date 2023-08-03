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
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	p2pbb "github.com/sourcenetwork/orbis-go/pkg/bulletin/p2p"
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

	tp, err := p2ptp.New(ctx, host, cfg.Transport)
	if err != nil {
		return nil, fmt.Errorf("create transport: %w", err)
	}

	bb, err := p2pbb.New(ctx, host, cfg.Bulletin)
	if err != nil {
		return nil, fmt.Errorf("create bulletin: %w", err)
	}

	// Services & Factory Options
	//
	// Services are global instances that are shared between all
	// consumers/callers (rings). They are instanciated once (like above).
	//
	// Factories are singletons that produce newly instanciated
	// objects for each new consumer/caller (rings)
	opts := []app.Option{
		app.DefaultOptions(),

		app.WithService[transport.Transport](tp),
		app.WithService[bulletin.Bulletin](bb),

		app.WithService(authz.NewAllow(authz.ALLOW_ALL)),
		app.WithService(did.NewResolver(key.Resolver{})),

		app.WithFactory[authn.CredentialService](jws.SelfSignedFactory),
		app.WithFactory[dkg.DKG](rabin.Factory),
		app.WithFactory[pre.PRE](elgamal.Factory),

		// Enable support the AVPSS, ECPSS, and CHURP based PSS systems.
		app.WithFactory[pss.PSS](avpss.Factory),

		// Also enable basic VSS for testing (no networking/bulleting required).
		// app.WithProactiveSecretSharing(vss.Provider),

		// mount DB Tables
		app.WithDBData(cfg.DB.Path),
	}

	app, err := app.New(ctx, host, opts...)
	if err != nil {
		return nil, fmt.Errorf("create app: %w", err)
	}

	// load state
	log.Info("Loading rings from state")
	err = app.LoadRings(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// 	rid := types.RingID("40b086ef") // cid/multihash of the ring config
// 	err = node.JoinRing(rid, config.Ring, ringPeers)

// 	ring := node.GetRing(rid)

// 	select {
// 	case <-time.NewTimer(time.Minute):
// 		// timeout
// 	case <-ring.WaitForState(ptypes.STATE_INITIALIZED):
// 		// ready
// 	}

// 	sid = types.SecretID("mySecretIdentifier")
// 	secretVal := []byte("My Super Secret Value or Private Key or Symmetric Key or w.e")
// 	ring.Store(ctx, sid, secretVal)

// 	secret, err := ring.Get(ctx, auth.NilToken, sid)
// 	share, err := ring.GetLocalShare(ctx, auth.NilToken, sid)
// 	shares, err := ring.GetShares(ctx, auth.NilToken, sid)

// 	recovered := orbis.RecoverFromShares(shares)

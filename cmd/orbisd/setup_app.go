package main

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/orbis-go/app"
	"github.com/sourcenetwork/orbis-go/config"
	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
	p2pbb "github.com/sourcenetwork/orbis-go/pkg/bulletin/p2p"
	"github.com/sourcenetwork/orbis-go/pkg/dkg/rabin"
	"github.com/sourcenetwork/orbis-go/pkg/host"
	"github.com/sourcenetwork/orbis-go/pkg/pre/elgamal"
	"github.com/sourcenetwork/orbis-go/pkg/pss/avpss"
	"github.com/sourcenetwork/orbis-go/pkg/ring"
	p2ptp "github.com/sourcenetwork/orbis-go/pkg/transport/p2p"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

func setupApp(ctx context.Context, cfg config.Config) (*app.App, error) {

	// testingStreamHandler := func(msg *transport.Message) error {
	// 	log.Infof("Transport Received message: %s", msg.Payload)
	// 	return nil
	// }
	// tp.AddHandler(p2ptp.ProtocolID, testingStreamHandler)

	opts := []app.Option{
		app.DefaultOptions(),
		app.WithHost(host.Factory{}),
		app.WithTransport(p2ptp.Factory),
		app.WithBulletin(p2pbb.Factory),

		app.WithDistKeyGenerator(rabin.Factory),
		app.WithProxyReencryption(elgamal.Factory),

		// Enable support the AVPSS, ECPSS, and CHURP based PSS systems.
		app.WithProactiveSecretSharing(avpss.Factory),

		// Also enable basic VSS for testing (no networking/bulleting required).
		// app.WithProactiveSecretSharing(vss.Provider),

		// mount DB Tables
		// app.WithTable()
	}

	app, err := app.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("create app: %w", err)
	}

	manifest := &types.Ring{
		Ring: ringv1alpha1.Ring{
			Id:        "40b086ef",
			N:         3,
			T:         2,
			Dkg:       "",
			Pss:       "",
			Pre:       "",
			Bulletin:  "p2pbb",
			Transport: "p2ptp",
			Nodes:     nil,
		},
	}

	// repo, err := db.New()
	// if err != nil {
	// 	return nil, fmt.Errorf("create ring repo: %w", err)
	// }

	rr, err := ring.NewRing(ctx, app.Injector(), manifest)
	if err != nil {
		return nil, fmt.Errorf("create ring: %w", err)
	}
	_ = rr

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

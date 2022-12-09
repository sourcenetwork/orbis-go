package main

import (
	"context"
	"time"

	"github.com/sourcenetwork/orbis-go/orbis"
	"github.com/sourcenetwork/orbis-go/pkg/types"
	"github.com/sourcenetwork/orbis-go/pkg/transport/p2p"
)

func main() {

	ctx := context.Background()

	host := libp2p2.New(...)
	transport := p2p.NewTransport(host)

	opts := []orbis.Option{
		orbis.DefaultOptions(),

		// Enable Tendermint based BulletinBoard Service
		orbis.WithBulletinService(tmbb.Provider),

		// Enable support the AVPSS, ECPSS, and CHURP based
		// PSS systems.
		orbis.WithSharingService(avpss.Provider),
		orbis.WithSharingService(ecpss.Provider),
		orbis.WithSharingService(churp.Provider),

		// Also enable basic VSS (no networking/bulleting required)
		// Primarily used for testing
		orbis.WithSharingService(vss.Provider),
	}

	node, err := orbis.NewNode(ctx, transport, opts)
	if err != nil {
	}

	err = node.Start()
	if err != nil {
	}

	rid := types.RingID("40b086ef") // cid/multihash of the ring config
	err = node.JoinRing(rid, config.Ring, ringPeers)
	if err != nil {
	}

	ring := node.GetRing(rid)

	select {
	case <-time.NewTimer(time.Minute):
		// timeout
	case <-ring.WaitForState(ptypes.STATE_INITIALIZED):
		// ready
	// } <== UNCOMMENT TO COMPILE - HACK TO HIDE IDE LINTER WARNINGS WHILE DESIGNING/PLAYING

	assert(ring.State() == ptypes.STATE_INITIALIZED)
	assert(ring.Type() == "AVPSS")

	sid = types.SecretID("mySecretIdentifier")
	secretVal := []byte("My Super Secret Value or Private Key or Symmetric Key or w.e")
	ring.Store(ctx, sid, secretVal)

	secret, err := ring.Get(ctx, auth.NilToken, sid)
	share, err := ring.GetLocalShare(ctx, auth.NilToken, sid)
	shares, err := ring.GetShares(ctx, auth.NilToken, sid)
	if err != nil {
	}

	recovered := orbis.RecoverFromShares(shares)
	assert(secret == recovered)
}

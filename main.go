package main

import (
	"time"

	"github.com/sourcenetwork/orbis-go/orbis"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)


func main() {

	opts := []orbis.Option{
		orbis.DefaultOptions(),
		orbis.WithP2PHost(libp2pHost),

		// Enable Tendermint based BulletinBoard Service
		orbis.WithBulletinService(tmbb.Provider),

		// Enable support the AVPSS, ECPSS, and CHURP based
		// PSS systems.
		orbis.WithSharingService(avpss.Provider),
		orbis.WithSharingService(ecpss.Provider),
		orbis.WithSharingService(churp.Provider),
	},

	node := orbis.NewNode(opts)
	err := on.Init()
	if err != nil {}

	rid := types.RingID("40b086ef") // cid/multihash of the ring config
	err = node.JoinRing(rid, config.Ring, ringPeers)
	if err != nil {}

	ring := node.GetRing(rid)
	
	select {
	case <- time.NewTimer(time.Minute):
		// timeout
	case <- ring.WaitForState(orbis.INITIALIZED):
		// ready

	assert(ring.State() == orbis.INITIALIZED)
	assert(ring.Type() == "AVPSS")

	sid = types.SecretID("mySecretIdentifier")
	secretVal := []byte("My Super Secret Value or Private Key or Symmetric Key or w.e")
	ring.Store(ctx, sid, secretVal)

	secret, err := ring.Get(ctx, auth.NilToken, sid)
	share, err := ring.GetLocalShare(ctx, auth.NilToken, sid)
	shares, err := ring.GetShares(ctx, auth.NilToken, sid)
	if err != nil {}

	recovered := orbis.RecoverFromShares(shares)
	assert(secret == recovered)

	

}
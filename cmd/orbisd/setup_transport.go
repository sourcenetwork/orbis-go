package main

import (
	"context"
	"time"

	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/infra/logger"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/transport/p2p"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	libp2phost "github.com/libp2p/go-libp2p/core/host"
	libp2prouting "github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
)

func setupTransport(ctx context.Context, lg logger.Logger, cfg config.Transport) (transport.Transport, error) {

	priv, _, err := libp2pcrypto.GenerateKeyPair(
		libp2pcrypto.Ed25519, // Select your key type. Ed25519 are nice short
		-1,                   // Select key length when possible (i.e. RSA).
	)
	if err != nil {
		return nil, err
	}

	var idht *dht.IpfsDHT

	connmgr, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		return nil, err
	}

	host, err := libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(priv),
		// Multiple listen addresses
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/9000",      // regular tcp connections
			"/ip4/0.0.0.0/udp/9000/quic", // a UDP endpoint for the QUIC transport
		),
		// support TLS connections
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		// support noise connections
		libp2p.Security(noise.ID, noise.New),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connmgr),
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		// Let this host use the DHT to find other hosts
		libp2p.Routing(func(h libp2phost.Host) (libp2prouting.PeerRouting, error) {
			idht, err = dht.New(ctx, h)
			return idht, err
		}),
		libp2p.EnableNATService(),
	)
	if err != nil {
		return nil, err
	}

	tp := p2p.NewTransport(host)

	return tp, nil
}

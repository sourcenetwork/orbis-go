package transport

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"math"
	mrand "math/rand"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	logging "github.com/ipfs/go-log"
	libp2p "github.com/libp2p/go-libp2p"
	libp2pdht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/event"
	libp2pevent "github.com/libp2p/go-libp2p/core/event"
	libp2phost "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	libp2pconnmgr "github.com/libp2p/go-libp2p/p2p/net/connmgr"
	libp2pnoise "github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	ma "github.com/multiformats/go-multiaddr"
	"google.golang.org/protobuf/proto"

	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

const (
	// wait a random amount of time from this interval
	// before dialing peers or reconnecting to help prevent DoS
	dialRandomizerIntervalMilliseconds = 3000

	// repeatedly try to reconnect for a few minutes
	// ie. 5 * 20 = 100s
	reconnectAttempts = 20
	reconnectInterval = 5 * time.Second

	// then move into exponential backoff mode for ~1day
	// ie. 3**10 = 16hrs
	reconnectBackOffAttempts    = 10
	reconnectBackOffBaseSeconds = 3
)

var log = logging.Logger("orbis/transport")

type Host struct {
	p2pHost libp2phost.Host
	privKey crypto.PrivateKey
	idht    *libp2pdht.IpfsDHT
	pubsub  *pubsub.PubSub

	reonnecting     sync.Map
	persistentPeers map[peer.ID]peer.AddrInfo
}

var _ Transport = (*Host)(nil)

func NewHost(ctx context.Context, cfg config.Host) (*Host, error) {

	// Convert string to libp2p crypto type.
	// Invalid types and/or bits are handled by libp2p.
	cryptoType := libp2pcrypto.RSA
	switch cfg.Crypto.Type {
	case "ed25519":
		cryptoType = libp2pcrypto.Ed25519
	case "secp256k1":
		cryptoType = libp2pcrypto.Secp256k1
	case "ecdsa":
		cryptoType = libp2pcrypto.ECDSA
	}

	randomness := rand.Reader
	if seed := cfg.Crypto.Seed; seed != 0 {
		randomness = mrand.New(mrand.NewSource(int64(seed)))
	}

	priv, _, err := libp2pcrypto.GenerateKeyPairWithReader(cryptoType, cfg.Crypto.Bits, randomness)
	if err != nil {
		return nil, fmt.Errorf("generate key pair: %w", err)
	}

	cpriv, err := crypto.PrivateKeyFromLibP2P(priv)
	if err != nil {
		return nil, fmt.Errorf("converting to crypto private key: %w", err)
	}

	connmgr, err := libp2pconnmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		libp2pconnmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("create connection manager: %w", err)
	}

	var dhtOptions []libp2pdht.Option
	if len(cfg.PersistentPeers) == 0 {
		log.Infof("Host running as a bootsrap node")
		dhtOptions = append(dhtOptions, libp2pdht.Mode(libp2pdht.ModeServer))
	}

	p2pHost, err := libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(priv),
		// Multiple listen addresses
		libp2p.ListenAddrStrings(cfg.ListenAddresses...),
		// support TLS connections
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		// support noise connections
		libp2p.Security(libp2pnoise.ID, libp2pnoise.New),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connmgr),
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		libp2p.EnableNATService(),
	)
	if err != nil {
		return nil, fmt.Errorf("create libp2p host: %w", err)
	}

	pubsubTracer := new(pubsubTracer)
	gossipSub, err := pubsub.NewGossipSub(ctx, p2pHost, pubsub.WithEventTracer(pubsubTracer))
	if err != nil {
		return nil, fmt.Errorf("create gossipsub: %w", err)
	}

	h := &Host{
		p2pHost:         p2pHost,
		pubsub:          gossipSub,
		privKey:         cpriv,
		persistentPeers: make(map[peer.ID]peer.AddrInfo),
	}

	for _, paddr := range cfg.PersistentPeers {
		pi, err := peer.AddrInfoFromString(paddr)
		if err != nil {
			return nil, fmt.Errorf("parse persistent peer addresses: %w", err)
		}
		h.persistentPeers[pi.ID] = *pi
	}

	go h.maintainPeers(ctx)

	return h, nil
}

func (h *Host) Name() string {
	return "p2p"
}

func (h *Host) PubSub() *pubsub.PubSub {
	return h.pubsub
}

func (h *Host) PublicKey() crypto.PublicKey {
	return h.privKey.GetPublic()
}

func (h *Host) Network() network.Network {
	return h.p2pHost.Network()
}

func (h *Host) NewMessage(rid types.RingID, id string, gossip bool, payload []byte, msgType string, target Node) (*Message, error) {

	pubkeyBytes, err := h.PublicKey().Raw()
	if err != nil {
		return nil, fmt.Errorf("get raw public key: %w", err)
	}

	// todo: Signature (should be done on send)
	// replay? nonce?
	msg := &Message{
		Timestamp:  time.Now().Unix(),
		Id:         id,
		RingId:     string(rid),
		NodeId:     h.ID().String(),
		NodePubKey: pubkeyBytes,
		Type:       msgType,
		Payload:    payload,
		Gossip:     gossip,
	}

	if target != nil {
		msg.TargetId = target.ID().String()
		pubkeyBuf, err := target.PublicKey().Raw()
		if err != nil {
			return nil, err
		}
		msg.TargetPubKey = pubkeyBuf
	}

	return msg, nil
}

func (h *Host) ID() peer.ID {
	return h.p2pHost.ID()
}

func (h *Host) Address() ma.Multiaddr {
	return h.p2pHost.Addrs()[0]
}

func (h *Host) EventBus() libp2pevent.Bus {
	return h.p2pHost.EventBus()
}

func (h *Host) PrivateKey() (crypto.PrivateKey, error) {
	sk := h.p2pHost.Peerstore().PrivKey(h.ID())
	k, err := crypto.PrivateKeyFromLibP2P(sk)
	if err != nil {
		return nil, fmt.Errorf("convert libp2p private key: %w", err)
	}
	return k, nil
}

func (h *Host) Send(ctx context.Context, node Node, msg *Message) error {

	// todo: telemetry
	// todo: verify msg is of type p2p.message
	// todo sign message

	// todo protocol formatting
	protocolID := protocol.ConvertFromStrings([]string{msg.GetType()})

	log.Infof("Send(): peerID:%s, ProtocolID:%v", node.ID(), protocolID)
	var stream network.Stream
	var err error
	newStream := func() error {
		stream, err = h.p2pHost.NewStream(ctx, node.ID(), protocolID...)
		return err
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 10 * time.Second
	bctx := backoff.WithContext(b, ctx)

	err = backoff.Retry(newStream, bctx)
	if err != nil {
		return fmt.Errorf("new stream: %v", err)
	}
	defer stream.Close()

	buf, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	_, err = stream.Write(buf)
	if err != nil {
		return fmt.Errorf("write stream: %w", err)
	}

	return nil
}

func (h *Host) AddHandler(pid protocol.ID, handler Handler) {
	streamHandler := streamHandlerFrom(handler)
	h.p2pHost.SetStreamHandler(pid, streamHandler)
}

func (h *Host) RemoveHandler(pid protocol.ID) {
	h.p2pHost.RemoveStreamHandler(pid)
}

func streamHandlerFrom(handler Handler) func(network.Stream) {
	return func(stream network.Stream) {

		log.Infof("new stream from %s", stream.Conn().RemotePeer())

		buf, err := io.ReadAll(stream)
		if err != nil {
			if err != io.EOF {
				log.Errorf("read stream: %s", err)
			}

			err = stream.Reset()
			if err != nil {
				log.Errorf("reset stream: %s", err)
			}

			return
		}

		err = stream.Close()
		if err != nil {
			log.Errorf("close stream: %s", err)
			return
		}

		data := &Message{}
		err = proto.Unmarshal(buf, data)
		if err != nil {
			log.Errorf("unmarshal data: %s", err)
			return
		}

		log.Infof("received message: id:%s, type: %s", data.Id, data.Type)
		err = handler(data)
		if err != nil {
			log.Errorf("handle data: %s", err)
			return
		}
	}
}

func (h *Host) maintainPeers(ctx context.Context) {
	go func() {
		for _, p := range h.persistentPeers {
			go h.reconnectToPeer(ctx, p)
		}
	}()

	subCh, err := h.EventBus().Subscribe(new(event.EvtPeerConnectednessChanged))
	if err != nil {
		log.Fatalf("Error subscribing to peer connectedness changes: %s", err)
	}
	defer subCh.Close()

	for {
		select {
		case ev, ok := <-subCh.Out():
			if !ok {
				return
			}

			evt := ev.(event.EvtPeerConnectednessChanged)
			if evt.Connectedness != network.NotConnected {
				continue
			}

			if _, ok := h.persistentPeers[evt.Peer]; !ok {
				continue
			}

			paddr := h.persistentPeers[evt.Peer]
			go h.reconnectToPeer(ctx, paddr)

		case <-ctx.Done():
			return
		}
	}
}

func (h *Host) reconnectToPeer(ctx context.Context, paddr peer.AddrInfo) {
	if _, ok := h.reonnecting.Load(paddr.ID.String()); ok {
		log.Infof("duplicate peer maintainence goroutine: %s", paddr.ID)
		return
	}

	h.reonnecting.Store(paddr.ID.String(), struct{}{})
	defer h.reonnecting.Delete(paddr.ID.String())

	start := time.Now()
	log.Infof("Reconnecting to peer %s", paddr)
	for i := 0; i < reconnectAttempts; i++ {
		select {
		case <-ctx.Done():
			log.Debug("peer maintainence goroutine context finished", paddr.ID)
			return
		default:
			// noop fallthrough
		}

		err := h.p2pHost.Connect(ctx, paddr)
		if err == nil {
			log.Infof("reconnected to peer %s during regular backoff", paddr.ID)
			return //success
		}

		log.Infof("Error reconnecting to peer %s: %s, Retrying %d/%d attemps", paddr, err, i, reconnectAttempts)
		randomSleep(reconnectInterval)
	}

	log.Errorf("Failed to reconnect to peer %s. Beginning exponential backoff. Elapsed %s", paddr, time.Since(start))
	for i := 0; i < reconnectBackOffAttempts; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			// noop fallthrough
		}

		// sleep an exponentially increasing amount
		sleepIntervalSeconds := math.Pow(reconnectBackOffBaseSeconds, float64(i))
		randomSleep(time.Duration(sleepIntervalSeconds) * time.Second)

		err := h.p2pHost.Connect(ctx, paddr)
		if err == nil {
			log.Infof("reconnected to peer %s during exponential backoff", paddr.ID)
			return //success
		}

		log.Infof("Error reconnecting to peer %s: %s, Retrying %d/%d attemps", paddr, err, i, reconnectAttempts)
	}
	log.Errorf("Failed to reconnect to peer %s. Giving up. Elapsed %s", paddr, time.Since(start))
}

func randomSleep(interval time.Duration) {
	r := time.Duration(mrand.Int63n(dialRandomizerIntervalMilliseconds)) * time.Millisecond
	time.Sleep(r + interval)
}

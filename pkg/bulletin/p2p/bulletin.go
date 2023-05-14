package p2p

import (
	"context"

	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/host"

	logging "github.com/ipfs/go-log"
	libp2phost "github.com/libp2p/go-libp2p/core/host"
	libp2pnetwork "github.com/libp2p/go-libp2p/core/network"
)

var log = logging.Logger("orbis/bulletin/p2p")

const (
	ProtocolID = "/orbis-bulletion/1.0.0"
	name       = "p2pbb"
)

type Bulletin struct {
	h libp2phost.Host
}

func New(ctx context.Context, host *host.Host, cfg config.Bulletin) (*Bulletin, error) {
	bb := &Bulletin{
		h: host,
	}

	host.SetStreamHandler(ProtocolID, bb.HandleStream)
	host.Discover(ctx, cfg.Rendezvous)

	return bb, nil
}

func (bb *Bulletin) Name() string {
	return name
}

func (bb *Bulletin) Post(ctx context.Context, path string, msg bulletin.Message) (bulletin.Response, error) {
	panic("implement me")
}

func (bb *Bulletin) Read(ctx context.Context, path string) (bulletin.Response, error) {
	panic("implement me")
}

func (bb *Bulletin) Query(ctx context.Context, query string) ([]bulletin.Response, error) {
	panic("implement me")
}

func (bb *Bulletin) Verify(context.Context, bulletin.Proof, string, bulletin.Message) bool {
	return true
}

// EventBus
// Events() eventbus.Bus

func (bb *Bulletin) HandleStream(stream libp2pnetwork.Stream) {
	log.Infof("Received stream: %s", stream.Conn().RemotePeer().Pretty())
}

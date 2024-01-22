package p2p

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/eventbus-go"
	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

func init() {
	logging.SetLogLevelRegex("orbis.*", "debug")
}

func newMessage(bb *Bulletin, rid string, typ string, buf []byte) (*transport.Message, error) {
	cid, err := types.CidFromBytes(buf)
	if err != nil {
		return nil, fmt.Errorf("cid from bytes: %w", err)
	}

	msg, err := bb.h.NewMessage(types.RingID(rid), cid.String(), false, buf, typ, nil)
	return msg, err
}

func newDefaultP2PHost(t *testing.T, ctx context.Context) transport.Transport {
	defaultHost, err := config.Default[config.Host]()
	require.NoError(t, err)

	defaultHost.Crypto.Seed = 1
	h, err := transport.NewHost(ctx, defaultHost)
	require.NoError(t, err)
	return h
}

func newRandomP2PHost(t *testing.T, ctx context.Context) transport.Transport {
	defaultHost, err := config.Default[config.Host]()
	require.NoError(t, err)
	// 0 port will result in random
	defaultHost.ListenAddresses = []string{"/ip4/0.0.0.0/tcp/0"}

	h, err := transport.NewHost(ctx, defaultHost)
	require.NoError(t, err)
	return h
}

func TestNewBulletin(t *testing.T) {
	ctx := context.Background()
	h := newDefaultP2PHost(t, ctx)
	defaultBulletinCfg, err := config.Default[config.Bulletin]()
	require.NoError(t, err)

	bb, err := New(ctx, h, defaultBulletinCfg)
	require.NoError(t, err)
	require.NotNil(t, bb)
}

func TestMultipleBulletinNetworkConnections(t *testing.T) {
	ctx := context.Background()
	h0 := newDefaultP2PHost(t, ctx)
	h1 := newRandomP2PHost(t, ctx)
	h2 := newRandomP2PHost(t, ctx)

	// bulletin setup
	cfg0, err := config.Default[config.Bulletin]()
	require.NoError(t, err)

	bb0, err := New(ctx, h0, cfg0)
	require.NoError(t, err)
	require.NotNil(t, bb0)

	cfg1, err := config.Default[config.Bulletin]()
	require.NoError(t, err)
	addr := h0.Address()
	cfg1.P2P.PersistentPeers = fmt.Sprintf("%s/p2p/%s", addr, h0.ID())

	bb1, err := New(ctx, h1, cfg1)
	require.NoError(t, err)
	require.NotNil(t, bb1)

	cfg2, err := config.Default[config.Bulletin]()
	require.NoError(t, err)
	cfg2.P2P.PersistentPeers = cfg1.P2P.PersistentPeers

	bb2, err := New(ctx, h2, cfg2)
	require.NoError(t, err)
	require.NotNil(t, bb2)

	// wait for the net connections
	time.Sleep(1 * time.Second)

	// test peer connections
	require.Len(t, h0.Network().Conns(), 2)
	require.Len(t, h1.Network().Conns(), 1)
	require.Len(t, h2.Network().Conns(), 1)

	ringTopic := "/ring/123"
	bb0.Register(ctx, ringTopic)
	bb1.Register(ctx, ringTopic)
	bb2.Register(ctx, ringTopic)

	time.Sleep(2 * time.Second)

	require.Len(t, bb0.h.PubSub().GetTopics(), 2)
	require.Len(t, bb1.h.PubSub().GetTopics(), 2)
	require.Len(t, bb1.h.PubSub().GetTopics(), 2)

	peers := []peer.ID{h2.ID(), h1.ID()}
	sort.Slice(peers, func(i, j int) bool {
		return strings.Compare(peers[i].String(), peers[j].String()) < 0
	})

	require.Equal(t, bb0.h.PubSub().ListPeers(ringTopic), peers)
	require.Equal(t, bb1.h.PubSub().ListPeers(ringTopic), []peer.ID{h0.ID()})
	require.Equal(t, bb2.h.PubSub().ListPeers(ringTopic), []peer.ID{h0.ID()})
}

func setupTestBulletins(t *testing.T, ctx context.Context) (*Bulletin, *Bulletin, *Bulletin) {
	h0 := newDefaultP2PHost(t, ctx)
	h1 := newRandomP2PHost(t, ctx)
	h2 := newRandomP2PHost(t, ctx)

	// bulletin setup
	cfg0, err := config.Default[config.Bulletin]()
	require.NoError(t, err)

	bb0, err := New(ctx, h0, cfg0)
	require.NoError(t, err)
	require.NotNil(t, bb0)

	cfg1, err := config.Default[config.Bulletin]()
	require.NoError(t, err)
	addr := h0.Address()
	cfg1.P2P.PersistentPeers = fmt.Sprintf("%s/p2p/%s", addr, h0.ID())

	bb1, err := New(ctx, h1, cfg1)
	require.NoError(t, err)
	require.NotNil(t, bb1)

	cfg2, err := config.Default[config.Bulletin]()
	require.NoError(t, err)
	cfg2.P2P.PersistentPeers = cfg1.P2P.PersistentPeers

	bb2, err := New(ctx, h2, cfg2)
	require.NoError(t, err)
	require.NotNil(t, bb2)

	// wait for the net connections
	time.Sleep(2 * time.Second)

	return bb0, bb1, bb2
}

func TestBulletinPostBroadcastLocalRead(t *testing.T) {
	ctx := context.Background()
	bb0, bb1, bb2 := setupTestBulletins(t, ctx)

	// test peer connections
	require.Len(t, bb0.h.Network().Conns(), 2)
	require.Len(t, bb1.h.Network().Conns(), 1)
	require.Len(t, bb2.h.Network().Conns(), 1)

	ringID := "123"
	ringTopic := "/ring/" + ringID
	bb0.Register(ctx, ringTopic)
	bb1.Register(ctx, ringTopic)
	bb2.Register(ctx, ringTopic)

	time.Sleep(2 * time.Second)

	msgType := ringTopic + "/dkg/rabin"
	msgID := msgType + "/1"
	msgBuf := []byte("helloworld")

	msg, err := newMessage(bb0, ringID, msgType, msgBuf)
	require.NoError(t, err)

	resp, err := bb0.Post(ctx, msgID, msg)
	require.NoError(t, err)
	require.NotEmpty(t, resp)
	require.Equal(t, ringID, resp.Data.RingId)
	require.Equal(t, msgBuf, resp.Data.Payload)
	require.Equal(t, msgType, resp.Data.Type)
	require.Equal(t, bb0.h.ID().String(), resp.Data.NodeId)

	// wait for post pubsub
	time.Sleep(2 * time.Second)

	// query the internal memory bulletin directly
	// which will avoid the pubsub read request, and
	// rely soley on the post request populating our
	// local in-memory bulletin state
	resp, err = bb1.mem.Read(ctx, msgID)
	require.NoError(t, err)
	require.NotEmpty(t, resp)
	require.Equal(t, ringID, resp.Data.RingId)
	require.Equal(t, msgBuf, resp.Data.Payload)
	require.Equal(t, msgType, resp.Data.Type)
	require.Equal(t, bb0.h.ID().String(), resp.Data.NodeId)
}

func TestBulletinPostBroadcastPubSubRead(t *testing.T) {
	ctx := context.Background()
	bb0, bb1, bb2 := setupTestBulletins(t, ctx)

	ringID := "123"
	ringTopic := "/ring/" + ringID
	bb0.Register(ctx, ringTopic)
	bb1.Register(ctx, ringTopic)
	bb2.Register(ctx, ringTopic)

	time.Sleep(2 * time.Second)

	msgType := ringTopic + "/dkg/rabin"
	msgID := msgType + "/1"
	msgBuf := []byte("helloworld")

	msg, err := newMessage(bb0, ringID, msgType, msgBuf)
	require.NoError(t, err)

	resp, err := bb0.Post(ctx, msgID, msg)
	require.NoError(t, err)
	require.NotEmpty(t, resp)
	require.Equal(t, ringID, resp.Data.RingId)
	require.Equal(t, msgBuf, resp.Data.Payload)
	require.Equal(t, msgType, resp.Data.Type)
	require.Equal(t, bb0.h.ID().String(), resp.Data.NodeId)

	// WE ARE NOT WAITING FOR PUBSUB
	//time.Sleep(2 * time.Second)

	// query the main store
	// todo: probably need a better way to gurantee
	// network read request
	resp, err = bb1.Read(ctx, msgID)
	require.NoError(t, err)
	require.NotEmpty(t, resp)
	require.Equal(t, ringID, resp.Data.RingId)
	require.Equal(t, msgBuf, resp.Data.Payload)
	require.Equal(t, msgType, resp.Data.Type)
	require.Equal(t, bb0.h.ID().String(), resp.Data.NodeId)
}

func TestBulletinPostBroadcastIndirectGossipRead(t *testing.T) {
	ctx := context.Background()
	bb0, bb1, bb2 := setupTestBulletins(t, ctx)

	ringID := "123"
	ringTopic := "/ring/" + ringID
	bb0.Register(ctx, ringTopic)
	bb1.Register(ctx, ringTopic)
	bb2.Register(ctx, ringTopic)

	time.Sleep(2 * time.Second)

	msgType := ringTopic + "/dkg/rabin"
	msgID := msgType + "/1"
	msgBuf := []byte("helloworld")

	msg, err := newMessage(bb1, ringID, msgType, msgBuf)
	require.NoError(t, err)

	resp, err := bb1.Post(ctx, msgID, msg)
	require.NoError(t, err)
	require.NotEmpty(t, resp)
	require.Equal(t, ringID, resp.Data.RingId)
	require.Equal(t, msgBuf, resp.Data.Payload)
	require.Equal(t, msgType, resp.Data.Type)
	require.Equal(t, bb1.h.ID().String(), resp.Data.NodeId)

	// WE ARE NOT WAITING FOR PUBSUB
	//time.Sleep(2 * time.Second)

	// NOTE: we are submitting the post request on node1
	// and reading from node2, which aren't directly connected
	// (at least initially), so were testing the gossip pubsub
	// todo: probably need a better way to gurantee
	// network read request
	resp, err = bb2.Read(ctx, msgID)
	require.NoError(t, err)
	require.NotEmpty(t, resp)
	require.Equal(t, ringID, resp.Data.RingId)
	require.Equal(t, msgBuf, resp.Data.Payload)
	require.Equal(t, msgType, resp.Data.Type)
	require.Equal(t, bb1.h.ID().String(), resp.Data.NodeId)
}

func TestBulletinEvents(t *testing.T) {
	ctx := context.Background()
	bb0, bb1, bb2 := setupTestBulletins(t, ctx)

	ringID := "123"
	ringTopic := "/ring/" + ringID
	bb0.Register(ctx, ringTopic)
	bb1.Register(ctx, ringTopic)
	bb2.Register(ctx, ringTopic)

	// get eventbus handle from node2
	bus := bb1.Events()
	require.NotNil(t, bus)

	subCh, err := eventbus.Subscribe[*transport.Message](bus)
	require.NoError(t, err)
	require.NotNil(t, subCh)

	time.Sleep(2 * time.Second)

	msgType := ringTopic + "/dkg/rabin"
	msgID := msgType + "/1"
	msgBuf := []byte("helloworld")

	// doneCh to track event completion
	doneCh := make(chan struct{})
	go func() {
		data := <-subCh
		require.NotEmpty(t, data)
		require.Equal(t, ringID, data.RingId)
		require.Equal(t, msgBuf, data.Payload)
		require.Equal(t, msgType, data.Type)
		require.Equal(t, bb0.h.ID().String(), data.NodeId)
		doneCh <- struct{}{}
	}()

	msg, err := newMessage(bb0, ringID, msgType, msgBuf)
	require.NoError(t, err)

	resp, err := bb0.Post(ctx, msgID, msg)
	require.NoError(t, err)
	require.NotEmpty(t, resp)
	require.Equal(t, ringID, resp.Data.RingId)
	require.Equal(t, msgBuf, resp.Data.Payload)
	require.Equal(t, msgType, resp.Data.Type)
	require.Equal(t, bb0.h.ID().String(), resp.Data.NodeId)

	select {
	case <-doneCh:
		return
	case <-time.After(5 * time.Second):
		// reached timeout before done, error
		t.Fatal("timeout reached before bulletin events subscription")
	}
}

func TestBulletinPostAndLocalQuery(t *testing.T) {
	ctx := context.Background()
	bb0, bb1, bb2 := setupTestBulletins(t, ctx)

	ringID := "123"
	ringTopic := "/ring/" + ringID
	bb0.Register(ctx, ringTopic)
	bb1.Register(ctx, ringTopic)
	bb2.Register(ctx, ringTopic)

	time.Sleep(2 * time.Second)

	msgType := ringTopic + "/dkg/rabin"
	msgID := msgType + "/1"
	msgBuf := []byte("helloworld")

	msg, err := newMessage(bb0, ringID, msgType, msgBuf)
	require.NoError(t, err)

	resp, err := bb0.Post(ctx, msgID, msg)
	require.NoError(t, err)
	require.NotEmpty(t, resp)
	require.Equal(t, ringID, resp.Data.RingId)
	require.Equal(t, msgBuf, resp.Data.Payload)
	require.Equal(t, msgType, resp.Data.Type)
	require.Equal(t, bb0.h.ID().String(), resp.Data.NodeId)

	// WAITING FOR PUBSUB
	time.Sleep(2 * time.Second)

	// query the internal local store
	// todo: probably need a better way to gurantee
	// network read request
	respCh, err := bb1.mem.Query(ctx, "*")
	require.NoError(t, err)
	require.NotNil(t, respCh)

	// should only be 1 event
	count := 0
	for evt := range respCh {
		require.NoError(t, evt.Err)
		require.NotEmpty(t, evt.Resp)
		require.Equal(t, ringID, evt.Resp.Data.RingId)
		require.Equal(t, msgBuf, evt.Resp.Data.Payload)
		require.Equal(t, msgType, evt.Resp.Data.Type)
		require.Equal(t, bb0.h.ID().String(), evt.Resp.Data.NodeId)
		count++
	}
	require.Equal(t, 1, count)

	// lets submit some more posts to make sure the query is working
	_, err = bb0.Post(ctx, msgType+"/2", msg)
	require.NoError(t, err)
	_, err = bb0.Post(ctx, msgType+"/3", msg)
	require.NoError(t, err)

	// WAITING FOR PUBSUB
	time.Sleep(2 * time.Second)

	// rerun the query
	respCh, err = bb1.mem.Query(ctx, "*")
	require.NoError(t, err)
	require.NotNil(t, respCh)

	// should be 3 events
	count = 0
	for range respCh {
		// TODO: actually verify all 3 event states
		count++
	}
	require.Equal(t, 3, count)
}

func TestBulletinRemoteQuery(t *testing.T) {
	ctx := context.Background()
	bb0, bb1, bb2 := setupTestBulletins(t, ctx)

	ringID := "123"
	ringTopic := "/ring/" + ringID
	bb0.Register(ctx, ringTopic)
	bb1.Register(ctx, ringTopic)
	bb2.Register(ctx, ringTopic)

	time.Sleep(2 * time.Second)

	msgType := ringTopic + "/dkg/rabin"
	msgBuf := []byte("helloworld")

	msg, err := newMessage(bb0, ringID, msgType, msgBuf)
	require.NoError(t, err)

	// were going to post everything locally and avoid the pubsub
	// so that we can test the net queries are actually calling out
	_, err = bb0.mem.Post(ctx, msgType+"/1", msg)
	require.NoError(t, err)
	_, err = bb0.mem.Post(ctx, msgType+"/2", msg)
	require.NoError(t, err)
	_, err = bb0.mem.Post(ctx, msgType+"/3", msg)
	require.NoError(t, err)

	respCh, err := bb1.Query(ctx, "*")
	require.NoError(t, err)

	// just count for now, and we can verify the local state afterwards
	count := 0
	for range respCh {
		count++
	}
	require.Equal(t, 3, count)

	assertEqualBulletinState(t, bb0, bb1)
}

func assertEqualBulletinState(t *testing.T, b0 *Bulletin, b1 *Bulletin) {
	ctx := context.Background()
	b0respCh, err := b0.mem.Query(ctx, "*")
	require.NoError(t, err)
	b1respCh, err := b1.mem.Query(ctx, "*")
	require.NoError(t, err)

	b0resp := channelToMap(b0respCh)
	b1resp := channelToMap(b1respCh)

	require.Equal(t, b0resp, b1resp)
}

// ChannelToSlice returns a slice built from channels items. Blocks until channel closes.
func channelToMap(ch <-chan bulletin.QueryResponse) map[string]*transport.Message {
	collection := make(map[string]*transport.Message)

	for item := range ch {
		collection[item.Resp.ID] = item.Resp.Data
	}

	return collection
}

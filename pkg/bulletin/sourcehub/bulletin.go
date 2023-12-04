package sourcehub

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	logging "github.com/ipfs/go-log"
	"google.golang.org/protobuf/proto"

	eventbus "github.com/sourcenetwork/eventbus-go"

	"github.com/sourcenetwork/orbis-go/config"
	gossipbulletinv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/gossipbulletin/v1alpha1"
	transportv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/transport/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/host"
	"github.com/sourcenetwork/orbis-go/pkg/transport"

	"sourcehub/x/bulletin/types"

	rpctypes "github.com/cometbft/cometbft/rpc/core/types"
	rpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"

	"github.com/ignite/cli/ignite/pkg/cosmosclient"
)

var log = logging.Logger("orbis/bulletin/sourcehub")

const (
	name          = "sourcehub"
	addressPrefix = "cosmos"
	accountName   = "alice"
	nodeAddress   = "http://host.docker.internal:26657"
	rpcAddress    = "tcp://host.docker.internal:26657"
)

var _ bulletin.Bulletin = (*Bulletin)(nil)

type Message = gossipbulletinv1alpha1.Message

type Bulletin struct {
	ctx context.Context

	client    cosmosclient.Client
	rpcClient *rpcclient.WSClient
	bus       eventbus.Bus
}

func New(ctx context.Context, host *host.Host, cfg config.Bulletin) (*Bulletin, error) {

	opts := cosmosclient.WithNodeAddress(nodeAddress)
	client, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix), opts)
	if err != nil {
		return nil, fmt.Errorf("new cosmos client: %w", err)
	}

	rpcClient, err := rpcclient.NewWS(rpcAddress, "/websocket")
	if err != nil {
		return nil, fmt.Errorf("new rpc client: %w", err)
	}

	err = rpcClient.Start()
	if err != nil {
		return nil, fmt.Errorf("rpc client start: %w", err)
	}

	err = rpcClient.Subscribe(ctx, "tm.event='Tx' AND NewPost.payload EXISTS")
	if err != nil {
		return nil, fmt.Errorf("subscribe to namespace: %w", err)
	}

	bus := eventbus.NewBus()

	bb := &Bulletin{
		ctx:       ctx,
		client:    client,
		rpcClient: rpcClient,
		bus:       bus,
	}

	go bb.HandleEvents()

	return bb, nil
}

func (bb *Bulletin) Name() string {
	return name
}

func (bb *Bulletin) Register(ctx context.Context, namespace string) error {
	if namespace == "" {
		return bulletin.ErrEmptyNamespace
	}

	return nil
}

func (bb *Bulletin) Post(ctx context.Context, id string, msg *transport.Message) (bulletin.Response, error) {
	var resp bulletin.Response

	account, err := bb.client.Account(accountName)
	if err != nil {
		return resp, fmt.Errorf("get account by name: %w", err)
	}

	addr, err := account.Address(addressPrefix)
	if err != nil {
		return resp, fmt.Errorf("get account address: %w", err)
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		return bulletin.Response{}, fmt.Errorf("marshal post message payload: %w", err)
	}

	hubMsg := &types.MsgCreatePost{
		Creator:   addr,
		Namespace: id,
		Payload:   payload,
		Proof:     nil,
	}

	resp.Data = msg
	resp.ID = id

	for retries := 20; retries > 0; retries-- {

		_, err = bb.client.BroadcastTx(ctx, account, hubMsg)
		if err == nil {
			log.Infof("Posted to bulletin, namespace: %s", id)
			return resp, nil
		}

		du := time.Duration(200+rand.Intn(1000)) * time.Millisecond
		log.Warnf("Broadcast tx: %s, retries(%d) in %s", err, retries, du)
		time.Sleep(du)
	}
	log.Errorf("Broadcast tx: %s", err)

	return resp, fmt.Errorf("broadcast tx: %w", err)
}

func (bb *Bulletin) Read(ctx context.Context, id string) (bulletin.Response, error) {
	var resp bulletin.Response

	queryClient := types.NewQueryClient(bb.client.Context())
	in := &types.QueryReadPostRequest{
		Namespace: id,
	}

	queryResp, err := queryClient.ReadPost(ctx, in)
	if err != nil {
		return resp, fmt.Errorf("query read post: %w", err)
	}

	var pbPayload transportv1alpha1.Message
	err = proto.Unmarshal(queryResp.Post.Payload, &pbPayload)
	if err != nil {
		return bulletin.Response{}, fmt.Errorf("unmarshal message payload: %w", err)
	}

	resp.Data = &pbPayload
	resp.ID = id

	return resp, nil
}

func (bb *Bulletin) Query(ctx context.Context, query string) (<-chan bulletin.QueryResponse, error) {
	if query == "" {
		return nil, fmt.Errorf("query can't be empty")
	}

	return nil, nil

}

func (bb *Bulletin) Verify(context.Context, bulletin.Proof, string, bulletin.Message) bool {
	return true
}

func (bb *Bulletin) Events() eventbus.Bus {
	return bb.bus
}

func (bb *Bulletin) HandleEvents() {

	for resp := range bb.rpcClient.ResponsesCh {
		result := &rpctypes.ResultEvent{}
		err := json.Unmarshal((resp.Result), result)
		if err != nil {
			log.Warnf("coud not unmarshal events resp: %v", err)
		}

		attrNamespace, ok := result.Events["NewPost.namespace"]
		if !ok {
			continue
		}
		attrPayload, ok := result.Events["NewPost.payload"]
		if !ok {
			continue
		}
		namespace := attrNamespace[0]
		b64Msg := attrPayload[0]
		rawMsg, err := base64.StdEncoding.DecodeString(b64Msg)
		if err != nil {
			log.Warnf("coud not decode base64 payload: %v", err)
			continue
		}

		var msg transportv1alpha1.Message
		if err := proto.Unmarshal(rawMsg, &msg); err != nil {
			log.Warnf("coud not unmarshal payload: %v", err)
			continue
		}

		evt := bulletin.Event{
			Message: &msg,
			ID:      namespace,
		}

		err = eventbus.Publish(bb.bus, evt)
		if err != nil {
			log.Warnf("failed to publish event to channel: %w", err)
			continue
		}
	}
}

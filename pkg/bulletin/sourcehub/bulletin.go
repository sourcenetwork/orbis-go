package sourcehub

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	logging "github.com/ipfs/go-log"
	"google.golang.org/protobuf/proto"

	eventbus "github.com/sourcenetwork/eventbus-go"

	"github.com/sourcenetwork/orbis-go/config"
	gossipbulletinv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/gossipbulletin/v1alpha1"
	transportv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/transport/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/host"
	"github.com/sourcenetwork/orbis-go/pkg/transport"

	"github.com/sourcenetwork/sourcehub/x/bulletin/types"

	rpctypes "github.com/cometbft/cometbft/rpc/core/types"
	rpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
)

var log = logging.Logger("orbis/bulletin/sourcehub")

const name = "sourcehub"

var _ bulletin.Bulletin = (*Bulletin)(nil)

type Message = gossipbulletinv1alpha1.Message

type Bulletin struct {
	ctx context.Context
	cfg config.Bulletin

	client    cosmosclient.Client
	account   cosmosaccount.Account
	address   string
	rpcClient *rpcclient.WSClient
	bus       eventbus.Bus
}

func New(ctx context.Context, host *host.Host, cfg config.Bulletin) (*Bulletin, error) {

	bb := &Bulletin{
		ctx: ctx,
		cfg: cfg,
	}

	return bb, nil
}

func (bb *Bulletin) Name() string {
	return name
}

func (bb *Bulletin) Init(ctx context.Context) error {

	opts := []cosmosclient.Option{
		cosmosclient.WithNodeAddress(bb.cfg.SourceHub.NodeAddress),
		cosmosclient.WithAddressPrefix(bb.cfg.SourceHub.AddressPrefix),
		cosmosclient.WithFees(bb.cfg.SourceHub.Fees),
	}
	client, err := cosmosclient.New(ctx, opts...)
	if err != nil {
		return fmt.Errorf("new cosmos client: %w", err)
	}

	account, err := client.Account(bb.cfg.SourceHub.AccountName)
	if err != nil {
		return fmt.Errorf("get account by name: %w", err)
	}

	address, err := account.Address(bb.cfg.SourceHub.AddressPrefix)
	if err != nil {
		return fmt.Errorf("get account address: %w", err)
	}

	rpcClient, err := rpcclient.NewWS(bb.cfg.SourceHub.RPCAddress, "/websocket")
	if err != nil {
		return fmt.Errorf("new rpc client: %w", err)
	}

	err = rpcClient.Start()
	if err != nil {
		return fmt.Errorf("rpc client start: %w", err)
	}

	err = rpcClient.Subscribe(ctx, "tm.event='Tx' AND NewPost.payload EXISTS")
	if err != nil {
		return fmt.Errorf("subscribe to namespace: %w", err)
	}

	bus := eventbus.NewBus()

	bb.ctx = ctx
	bb.client = client
	bb.account = account
	bb.address = address
	bb.rpcClient = rpcClient
	bb.bus = bus

	go bb.HandleEvents()

	return nil
}

func (bb *Bulletin) Register(ctx context.Context, namespace string) error {
	if namespace == "" {
		return bulletin.ErrEmptyNamespace
	}

	return nil
}

func (bb *Bulletin) Post(ctx context.Context, namespace, id string, msg *transport.Message) (bulletin.Response, error) {
	var resp bulletin.Response

	payload, err := proto.Marshal(msg)
	if err != nil {
		return bulletin.Response{}, fmt.Errorf("marshal post message payload: %w", err)
	}

	id = namespace + id
	hubMsg := &types.MsgCreatePost{
		Creator:   bb.address,
		Namespace: id,
		Payload:   payload,
		Proof:     nil,
	}

	resp.Data = msg
	resp.ID = id

	_, err = bb.client.BroadcastTx(ctx, bb.account, hubMsg)
	if err != nil {
		return resp, fmt.Errorf("broadcast tx: %w", err)
	}
	log.Infof("Posted to bulletin, namespace: %s", id)

	return resp, nil
}

func (bb *Bulletin) Read(ctx context.Context, namespace, id string) (bulletin.Response, error) {
	var resp bulletin.Response

	queryClient := types.NewQueryClient(bb.client.Context())
	id = namespace + id
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

func (bb *Bulletin) Query(ctx context.Context, namespace, query string) (<-chan bulletin.QueryResponse, error) {
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

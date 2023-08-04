package grpcserver

import (
	"context"

	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/app"
	transportv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/transport/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/transport/p2p"

	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto/pb"
	"github.com/multiformats/go-multiaddr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type transportService struct {
	transportv1alpha1.UnimplementedTransportServiceServer
	app *app.App
}

func newTransportService(app *app.App) *transportService {
	return &transportService{
		app: app,
	}
}

func (s *transportService) GetHost(ctx context.Context, req *transportv1alpha1.GetHostRequest) (*transportv1alpha1.GetHostResponse, error) {

	tp, ok := getTransport(s.app.Injector(), req.Transport)
	if !ok {
		return nil, status.Error(codes.NotFound, "transport doens't exist")
	}

	raw, err := tp.Host().PublicKey().Raw()
	if err != nil {
		return nil, err
	}

	resp := &transportv1alpha1.GetHostResponse{
		Node: &transportv1alpha1.Node{
			Id:      tp.Host().ID(),
			Address: tp.Host().Address().String(),
			PublicKey: &libp2pcrypto.PublicKey{
				Type: libp2pcrypto.KeyType_Ed25519.Enum(),
				Data: raw,
			},
		},
	}

	return resp, nil
}

func (s *transportService) Send(ctx context.Context, req *transportv1alpha1.SendRequest) (*transportv1alpha1.SendResponse, error) {

	resp := &transportv1alpha1.SendResponse{}

	return resp, errUnimplemented
}

func (s *transportService) Connect(ctx context.Context, req *transportv1alpha1.ConnectRequest) (*transportv1alpha1.ConnectResponse, error) {
	tp, ok := getTransport(s.app.Injector(), req.Transport)
	if !ok {
		return nil, status.Error(codes.NotFound, "transport doens't exist")
	}

	addr, err := multiaddr.NewMultiaddr(req.GetAddress())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	n := p2p.NewNode(req.GetId(), nil, addr)

	err = tp.Connect(ctx, n)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &transportv1alpha1.ConnectResponse{}

	return resp, nil
}

func (s *transportService) Gossip(ctx context.Context, req *transportv1alpha1.GossipRequest) (*transportv1alpha1.GossipResponse, error) {
	tp, ok := getTransport(s.app.Injector(), req.Transport)
	if !ok {
		return nil, status.Error(codes.NotFound, "transport doens't exist")
	}

	err := tp.Gossip(ctx, req.GetTopic(), req.GetMessage())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &transportv1alpha1.GossipResponse{}

	return resp, nil
}

func (s *transportService) NewMessage(ctx context.Context, req *transportv1alpha1.NewMessageRequest) (*transportv1alpha1.NewMessageResponse, error) {
	tp, ok := getTransport(s.app.Injector(), req.Transport)
	if !ok {
		return nil, status.Error(codes.NotFound, "transport doens't exist")
	}

	pubkey, err := tp.Host().PublicKey().Raw()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// TODO: Check what else needs to be signed.
	sig, err := tp.Host().Sign(req.Payload)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	msg := &transportv1alpha1.Message{
		Id:      req.Id,
		Type:    req.Type,
		Payload: req.Payload,
		Gossip:  req.Gossip,
		RingId:  req.RingId,

		NodeId:     tp.Host().ID(),
		NodePubKey: pubkey,
		Signature:  sig,
	}

	resp := &transportv1alpha1.NewMessageResponse{
		Message: msg,
	}

	return resp, nil
}

func getTransport(inj *do.Injector, transportName string) (transport.Transport, bool) {
	tp, err := do.InvokeNamed[transport.Transport](inj, transportName)
	if err != nil {
		return nil, false
	}
	return tp, true
}

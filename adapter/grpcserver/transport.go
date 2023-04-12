package grpcserver

import (
	"context"

	cryptov1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/libp2p/crypto/v1alpha1"
	transportv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/transport/v1alpha1"
	"github.com/sourcenetwork/orbis-go/infra/logger"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

type transportService struct {
	transportv1alpha1.UnimplementedTransportServiceServer
	lg logger.Logger
	tp transport.Transport
}

func newTransportService(lg logger.Logger, tp transport.Transport) *transportService {
	return &transportService{
		lg: lg,
		tp: tp,
	}
}

func (s *transportService) GetHost(ctx context.Context, req *transportv1alpha1.GetHostRequest) (*transportv1alpha1.GetHostResponse, error) {

	h := s.tp.Host()

	raw, err := h.PublicKey().Raw()
	if err != nil {
		return nil, err
	}
	resp := &transportv1alpha1.GetHostResponse{
		Node: &transportv1alpha1.Node{
			Id:      h.ID(),
			Address: h.Address().String(),
			PublicKey: &libp2pcrypto.PublicKey{
				Type: libp2pcrypto.KeyType_Ed25519.Enum(),
				Data: raw,
			},
		},
	}

	return resp, nil
}

func (s *transportService) ListNodes(ctx context.Context, req *transportv1alpha1.ListNodesRequest) (*transportv1alpha1.ListNodesResponse, error) {

	s.lg.Debugf("ListNode()")

	Nodes := []*transportv1alpha1.Node{
		{
			Id: "Node1",
		},
		{
			Id: "Node2",
		},
	}

	resp := &transportv1alpha1.ListNodesResponse{
		Nodes: Nodes,
	}

	return resp, nil
}

func (s *transportService) Send(ctx context.Context, req *transportv1alpha1.SendRequest) (*transportv1alpha1.SendResponse, error) {

	s.lg.Debugf("Send()")

	resp := &transportv1alpha1.SendResponse{}

	return resp, errUnimplemented
}

func (s *transportService) Connect(ctx context.Context, req *transportv1alpha1.ConnectRequest) (*transportv1alpha1.ConnectResponse, error) {

	s.lg.Debugf("Connect()")

	resp := &transportv1alpha1.ConnectResponse{}

	return resp, errUnimplemented
}

func (s *transportService) Gossip(ctx context.Context, req *transportv1alpha1.GossipRequest) (*transportv1alpha1.GossipResponse, error) {

	s.lg.Debugf("Gossip()")

	resp := &transportv1alpha1.GossipResponse{}

	return resp, errUnimplemented
}

func (s *transportService) NewMessage(ctx context.Context, req *transportv1alpha1.NewMessageRequest) (*transportv1alpha1.NewMessageResponse, error) {

	s.lg.Debugf("NewMessage()")

	return nil, errUnimplemented
}

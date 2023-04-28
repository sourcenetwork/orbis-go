package grpcserver

import (
	"context"

	p2pv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/p2p/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/p2p"

	libp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type p2pService struct {
	p2pv1alpha1.UnimplementedP2PServiceServer

	h *p2p.Host
}

func newP2PService(h *p2p.Host) *p2pService {
	return &p2pService{
		h: h,
	}
}

func (s *p2pService) Host(ctx context.Context, req *p2pv1alpha1.HostRequest) (*p2pv1alpha1.HostResponse, error) {

	pi := libp2ppeer.AddrInfo{
		ID:    s.h.ID(),
		Addrs: s.h.Addrs(),
	}

	resp := &p2pv1alpha1.HostResponse{
		Id: pi.String(),
	}

	return resp, nil
}

func (s *p2pService) Peers(ctx context.Context, req *p2pv1alpha1.PeersRequest) (*p2pv1alpha1.PeersResponse, error) {

	resp := &p2pv1alpha1.PeersResponse{
		Ids: s.h.Peers(),
	}

	return resp, nil
}

func (s *p2pService) Send(ctx context.Context, req *p2pv1alpha1.SendRequest) (*p2pv1alpha1.SendResponse, error) {

	str := req.PeerInfo
	pi, err := libp2ppeer.AddrInfoFromString(str)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.h.Send(ctx, *pi, req.Protocol, req.Data)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &p2pv1alpha1.SendResponse{}

	return resp, nil
}

func (s *p2pService) Connect(ctx context.Context, req *p2pv1alpha1.ConnectRequest) (*p2pv1alpha1.ConnectResponse, error) {

	pi, err := libp2ppeer.AddrInfoFromString(req.PeerInfo)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.h.Connect(ctx, *pi)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &p2pv1alpha1.ConnectResponse{}

	return resp, nil
}

func (s *p2pService) Publish(ctx context.Context, req *p2pv1alpha1.PublishRequest) (*p2pv1alpha1.PublishResponse, error) {

	err := s.h.Publish(ctx, req.Topic, req.Data)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &p2pv1alpha1.PublishResponse{}

	return resp, nil
}

func (s *p2pService) Subscribe(req *p2pv1alpha1.SubscribeRequest, srv p2pv1alpha1.P2PService_SubscribeServer) error {

	ctx := srv.Context()
	sub, err := s.h.Subscribe(ctx, req.Topic)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		err = srv.Send(&p2pv1alpha1.SubscribeResponse{
			Data: msg.Data,
		})
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}

package grpcserver

import (
	"context"

	hostv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/host/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/host"

	libp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type hostService struct {
	hostv1alpha1.UnimplementedHostServiceServer

	h *host.Host
}

func newHostService(h *host.Host) *hostService {
	return &hostService{
		h: h,
	}
}

func (s *hostService) Host(ctx context.Context, req *hostv1alpha1.HostRequest) (*hostv1alpha1.HostResponse, error) {

	pi := libp2ppeer.AddrInfo{
		ID:    s.h.ID(),
		Addrs: s.h.Addrs(),
	}

	resp := &hostv1alpha1.HostResponse{
		Id: pi.String(),
	}

	return resp, nil
}

func (s *hostService) Peers(ctx context.Context, req *hostv1alpha1.PeersRequest) (*hostv1alpha1.PeersResponse, error) {

	resp := &hostv1alpha1.PeersResponse{
		Ids: s.h.Peers(),
	}

	return resp, nil
}

func (s *hostService) Send(ctx context.Context, req *hostv1alpha1.SendRequest) (*hostv1alpha1.SendResponse, error) {

	str := req.PeerInfo
	pi, err := libp2ppeer.AddrInfoFromString(str)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.h.Send(ctx, *pi, req.Protocol, req.Data)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &hostv1alpha1.SendResponse{}

	return resp, nil
}

func (s *hostService) Connect(ctx context.Context, req *hostv1alpha1.ConnectRequest) (*hostv1alpha1.ConnectResponse, error) {

	pi, err := libp2ppeer.AddrInfoFromString(req.PeerInfo)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.h.Connect(ctx, *pi)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &hostv1alpha1.ConnectResponse{}

	return resp, nil
}

func (s *hostService) Publish(ctx context.Context, req *hostv1alpha1.PublishRequest) (*hostv1alpha1.PublishResponse, error) {

	err := s.h.Publish(ctx, req.Topic, req.Data)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &hostv1alpha1.PublishResponse{}

	return resp, nil
}

func (s *hostService) Subscribe(req *hostv1alpha1.SubscribeRequest, srv hostv1alpha1.HostService_SubscribeServer) error {

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
		err = srv.Send(&hostv1alpha1.SubscribeResponse{
			Data: msg.Data,
		})
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}

package grpcserver

import (
	"context"

	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
	"github.com/sourcenetwork/orbis-go/infra/logger"

	"google.golang.org/protobuf/types/known/emptypb"
)

// ringService wraps application to provides gRPCs.
type ringService struct {
	ringv1alpha1.UnimplementedRingServiceServer
	lg logger.Logger
}

func newRingService(lg logger.Logger) *ringService {
	return &ringService{
		lg: lg,
	}
}

func (s *ringService) ListRings(ctx context.Context, req *ringv1alpha1.ListRingsRequest) (*ringv1alpha1.ListRingsResponse, error) {

	rings := []*ringv1alpha1.Ring{
		{
			Id: "ring1",
		},
		{
			Id: "ring2",
		},
	}

	resp := &ringv1alpha1.ListRingsResponse{
		Rings: rings,
	}

	return resp, nil
}

func (s *ringService) CreateRing(ctx context.Context, req *ringv1alpha1.CreateRingRequest) (*ringv1alpha1.CreateRingResponse, error) {

	// err := s.app.CreatRing(...)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	resp := &ringv1alpha1.CreateRingResponse{}

	return resp, errUnimplemented
}

func (s *ringService) GetRing(ctx context.Context, req *ringv1alpha1.GetRingRequest) (*ringv1alpha1.GetRingResponse, error) {

	// value, err := s.app.GetRing()
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	resp := &ringv1alpha1.GetRingResponse{}

	return resp, errUnimplemented
}

func (s *ringService) DeleteRing(ctx context.Context, req *ringv1alpha1.DeleteRingRequest) (*emptypb.Empty, error) {

	// err := s.app.DeleteRing()
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	return nil, errUnimplemented
}

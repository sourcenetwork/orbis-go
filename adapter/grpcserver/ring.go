package grpcserver

import (
	"context"

	"github.com/sourcenetwork/orbis-go/app"
	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/ring/v1alpha1"

	"google.golang.org/protobuf/types/known/emptypb"
)

// ringService wraps application to provides gRPCs.
type ringService struct {
	ringv1alpha1.UnimplementedRingServiceServer
	app *app.App
}

func newRingService(app *app.App) *ringService {
	return &ringService{
		app: app,
	}
}

func (s *ringService) ListRings(ctx context.Context, req *ringv1alpha1.ListRingsRequest) (*ringv1alpha1.ListRingsResponse, error) {

	s.app.Logger().Debugf("ListRing()")

	// err := s.app.ListRings(...)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	// For demo purpose.
	// All processing should be done in the app.
	// One day, all these stubs code/files will be generated.
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

	s.app.Logger().Debugf("CreateRing()")

	// err := s.app.CreatRing(...)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	resp := &ringv1alpha1.CreateRingResponse{}

	return resp, errUnimplemented
}

func (s *ringService) GetRing(ctx context.Context, req *ringv1alpha1.GetRingRequest) (*ringv1alpha1.GetRingResponse, error) {

	s.app.Logger().Debugf("GetRing()")

	// value, err := s.app.GetRing()
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	resp := &ringv1alpha1.GetRingResponse{}

	return resp, errUnimplemented
}

func (s *ringService) DeleteRing(ctx context.Context, req *ringv1alpha1.DeleteRingRequest) (*emptypb.Empty, error) {

	s.app.Logger().Debugf("DeleteRing()")

	// err := s.app.DeleteRing()
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	return nil, errUnimplemented
}

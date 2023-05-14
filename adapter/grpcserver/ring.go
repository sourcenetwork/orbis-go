package grpcserver

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/orbis-go/app"
	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/types"

	"google.golang.org/protobuf/types/known/emptypb"
)

// ringService wraps application to provides gRPCs.
type ringService struct {
	app *app.App
	ringv1alpha1.UnimplementedRingServiceServer
}

func newRingService(app *app.App) *ringService {
	return &ringService{
		app: app,
	}
}

func (s *ringService) ListRings(ctx context.Context, req *ringv1alpha1.ListRingsRequest) (*ringv1alpha1.ListRingsResponse, error) {

	resp := &ringv1alpha1.ListRingsResponse{}

	return resp, nil
}

func (s *ringService) CreateRing(ctx context.Context, req *ringv1alpha1.CreateRingRequest) (*ringv1alpha1.CreateRingResponse, error) {

	// err := s.app.CreatRing(...)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	manifest := &types.Ring{
		Ring: ringv1alpha1.Ring{
			Id:        "40b086ef",
			N:         3,
			T:         2,
			Dkg:       "rabin",
			Pss:       "avpss",
			Pre:       "elgamal",
			Bulletin:  "p2pbb",
			Transport: "p2p",
			Nodes:     nil,
		},
	}

	// repo, err := db.New()
	// if err != nil {
	// 	return nil, fmt.Errorf("create ring repo: %w", err)
	// }

	rr, err := s.app.NewRing(ctx, manifest)
	if err != nil {
		return nil, fmt.Errorf("create ring: %w", err)
	}
	_ = rr
	resp := &ringv1alpha1.CreateRingResponse{}

	return resp, nil
}

func (s *ringService) GetRing(ctx context.Context, req *ringv1alpha1.GetRingRequest) (*ringv1alpha1.GetRingResponse, error) {

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

func (s *ringService) PublicKey(ctx context.Context, req *ringv1alpha1.PublicKeyRequest) (*ringv1alpha1.PublicKeyResponse, error) {

	resp := &ringv1alpha1.PublicKeyResponse{}

	return resp, errUnimplemented
}

func (s *ringService) Refresh(ctx context.Context, req *ringv1alpha1.RefreshRequest) (*ringv1alpha1.RefreshResponse, error) {

	resp := &ringv1alpha1.RefreshResponse{}

	return resp, errUnimplemented
}

func (s *ringService) State(ctx context.Context, req *ringv1alpha1.StateRequest) (*ringv1alpha1.StateResponse, error) {

	resp := &ringv1alpha1.StateResponse{}

	return resp, errUnimplemented
}
func (s *ringService) Nodes(ctx context.Context, req *ringv1alpha1.NodesRequest) (*ringv1alpha1.NodesResponse, error) {

	resp := &ringv1alpha1.NodesResponse{}

	return resp, errUnimplemented
}

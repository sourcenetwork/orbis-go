package grpcserver

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/orbis-go/app"
	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"

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

	manifest := &ringv1alpha1.Ring{
		Id:        "40b086ef",
		N:         3,
		T:         2,
		Dkg:       "rabin",
		Pss:       "avpss",
		Pre:       "elgamal",
		Bulletin:  "p2pbb",
		Transport: "p2p",
		Nodes:     req.Ring.Nodes,
	}

	rr, err := s.app.JoinRing(ctx, manifest)
	if err != nil {
		return nil, fmt.Errorf("create ring: %w", err)
	}

	err = rr.Start(ctx)
	if err != nil {
		return nil, fmt.Errorf("starting ring: %w", err)
	}
	resp := &ringv1alpha1.CreateRingResponse{}

	return resp, nil
}

func (s *ringService) GetRing(ctx context.Context, req *ringv1alpha1.GetRingRequest) (*ringv1alpha1.GetRingResponse, error) {

	resp := &ringv1alpha1.GetRingResponse{}

	return resp, errUnimplemented
}

func (s *ringService) DeleteRing(ctx context.Context, req *ringv1alpha1.DeleteRingRequest) (*emptypb.Empty, error) {

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

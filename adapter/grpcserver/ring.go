package grpcserver

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/orbis-go/app"
	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ringService wraps application to provides gRPCs.
type ringService struct {
	ringv1alpha1.UnimplementedRingServiceServer

	app       *app.App
	rings     map[types.RingID]*app.Ring
	manifests map[types.RingID]*ringv1alpha1.Ring
}

func newRingService(a *app.App) *ringService {
	return &ringService{
		app:       a,
		rings:     map[types.RingID]*app.Ring{},
		manifests: map[types.RingID]*ringv1alpha1.Ring{},
	}
}

func (s *ringService) ListRings(ctx context.Context, req *ringv1alpha1.ListRingsRequest) (*ringv1alpha1.ListRingsResponse, error) {

	resp := &ringv1alpha1.ListRingsResponse{}

	return resp, nil
}

func (s *ringService) CreateRing(ctx context.Context, req *ringv1alpha1.CreateRingRequest) (*ringv1alpha1.CreateRingResponse, error) {

	r, err := s.app.JoinRing(ctx, req.Ring)
	if err != nil {
		return nil, fmt.Errorf("create ring: %w", err)
	}

	s.rings[r.ID] = r
	s.manifests[r.ID] = req.Ring

	resp := &ringv1alpha1.CreateRingResponse{
		Id: string(r.ID),
	}

	err = r.Start(ctx)
	if err != nil {
		return nil, fmt.Errorf("start ring: %w", err)
	}

	return resp, nil
}

func (s *ringService) GetRing(ctx context.Context, req *ringv1alpha1.GetRingRequest) (*ringv1alpha1.GetRingResponse, error) {

	r, ok := s.manifests[types.RingID(req.Id)]
	if !ok {
		return nil, status.Error(codes.NotFound, "ring not found")
	}

	resp := &ringv1alpha1.GetRingResponse{
		Ring: r,
	}

	return resp, nil
}

func (s *ringService) DeleteRing(ctx context.Context, req *ringv1alpha1.DeleteRingRequest) (*emptypb.Empty, error) {

	return nil, errUnimplemented
}

func (s *ringService) PublicKey(ctx context.Context, req *ringv1alpha1.PublicKeyRequest) (*ringv1alpha1.PublicKeyResponse, error) {

	r, ok := s.rings[types.RingID(req.Id)]
	if !ok {
		return nil, status.Error(codes.NotFound, "ring not found")
	}

	pub, err := r.PublicKey()
	if err != nil {
		return nil, status.Error(codes.Internal, "can't get public key")
	}

	pubProto, err := crypto.PublicKeyToProto(pub)
	if err != nil {
		return nil, status.Error(codes.Internal, "can't get public key")
	}

	resp := &ringv1alpha1.PublicKeyResponse{
		PublicKey: pubProto,
	}

	return resp, nil
}

func (s *ringService) Refresh(ctx context.Context, req *ringv1alpha1.RefreshRequest) (*ringv1alpha1.RefreshResponse, error) {

	resp := &ringv1alpha1.RefreshResponse{}

	return resp, errUnimplemented
}

func (s *ringService) State(ctx context.Context, req *ringv1alpha1.StateRequest) (*ringv1alpha1.StateResponse, error) {

	r, ok := s.rings[types.RingID(req.Id)]
	if !ok {
		return nil, status.Error(codes.NotFound, "ring not found")
	}
	r.State()

	resp := &ringv1alpha1.StateResponse{
		DkgState: r.State().DKG.String(),
		PssState: r.State().PSS.String(),
	}

	return resp, nil
}

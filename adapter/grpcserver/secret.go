package grpcserver

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/authz"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

func (s *ringService) ListSecrets(ctx context.Context, req *ringv1alpha1.ListSecretsRequest) (*ringv1alpha1.ListSecretsResponse, error) {

	// err := s.app.ListSecrets(...)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	resp := &ringv1alpha1.ListSecretsResponse{}

	return resp, errUnimplemented
}

func (s *ringService) StoreSecret(ctx context.Context, req *ringv1alpha1.StoreSecretRequest) (*ringv1alpha1.StoreSecretResponse, error) {

	// err := s.app.StoreSecret(...)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	resp := &ringv1alpha1.StoreSecretResponse{}

	return resp, errUnimplemented
}

func (s *ringService) GetSecret(ctx context.Context, req *ringv1alpha1.GetSecretRequest) (*ringv1alpha1.GetSecretResponse, error) {
	resp := &ringv1alpha1.GetSecretResponse{}

	ring, err := s.app.GetRing(ctx, types.RingID(req.RingId))
	if err != nil {
		return resp, err
	}

	authInfo, err := ring.Authn.GetAndVerifyRequestMetadata(ctx)
	ok, err := ring.Authz.Check(ctx, req.SecretId, authz.READ, authInfo.Subject)
	if err != nil {
		return resp, err
	}

	if !ok {
		return resp, errUnAuthorized
	}

	return resp, nil
}

func (s *ringService) DeleteSecret(ctx context.Context, req *ringv1alpha1.DeleteSecretRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, errUnimplemented
}

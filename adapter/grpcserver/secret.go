package grpcserver

import (
	"context"

	secretv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/secret/v1alpha1"

	"google.golang.org/protobuf/types/known/emptypb"
)

// secretService wraps application to provides gRPCs.
type secretService struct {
	secretv1alpha1.UnimplementedSecretServiceServer
}

func newSecretService() *secretService {
	return &secretService{}
}

func (s *secretService) ListSecrets(ctx context.Context, req *secretv1alpha1.ListSecretsRequest) (*secretv1alpha1.ListSecretsResponse, error) {

	// err := s.app.ListSecrets(...)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	resp := &secretv1alpha1.ListSecretsResponse{}

	return resp, errUnimplemented
}

func (s *secretService) StoreSecret(ctx context.Context, req *secretv1alpha1.StoreSecretRequest) (*secretv1alpha1.StoreSecretResponse, error) {

	// err := s.app.StoreSecret(...)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	resp := &secretv1alpha1.StoreSecretResponse{}

	return resp, errUnimplemented
}

func (s *secretService) GetSecret(ctx context.Context, req *secretv1alpha1.GetSecretRequest) (*secretv1alpha1.GetSecretResponse, error) {

	resp := &secretv1alpha1.GetSecretResponse{}

	return resp, errUnimplemented
}

func (s *secretService) DeleteSecret(ctx context.Context, req *secretv1alpha1.DeleteSecretRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, errUnimplemented
}

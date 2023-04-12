package grpcserver

import (
	"context"

	secretv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/secret/v1alpha1"
	"github.com/sourcenetwork/orbis-go/infra/logger"

	"google.golang.org/protobuf/types/known/emptypb"
)

// secretService wraps application to provides gRPCs.
type secretService struct {
	secretv1alpha1.UnimplementedSecretServiceServer
	lg logger.Logger
}

func newSecretService(lg logger.Logger) *secretService {
	return &secretService{
		lg: lg,
	}
}

func (s *secretService) ListSecrets(ctx context.Context, req *secretv1alpha1.ListSecretsRequest) (*secretv1alpha1.ListSecretsResponse, error) {

	s.lg.Debugf("ListSecret()")

	// err := s.app.ListSecrets(...)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	resp := &secretv1alpha1.ListSecretsResponse{}

	return resp, errUnimplemented
}

func (s *secretService) StoreSecret(ctx context.Context, req *secretv1alpha1.StoreSecretRequest) (*secretv1alpha1.StoreSecretResponse, error) {

	s.lg.Debugf("StoreSecret()")

	// err := s.app.StoreSecret(...)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	resp := &secretv1alpha1.StoreSecretResponse{}

	return resp, errUnimplemented
}

func (s *secretService) GetSecret(ctx context.Context, req *secretv1alpha1.GetSecretRequest) (*secretv1alpha1.GetSecretResponse, error) {

	s.lg.Debugf("GetSecret()")

	// value, err := s.app.GetSecret(...)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	resp := &secretv1alpha1.GetSecretResponse{}

	return resp, errUnimplemented
}

func (s *secretService) DeleteSecret(ctx context.Context, req *secretv1alpha1.DeleteSecretRequest) (*emptypb.Empty, error) {

	s.lg.Debugf("DeleteSecret()")

	// err := s.app.DeleteSecret(...)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	return &emptypb.Empty{}, errUnimplemented
}

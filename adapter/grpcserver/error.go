package grpcserver

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errUnimplemented = status.Error(codes.Unimplemented, "not implemented yet")
	errUnAuthorized  = status.Error(codes.PermissionDenied, "not authorized")
	errNotFound      = status.Error(codes.NotFound, "not found")
)

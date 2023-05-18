package grpcserver

import (
	"context"

	"google.golang.org/grpc"
)

func loggingInterceptor() grpc.ServerOption {

	interceptor := func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		resp, err := handler(ctx, req)

		// Output format favors readability in console over parsability.
		// Might be changed in the future for ingestors.
		if err != nil {
			log.Errorf("gRPC %s(%+v), error: %v", info.FullMethod, req, err)
		} else {
			log.Infof("gRPC %s(%+v): %+v", info.FullMethod, req, resp)
		}

		return resp, err
	}

	return grpc.UnaryInterceptor(interceptor)
}

package grpcserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sourcenetwork/orbis-go/app"
	"github.com/sourcenetwork/orbis-go/config"
	hostv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/host/v1alpha1"
	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
	transportv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/transport/v1alpha1"
	utilityv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/utility/v1alpha1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("orbis/grpc/server")

var (
	errUnimplemented = status.Error(codes.Unimplemented, "not implemented yet")
	errUnAuthorized  = status.Error(codes.PermissionDenied, "not authorized")
)

func NewGRPCServer(cfg config.GRPC, a *app.App) *grpc.Server {

	var opts []grpc.ServerOption
	if cfg.Logging {
		opts = append(opts, loggingInterceptor())
	}

	s := grpc.NewServer(opts...)

	// Setup orbis service handlers to the server.
	utilityv1alpha1.RegisterUtilityServiceServer(s, newUtilService(a))
	hostv1alpha1.RegisterHostServiceServer(s, newHostService(a.Host()))
	transportv1alpha1.RegisterTransportServiceServer(s, newTransportService(a))
	ringv1alpha1.RegisterRingServiceServer(s, newRingService(a))

	return s
}

func NewGRPCGatewayServer(cfg config.GRPC) (*http.Server, error) {

	ctx := context.Background()
	// Create a client connection to the gRPC server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	conn, err := grpc.DialContext(
		ctx,
		cfg.GRPCURL,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial to gRPC server %s, %w", cfg.GRPCURL, err)
	}

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err = ringv1alpha1.RegisterRingServiceHandler(ctx, mux, conn)
	if err != nil {
		return nil, fmt.Errorf("register ring service handler, %w", err)
	}

	err = utilityv1alpha1.RegisterUtilityServiceHandler(ctx, mux, conn)
	if err != nil {
		return nil, fmt.Errorf("register utility service handler, %w", err)
	}

	// Register Orbis Services.
	err = ringv1alpha1.RegisterRingServiceHandlerFromEndpoint(ctx, mux, cfg.GRPCURL, opts)
	if err != nil {
		return nil, fmt.Errorf("register ring service handler from endpoint, %w", err)
	}

	err = utilityv1alpha1.RegisterUtilityServiceHandlerFromEndpoint(ctx, mux, cfg.GRPCURL, opts)
	if err != nil {
		return nil, fmt.Errorf("register utility service handler from endpoint, %w", err)
	}

	gw := &http.Server{
		Addr:    cfg.RESTURL,
		Handler: mux,
	}

	return gw, nil
}

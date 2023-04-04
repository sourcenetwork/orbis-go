package grpcserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sourcenetwork/orbis-go/app"
	"github.com/sourcenetwork/orbis-go/config"
	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/ring/v1alpha1"
	secretv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/secret/v1alpha1"
	transportv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/transport/v1alpha1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	errUnimplemented = status.Error(codes.Unimplemented, "not implemented yet")
)

func NewGRPCServer(app *app.App) *grpc.Server {

	lg := app.Logger()
	s := grpc.NewServer()

	// Setup orbis service handlers to the server.
	transportv1alpha1.RegisterTransportServiceServer(s, newTransportService(lg, app.Transport()))
	ringv1alpha1.RegisterRingServiceServer(s, newRingService(lg))
	secretv1alpha1.RegisterSecretServiceServer(s, newSecretService(lg))

	return s
}

func NewGRPCGatewayServer(cfg config.GRPC) (*http.Server, error) {

	// Create a client connection to the gRPC server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	conn, err := grpc.DialContext(
		context.Background(),
		cfg.GRPCURL,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial to gRPC server %s, %w", cfg.GRPCURL, err)
	}

	mux := runtime.NewServeMux()

	// Register Orbis Services.
	err = ringv1alpha1.RegisterRingServiceHandler(context.Background(), mux, conn)
	if err != nil {
		return nil, fmt.Errorf("register ringService: %w", err)
	}

	err = secretv1alpha1.RegisterSecretServiceHandler(context.Background(), mux, conn)
	if err != nil {
		return nil, fmt.Errorf("register secretService: %w", err)
	}

	gw := &http.Server{
		Addr:    cfg.RESTURL,
		Handler: mux,
	}

	return gw, nil
}

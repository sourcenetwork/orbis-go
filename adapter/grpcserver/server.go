package grpcserver

import (
	"net/http"

	"github.com/sourcenetwork/orbis-go/app"
	"github.com/sourcenetwork/orbis-go/config"
	p2pv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/p2p/v1alpha1"
	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
	secretv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/secret/v1alpha1"
	transportv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/transport/v1alpha1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("orbis/grpc/server")

var (
	errUnimplemented = status.Error(codes.Unimplemented, "not implemented yet")
)

func NewGRPCServer(cfg config.GRPC, a *app.App) *grpc.Server {

	var opts []grpc.ServerOption
	if cfg.Logging {
		opts = append(opts, loggingInterceptor())
	}

	s := grpc.NewServer(opts...)

	// Setup orbis service handlers to the server.
	p2pv1alpha1.RegisterP2PServiceServer(s, newP2PService(a.Host()))
	transportv1alpha1.RegisterTransportServiceServer(s, newTransportService(a.Transport()))
	ringv1alpha1.RegisterRingServiceServer(s, newRingService())
	secretv1alpha1.RegisterSecretServiceServer(s, newSecretService())

	return s
}

func NewGRPCGatewayServer(cfg config.GRPC) (*http.Server, error) {

	// Create a client connection to the gRPC server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	// conn, err := grpc.DialContext(
	// 	context.Background(),
	// 	cfg.GRPCURL,
	// 	grpc.WithBlock(),
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("dial to gRPC server %s, %w", cfg.GRPCURL, err)
	// }

	mux := runtime.NewServeMux()

	// Register Orbis Services.
	// err = ringv1alpha1.RegisterRingServiceServer(context.Background(), mux, conn)
	// if err != nil {
	// 	return nil, fmt.Errorf("register ringService: %w", err)
	// }

	// err = secretv1alpha1.RegisterSecretServiceHandler(context.Background(), mux, conn)
	// if err != nil {
	// 	return nil, fmt.Errorf("register secretService: %w", err)
	// }

	gw := &http.Server{
		Addr:    cfg.RESTURL,
		Handler: mux,
	}

	return gw, nil
}

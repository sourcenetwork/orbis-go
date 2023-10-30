package main

import (
	"context"
	"fmt"
	"net"

	"github.com/NathanBaulch/protoc-gen-cobra/client"
	"github.com/NathanBaulch/protoc-gen-cobra/naming"
	"github.com/samber/do"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"github.com/sourcenetwork/orbis-go/adapter/grpcserver"
	"github.com/sourcenetwork/orbis-go/app"
	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/pkg/authn"
	"github.com/sourcenetwork/orbis-go/pkg/util/cleaner"
)

func init() {

	var key string

	client.RegisterFlagBinder(func(fs *pflag.FlagSet, namer naming.Namer) {
		fs.StringVar(&key, namer("JWT"), key, "JWT")
	})

	client.RegisterPreDialer(func(_ context.Context, opts *[]grpc.DialOption) error {

		if key != "" {
			cred, err := newJWTCred(key)
			if err != nil {
				return fmt.Errorf("jwt key: %v", err)
			}
			*opts = append(*opts, grpc.WithPerRPCCredentials(cred))
		}

		return nil
	})
}

// This is a temporary solution to get the JWT key into the gRPC client.
// Once we have the TLS setup, we can remove this, and take advantage of the oauth package.
// "google.golang.org/grpc/credentials/oauth"
func newJWTCred(token string) (credentials.PerRPCCredentials, error) {
	return jwtAccess{token}, nil
}

type jwtAccess struct {
	token string
}

func (j jwtAccess) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + j.token,
	}, nil
}

func (j jwtAccess) RequireTransportSecurity() bool {
	return false
}

func setupGRPCServer(cfg config.GRPC, errGrp *errgroup.Group, clnr *cleaner.Cleaner, a *app.App) error {

	// inject Metadata parser
	// NOTE(jsimnz): I dont really like this location for this, as I like to keep
	// all the dependency injection "magic" together, so its easier to track/follow.
	// But at least for now its going to live here since I wanted this particular
	// dependency to live close to the request/server intialization flow.
	do.ProvideValue[authn.RequestMetadataParser](a.Injector(), GRPCMetadataParser{})

	// Create a gRPC server object
	s := grpcserver.NewGRPCServer(cfg, a)

	reflection.Register(s)

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", cfg.GRPCURL)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	// Start the gRPC Server.
	errGrp.Go(func() error {
		log.Infof("Serving gRPC on %s", cfg.GRPCURL)
		return s.Serve(lis)
	})

	// Serve REST with gRPC-Gateway.
	gwServer, err := grpcserver.NewGRPCGatewayServer(cfg)
	if err != nil {
		return fmt.Errorf("create gRPC gateway server: %w", err)
	}
	errGrp.Go(func() error {
		log.Infof("Serving gRPC-Gateway on http://%s", cfg.RESTURL)
		return gwServer.ListenAndServe()
	})

	clnr.Regster(func() {

		log.Infof("Shutting down gRPC-Gateway server")
		err := gwServer.Shutdown(context.Background())
		if err != nil {
			log.Errorf("Shutting down gRPC-Gateway, server error: %s", err)
		}

		log.Infof("Shutting down gRPC server")
		s.GracefulStop()
	})

	return nil
}

var _ authn.RequestMetadataParser = (*GRPCMetadataParser)(nil)

type GRPCMetadataParser struct{}

func (GRPCMetadataParser) Parse(ctx context.Context) (authn.Metadata, bool) {
	return metadata.FromIncomingContext(ctx)
}

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/pkg/cleaner"

	"golang.org/x/sync/errgroup"
)

func setupServer(cfg config.Config) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clnr := cleaner.New()
	defer clnr.CleanUp()

	app, err := setupApp(ctx, cfg)
	if err != nil {
		return fmt.Errorf("setup app: %w", err)
	}

	// Errgroup tracks long running goroutines.
	// Any of the goroutines returns an error, the errgroup will return the error.
	errGrp, errGrpCtx := errgroup.WithContext(ctx)

	// Expose app services via gRPC server.
	err = setupGRPCServer(cfg.GRPC, errGrp, clnr, app)
	if err != nil {
		return fmt.Errorf("setup gRPC server: %w", err)
	}

	// Catch and handle signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		var sig os.Signal
		select {
		case sig = <-sigs:
			app.Logger().Infof("Received signal %q", sig)
		case <-errGrpCtx.Done():
			// At least 1 managed goroutines returns an error.
		}
		cancel()
		clnr.CleanUp()
	}()

	// Wait for all goroutines to finish.
	return errGrp.Wait()
}

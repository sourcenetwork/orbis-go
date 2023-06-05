package main

import (
	"time"

	hostv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/host/v1alpha1"
	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
	secretv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/secret/v1alpha1"
	transportv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/transport/v1alpha1"

	"github.com/NathanBaulch/protoc-gen-cobra/client"
	"github.com/spf13/cobra"
)

func main() {

	rootCmd := &cobra.Command{
		Use:          "orbisctl",
		Long:         "orbisctl controls the Orbis server",
		SilenceUsage: true,
	}

	opts := []client.Option{
		client.WithTimeout(1 * time.Second),
	}

	// Setup client commands for the Orbis client.
	rootCmd.AddCommand(
		hostv1alpha1.HostServiceClientCommand(opts...),
		ringv1alpha1.RingServiceClientCommand(opts...),
		secretv1alpha1.SecretServiceClientCommand(opts...),
		transportv1alpha1.TransportServiceClientCommand(opts...),
	)

	rootCmd.Execute() // nolint:errcheck
}

package main

import (
	"time"

	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/ring/v1alpha1"
	secretv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/secret/v1alpha1"
	transportv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/transport/v1alpha1"

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
		ringv1alpha1.RingServiceClientCommand(opts...),
		secretv1alpha1.SecretServiceClientCommand(opts...),
		transportv1alpha1.TransportServiceClientCommand(opts...),
	)

	rootCmd.Execute()
}

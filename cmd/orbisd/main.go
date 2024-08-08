package main

import (
	"time"

	"github.com/sourcenetwork/orbis-go/adapter/cobracli"
	oclient "github.com/sourcenetwork/orbis-go/adapter/cobracli/client"
	"github.com/sourcenetwork/orbis-go/adapter/cobracli/keys"

	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
	transportv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/transport/v1alpha1"
	utilityv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/utility/v1alpha1"

	"github.com/NathanBaulch/protoc-gen-cobra/client"
	logging "github.com/ipfs/go-log"
	"github.com/spf13/cobra"
)

var log = logging.Logger("orbis/orbisd")

func main() {

	rootCmd := &cobra.Command{
		Use:          "orbisd",
		Long:         "Orbis is a hybrid secrets management engine designed as a decentralized custodial system.",
		SilenceUsage: true,
	}

	// Setup the start command for the Orbis server.
	startCmd, err := cobracli.StartCmd(setupServer)
	if err != nil {
		log.Fatalf("Start command: %s", err)
	}

	// Setup client commands for the Orbis client.
	rootCmd.AddCommand(oclient.ClientCmd())

	// Setup server start commands
	rootCmd.AddCommand(
		startCmd,
		keys.KeyCmd(),
	)

	opts := []client.Option{
		client.WithTimeout(1 * time.Second),
		client.WithEnvVars("orbis"),
	}

	rootCmd.AddCommand(
		utilityv1alpha1.UtilityServiceClientCommand(opts...),
		ringv1alpha1.RingServiceClientCommand(opts...),
		transportv1alpha1.TransportServiceClientCommand(opts...),
	)

	rootCmd.Execute() // nolint:errcheck
}

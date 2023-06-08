package main

import (
	"github.com/sourcenetwork/orbis-go/adapter/cobracli"

	logging "github.com/ipfs/go-log"
	"github.com/spf13/cobra"
)

var log = logging.Logger("orbis/orbisd")

func main() {

	logging.SetAllLoggers(logging.LevelDPanic)
	logging.SetLogLevelRegex("orbis.*", "info")

	err := logging.SetLogLevelRegex("dht/.*", "error")
	if err != nil {
		log.Fatalf("Set log level: %s", err)
	}

	err = logging.SetLogLevelRegex("orbis/transport/.*", "error")
	if err != nil {
		log.Fatalf("Set log level: %s", err)
	}

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
	rootCmd.AddCommand(startCmd)

	rootCmd.Execute() // nolint:errcheck
}

package main

import (
	"log"

	"github.com/sourcenetwork/orbis-go/adapter/cobracli"

	"github.com/spf13/cobra"
)

func main() {

	rootCmd := &cobra.Command{
		Use:          "orbisd",
		Long:         "Orbis is a hybrid secrets management engine designed as a decentralized custodial system.",
		SilenceUsage: true,
	}

	// Setup the start command for the Orbis server.
	startCmd, err := cobracli.StartCmd(setupServer)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// Setup client commands for the Orbis client.
	rootCmd.AddCommand(startCmd)

	rootCmd.Execute()
}

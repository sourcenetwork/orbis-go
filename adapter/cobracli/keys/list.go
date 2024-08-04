package keys

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ListCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all keys",
		RunE:  runListCmd,
	}

	return cmd
}

func runListCmd(cmd *cobra.Command, _ []string) error {
	fmt.Println("hi from keys list")
	return nil
}

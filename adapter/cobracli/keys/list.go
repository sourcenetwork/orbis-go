package keys

import (
	"fmt"

	"github.com/sourcenetwork/orbis-go/adapter/cobracli"

	"github.com/spf13/cobra"
)

func ListCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all keys",
		RunE:  runListCmd(cfg),
	}

	return cmd
}

func runListCmd(cfg *Config) func(cmd *cobra.Command, _ []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		ctx, ok := cobracli.FromContext(cmd.Context())
		if !ok {
			return fmt.Errorf("invalid keyring context")
		}
		kr := ctx.Keyring()

		keys, err := kr.List()
		if err != nil {
			return fmt.Errorf("getting keys: %w", err)
		}

		p := ctx.Output()

		outset := make([]keyOutput, len(keys))
		for i, keyInfo := range keys {
			out, err := keyToOutput(keyInfo.Name, keyInfo.Key)
			if err != nil {
				return err
			}
			outset[i] = out
		}
		p.Print(outset)

		return nil
	}
}

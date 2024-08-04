package keys

import (
	"github.com/sourcenetwork/orbis-go/pkg/util/flag"
	"github.com/spf13/cobra"
)

func Commands() *cobra.Command {
	cfg := DefaultConfig
	var cmd *cobra.Command // separate variable defition is required!
	cmd = &cobra.Command{
		Use:   "keys",
		Short: "Manage your keys",
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			if cfg.UseEnvVars {
				err := flag.SetFlagsFromEnv(cmd.PersistentFlags(), false, cfg.EnvVarNamer, cfg.EnvVarPrefix, "keys")
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	cfg.BindFlags(cmd.PersistentFlags())
	cmd.AddCommand()

	return cmd
}

package client

import (
	"os"

	"github.com/segmentio/cli"
	"github.com/sourcenetwork/orbis-go/adapter/cobracli"
	"github.com/sourcenetwork/orbis-go/adapter/cobracli/keys"
	"github.com/sourcenetwork/orbis-go/pkg/keyring"
	"github.com/sourcenetwork/orbis-go/pkg/util/flag"
	"github.com/spf13/cobra"
)

func ClientCmd() *cobra.Command {
	cfg := DefaultConfig
	var cmd *cobra.Command // separate variable defition is required!
	cmd = &cobra.Command{
		Use:   "client",
		Short: "Orbis client",
		Long: `Stateful Orbis client that simplifies the end-to-end flow 
for interacting with an Orbis Ring.`,
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			if cfg.UseEnvVars {
				err := flag.SetFlagsFromEnv(cmd.PersistentFlags(), false, cfg.EnvVarNamer, cfg.EnvVarPrefix, "client")
				if err != nil {
					return err
				}
			}

			kr, err := keys.KeyringFromConfig(cfg.Keyring)
			if err != nil {
				return err
			}
			c.SetContext(keyring.WithKeyring(c.Context(), kr))

			pf, err := cli.Format(cfg.Output, os.Stdout)
			if err != nil {
				return err
			}
			c.SetContext(cobracli.WithContext(c.Context(), kr, pf))
			return nil
		},
	}
	cfg.BindFlags(cmd.PersistentFlags())
	cmd.AddCommand(
		PolicyCmd(cfg),
		GetSecretClientCmd(cfg),
		PutSecretClientCmd(cfg),
	)

	return cmd
}

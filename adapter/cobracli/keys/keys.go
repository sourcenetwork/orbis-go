package keys

import (
	"fmt"

	"github.com/sourcenetwork/orbis-go/pkg/keyring"
	"github.com/sourcenetwork/orbis-go/pkg/util/flag"

	"github.com/spf13/cobra"
)

func KeyCmd() *cobra.Command {
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

			kr, err := KeyringFromConfig(cfg)
			if err != nil {
				return err
			}
			c.SetContext(keyring.WithKeyring(c.Context(), kr))

			return nil
		},
	}
	cfg.BindFlags(cmd.PersistentFlags())
	cmd.AddCommand(
		ListCmd(cfg),
	)

	return cmd
}

func KeyringFromConfig(cfg *Config) (keyring.Keyring, error) {
	switch cfg.KeyringBackend {
	case "file":
		return keyring.New("file", cfg.KeyringPath)
	case "os":
		return keyring.New("os", cfg.KeyringPath, cfg.KeyringService)
	case "test":
		return keyring.New("test", cfg.KeyringPath)
	}

	return nil, fmt.Errorf("invalid keyring backend")
}

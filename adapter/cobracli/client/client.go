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

/*
# Env Vars
export ORBIS_CLIENT_FROM="alice"
export ORBIS_CLIENT_SERVER_ADDR=":8081"
export ORBIS_CLIENT_AUTHZ_ADDR=":8080"
export ORBIS_CLIENT_RING_ID="zQ123"

# Create Policy
orbisd client policy create -f policy.yaml => policy-id=0x123
orbisd client policy describe 0x123 => policy-data=...

# Create Secret (managed authorization)
orbisd client put "mysecret" --authz managed --policy 0x123 --resource secret --permission read
orbisd client get ABC123

# Create Secret (unmanaged authorization)
orbisd client policy register 0x123 secret mysecret
orbisd client put "mysecret" --authz unmanaged --permission "0x123/secret:mysecret#read"
orbisd client get ABC123

# Add Bob as a reader
orbisd client policy set 0x123 secret mysecret collaborator did:key:bob
orbisd client policy check 0x123 did:key:bob secret:mysecret#read => valid=true/false

*/

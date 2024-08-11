package keys

import (
	"fmt"
	"os"

	scrypto "github.com/TBD54566975/ssi-sdk/crypto"
	"github.com/TBD54566975/ssi-sdk/did/key"
	"github.com/segmentio/cli"
	"github.com/sourcenetwork/orbis-go/adapter/cobracli"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
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

			p, err := cli.Format(cfg.Output, os.Stdout)
			if err != nil {
				return err
			}
			// we flush the cli output in the PersistentPostRun

			c.SetContext(cobracli.WithContext(c.Context(), kr, p))

			return nil
		},
		PersistentPostRun: func(c *cobra.Command, args []string) {
			clientCtx, ok := cobracli.FromContext(c.Context())
			if !ok {
				return
			}
			clientCtx.Output().Flush()
		},
	}
	cfg.BindFlags(cmd.PersistentFlags())
	// output needs to be handled separetly because its duplicated on client subcommand
	cmd.PersistentFlags().StringVarP(&cfg.Output, namer("Output"), "", cfg.Output, "output format (text|json|yaml)")
	cmd.AddCommand(
		ListCmd(cfg),
		AddCmd(cfg),
		ShowCmd(cfg),
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

type keyOutput struct {
	Name   string `json:"name"`
	DID    string `json:"did"`
	Type   string `json:"type"`
	PubKey struct {
		Type string `json:"type"`
		Key  string `json:"key"`
	} `json:"pubkey"`
}

func cryptoKeyTypeToDID(kt crypto.KeyType) (scrypto.KeyType, error) {
	switch kt {
	case crypto.Ed25519:
		return scrypto.Ed25519, nil
	case crypto.Secp256k1:
		return scrypto.SECP256k1, nil
	}

	return "", fmt.Errorf("invalid key type")
}

func keyToOutput(name string, k crypto.Key) (keyOutput, error) {
	var out keyOutput
	out.Name = name
	if crypto.IsAsymmetric(k) {
		out.Type = "asymmetric (keypair)"
		pubkey, err := crypto.GetPublic(k)
		if err != nil {
			return keyOutput{}, err
		}
		out.PubKey.Type = pubkey.Type().String()
		out.PubKey.Key = pubkey.String()

		raw, err := pubkey.Raw()
		if err != nil {
			return keyOutput{}, err
		}
		didKeyType, err := cryptoKeyTypeToDID(pubkey.Type())
		if err != nil {
			return keyOutput{}, err
		}
		did, err := key.CreateDIDKey(didKeyType, raw)
		if err != nil {
			return keyOutput{}, err
		}

		out.DID = did.String()
	} else {
		out.Type = "symmetric (encryption key)"
	}

	return out, nil
}

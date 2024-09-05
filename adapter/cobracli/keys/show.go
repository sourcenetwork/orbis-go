package keys

import (
	"encoding/json"
	"fmt"

	"github.com/sourcenetwork/orbis-go/adapter/cobracli"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/spf13/cobra"
)

const (
	flagShowDid    = "did"
	flagShowPubkey = "pubkey"
)

func ShowCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show and existing key info",
		Args:  cobra.ExactArgs(1),
		RunE:  runShowCmd(cfg),
	}

	cmd.PersistentFlags().BoolP(flagShowDid, "d", false, "Output the DID identifier only")
	cmd.PersistentFlags().BoolP(flagShowPubkey, "p", false, "Output the public key only")

	return cmd
}

func runShowCmd(cfg *Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx, ok := cobracli.FromContext(cmd.Context())
		if !ok {
			return fmt.Errorf("invalid client context")
		}
		name := args[0]
		key, err := ctx.Keyring().Get(name)
		if err != nil {
			return fmt.Errorf("getting key: %w", err)
		}

		showDID, err := cmd.Flags().GetBool(flagShowDid)
		if err != nil {
			return err
		}
		showPubKey, err := cmd.Flags().GetBool(flagShowPubkey)
		if err != nil {
			return err
		}

		if showDID && showPubKey && !crypto.IsAsymmetric(key) {
			return fmt.Errorf("can't get public key of a non asymmetric key")
		}

		keyOutput, err := keyToOutput(name, key)
		if err != nil {
			return err
		}

		if showDID {
			fmt.Println(keyOutput.DID)
		} else if showPubKey {
			buf, err := json.Marshal(keyOutput.PubKey)
			if err != nil {
				return err
			}
			fmt.Println(string(buf))
		} else {
			ctx.Output().Print(keyOutput)
		}

		return nil
	}
}

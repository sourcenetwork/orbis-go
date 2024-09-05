package keys

import (
	"crypto/rand"
	"fmt"

	"github.com/sourcenetwork/orbis-go/adapter/cobracli"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/keyring"
	"github.com/spf13/cobra"
)

const (
	flagAddType   = "type"
	flagAddPubkey = "pubkey"
)

func AddCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new key",
		Args:  cobra.ExactArgs(1),
		RunE:  runAddCmd(cfg),
	}

	cmd.PersistentFlags().StringP(flagAddType, "t", "secp256k1", "Type of key to create (secp256k1|ed25519|symmetric)")
	cmd.PersistentFlags().BytesBase64(flagAddPubkey, nil, "import a public key (base64 encoded)")

	return cmd
}

func runAddCmd(cfg *Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx, ok := cobracli.FromContext(cmd.Context())
		if !ok {
			return fmt.Errorf("invalid keyring context")
		}
		kr := ctx.Keyring()

		name := args[0]

		flagType, err := cmd.Flags().GetString(flagAddType)
		if err != nil {
			return err
		}
		keyType, err := crypto.KeyPairTypeFromString(flagType)
		if err != nil {
			return err
		}

		var addedKey crypto.Key
		if cmd.Flags().Lookup(flagAddPubkey).Changed {
			pubBytes, err := cmd.Flags().GetBytesBase64(flagAddPubkey)
			if err != nil {
				return err
			}

			addedKey, err = importPublicKey(kr, keyType, name, pubBytes)
			if err != nil {
				return err
			}
		} else {
			addedKey, err = generateKey(kr, keyType, name)
			if err != nil {
				return err
			}
		}

		p := ctx.Output()

		outKey, err := keyToOutput(name, addedKey)
		if err != nil {
			return err
		}
		p.Print(outKey)
		return nil
	}
}

func generateKey(kr keyring.Keyring, keyType crypto.KeyType, name string) (crypto.Key, error) {
	if keyType == crypto.Octet {
		// symmetric
		panic("todo symmetric keys")
	}
	priv, _, err := crypto.GenerateKeyPair(keyType, rand.Reader)
	if err != nil {
		return nil, err
	}

	err = kr.Set(name, priv)
	return priv, err
}

func importPublicKey(kr keyring.Keyring, keyType crypto.KeyType, name string, pubkey []byte) (crypto.Key, error) {
	pk, err := crypto.PublicKeyFromBytes(keyType, pubkey)
	if err != nil {
		return nil, err
	}

	err = kr.Set(name, pk)
	return pk, err
}

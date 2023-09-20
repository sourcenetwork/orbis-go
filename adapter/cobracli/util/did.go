package util

import (
	"fmt"

	"github.com/TBD54566975/ssi-sdk/did/key"
	"github.com/spf13/cobra"
)

func createDIDCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create-did",
		Short: "create-did",
		RunE: func(cmd *cobra.Command, args []string) error {

			ste, err := suiteFromFlag(cmd.Flag("suite"))
			if err != nil {
				return fmt.Errorf("find suite: %w", err)
			}
			keyType, err := kyberSuiteToSSIKeyType(ste)
			if err != nil {
				return fmt.Errorf("convert suite to key type: %w", err)
			}

			rawPub, err := bytesFromFlag(cmd.Flag("pk"))
			if err != nil {
				return fmt.Errorf("encode public key to base64: %w", err)
			}

			didKey, err := key.CreateDIDKey(keyType, rawPub)
			if err != nil {
				return fmt.Errorf("create did key: %w", err)
			}

			return result{
				"did": didKey.String(),
			}.Output()
		},
	}

	cmd.Flags().String("suite", "ed25519", "Crypto suite. Must be one of ed25519, secp256k1, rsa, and ecdsa")
	cmd.Flags().String("pk", "", "Public key for the DID key (in base64)")

	return cmd
}

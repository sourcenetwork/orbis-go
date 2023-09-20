package util

import (
	"fmt"

	"github.com/sourcenetwork/orbis-go/pkg/pre/elgamal"

	"github.com/spf13/cobra"
)

func encryptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encrypt-secret",
		Short: "encrypt-secret",
		RunE: func(cmd *cobra.Command, args []string) error {

			scrt, err := bytesFromFlag(cmd.Flag("scrt"))
			if err != nil {
				return fmt.Errorf("decode secret: %w", err)
			}

			ste, err := suiteFromFlag(cmd.Flag("suite"))
			if err != nil {
				return fmt.Errorf("find suite: %w", err)
			}

			dkgPk, err := pointFromFlag(ste, cmd.Flag("dkg-pk"))
			if err != nil {
				return fmt.Errorf("unmarshal dkg-pk: %w", err)
			}

			encCmt, encScrt := elgamal.EncryptSecret(ste, dkgPk, scrt)

			b64EncCmt, err := pointToB64(encCmt)
			if err != nil {
				return fmt.Errorf("marshal enc-cmt: %w", err)
			}

			b64EncScrt := make([]string, len(encScrt))
			for i, encScrti := range encScrt {

				b64EncScrti, err := pointToB64(encScrti)
				if err != nil {
					return fmt.Errorf("marshal enc-scrt: %w", err)
				}
				b64EncScrt[i] = b64EncScrti
			}

			return result{
				"encCmt":  b64EncCmt,
				"encScrt": b64EncScrt,
			}.Output()
		},
	}

	cmd.Flags().String("suite", "ed25519", "Crypto suite. Must be one of ed25519, secp256k1, rsa, and ecdsa")
	cmd.Flags().String("dkg-pk", "", "DKG ring's shared public key (in base64)")
	cmd.Flags().String("scrt", "", "Secret to encrypt (in base64))")

	return cmd
}

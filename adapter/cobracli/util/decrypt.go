package util

import (
	"encoding/base64"
	"fmt"

	"github.com/sourcenetwork/orbis-go/pkg/pre/elgamal"

	"github.com/spf13/cobra"
	"go.dedis.ch/kyber/v3"
)

func decryptCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "decrypt-secret",
		Short: "decrypt-secret",
		RunE: func(cmd *cobra.Command, args []string) error {

			ste, err := suiteFromFlag(cmd.Flag("suite"))
			if err != nil {
				return fmt.Errorf("find suite: %w", err)
			}

			dkgPk, err := pointFromFlag(ste, cmd.Flag("dkg-pk"))
			if err != nil {
				return fmt.Errorf("unmarshal dkg-pk: %w", err)
			}

			encScrt, err := pointFromFlag(ste, cmd.Flag("enc-scrt"))
			if err != nil {
				return fmt.Errorf("unmarshal enc-scrt: %w", err)
			}

			xncCmt, err := pointFromFlag(ste, cmd.Flag("xnc-cmt"))
			if err != nil {
				return fmt.Errorf("unmarshal xnc-cmt: %w", err)
			}

			encScrts := []kyber.Point{encScrt}

			rdrSk, err := skFromFlag(ste, cmd.Flag("rdr-sk"))
			if err != nil {
				return fmt.Errorf("unmarshal rdr-sk: %w", err)
			}

			scrt, err := elgamal.DecryptSecret(ste, encScrts, dkgPk, xncCmt, rdrSk)
			if err != nil {
				return fmt.Errorf("decrypt secret: %w", err)
			}

			return result{
				"scrt": base64.StdEncoding.EncodeToString(scrt),
			}.Output()
		},
	}

	cmd.Flags().String("suite", "ed25519", "Crypto suite. Must be one of ed25519, secp256k1, rsa, and ecdsa")
	cmd.Flags().String("enc-scrt", "", "Encrypted secret (in base64)")
	cmd.Flags().String("xnc-cmt", "", "Reencrypted commitment (in base64)")
	cmd.Flags().String("rdr-sk", "", "Secret key to decrypt secret (in base64)")
	cmd.Flags().String("dkg-pk", "", "DKG ring's shared public key (in base64)")

	return cmd
}

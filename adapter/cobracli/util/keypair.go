package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"

	"github.com/spf13/cobra"
)

type zeroReader struct {
}

func (zeroReader) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p); i++ {
		p[i] = 0
	}

	return len(p), nil
}

func createKeypair() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create-keypair",
		Short: "create-keypair",
		RunE: func(cmd *cobra.Command, args []string) error {

			ste, err := suiteFromFlag(cmd.Flag("suite"))
			if err != nil {
				return fmt.Errorf("find suite: %w", err)
			}

			privateKey, publicKey, err := crypto.GenerateKeyPair(ste, rand.Reader)
			if err != nil {
				return fmt.Errorf("generate key pair: %w", err)
			}

			b64PublicKey, err := pointToB64(publicKey.Point())
			if err != nil {
				return fmt.Errorf("encode public key to base64: %w", err)
			}

			rawPrivateKey, err := privateKey.Raw()
			if err != nil {
				return fmt.Errorf("marshal private key: %w", err)
			}

			b64PrivateKey := base64.StdEncoding.EncodeToString(rawPrivateKey)
			if err != nil {
				return fmt.Errorf("encode private key to base64: %w", err)
			}

			return result{
				"privateKey": b64PrivateKey,
				"publicKey":  b64PublicKey,
			}.Output()
		},
	}

	cmd.Flags().String("suite", "ed25519", "Crypto suite. Must be one of ed25519, secp256k1, rsa, and ecdsa")
	cmd.Flags().String("rand", "rand", "random reader")

	return cmd
}

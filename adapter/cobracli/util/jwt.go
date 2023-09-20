package util

import (
	"crypto/ed25519"
	"fmt"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/spf13/cobra"
)

func createJWTCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-jwt",
		Short: "create-jwt",
		RunE: func(cmd *cobra.Command, args []string) error {

			kid := cmd.Flag("kid").Value.String()

			claims, err := claimsFromFlag(cmd.Flag("claims"))
			if err != nil {
				return fmt.Errorf("parse claims: %w", err)
			}

			rawSk, err := bytesFromFlag(cmd.Flag("sk"))
			if err != nil {
				return fmt.Errorf("decode: %w", err)
			}

			sk := ed25519.PrivateKey(rawSk)

			opts := new(jose.SignerOptions)
			opts.WithHeader(jose.HeaderKey("kid"), kid)
			signer, err := jose.NewSigner(
				jose.SigningKey{
					Algorithm: jose.EdDSA,
					Key:       sk,
				},
				opts,
			)
			if err != nil {
				return fmt.Errorf("create signer: %w", err)
			}

			signedJWT, err := jwt.Signed(signer).Claims(claims).CompactSerialize()
			if err != nil {
				return fmt.Errorf("sign jwt: %w", err)
			}

			return result{
				"jwt": signedJWT,
			}.Output()
		},
	}

	cmd.Flags().String("kid", "", "Key ID for JWT header")
	cmd.Flags().String("suite", "ed25519", "Crypto suite. Must be one of ed25519, secp256k1, rsa, and ecdsa")
	cmd.Flags().String("claims", "", "DKG ring's shared public key (in base64)")
	cmd.Flags().String("sk", "", "Secret key to sign the token (in base64))")

	return cmd
}

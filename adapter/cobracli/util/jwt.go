package util

import (
	"crypto/ed25519"
	"fmt"
	"time"

	ssicrypto "github.com/TBD54566975/ssi-sdk/crypto"
	"github.com/TBD54566975/ssi-sdk/did/key"
	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/spf13/cobra"
)

var jwtExpiry = time.Second * 10

func createJWTCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-jwt",
		Short: "create-jwt",
		RunE: func(cmd *cobra.Command, args []string) error {

			// claims, err := claimsFromFlag(cmd.Flag("claims"))
			// if err != nil {
			// 	return fmt.Errorf("parse claims: %w", err)
			// }

			rawSk, err := bytesFromFlag(cmd.Flag("sk"))
			if err != nil {
				return fmt.Errorf("decode: %w", err)
			}

			sk := ed25519.PrivateKey(rawSk)
			pk := sk.Public()

			didKey, err := key.CreateDIDKey(ssicrypto.Ed25519, pk.(ed25519.PublicKey))
			if err != nil {
				return fmt.Errorf("did key: %w", err)
			}

			suffix, err := didKey.Suffix()
			if err != nil {
				return fmt.Errorf("did suffix: %w", err)
			}
			kid := fmt.Sprintf("%s#%s", didKey.String(), suffix)

			claims := jwt.Claims{
				Subject:  didKey.String(),
				Issuer:   didKey.String(),
				Expiry:   jwt.NewNumericDate(time.Now().Add(jwtExpiry)),
				IssuedAt: jwt.NewNumericDate(time.Now()),
				Audience: jwt.Audience{"orbis"},
			}

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

	cmd.Flags().String("suite", "ed25519", "Crypto suite. Must be one of ed25519, secp256k1, rsa, and ecdsa")
	cmd.Flags().String("sk", "", "Secret key to sign the token (in base64))")

	return cmd
}

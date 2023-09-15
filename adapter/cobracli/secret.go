package cobracli

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/pre/elgamal"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
)

func SecretCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "secret",
		Short: "secret",
	}

	cmd.AddCommand(encryptCmd())
	cmd.AddCommand(decryptCmd())
	cmd.AddCommand(keypairCmd())

	return cmd
}

func encryptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encrypt",
		Short: "encrypt",
		RunE: func(cmd *cobra.Command, args []string) error {

			b64Scrt := cmd.Flag("scrt").Value.String()
			scrt, err := base64.StdEncoding.DecodeString(b64Scrt)
			if err != nil {
				return fmt.Errorf("decode secret: %w", err)
			}

			ste, err := suiteFromFlag(cmd.Flag("suite"))
			if err != nil {
				return fmt.Errorf("find suite: %w", err)
			}

			dkgPk, err := pointFromFlag(ste, cmd.Flag("dkg-pk"))
			if err != nil {
				return fmt.Errorf("unmarshal dkgPk: %w", err)
			}

			encCmt, encScrt := elgamal.EncryptSecret(ste, dkgPk, scrt)

			b64EncCmt, err := pointToB64(encCmt)
			if err != nil {
				return fmt.Errorf("marshal encCmt: %w", err)
			}

			b64EncScrt := make([]string, len(encScrt))
			for i, encScrti := range encScrt {

				b64EncScrti, err := pointToB64(encScrti)
				if err != nil {
					return fmt.Errorf("marshal encScrt: %w", err)
				}
				b64EncScrt[i] = b64EncScrti
			}

			result := struct {
				EncCmt  string   `json:"encCmt"`
				EncScrt []string `json:"encScrt"`
			}{
				EncCmt:  b64EncCmt,
				EncScrt: b64EncScrt,
			}

			j, err := json.Marshal(result)
			if err != nil {
				return fmt.Errorf("marshal json: %w", err)
			}

			fmt.Printf("%s\n", j)

			return nil
		},
	}

	cmd.Flags().String("suite", "ed25519", "Crypto suite. Must be one of ed25519, secp256k1, rsa, and ecdsa")
	cmd.Flags().String("dkg-pk", "", "DKG ring's shared public key (in base64)")
	cmd.Flags().String("scrt", "", "Secret to encrypt (in base64))")

	return cmd
}

func decryptCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "decrypt",
		Short: "decrypt",
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

			rdrSk, err := scalarFromFlag(ste, cmd.Flag("rdr-sk"))
			if err != nil {
				return fmt.Errorf("unmarshal rdr-sk: %w", err)
			}

			scrt, err := elgamal.DecryptSecret(ste, encScrts, dkgPk, xncCmt, rdrSk)
			if err != nil {
				return fmt.Errorf("decrypt secret: %w", err)
			}

			jsonScrt := struct {
				Secret string `json:"scrt"`
			}{
				Secret: base64.StdEncoding.EncodeToString(scrt),
			}

			j, err := json.Marshal(jsonScrt)
			if err != nil {
				return fmt.Errorf("marshal json: %w", err)
			}

			fmt.Printf("%s\n", j)

			return nil
		},
	}

	cmd.Flags().String("suite", "ed25519", "Crypto suite. Must be one of ed25519, secp256k1, rsa, and ecdsa")
	cmd.Flags().String("enc-scrt", "", "Encrypted secret (in base64)")
	cmd.Flags().String("xnc-cmt", "", "Reencrypted commitment (in base64)")
	cmd.Flags().String("rdr-sk", "", "Secret key to decrypt secret (in base64)")
	cmd.Flags().String("dkg-pk", "", "DKG ring's shared public key (in base64)")

	return cmd
}

func keypairCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "keypair",
		Short: "keypair",
		RunE: func(cmd *cobra.Command, args []string) error {

			ste, err := suiteFromFlag(cmd.Flag("suite"))
			if err != nil {
				return fmt.Errorf("find suite: %w", err)
			}

			privateKey, publicKey, err := crypto.GenerateKeyPair(ste)
			if err != nil {
				return fmt.Errorf("generate key pair: %w", err)
			}

			b64PublicKey, err := pointToB64(publicKey.Point())
			if err != nil {
				return fmt.Errorf("marshal public key: %w", err)
			}

			b64PrivateKey, err := scalarToB64(privateKey.Scalar())
			if err != nil {
				return fmt.Errorf("marshal private key: %w", err)
			}

			keyPair := struct {
				PrivateKey string `json:"privateKey"`
				PublicKey  string `json:"publicKey"`
			}{
				PrivateKey: b64PrivateKey,
				PublicKey:  b64PublicKey,
			}

			j, err := json.Marshal(keyPair)
			if err != nil {
				return fmt.Errorf("marshal json: %w", err)
			}

			fmt.Printf("%s\n", j)
			return nil
		},
	}

	cmd.Flags().String("suite", "ed25519", "Crypto suite. Must be one of ed25519, secp256k1, rsa, and ecdsa")

	return cmd
}

func suiteFromFlag(flag *pflag.Flag) (suites.Suite, error) {

	keyType := flag.Value.String()
	ste, err := suites.Find(keyType)
	if err != nil {
		return nil, fmt.Errorf("find suite: %w", err)
	}

	return ste, nil
}

func pointFromFlag(ste suites.Suite, flag *pflag.Flag) (kyber.Point, error) {

	b64 := flag.Value.String()
	p := ste.Point()

	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return p, fmt.Errorf("decode: %w", err)
	}

	err = p.UnmarshalBinary(raw)
	if err != nil {
		return p, fmt.Errorf("unmarshal: %w", err)
	}

	return p, nil
}

func pointToB64(p kyber.Point) (string, error) {

	raw, err := p.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	b64 := base64.StdEncoding.EncodeToString(raw)

	return b64, nil
}

func scalarFromFlag(ste suites.Suite, flag *pflag.Flag) (kyber.Scalar, error) {

	b64 := flag.Value.String()
	s := ste.Scalar()

	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return s, fmt.Errorf("decode: %w", err)
	}

	err = s.UnmarshalBinary(raw)
	if err != nil {
		return s, fmt.Errorf("unmarshal: %w", err)
	}

	return s, nil
}

func scalarToB64(s kyber.Scalar) (string, error) {

	raw, err := s.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	b64 := base64.StdEncoding.EncodeToString(raw)

	return b64, nil
}

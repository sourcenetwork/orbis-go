package keys

import (
	"fmt"
	"os"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/keyring"

	scrypto "github.com/TBD54566975/ssi-sdk/crypto"
	"github.com/TBD54566975/ssi-sdk/did/key"
	"github.com/segmentio/cli"
	"github.com/spf13/cobra"
)

func ListCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all keys",
		RunE:  runListCmd,
	}

	return cmd
}

type keyOutput struct {
	Name   string
	DID    string
	Type   string
	PubKey struct {
		Type string
		Key  string
	}
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
		var pubkey crypto.PublicKey
		switch kt := k.(type) {
		case crypto.PublicKey:
			pubkey = kt
		case crypto.PrivateKey:
			pubkey = kt.GetPublic()
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

func runListCmd(cmd *cobra.Command, _ []string) error {
	kr, ok := keyring.FromContext(cmd.Context())
	if !ok {
		return fmt.Errorf("invalid keyring context")
	}

	keys, err := kr.List()
	if err != nil {
		return err
	}

	p, err := cli.Format("yaml", os.Stdout)
	if err != nil {
		return err
	}
	defer p.Flush()

	for _, keyInfo := range keys {
		out, err := keyToOutput(keyInfo.Name, keyInfo.Key)
		if err != nil {
			return err
		}
		p.Print(out)
	}

	return nil
}

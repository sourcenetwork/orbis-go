package client

import (
	b64 "encoding/base64"
	"fmt"
	"time"

	"github.com/TBD54566975/ssi-sdk/did/key"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/sourcenetwork/orbis-go/adapter/cobracli"
	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/pre/elgamal"
	"github.com/sourcenetwork/orbis-go/pkg/types"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/spf13/cobra"
)

var (
	flagPutAuthz = "authz"
)

func GetSecretClientCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get secret-id",
		Short: "Get a secret",
		Args:  cobra.ExactArgs(1),
		RunE:  runGetSecretClientCmd(cfg),
	}
	return cmd
}

func runGetSecretClientCmd(cfg *Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		secretId := args[0]
		return RoundTrip(cmd.Context(), cfg, cfg.ServerAddr, func(conn grpc.ClientConnInterface) error {
			clientCtx, ok := cobracli.FromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("invalid client context")
			}

			ringClient := ringv1alpha1.NewRingServiceClient(conn)
			from := cfg.From
			if from == "" {
				return fmt.Errorf("missing required flag 'from'")
			}
			ringId := cfg.RingId
			if ringId == "" {
				return fmt.Errorf("missing requied flag 'ring-id'")
			}

			ringPkResp, err := ringClient.PublicKey(cmd.Context(), &ringv1alpha1.PublicKeyRequest{
				Id: ringId,
			})
			if err != nil {
				return fmt.Errorf("getting ring public-key: %w", err)
			}

			fromKey, err := clientCtx.Keyring().Get(from)
			if err != nil {
				return fmt.Errorf("missing from key: %w", err)
			}

			fromSk, err := crypto.GetPrivate(fromKey)
			if err != nil {
				return fmt.Errorf("invalid from key: %w", err)
			}

			pkProto, err := crypto.PublicKeyToProto(fromSk.GetPublic())
			if err != nil {
				return err
			}

			token, err := createJWT(fromSk)
			if err != nil {
				return fmt.Errorf("create jwt: %w", err)
			}
			md := metadata.New(map[string]string{
				"authorization": "Bearer " + token,
			})
			ctx := metadata.NewOutgoingContext(cmd.Context(), md)

			reencryptReq := &ringv1alpha1.ReencryptSecretRequest{
				RingId:   ringId,
				SecretId: secretId,
				RdrPk:    pkProto,
			}
			resp, err := ringClient.ReencryptSecret(ctx, reencryptReq)
			if err != nil {
				return fmt.Errorf("reencrypt secret: %w", err)
			}

			secret, err := decryptSecret(fromSk.Type(), fromSk, ringPkResp.PublicKey.Data, resp.XncCmt, resp.EncScrt)
			if err != nil {
				return fmt.Errorf("decrypt secret: %w", err)
			}

			fmt.Println(string(secret))
			return nil
		})
	}
}

func PutSecretClientCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "put secret",
		Short: "Store a secret",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			managed, err := cmd.Flags().GetBool("managed")
			if err != nil {
				return err
			}

			if managed {
				return doManagedSecretPut(cfg, cmd, args)
			}
			return doUnmanagedSecretPut(cfg, cmd, args)
		},
	}

	cmd.PersistentFlags().Bool("managed", false, "Managed authorization to automatically create permissions")
	cmd.PersistentFlags().StringP("permission", "p", "", "Permission to access the secret")
	cmd.PersistentFlags().BoolP("base64", "b", false, "Secret is base64 encoded")
	cmd.PersistentFlags().StringP("policy", "", "", "Policy for managed authorization")
	cmd.PersistentFlags().StringP("resource", "r", "", "Resource for managed authorization")
	return cmd
}

func doManagedSecretPut(cfg *Config, cmd *cobra.Command, args []string) error {
	secret := []byte(args[0])
	isBase64, err := cmd.Flags().GetBool("base64")
	if err != nil {
		return err
	}

	if len(secret) == 0 {
		return fmt.Errorf("secret zero lenght")
	}

	if isBase64 {
		var secretBuf []byte
		b64.StdEncoding.Decode(secretBuf, secret)
		secret = secretBuf
	}

	policyID, err := cmd.Flags().GetString("policy")
	if err != nil {
		return err
	}
	if policyID == "" {
		return fmt.Errorf("managed authorization requires 'policy' flag")
	}
	permission, err := cmd.Flags().GetString("permission")
	if err != nil {
		return err
	}
	if permission == "" {
		return fmt.Errorf("managed authorization requires 'permission' flag")
	}
	resource, err := cmd.Flags().GetString("resource")
	if err != nil {
		return err
	}
	if resource == "" {
		return fmt.Errorf("managed authorization requires 'resource' flag")
	}

	return RoundTrip(cmd.Context(), cfg, cfg.ServerAddr, func(conn grpc.ClientConnInterface) error {
		ringClient := ringv1alpha1.NewRingServiceClient(conn)

		pkReq := &ringv1alpha1.PublicKeyRequest{
			Id: cfg.RingId,
		}
		pkResp, err := ringClient.PublicKey(cmd.Context(), pkReq)
		if err != nil {
			return fmt.Errorf("get ring public key: %w", err)
		}

		encCmt, encScrt, err := encryptSecret(
			pkResp.PublicKey.Type.String(),
			pkResp.PublicKey.Data,
			secret)
		if err != nil {
			return fmt.Errorf("encrypt secret: %w", err)
		}

		// we can omit the authz context since the secret ID
		// doesn't depend on it
		secretType := types.NewSecret(encCmt, encScrt, "")
		sid, err := secretType.ID()
		if err != nil {
			return fmt.Errorf("secret ID: %w", err)
		}

		ctx, ok := cobracli.FromContext(cmd.Context())
		if !ok {
			return fmt.Errorf("couldn't get client context")
		}
		did, err := fromDID(ctx.Keyring(), cfg.From)
		if err != nil {
			return fmt.Errorf("getting key DID identifier: %w", err)
		}

		_, err = doRelationshipRequest(cmd.Context(), cfg, policyID, resource,
			string(sid), did, "owner")
		if err != nil {
			return fmt.Errorf("policy register object: %w", err)
		}

		// build authz context
		// <policyID>/resourceName:resourceID#permission
		authzCtx := fmt.Sprintf("%s/%s:%s#%s", policyID, resource, sid, permission)

		storeSecretReq := &ringv1alpha1.StoreSecretRequest{
			RingId: cfg.RingId,
			Secret: &ringv1alpha1.Secret{
				EncCmt:   encCmt,
				EncScrt:  encScrt,
				AuthzCtx: authzCtx,
			},
		}
		storeSecretResp, err := ringClient.StoreSecret(cmd.Context(), storeSecretReq)
		if err != nil {
			return fmt.Errorf("store secret: %w", err)
		}

		fmt.Println(storeSecretResp.SecretId)
		return nil
	})
}

func doUnmanagedSecretPut(cfg *Config, cmd *cobra.Command, args []string) error {
	secret := []byte(args[0])
	isBase64, err := cmd.Flags().GetBool("base64")
	if err != nil {
		return err
	}

	if len(secret) == 0 {
		return fmt.Errorf("secret zero lenght")
	}

	if isBase64 {
		var secretBuf []byte
		b64.StdEncoding.Decode(secretBuf, secret)
		secret = secretBuf
	}

	permission, err := cmd.Flags().GetString("permission")
	if err != nil {
		return err
	}

	return RoundTrip(cmd.Context(), cfg, cfg.ServerAddr, func(conn grpc.ClientConnInterface) error {
		ringClient := ringv1alpha1.NewRingServiceClient(conn)

		pkReq := &ringv1alpha1.PublicKeyRequest{
			Id: cfg.RingId,
		}
		pkResp, err := ringClient.PublicKey(cmd.Context(), pkReq)
		if err != nil {
			return fmt.Errorf("get ring public key: %w", err)
		}

		encCmt, encScrt, err := encryptSecret(
			pkResp.PublicKey.Type.String(),
			pkResp.PublicKey.Data,
			secret)
		if err != nil {
			return fmt.Errorf("encrypt secret: %w", err)
		}

		storeSecretReq := &ringv1alpha1.StoreSecretRequest{
			RingId: cfg.RingId,
			Secret: &ringv1alpha1.Secret{
				EncCmt:   encCmt,
				EncScrt:  encScrt,
				AuthzCtx: permission,
			},
		}
		storeSecretResp, err := ringClient.StoreSecret(cmd.Context(), storeSecretReq)
		if err != nil {
			return fmt.Errorf("store secret: %w", err)
		}

		fmt.Println(storeSecretResp.SecretId)
		return nil
	})
}

func encryptSecret(keyType string, ringPk []byte, secret []byte) ([]byte, [][]byte, error) {
	ste, err := crypto.FindSuite(keyType)
	if err != nil {
		return nil, nil, fmt.Errorf("unsupported key type: %w", err)
	}

	dkgPk := ste.Point()
	err = dkgPk.UnmarshalBinary(ringPk)
	if err != nil {
		return nil, nil, fmt.Errorf("unmarshal dkgPk: %w", err)
	}

	encCmt, encScrt := elgamal.EncryptSecret(ste, dkgPk, secret)

	rawEncCmt, err := encCmt.MarshalBinary()
	if err != nil {
		return nil, nil, fmt.Errorf("marshal encCmt: %w", err)
	}

	rawEncScrt := make([][]byte, len(encScrt))
	for i, encScrti := range encScrt {

		rawEncScrti, err := encScrti.MarshalBinary()
		if err != nil {
			return nil, nil, fmt.Errorf("marshal encScrt: %w", err)
		}
		rawEncScrt[i] = rawEncScrti
	}

	return rawEncCmt, rawEncScrt, nil
}

// func encryptSecret(keyType string, ringPk []byte, secret []byte) ([]byte, [][]byte, error) {
func decryptSecret(keyType crypto.KeyType, privKey crypto.PrivateKey, ringPk []byte, xncCmt []byte, encScrt [][]byte) ([]byte, error) {
	suite, err := crypto.SuiteForType(keyType)
	if err != nil {
		return nil, fmt.Errorf("find suite: %w", err)
	}

	dkgPk := suite.Point()
	err = dkgPk.UnmarshalBinary(ringPk)
	if err != nil {
		return nil, fmt.Errorf("unmarshal ring public-key: %w", err)
	}

	encScrts := make([]kyber.Point, len(encScrt))
	for i, rawEncScrt := range encScrt {
		encScrtPoint := suite.Point()
		err = encScrtPoint.UnmarshalBinary(rawEncScrt)
		if err != nil {
			return nil, fmt.Errorf("unmarshal encrypted secret: %w", err)
		}

		encScrts[i] = encScrtPoint
	}

	xncCmtPoint, err := unmarshalPoint(suite, xncCmt)
	if err != nil {
		return nil, fmt.Errorf("unmarshal commitment: %w", err)
	}

	privKeyScalar := privKey.Scalar()

	scrt, err := elgamal.DecryptSecret(suite, encScrts, dkgPk, xncCmtPoint, privKeyScalar)
	if err != nil {
		return nil, fmt.Errorf("decrypt secret: %w", err)
	}

	return scrt, nil
}

func unmarshalPoint(suite suites.Suite, buf []byte) (kyber.Point, error) {
	point := suite.Point()
	err := point.UnmarshalBinary(buf)
	return point, err
}

func createJWT(privKey crypto.PrivateKey) (string, error) {
	did, err := privKey.GetPublic().DID()
	if err != nil {
		return "", fmt.Errorf("get DID: %w", err)
	}
	// claims := jwt.Claims{
	// 	Audience: jwt.Audience{"orbis"},
	// 	Issuer:   did,
	// 	Subject:  did,
	// 	IssuedAt: jwt.NewNumericDate(time.Now()),
	// 	Expiry:   jwt.NewNumericDate(time.Now().Add(10 * time.Second)),
	// }
	tokenUnsigned, err := jwt.NewBuilder().
		Audience([]string{"orbis"}).
		Issuer(did).
		Subject(did).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(10 * time.Second)).Build()

	var alg jwa.KeyAlgorithm
	switch privKey.Type() {
	case crypto.Ed25519:
		alg = jwa.EdDSA
	case crypto.Secp256k1:
		alg = jwa.ES256K
	default:
		return "", fmt.Errorf("unknown key type: %w", err)
	}

	goKey, err := privKey.Std()
	if err != nil {
		return "", fmt.Errorf("getting private key: %w", err)
	}
	switch kt := goKey.(type) {
	case (*secp256k1.PrivateKey):
		goKey = kt.ToECDSA()
	}

	didKey := key.DIDKey(did)
	suffix, err := didKey.Suffix()
	if err != nil {
		return "", fmt.Errorf("did key suffix: %w", err)
	}
	hdr := jws.NewHeaders()
	if err := hdr.Set("kid", did+"#"+suffix); err != nil {
		return "", fmt.Errorf("jwt header KeyID: %w", err)
	}
	payload, err := jwt.Sign(tokenUnsigned,
		jwt.WithKey(alg, goKey, jws.WithProtectedHeaders(hdr)))
	if err != nil {
		return "", fmt.Errorf("jwt sign: %w", err)
	}

	return string(payload), nil
}

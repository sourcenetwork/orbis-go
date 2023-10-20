package grpcserver

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sourcenetwork/orbis-go/app"
	utilityv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/utility/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/pre/elgamal"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ssicrypto "github.com/TBD54566975/ssi-sdk/crypto"
	"github.com/TBD54566975/ssi-sdk/did/key"
	ic "github.com/libp2p/go-libp2p/core/crypto"
)

type utilService struct {
	utilityv1alpha1.UnimplementedUtilityServiceServer
	app *app.App
}

func newUtilService(app *app.App) *utilService {
	return &utilService{
		app: app,
	}
}

func (s *utilService) CreateDID(ctx context.Context, req *utilityv1alpha1.CreateDIDRequest) (*utilityv1alpha1.CreateDIDResponse, error) {

	if req.PublicKey == nil || req.KeyType == "" {
		return nil, status.Error(codes.InvalidArgument, "public key and key type are required")
	}

	ssiKeyType, err := kyberSuiteToSSIKeyType(req.KeyType)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "unsupported public key type")
	}

	didKey, err := key.CreateDIDKey(ssiKeyType, req.PublicKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "create did key")
	}

	fields := strings.SplitN(didKey.String(), ":", 3)

	resp := &utilityv1alpha1.CreateDIDResponse{
		Did:        didKey.String(),
		Scheme:     fields[0],
		Method:     fields[1],
		Identifier: fields[2],
	}

	return resp, nil
}

func (s *utilService) CreateJWT(ctx context.Context, req *utilityv1alpha1.CreateJWTRequest) (*utilityv1alpha1.CreateJWTResponse, error) {

	claims := jwt.Claims{}
	err := json.Unmarshal([]byte(req.Claims), &claims)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unmarshal claims: %s", err)
	}

	privateKey := ed25519.PrivateKey(req.PrivateKey)

	opts := new(jose.SignerOptions)
	opts.WithHeader(jose.HeaderKey("kid"), req.Kid)
	signer, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.EdDSA,
			Key:       privateKey,
		},
		opts,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create signer: %s", err)
	}

	signedJWT, err := jwt.Signed(signer).Claims(claims).CompactSerialize()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "sign jwt: %s", err)
	}

	resp := &utilityv1alpha1.CreateJWTResponse{
		Jwt: signedJWT,
	}

	return resp, nil
}

func (s *utilService) CreateKeypair(ctx context.Context, req *utilityv1alpha1.CreateKeypairRequest) (*utilityv1alpha1.CreateKeypairResponse, error) {

	ste, err := suites.Find(req.KeyType)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported key type: %s", err)
	}

	privateKey, publicKey, err := crypto.GenerateKeyPair(ste, rand.Reader)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate key pair: %s", err)
	}

	rawPrivateKey, err := privateKey.Raw()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "marshal private key: %s", err)
	}

	rawPublicKey, err := publicKey.Raw()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "marshal public key: %s", err)
	}

	resp := &utilityv1alpha1.CreateKeypairResponse{
		PrivateKey: rawPrivateKey,
		PublicKey:  rawPublicKey,
	}

	return resp, nil
}

func (s *utilService) EncryptSecret(ctx context.Context, req *utilityv1alpha1.EncryptSecretRequest) (*utilityv1alpha1.EncryptSecretResponse, error) {

	ste, err := suites.Find(req.KeyType)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported key type: %s", err)
	}

	dkgPk := ste.Point()
	err = dkgPk.UnmarshalBinary(req.DkgPk)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unmarshal dkgPk: %s", err)
	}

	encCmt, encScrt := elgamal.EncryptSecret(ste, dkgPk, req.Scrt)

	rawEncCmt, err := encCmt.MarshalBinary()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "marshal encCmt: %s", err)
	}

	rawEncScrt := make([][]byte, len(encScrt))
	for i, encScrti := range encScrt {

		rawEncScrti, err := encScrti.MarshalBinary()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "marshal encScrt: %s", err)
		}
		rawEncScrt[i] = rawEncScrti
	}

	resp := &utilityv1alpha1.EncryptSecretResponse{
		EncCmt:  rawEncCmt,
		EncScrt: rawEncScrt,
	}
	return resp, nil
}

func (s *utilService) DecryptSecret(ctx context.Context, req *utilityv1alpha1.DecryptSecretRequest) (*utilityv1alpha1.DecryptSecretResponse, error) {

	ste, err := suites.Find(req.KeyType)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported key type: %s", err)
	}

	dkgPk := ste.Point()
	err = dkgPk.UnmarshalBinary(req.DkgPk)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unmarshal dkgPk: %s", err)
	}

	var encScrts []kyber.Point
	for _, rawEncScrt := range req.EncScrt {
		encScrt := ste.Point()
		err = encScrt.UnmarshalBinary(rawEncScrt)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "unmarshal encScrt: %s", err)
		}
		encScrts = append(encScrts, encScrt)
	}

	xncCmt := ste.Point()
	err = xncCmt.UnmarshalBinary(req.XncCmt)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unmarshal xncCmt: %s", err)
	}

	icRdrSk, err := ic.UnmarshalEd25519PrivateKey(req.RdrSk)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unmarshal rdrSk: %s", err)
	}
	sk, err := crypto.PrivateKeyFromLibP2P(icRdrSk)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unmarshal rdrSk: %s", err)
	}
	rdrSk := sk.Scalar()

	scrt, err := elgamal.DecryptSecret(ste, encScrts, dkgPk, xncCmt, rdrSk)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "decrypt secret: %s", err)
	}

	resp := &utilityv1alpha1.DecryptSecretResponse{
		Scrt: scrt,
	}

	return resp, nil
}

func kyberSuiteToSSIKeyType(kt string) (ssicrypto.KeyType, error) {
	var keyType ssicrypto.KeyType

	switch kt {
	case "ed25519":
		keyType = ssicrypto.Ed25519
	case "secp256k1":
		keyType = ssicrypto.SECP256k1
	case "rsa":
		keyType = ssicrypto.RSA
	case "ecdsa":
		keyType = ssicrypto.SECP256k1ECDSA
	default:
		return keyType, fmt.Errorf("unsupported key type: %s", kt)
	}

	return keyType, nil
}

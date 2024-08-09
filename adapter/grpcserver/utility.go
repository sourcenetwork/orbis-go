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
	"github.com/sourcenetwork/orbis-go/pkg/cosmos"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/pre/elgamal"

	"github.com/sourcenetwork/sourcehub/x/acp/types"

	ssicrypto "github.com/TBD54566975/ssi-sdk/crypto"
	"github.com/TBD54566975/ssi-sdk/did/key"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/samber/do"
	"go.dedis.ch/kyber/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *utilService) CreateBech32Address(ctx context.Context, req *utilityv1alpha1.CreateBech32AddressRequest) (*utilityv1alpha1.CreateBech32AddressResponse, error) {
	if len(req.PublicKey) == 0 {
		return nil, status.Error(codes.InvalidArgument, "public key and key type are required")
	}

	if req.Prefix == "" {
		return nil, status.Error(codes.InvalidArgument, "bech32 prefix required")
	}

	addr, err := crypto.PublicKeyBytesToBech32(req.PublicKey)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bech32 encoding: %s", err)
	}

	return &utilityv1alpha1.CreateBech32AddressResponse{
		Bech32Address: addr,
	}, nil
}

func (s *utilService) CreateJWT(ctx context.Context, req *utilityv1alpha1.CreateJWTRequest) (*utilityv1alpha1.CreateJWTResponse, error) {

	claims := jwt.Claims{}
	err := json.Unmarshal([]byte(req.Claims), &claims)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unmarshal claims: %s", err)
	}

	var privateKey any
	var alg jose.SignatureAlgorithm
	switch strings.ToLower(req.KeyType) {
	case "ed25519":
		privateKey = ed25519.PrivateKey(req.PrivateKey)
		alg = jose.EdDSA
	case "secp256k1":
		privateKey = secp256k1.PrivKeyFromBytes(req.PrivateKey).ToECDSA()
		alg = jose.ES256K
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid key type '%s'", req.KeyType)
	}

	opts := new(jose.SignerOptions)
	opts.WithHeader(jose.HeaderKey("kid"), req.Kid)
	signer, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: alg,
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
	keyType, err := crypto.KeyPairTypeFromString(req.KeyType)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported key type: %s", err)
	}

	privateKey, publicKey, err := crypto.GenerateKeyPair(keyType, rand.Reader)
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

	ste, err := crypto.FindSuite(req.KeyType)
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

	ste, err := crypto.FindSuite(req.KeyType)
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

	kt, err := crypto.KeyPairTypeFromString(req.KeyType)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unknown key type: %s", err)
	}
	sk, err := crypto.PrivateKeyFromBytes(kt, req.RdrSk)
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

func (s *utilService) ACPCreatePolicy(ctx context.Context, req *utilityv1alpha1.ACPCreatePolicyRequest) (*utilityv1alpha1.ACPCreatePolicyResponse, error) {
	cc, err := do.Invoke[*cosmos.Client](s.app.Injector())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "missing cosmos client")
	}

	if req.PolicyYaml == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing policy ID")
	}
	if req.Creator == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing creator")
	}

	msgCreatePolicy := types.NewMsgCreatePolicyNow(req.Creator, req.PolicyYaml, types.PolicyMarshalingType_SHORT_YAML)
	acct, err := cc.AccountRegistry.GetByAddress(req.Creator)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't get account for creator %s: %s", req.Creator, err)
	}
	resp, err := cc.BroadcastTx(ctx, acct, msgCreatePolicy)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "broadcast tx: %s", err)
	}

	policyID := ""
	for _, evt := range resp.Events {
		if evt.Type == "sourcehub.acp.EventPolicyCreated" {
			for _, attr := range evt.Attributes {
				if attr.Key == "policy_id" {
					policyID = strings.ReplaceAll(attr.Value, "\"", "")
				}
			}
		}
	}

	return &utilityv1alpha1.ACPCreatePolicyResponse{
		PolicyId: policyID,
	}, nil
}

func (s *utilService) ACPRegisterObject(ctx context.Context, req *utilityv1alpha1.ACPRegisterObjectRequest) (*utilityv1alpha1.ACPRegisterObjectResponse, error) {
	cc, err := do.Invoke[*cosmos.Client](s.app.Injector())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "missing cosmos client")
	}

	if req.ObjectId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing object ID")
	}
	if req.ObjectResource == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing object resource")
	}
	if req.PolicyId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing policy ID")
	}
	if req.Creator == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing creator")
	}

	obj := &types.Object{
		Resource: req.ObjectResource,
		Id:       req.ObjectId,
	}
	registerMsg := types.NewMsgRegisterObjectNow(req.Creator, req.PolicyId, obj)
	// registerMsg := &types.MsgRegisterObject{
	// 	Creator:  req.Creator,
	// 	PolicyId: req.PolicyId,
	// 	Object: &types.Object{
	// 		Resource: req.ObjectResource,
	// 		Id:       req.ObjectId,
	// 	},
	// 	CreationTime: prototypes.TimestampNow(),
	// }

	acct, err := cc.AccountRegistry.GetByAddress(req.Creator)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't get account for creator %s: %s", req.Creator, err)
	}
	_, err = cc.BroadcastTx(ctx, acct, registerMsg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "broadcast tx: %s", err)
	}

	return &utilityv1alpha1.ACPRegisterObjectResponse{}, nil
}

func (s *utilService) ACPSetRelationship(ctx context.Context, req *utilityv1alpha1.ACPSetRelationshipRequest) (*utilityv1alpha1.ACPSetRelationshipResponse, error) {
	cc, err := do.Invoke[*cosmos.Client](s.app.Injector())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "missing cosmos client")
	}

	if req.ObjectId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing object ID")
	}
	if req.ObjectResource == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing object resource")
	}
	if req.Creator == "" {
		return nil, status.Error(codes.InvalidArgument, "missing creator")
	}
	if req.Actor == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing actor")
	}
	if req.Relation == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing relation")
	}
	if req.PolicyId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing policy ID")
	}

	relationship := types.NewActorRelationship(req.ObjectResource, req.ObjectId, req.Relation, req.Actor)
	setRelmsg := types.NewMsgSetRelationshipNow(
		req.Creator,
		req.PolicyId,
		relationship,
	)

	acct, err := cc.AccountRegistry.GetByAddress(req.Creator)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't get account for creator %s: %s", req.Creator, err)
	}
	_, err = cc.BroadcastTx(ctx, acct, setRelmsg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "broadcast tx: %s", err)
	}

	return &utilityv1alpha1.ACPSetRelationshipResponse{}, nil
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

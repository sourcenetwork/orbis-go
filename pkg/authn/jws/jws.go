package jws

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-jose/go-jose"
	"github.com/go-jose/go-jose/jwt"
	"golang.org/x/crypto/ed25519"

	"github.com/sourcenetwork/orbis-go/pkg/authn"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
)

const (
	OrbisJWSAudience = "orbis"
	tokenMetadataKey = "authorization"
)

var (
	verifyLeewayTime = 10 * time.Second
)

var _ authn.CredentialService = (*credentialSrv)(nil)

type credentialSrv struct {
	resolver       authn.KeyResolver
	metadataParser authn.RequestMetadataParser
}

func NewSelfSignedCredentialService(resolver authn.KeyResolver, metadataFn authn.RequestMetadataParser) authn.CredentialService {
	return credentialSrv{
		resolver:       resolver,
		metadataParser: metadataFn,
	}
}

func (c credentialSrv) GetAndVerifyRequestMetadata(ctx context.Context) (authn.SubjectInfo, error) {
	// grab the target expiration time at the start of the function
	// as there may be a non trivial amount of time that has passed
	// due parsing and resolving state by the time we actually want
	// to verify the token expiration
	expTime := time.Now()

	// parse content from request context
	md, ok := c.metadataParser.Parse(ctx)
	if !ok {
		return authn.SubjectInfo{}, fmt.Errorf("missing request metadata")
	}

	vals := md.Get(tokenMetadataKey)
	if len(vals) == 0 {
		return authn.SubjectInfo{}, fmt.Errorf("missing authorization token")
	}
	token := vals[0]
	jws, err := jose.ParseSigned(token)
	if err != nil {
		return authn.SubjectInfo{}, fmt.Errorf("parsing jws token: %w", err)
	}

	// this is likely impossible because ParseSigned above
	// will catch it, but just for sanity
	if len(jws.Signatures) == 0 {
		return authn.SubjectInfo{}, fmt.Errorf("missing jws signatures")
	}

	sig := jws.Signatures[0]
	if sig.Protected.KeyID == "" {
		return authn.SubjectInfo{}, fmt.Errorf("missing either JWK or KeyID")
	}

	// otherwise resolve the JWK from the KeyID
	kid := sig.Protected.KeyID
	userInfo, err := c.resolver.Resolve(ctx, kid)
	if err != nil {
		return authn.SubjectInfo{}, fmt.Errorf("resolving kid: %w", err)
	}
	key, err := crypto.PublicKeyToStdKey(userInfo.PubKey)
	if err != nil {
		return authn.SubjectInfo{}, fmt.Errorf("extracting JWK from resolved public key: %w", err)
	}

	// verify jws
	payload, err := jws.Verify(key)
	if err != nil {
		return authn.SubjectInfo{}, fmt.Errorf("verifying JWS: %w", err)
	}

	claims := claims{}
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return authn.SubjectInfo{}, fmt.Errorf("unmarshaling JWS payload: %w", err)
	}

	// verify claims
	expected := jwt.Expected{
		Audience: jwt.Audience{OrbisJWSAudience},
		Issuer:   userInfo.Subject,
		Subject:  userInfo.Subject,
		Time:     expTime,
	}

	err = claims.ValidateWithLeeway(expected, verifyLeewayTime)
	if err != nil {
		return authn.SubjectInfo{}, fmt.Errorf("JWS claim failed validation: %w", err)
	}

	return authn.SubjectInfo{
		Type:    "JWS",
		Subject: userInfo.Subject,
		PubKey:  userInfo.PubKey,
	}, nil
}

type dummyResolver struct{}

func (dummyResolver) Resolve(_ context.Context, _ string) (authn.SubjectInfo, error) {
	return authn.SubjectInfo{}, nil
}

// Converts a Public Key to a JWK
func JWKFromPublicKey(pk crypto.PublicKey) (*jose.JSONWebKey, error) {
	if pk == nil {
		return nil, fmt.Errorf("empty public key")
	}

	var key interface{}
	switch pk.Type() {
	case crypto.Ed25519:
		buf, err := pk.Raw()
		if err != nil {
			return nil, fmt.Errorf("extrating pubkey bytes: %w", err)
		}
		key = ed25519.PublicKey(buf)
	default:
		// invalid
		return nil, fmt.Errorf("invalid key type %s", pk.Type())
	}

	return &jose.JSONWebKey{
		Key: key,
	}, nil
}

type claims struct {
	jwt.Claims
}

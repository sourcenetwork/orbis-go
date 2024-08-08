package jws

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cometbft/cometbft/crypto/tmhash"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	logging "github.com/ipfs/go-log"
	"golang.org/x/crypto/ed25519"

	"github.com/sourcenetwork/orbis-go/pkg/authn"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
)

const (
	tokenPrefix = "Bearer "

	OrbisJWSAudience = "orbis"
	TokenMetadataKey = "authorization"
)

var (
	verifyLeewayTime = 10 * time.Second
)

var log = logging.Logger("orbis/authn/jws")

var _ authn.CredentialService = (*credentialSrv)(nil)

type credentialSrv struct {
	resolver       authn.KeyResolver
	metadataParser authn.RequestMetadataParser
	addressPrefix  string
}

func NewSelfSignedCredentialService(resolver authn.KeyResolver, metadataParser authn.RequestMetadataParser, addressPrefix string) authn.CredentialService {
	return credentialSrv{
		resolver:       resolver,
		metadataParser: metadataParser,
		addressPrefix:  addressPrefix,
	}
}

func (c credentialSrv) GetRequestToken(ctx context.Context) ([]byte, error) {
	// parse content from request context
	md, ok := c.metadataParser.Parse(ctx)
	if !ok {
		return nil, fmt.Errorf("missing request metadata")
	}

	vals := md.Get(TokenMetadataKey)
	if len(vals) == 0 {
		return nil, fmt.Errorf("missing authorization token")
	}
	token, found := strings.CutPrefix(vals[0], tokenPrefix)
	if !found {
		return nil, fmt.Errorf("missing token prefix %q", tokenPrefix)
	}

	return []byte(token), nil
}

func (c credentialSrv) VerifyRequestSubject(ctx context.Context, token []byte) (authn.SubjectInfo, error) {
	jws, err := jose.ParseSigned(string(token))
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
	key, err := userInfo.PubKey.Std()
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
		// Subject:  userInfo.Subject,
		Time: claims.Expiry.Time().Add(verifyLeewayTime),
	}

	err = claims.ValidateWithLeeway(expected, verifyLeewayTime)
	if err != nil {
		return authn.SubjectInfo{}, fmt.Errorf("JWS claim failed validation: %w", err)
	}

	// validate subject, either its a did:key that *must* match
	// the Issuer, or its a cosmos bech32 address that *must* also
	// match the Issuer
	subject := claims.Subject
	if strings.HasPrefix(subject, "did:key") {
		if subject != claims.Issuer {
			return authn.SubjectInfo{}, fmt.Errorf("JWS subject did:key verification failed")
		}
	} else {
		if c.addressPrefix == "" {
			return authn.SubjectInfo{}, fmt.Errorf("bech32 verification missing prefix")
		}
		// verify bech32 address
		addr, err := crypto.PublicKeyToBech32(userInfo.PubKey)
		if err != nil {
			return authn.SubjectInfo{}, fmt.Errorf("JWS subject bech32 verification failed: %w", err)
		}
		if addr != subject {
			return authn.SubjectInfo{}, fmt.Errorf("JWS subject bech32 mismatch %s != %s", addr, subject)
		}
	}

	return authn.SubjectInfo{
		Type:    "JWS",
		Subject: subject,
		PubKey:  userInfo.PubKey,
	}, nil
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

func publicKeyToBech32(prefix string, publicKey crypto.PublicKey) (string, error) {
	keyraw, err := publicKey.Raw()
	if err != nil {
		return "", fmt.Errorf("raw public key: %w", err)
	}
	buf := tmhash.SumTruncated(keyraw)
	return bech32.ConvertAndEncode(prefix, buf)

}

type claims struct {
	jwt.Claims
}

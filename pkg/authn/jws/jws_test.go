package jws

import (
	"context"
	gocrypto "crypto"
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"

	ssicrypto "github.com/TBD54566975/ssi-sdk/crypto"
	"github.com/TBD54566975/ssi-sdk/did/key"
	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"

	"github.com/sourcenetwork/orbis-go/pkg/authn"
	"github.com/sourcenetwork/orbis-go/pkg/authn/mocks"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/did"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testKID     = "alice#key"
	testSubject = "alice"

	ed25519PrivKey, ed25519PubKey = mustGeneratePrivateKey()
	cryptoPubKey                  = mustGetPublicKey(ed25519PubKey)
	ed25519Signer                 = mustMakeSigner(jose.EdDSA, ed25519PrivKey, testKID)
)

func mustGeneratePrivateKey() (ed25519.PrivateKey, ed25519.PublicKey) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	return priv, pub
}

func mustGetPublicKey(pub gocrypto.PublicKey) crypto.PublicKey {
	cpub, err := crypto.PublicKeyFromStdPublicKey(pub)
	if err != nil {
		panic(err)
	}
	return cpub
}

func mustMakeSigner(alg jose.SignatureAlgorithm, k interface{}, kid string) jose.Signer {
	opts := new(jose.SignerOptions)
	opts.WithHeader(jose.HeaderKey("kid"), kid)
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: alg, Key: k}, opts)
	if err != nil {
		panic("failed to create signer:" + err.Error())
	}

	return sig
}

func TestJWSCredentialService(t *testing.T) {
	ed25519PrivKey, ed25519PubKey = mustGeneratePrivateKey()
	cryptoPubKey = mustGetPublicKey(ed25519PubKey)
	ed25519Signer = mustMakeSigner(jose.EdDSA, ed25519PrivKey, testKID)

	// create signed JWT token
	claims := claims{
		Claims: jwt.Claims{
			Subject:  testSubject,
			Issuer:   testSubject,
			Expiry:   jwt.NewNumericDate(time.Now()),
			Audience: jwt.Audience{"orbis"},
		},
	}
	signedJWT, err := jwt.Signed(ed25519Signer).
		Claims(claims).
		CompactSerialize()
	if err != nil {
		t.Fatal(err)
	}

	// setup mocks

	// mock the request metadata and inject the test JWT into it
	// Get("authorization") => []string{signedJWT}}
	mockMD := mocks.NewMetadata(t)
	mockMD.EXPECT().Get(tokenMetadataKey).Return([]string{signedJWT})

	// mock the request parser and inject the above mocked metadata
	// Parse(ctx) => Metadata{"authorization": []string{signedJWT}}
	mockReqParser := mocks.NewRequestMetadataParser(t)
	mockReqParser.EXPECT().Parse(mock.Anything).Return(mockMD, true)

	// mock the key resolver to return our generated keys
	// Resolve(ctx, "alice#key") => SubjectInfo{"alice", publicKey}
	mockResolver := mocks.NewKeyResolver(t)
	mockResolver.EXPECT().Resolve(mock.Anything, testKID).Return(authn.SubjectInfo{
		Subject: testSubject,
		PubKey:  cryptoPubKey,
	}, nil)

	// Actual test block
	credService := NewSelfSignedCredentialService(mockResolver, mockReqParser)
	info, err := credService.GetAndVerifyRequestMetadata(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, info)
	require.Equal(t, authn.SubjectInfo{
		Subject: testSubject,
		PubKey:  cryptoPubKey,
		Type:    "JWS",
	}, info)
}

func TestDIDKeyJWSCredentialService(t *testing.T) {
	ed25519PrivKey, ed25519PubKey = mustGeneratePrivateKey()
	cryptoPubKey = mustGetPublicKey(ed25519PubKey)
	didKey, err := key.CreateDIDKey(ssicrypto.Ed25519, ed25519PubKey)
	if err != nil {
		t.Fatal(err)
	}
	suffix, err := didKey.Suffix()
	if err != nil {
		t.Fatal(err)
	}
	ed25519Signer = mustMakeSigner(jose.EdDSA, ed25519PrivKey, didKey.String()+"#"+suffix)

	// create signed JWT token
	claims := claims{
		Claims: jwt.Claims{
			Subject:  didKey.String(),
			Issuer:   didKey.String(),
			Expiry:   jwt.NewNumericDate(time.Now()),
			Audience: jwt.Audience{"orbis"},
		},
	}
	signedJWT, err := jwt.Signed(ed25519Signer).
		Claims(claims).
		CompactSerialize()
	if err != nil {
		t.Fatal(err)
	}

	// setup mocks

	// mock the request metadata and inject the test JWT into it
	// Get("authorization") => []string{signedJWT}}
	mockMD := mocks.NewMetadata(t)
	mockMD.EXPECT().Get(tokenMetadataKey).Return([]string{signedJWT})

	// mock the request parser and inject the above mocked metadata
	// Parse(ctx) => Metadata{"authorization": []string{signedJWT}}
	mockReqParser := mocks.NewRequestMetadataParser(t)
	mockReqParser.EXPECT().Parse(mock.Anything).Return(mockMD, true)

	// we'll use the actual did resolver this time instead of mocking
	resolver := did.NewResolver(key.Resolver{})

	// Actual test block
	credService := NewSelfSignedCredentialService(resolver, mockReqParser)
	info, err := credService.GetAndVerifyRequestMetadata(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, info)
	require.Equal(t, authn.SubjectInfo{
		Subject: didKey.String(),
		PubKey:  cryptoPubKey,
		Type:    "JWS",
	}, info)

}

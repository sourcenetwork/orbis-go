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
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
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
	testKID     = "did:key:alice#key"
	testSubject = "did:key:alice"

	// ed25519PrivKey, ed25519PubKey     = mustGenerateEd25519PrivateKey()
	// secp256k1PrivKey, secp256k1PubKey = mustGenerateSecp256k1PrivateKey()
	// cryptoPubKey                      = mustGetPublicKey(ed25519PubKey)
	// ed25519Signer                     = mustMakeSigner(jose.EdDSA, ed25519PrivKey, testKID)
	// secp256k1Signer                   = mustMakeSigner(jose.ES256K, secp256k1PrivKey, testKID)

	addressPresfix = "source"
	// testAddress    = mustGenerateBech32Address(addressPresfix, cryptoPubKey)
)

func mustGenerateBech32Address(pk crypto.PublicKey) string {
	addr, err := crypto.PublicKeyToBech32(pk)
	if err != nil {
		panic(err)
	}
	return addr
}

func mustGenerateEd25519PrivateKey() (ed25519.PrivateKey, ed25519.PublicKey) {
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

func mustGenerateSecp256k1PrivateKey() (*secp256k1.PrivateKey, *secp256k1.PublicKey) {
	priv, err := secp256k1.GeneratePrivateKeyFromRand(rand.Reader)
	if err != nil {
		panic(err)
	}
	return priv, priv.PubKey()
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

func TestJWSCredentialServicDID(t *testing.T) {
	ed25519PrivKey, ed25519PubKey := mustGenerateEd25519PrivateKey()
	cryptoPubKey := mustGetPublicKey(ed25519PubKey)
	ed25519Signer := mustMakeSigner(jose.EdDSA, ed25519PrivKey, testKID)

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
	mockMD.EXPECT().Get(TokenMetadataKey).Return([]string{tokenPrefix + signedJWT})

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
	ctx := context.Background()
	credService := NewSelfSignedCredentialService(mockResolver, mockReqParser, "")
	token, err := credService.GetRequestToken(ctx)
	require.NoError(t, err)
	require.NotNil(t, token)
	info, err := credService.VerifyRequestSubject(ctx, token)
	require.NoError(t, err)
	require.NotEmpty(t, info)
	require.Equal(t, authn.SubjectInfo{
		Subject: testSubject,
		PubKey:  cryptoPubKey,
		Type:    "JWS",
	}, info)
}

func TestJWSCredentialServicDIDSecp256k1(t *testing.T) {
	secp256k1PrivKey, secp256k1PubKey := mustGenerateSecp256k1PrivateKey()
	cryptoPubKey := mustGetPublicKey(secp256k1PubKey)
	secp256k1Signer := mustMakeSigner(jose.ES256K, secp256k1PrivKey, testKID)

	// create signed JWT token
	claims := claims{
		Claims: jwt.Claims{
			Subject:  testSubject,
			Issuer:   testSubject,
			Expiry:   jwt.NewNumericDate(time.Now()),
			Audience: jwt.Audience{"orbis"},
		},
	}
	signedJWT, err := jwt.Signed(secp256k1Signer).
		Claims(claims).
		CompactSerialize()
	if err != nil {
		t.Fatal(err)
	}

	// setup mocks

	// mock the request metadata and inject the test JWT into it
	// Get("authorization") => []string{signedJWT}}
	mockMD := mocks.NewMetadata(t)
	mockMD.EXPECT().Get(TokenMetadataKey).Return([]string{tokenPrefix + signedJWT})

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
	ctx := context.Background()
	credService := NewSelfSignedCredentialService(mockResolver, mockReqParser, "")
	token, err := credService.GetRequestToken(ctx)
	require.NoError(t, err)
	require.NotNil(t, token)
	info, err := credService.VerifyRequestSubject(ctx, token)
	require.NoError(t, err)
	require.NotEmpty(t, info)
	require.Equal(t, authn.SubjectInfo{
		Subject: testSubject,
		PubKey:  cryptoPubKey,
		Type:    "JWS",
	}, info)
}

func TestJWSCredentialServicBech32(t *testing.T) {
	ed25519PrivKey, ed25519PubKey := mustGenerateEd25519PrivateKey()
	cryptoPubKey := mustGetPublicKey(ed25519PubKey)
	ed25519Signer := mustMakeSigner(jose.EdDSA, ed25519PrivKey, testKID)
	testAddress := mustGenerateBech32Address(cryptoPubKey)

	// create signed JWT token
	claims := claims{
		Claims: jwt.Claims{
			Subject:  testAddress,
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
	mockMD.EXPECT().Get(TokenMetadataKey).Return([]string{tokenPrefix + signedJWT})

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
	ctx := context.Background()
	credService := NewSelfSignedCredentialService(mockResolver, mockReqParser, addressPresfix)
	token, err := credService.GetRequestToken(ctx)
	require.NoError(t, err)
	require.NotNil(t, token)
	info, err := credService.VerifyRequestSubject(ctx, token)
	require.NoError(t, err)
	require.NotEmpty(t, info)
	require.Equal(t, authn.SubjectInfo{
		Subject: testAddress,
		PubKey:  cryptoPubKey,
		Type:    "JWS",
	}, info)
}

func TestDIDKeyJWSCredentialService(t *testing.T) {
	ed25519PrivKey, ed25519PubKey := mustGenerateEd25519PrivateKey()
	cryptoPubKey := mustGetPublicKey(ed25519PubKey)
	didKey, err := key.CreateDIDKey(ssicrypto.Ed25519, ed25519PubKey)
	if err != nil {
		t.Fatal(err)
	}
	suffix, err := didKey.Suffix()
	if err != nil {
		t.Fatal(err)
	}
	ed25519Signer := mustMakeSigner(jose.EdDSA, ed25519PrivKey, didKey.String()+"#"+suffix)

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
	mockMD.EXPECT().Get(TokenMetadataKey).Return([]string{tokenPrefix + signedJWT})

	// mock the request parser and inject the above mocked metadata
	// Parse(ctx) => Metadata{"authorization": []string{signedJWT}}
	mockReqParser := mocks.NewRequestMetadataParser(t)
	mockReqParser.EXPECT().Parse(mock.Anything).Return(mockMD, true)

	// we'll use the actual did resolver this time instead of mocking
	resolver := did.NewResolver(key.Resolver{})

	// Actual test block
	ctx := context.Background()
	credService := NewSelfSignedCredentialService(resolver, mockReqParser, "")
	token, err := credService.GetRequestToken(ctx)
	require.NoError(t, err)
	require.NotNil(t, token)
	info, err := credService.VerifyRequestSubject(ctx, token)
	require.NoError(t, err)
	require.NotEmpty(t, info)
	require.Equal(t, authn.SubjectInfo{
		Subject: didKey.String(),
		PubKey:  cryptoPubKey,
		Type:    "JWS",
	}, info)

}

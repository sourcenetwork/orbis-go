package keyring

import (
	"crypto/rand"
	"testing"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/stretchr/testify/require"
)

func TestKeyringAsymmetricBasicEd25519(t *testing.T) {
	keyring, err := New("test", t.TempDir())
	require.NoError(t, err)

	// generate our test keypair
	testKeyringAsymmetricBasic(t, keyring, crypto.Ed25519)
}

func TestKeyringAsymmetricBasicSecp256k1(t *testing.T) {
	keyring, err := New("test", t.TempDir())
	require.NoError(t, err)

	// generate our test keypair
	testKeyringAsymmetricBasic(t, keyring, crypto.Secp256k1)
}

func testKeyringAsymmetricBasic(t *testing.T, keyring Keyring, keyType crypto.KeyType) {
	testKeyName := "testKey1"

	// nil arg is to ensure a random reader
	priv, pub, err := crypto.GenerateKeyPair(keyType, rand.Reader)
	require.NoError(t, err)

	err = keyring.Set(testKeyName, priv)
	require.NoError(t, err)

	priv2, err := keyring.Get(testKeyName)
	require.NoError(t, err)
	require.True(t, priv.Equals(priv2))

	pv2 := priv2.(crypto.PrivateKey)
	require.True(t, pv2.GetPublic().Equals(pub))

	// public keys
	testKeyName2 := "testKey2"
	err = keyring.Set(testKeyName2, pub)
	require.NoError(t, err)

	pk2, err := keyring.Get(testKeyName2)
	require.NoError(t, err)
	require.True(t, pk2.Equals(pub))

	keyInfos, err := keyring.List()
	require.NoError(t, err)

	require.Equal(t, testKeyName, keyInfos[0].Name)
	require.True(t, priv.Equals(keyInfos[0].Key))

	require.Equal(t, testKeyName2, keyInfos[1].Name)
	require.True(t, pub.Equals(keyInfos[1].Key))

	err = keyring.Delete(testKeyName)
	require.NoError(t, err)

	_, err = keyring.Get(testKeyName)
	require.Error(t, err)
}

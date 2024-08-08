package crypto

import (
	"crypto/rand"
	"fmt"
	mrand "math/rand"
	"testing"

	"github.com/sourcenetwork/orbis-go/pkg/crypto/suites/secp256k1"

	ic "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v3/group/edwards25519"
)

func TestKeyGeneration(t *testing.T) {
	// libp2p key
	p2pPriv, _, err := ic.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)

	buf, err := p2pPriv.Raw()
	require.NoError(t, err)
	fmt.Printf("%x\n", buf)

	suite := edwards25519.NewBlakeSHA256Ed25519()
	kyberKey := suite.NewKey(suite.RandomStream())
	pubKey := suite.Point().Mul(kyberKey, nil)

	buf2, err := kyberKey.MarshalBinary()
	require.NoError(t, err)

	buf3, err := pubKey.MarshalBinary()
	require.NoError(t, err)
	fmt.Printf("%x%x\n", buf2, buf3)

}

func TestSecp256PublicKeyMarshalling(t *testing.T) {
	// libp2p key
	p2pPriv, _, err := ic.GenerateSecp256k1Key(rand.Reader)
	require.NoError(t, err)

	buf, err := p2pPriv.GetPublic().Raw()
	require.NoError(t, err)
	fmt.Printf("%x\n", buf)
	fmt.Println("serialized length:", len(buf))

	suite := secp256k1.NewBlakeKeccackSecp256k1()
	point := suite.Point()
	err = point.UnmarshalBinary(buf)
	require.NoError(t, err)
	// fmt.Println(point)
	// t.Fail()
	// kyberKey := suite.NewKey(suite.RandomStream())
	// pubKey := suite.Point().Mul(kyberKey, nil)

	// buf2, err := kyberKey.MarshalBinary()
	// require.NoError(t, err)

	// buf3, err := pubKey.MarshalBinary()
	// require.NoError(t, err)
	// fmt.Printf("%x%x\n", buf2, buf3)
	// t.Fail()
}

func TestPeerIDToPublicKeySecp256k1(t *testing.T) {
	seed := 1
	randomness := mrand.New(mrand.NewSource(int64(seed)))
	suite := secp256k1.NewBlakeKeccackSecp256k1()
	p2pPriv, _, err := ic.GenerateSecp256k1Key(randomness)
	require.NoError(t, err)

	pid, err := peer.IDFromPublicKey(p2pPriv.GetPublic())
	require.NoError(t, err)
	fmt.Println(pid)

	cprivKey, err := PrivateKeyFromLibP2P(p2pPriv)
	require.NoError(t, err)
	require.NotNil(t, cprivKey)

	pubKey, err := pid.ExtractPublicKey()
	require.NoError(t, err)
	require.True(t, pubKey.Equals(p2pPriv.GetPublic()))

	cpubKey, err := PublicKeyFromLibP2P(pubKey)
	require.NoError(t, err)
	require.True(t, pubKey.Equals(cpubKey.(*pubKeyLibP2P).PubKey))

	/////////////////////////////////////////

	buf, err := cpubKey.Raw()
	require.NoError(t, err)
	p1 := suite.Point()
	err = p1.UnmarshalBinary(buf)
	require.NoError(t, err)

	require.True(t, p1.Equal(cpubKey.Point()))
	// fmt.Println(point.String())
	// fmt.Println(cpubKey.Point().String())
	// t.Fail()

	scalar := cprivKey.Scalar() // private Key scalar
	// suite := secp256k1.NewBlakeKeccackSecp256k1()
	p2 := suite.Point().Mul(scalar, nil) // public key point
	require.True(t, p2.Equal(cpubKey.Point()))
	t.Log("computed point:", p2.String())
	t.Log("extracted point:", cpubKey.Point().String())

}

func TestPeerIDToPublicKeyEd25519(t *testing.T) {
	seed := 1
	randomness := mrand.New(mrand.NewSource(int64(seed)))
	suite := edwards25519.NewBlakeSHA256Ed25519()
	p2pPriv, _, err := ic.GenerateEd25519Key(randomness)
	require.NoError(t, err)

	pid, err := peer.IDFromPublicKey(p2pPriv.GetPublic())
	require.NoError(t, err)
	fmt.Println(pid)

	cprivKey, err := PrivateKeyFromLibP2P(p2pPriv)
	require.NoError(t, err)
	require.NotNil(t, cprivKey)

	pubKey, err := pid.ExtractPublicKey()
	require.NoError(t, err)
	require.True(t, pubKey.Equals(p2pPriv.GetPublic()))

	cpubKey, err := PublicKeyFromLibP2P(pubKey)
	require.NoError(t, err)
	require.True(t, pubKey.Equals(cpubKey.(*pubKeyLibP2P).PubKey))

	scalar := cprivKey.Scalar()             // private Key scalar
	point := suite.Point().Mul(scalar, nil) // public key point
	require.True(t, point.Equal(cpubKey.Point()))
	t.Log("computed point:", point.String())
	t.Log("extracted point:", cpubKey.Point().String())

	// t.Fail()
}

package crypto

import (
	"crypto/rand"
	"fmt"
	"testing"

	ic "github.com/libp2p/go-libp2p/core/crypto"
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
	t.Fail()
}

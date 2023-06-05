package rabin

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	rabindkg "go.dedis.ch/kyber/v3/share/dkg/rabin"
	"go.dedis.ch/kyber/v3/util/random"
)

func TestProtoDealSerialization(t *testing.T) {
	rand := random.New()
	suite := edwards25519.NewBlakeSHA256Ed25519()

	s1 := suite.Scalar().Pick(rand)
	p1 := suite.Point().Mul(s1, nil)
	p2 := suite.Point().Pick(rand)
	p3 := suite.Point().Pick(rand)

	points := []kyber.Point{p1, p2, p3}
	dkg, err := rabindkg.NewDistKeyGenerator(suite, s1, points, 2)
	if err != nil {
		panic(err)
	}

	deals, err := dkg.Deals()
	if err != nil {
		panic(err)
	}

	d1, err := dealToProto(deals[1])
	if err != nil {
		panic(err)
	}

	d1r, err := dealFromProto(suite, d1)
	require.NoError(t, err)

	require.Equal(t, deals[1].Index, d1r.Index)
	require.Equal(t, deals[1].Deal.Nonce, d1r.Deal.Nonce)
	require.Equal(t, deals[1].Deal.Signature, d1r.Deal.Signature)
	require.Equal(t, deals[1].Deal.Cipher, d1r.Deal.Cipher)
	require.True(t, deals[1].Deal.DHKey.Equal(d1r.Deal.DHKey))

	// require.Equal(t, deals[1].Deal.DHKey, d1r.Deal.DHKey)
	// require.Equal(t, deals[1].Index, d1r.Index)
	// require.Equal(t, deals[1].Index, d1r.Index)

	require.Equal(t, deals[1], d1r)
}

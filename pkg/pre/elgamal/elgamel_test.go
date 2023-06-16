package elgamal

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/suites"
	"go.dedis.ch/kyber/v3/util/random"
)

func TestReencryptAndVerify(t *testing.T) {

	var (
		n       = 5
		th      = 3
		ste     = suites.MustFind("ed25519")
		s       = ste.Scalar().Pick(ste.RandomStream())
		priPoly = share.NewPriPoly(ste, th, s, ste.RandomStream())
		pubPoly = priPoly.Commit(nil)
		dkgPk   = pubPoly.Commit()

		rdrSk = ste.Scalar().Pick(ste.RandomStream())
		rdrPk = ste.Point().Mul(rdrSk, nil)
	)

	var pubShares []*share.PubShare

	// Generate a random secret.
	scrt := make([]byte, 32)
	random.Bytes(scrt, random.New())

	// 1. Encrypt the secret under the DKG public key.
	encCmt, encScrt := EncryptSecret(ste, dkgPk, scrt)

	for idx := 0; idx < n; idx++ {

		dkgSki := priPoly.Eval(idx).V
		dkgCmt := pubPoly.Eval(idx).V

		// 2. Re-encrypt the key under the reader's public key.
		xncSki, chlgi, proofi, err := reencrypt(ste, dkgSki, rdrPk, encCmt)
		require.NoErrorf(t, err, "failed to reencrypt for share %d", idx)

		// 3. Verify the re-encryption from other nodes.
		err = verify(ste, rdrPk, encCmt, xncSki, chlgi, proofi, dkgCmt)
		require.NoErrorf(t, err, "failed to verify reencryption for share %d", idx)

		pubShare := &share.PubShare{I: idx, V: xncSki}
		pubShares = append(pubShares, pubShare)
	}

	// 4 - Recover re-encrypted commmitment using Lagrange interpolation.
	// ski * (xG + rG) => rsG + xsG
	xncCmt, err := share.RecoverCommit(ste, pubShares, th, n)
	require.NoErrorf(t, err, "failed to recover commit")

	// 5 - Decode encrypted key with re-encrypted commitment and reader's privatekey.
	scrtHat, err := DecryptSecret(ste, encScrt, dkgPk, xncCmt, rdrSk)
	require.NoErrorf(t, err, "failed to decode key")
	require.Equal(t, scrt, scrtHat)
}

package elgamal

import (
	"crypto/sha256"
	"fmt"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/suites"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	"github.com/sourcenetwork/orbis-go/pkg/pre"
	"github.com/sourcenetwork/orbis-go/pkg/types"

	logging "github.com/ipfs/go-log"
)

const name = "elgamal"

var (
	// For Bulletin handlers
	InvalidReplyNamespace string = "invalidreply"

	// For P2P handlers
	EncryptedSecretRequest string = "encscrtrequest"
	EncryptedSecretReply   string = "encscrtreply"
)

var log = logging.Logger("orbis/pre/elgamal")

var (
	_ pre.PRE = (*ThesholdDealer)(nil)
)

// ThesholdDealer is a semi-trusted dealer implementations
// of the elgamal threshold re-encryption algorithm.
//
// In this instance, the semi-trusted dealer is responsible
// only for aggregating the respective node shares
// but is not able to recover the underlying secret. This
// helps in computation complexity and communication complexity
// between the server and the requesting client.
//
// Given the terms defined:
//
//	scrt    (k)         - A secret to be re-encrypted.
//	dkgPk   (sG)        - Aggregate public key of DKG.
//	dkgSki  (ski)       - Private share of secret key of DKG.
//	dkgCmt  (ci)        - Commitment (public polynomial) of DKG at index i
//	encCmt  (rG)        - Schnorr commitment of encoded keys.
//	encScrt (rsG + K)   - Encrypted key-slices.
//	xncCmt  (rsG + xsG) - Re-encrypted schnorr commitment.
//	xncSki  (Ui)        - Re-encrypted secret share.
//	rdrSk   (x)         - Reader's secret key.
//	rdrPk   (xG)        - Reader's public key.
//	chlgi   (ei)        - Random oracle challenge at index i.
//	proofi  (fi)        - NIZK proofi of re-encryption at index i. (ri + ei * ski)
//	rndi    (ri)        - Random number.
//
// A minimal flow of secret re-encryption is as follows:
//
// 1. Encrypt:
//
//	encCmt, encScrt = EncryptSecret(dkgPk, scrt)
//
//	- A document owner encrypts a document with a symetric key (secret).
//	- The owner encrypts the secrets with the aggregate public key of the DKG (dkgPk).
//	- The owner generates a schnorr commitment (encCmt) and the encrypted secret (encScrt).
//
// 2. Re-encrypt:
//
//	xncSki, chlgi, proofi = Reencrypt(dkgSki, rdrPk, encCmt)
//
//	- A document reader generates a pair of secret key (rdrSk) and public key (rdrPk).
//	- Each DKG node:
//	  - Re-encrypts the commitment (encCmt) into re-enncrypted secret share (xncSki/Ui)
//	    using its private share of DKG secret key (dkgSki) and the the reader's public key (rdrPk).
//	  - Generates a random oracle challenge () and a NIZK proofi (proofi).
//	    - Generates a random number (ri).
//	    - Generates a random oracle challenge (chlgi/ei)
//	      - UiHat = ri * (xG + rG)
//	      - HiHat = ri * G
//	      - chlgi = Hash(Ui, UiHat, HiHat)
//	  - Generates proofi(fi) as ri + ei * si
//
// 3. Verify:
//
//	Each DKG nodes recieves and verifies (xncSki, chlgi, proofi) from other node re-encrypted.
//
//	Verify(rdrPk, encCmt, xncSki, chlgi, proofi, dkgCmt)
//
//	- Reconstruct the UiHat and HiHat.
//	- Verify the reconstructed challenge (ei) matches the challenge ().
//
//	  UiHat(verifier) = fi              * (xG + rG) - ei * Ui
//	                  = (ri + ski * ei) * (xG + rG) - ei * ski * (xG + rG)
//	                  = (ri + ei * ski - ei * ski) * (xG + rG)
//	                  = ri * (xG + rG)
//	                  = UiHat(re-encreptor)
//
//	  hiHat(verifier) = fi * G                - ei * ci
//	                  = (ri + ei * ski) * G   - ei * ci
//	                  = ri * G + ei * ski * G - ei * ci
//	                  = ri * G + ei * ci      - ei * ci
//	                  = ri * G
//	                  = hiHat(re-encreptor)
//
//	- Verify the reconstructed (ei) matches the one from the re-encryptor.
//	  Hash(Ui, UiHatVerifier, CiHatVerifier) = Hash(Ui, UiHatReEncryptor, CiHatReEncryptor)
//
// 4. Recover:
//
//	  xncCmt = RecoverCommit([]xncSk, th, n)
//
//	- Given at least thrashold(th) number of verified xncSki,we can recover the
//	  re-encrypted commitment (xncCmt) using lagrange interpolation.
//
//	  ski * (xG + rG) ==> rsG + xsG
//
// 5. Decode encrypted key (encScrt) with dkgPk, xncCmt, rdrSk.
//
//	key = DecryptSecret(encScrt, dkgPk, xncCmt, rdrSk)
//
//	- The encrypted key encScrt consists of key K and an encrypting point rsG.
//	- The re-encrypted commitment xncCmt consists of rsG and xsG.
//	- The reader can recover the K as follows:
//
//	  encScrt = rsG + K
//	  xncCmt = rsG + xsG
//	  rdrSk  = x
//	  dkgPk  = sG
//
//	  encScrt - xncCmt + rdrSk * dkgPk
//	  = rsG + K - (rsG + xsG) + xsG
//	  = K
type ThesholdDealer struct {
}

func New(db *db.DB, repoKey []db.RepoKey) (pre.PRE, error) {

	return &ThesholdDealer{}, nil
}

func (e *ThesholdDealer) Init(rid types.RingID, n int32, t int32) error {
	return nil
}

func (e *ThesholdDealer) Name() string {
	return name
}

func (e *ThesholdDealer) Reencrypt(distKeyShare crypto.DistKeyShare, scrt *types.Secret, rdrPk crypto.PublicKey) (pre.ReencryptReply, error) {

	var reply pre.ReencryptReply
	ste, err := crypto.SuiteForType(rdrPk.Type())
	if err != nil {
		return reply, fmt.Errorf("get suite for type: %w", err)
	}

	idx := distKeyShare.PriShare.I
	ski := distKeyShare.PriShare.V

	encCmt := ste.Point()
	err = encCmt.UnmarshalBinary(scrt.EncCmt)
	if err != nil {
		return reply, fmt.Errorf("unmarshal encCmt: %w", err)
	}

	xncSki, chlgi, proofi, err := reencrypt(ste, ski, rdrPk.Point(), encCmt)
	if err != nil {
		return reply, err
	}

	reply = pre.ReencryptReply{
		Share: share.PubShare{
			I: idx,
			V: xncSki,
		},
		Challenge: chlgi,
		Proof:     proofi,
	}

	return reply, nil
}

// Verify verifies an incoming re-encryption reply from another node.
func (e *ThesholdDealer) Verify(rdrPk crypto.PublicKey, dkgCmt crypto.PubPoly, encCmt kyber.Point, r pre.ReencryptReply) error {

	ste, err := crypto.SuiteForType(rdrPk.Type())
	if err != nil {
		return fmt.Errorf("get suite for type: %w", err)
	}

	xncSki := r.Share.V
	idx := r.Share.I
	err = verify(ste,
		rdrPk.Point(),
		encCmt,
		xncSki,
		r.Challenge,
		r.Proof,
		dkgCmt.PubPoly.Eval(idx).V,
	)
	if err != nil {
		return fmt.Errorf("verification: %w", err)
	}

	return nil
}

func (e *ThesholdDealer) Recover(ste suites.Suite, xncSki []*share.PubShare, t int, n int) (kyber.Point, error) {
	if len(xncSki) < t {
		return nil, nil
	}

	return share.RecoverCommit(ste, xncSki, t, n)
}

// reencrypt re-encrypts a secret share using the reciever public key.
//
// Input:
//
//	ste          - Crypto suite
//	dkgSki (ski) - Private share of secret key of DKG.
//	rdrPk  (xG)  - Public key of of the reader.
//	encCmt (rG)  - Schnorr commit of encoded keys.
//
// Output:
//
//	xncSki (Ui) - Re-encrypted secret share.
//	chlgi  (ei) - Random oracle challenge.
//	proofi (fi) - NIZK proofi of re-encryption.
//	err         - Error if re-encryption fails.
func reencrypt(
	ste suites.Suite,
	dkgSki kyber.Scalar,
	rdrPk kyber.Point,
	encCmt kyber.Point,
) (
	xncSki kyber.Point,
	chlgi kyber.Scalar,
	proofi kyber.Scalar,
	err error,
) {
	// Re-encrypted secret share (Ui)
	xrG := ste.Point().Add(rdrPk, encCmt) // xrG = xG + rG
	xncSki = ste.Point().Mul(dkgSki, xrG) // Ui  = ski * (xG + rG)

	// Produce random oracle challenge (ei)
	// ei = Hash(Ui + UiHat + HiHat)
	ri := ste.Scalar().Pick(ste.RandomStream()) // ri    = Random scalar
	uiHat := ste.Point().Mul(ri, xrG)           // UiHat = ri * (xG + rG)
	hiHat := ste.Point().Mul(ri, nil)           // HiHat = ri * G

	b, err := hashPoints(xncSki, uiHat, hiHat)
	if err != nil {
		return xncSki, chlgi, proofi, fmt.Errorf("marshal Ui: %v", err)
	}
	chlgi = ste.Scalar().SetBytes(b)

	// Produce NIZK proofi of re-encryption (fi)
	// fi = ri + ei * ski
	proofi = ste.Scalar().Add(ri, ste.Scalar().Mul(chlgi, dkgSki))

	return xncSki, chlgi, proofi, nil
}

// Input:
//
//	ste          - Crypto suite
//	rdrPk  (xG)  - Public key of of the reader.
//	encCmt (rG)  - Schnorr commit of encoded keys.
//	dkgSki (ski) - Re-encrepyred share of commitment.
//	chlgi  (ei)  - Random oracle challenge at index i.
//	proofi (fi)  - NIZK proofi of re-encryption at index i.
//	dkgCmt (ci)  - Commitment (public polynomial) of DKG at index i.
//
// Output:
//
//	err - Error if verification fails.
func verify(
	ste suites.Suite,
	rdrPk kyber.Point,
	encCmt kyber.Point,
	dkgSki kyber.Point,
	chlgi kyber.Scalar,
	proofi kyber.Scalar,
	dkgCmt kyber.Point,
) error {

	// Reconstruct UiHat.
	fixrG := ste.Point().Mul(proofi, ste.Point().Add(rdrPk, encCmt)) // fi * (xG + rG)
	eiui := ste.Point().Mul(chlgi, dkgSki)                           // ei * Ui
	uiHat := ste.Point().Sub(fixrG, eiui)                            // UiHat = fi * (xG + rG) - ei * Ui

	// Reconstruct HiHat.
	fig := ste.Point().Mul(proofi, nil)    // FiG   = fi * G
	eici := ste.Point().Mul(chlgi, dkgCmt) // EiHi  = ei * ci
	hiHat := ste.Point().Sub(fig, eici)    // HiHat = fi * G - ei * ci

	// Reconstruct random oracle challenge (ei).
	// ei = Hash(Ui + UiHat + HiHat)
	b, err := hashPoints(dkgSki, uiHat, hiHat)
	if err != nil {
		return fmt.Errorf("failed to marshal Ui: %v", err)
	}
	chlg := ste.Scalar().SetBytes(b)

	// Verify local challenge
	if !chlg.Equal(chlgi) {
		return fmt.Errorf("failed verification")
	}

	return nil
}

// EncryptSecret encrypts a secret using the aggregate public key of the DKG.
//
// Input:
//
//	ste        - Crypto suite.
//	dkgPk (sG) - Aggregate public key of the DKG.
//	scrt  (k)  - Secret to be encrypted.
//
// Output:
//
//	encCmt - Schnorr commit (rG)
//	encScrt - Encrypted key-slices (rsG + Ki)
func EncryptSecret(
	ste suites.Suite,
	dkgPk kyber.Point,
	scrt []byte,
) (
	encCmt kyber.Point,
	encScrt []kyber.Point,
) {

	r := ste.Scalar().Pick(ste.RandomStream())
	encCmt = ste.Point().Mul(r, nil) // rG = r * G
	rsG := ste.Point().Mul(r, dkgPk) // rsG = r * sG

	for len(scrt) > 0 {
		k := ste.Point().Embed(scrt, ste.RandomStream())
		scrt = scrt[min(len(scrt), k.EmbedLen()):]

		keyi := ste.Point().Add(rsG, k)
		encScrt = append(encScrt, keyi)
	}

	return encCmt, encScrt
}

// DecryptSecret decrypts a secret using the reader's secret key.
//
// Input:
//
//	ste                 - Crypto suite.
//	encScrt (rsG + K) - Encrypted key-slices.
//	dkgPk  (sG)         - Aggregate public key of DKG.
//	xncCmt (rsG + xsG)  - Re-encrypted schnorr-commit.
//	rdrSk  (x)          - Secret key of the reader.
//
// Output:
//
//	scrt - Recovered secret.
//	err - Error if decryption failed.
func DecryptSecret(
	ste suites.Suite,
	encScrt []kyber.Point,
	dkgPk kyber.Point,
	xncCmt kyber.Point,
	rdrSk kyber.Scalar,
) (
	scrt []byte,
	err error,
) {
	// To retrieve each key slice (Ki) from the encrypted key point (Ki + rsG),
	// we must deduct the encryption point (rsG). This can be inferred from
	// the re-encrypted schnorr-commit (rsG + xsG) by removing the product of
	// the reader's secret key (x) and the aggregate public key from the DKG (sG).

	xsG := ste.Point().Mul(rdrSk, dkgPk) // xsG = x * sG
	rsG := ste.Point().Sub(xncCmt, xsG)  // rsG = (rsG + xsG) - xsG

	for _, encKey := range encScrt {
		k := ste.Point().Sub(encKey, rsG) // K = (rsG + K) - rsG
		keyi, err := k.Data()
		if err != nil {
			return nil, fmt.Errorf("extract key share from key point: %w", err)
		}
		scrt = append(scrt, keyi...)
	}

	return scrt, err
}

func hashPoints(points ...kyber.Point) ([]byte, error) {
	hash := sha256.New()
	for _, p := range points {
		_, err := p.MarshalTo(hash)
		if err != nil {
			return nil, fmt.Errorf("marshal point: %v", err)
		}
	}
	return hash.Sum(nil), nil
}

// TODO: remove after Go 1.21
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

package elgamal

import (
	"crypto/sha256"
	"fmt"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/suites"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/pre"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
)

var (
	_ pre.ReencryptReply = (*ReencryptReply)(nil)
	_ pre.PRE            = (*ThesholdDealer)(nil)
)

type ReencryptReply struct {
	Ui share.PubShare // nodes re-encrypted secret share
	Ei kyber.Scalar   // random oracle challenge
	Fi kyber.Scalar   // nizk proof of re-encryption
}

// Share
func (rr ReencryptReply) Share() share.PubShare {
	return rr.Ui
}

// Challenge
func (rr ReencryptReply) Challenge() kyber.Scalar {
	return rr.Ei
}

// Proof
func (rr ReencryptReply) Proof() kyber.Scalar {
	return rr.Fi
}

// ThesholdDealer is a semi-trusted dealer implementations
// of the elgamal threshold re-encryption algorithm.
//
// In this instance, the semi-trusted dealer is responsible
// only for aggregating the respective node shares
// but is not able to recover the underlying secret. This
// helps in computation complexity and communication complexity
// between the server and the requesting client.
type ThesholdDealer struct {
	// Embedded PSS Service
	pss   pss.PSS
	suite suites.Suite

	u              kyber.Point       // Encrypted Secret
	pk             crypto.PublicKey  // Receiver Public Key
	uis            []*share.PubShare // reencrypted secret shares
	verifiedShares *int64

	finished bool
}

func (e *ThesholdDealer) Name() string {
	return "elgamal"
}

// pk: reciever public key
// u: key share commit
func (e *ThesholdDealer) Reencrypt(pk crypto.PublicKey, u kyber.Point) (pre.ReencryptReply, error) {
	// ui = (u ^ ski) * (pk ^ ski)
	pkPoint := pk.Point()
	ski := e.pss.Share().V
	uski := e.suite.Point().Mul(ski, u)
	pkski := e.suite.Point().Mul(ski, pkPoint)
	ui := uski.Add(uski, pkski)

	si := e.suite.Scalar().Pick(e.suite.RandomStream())               // si = random scalar
	uiHat := e.suite.Point().Mul(si, e.suite.Point().Add(u, pkPoint)) // uiHat = (u * pk) ^ si
	hiHat := e.suite.Point().Mul(si, nil)                             // hiHat = g^si (nil implies default base g)

	// ei = H(ui + uiHat + hiHat)
	hash := sha256.New()
	ui.MarshalTo(hash)
	uiHat.MarshalTo(hash)
	hiHat.MarshalTo(hash)
	ei := e.suite.Scalar().SetBytes(hash.Sum(nil))

	fi := e.suite.Scalar().Add(si, e.suite.Scalar().Mul(ei, ski)) // fi = si + (ski * ei)

	return ReencryptReply{
		Ui: share.PubShare{
			V: ui,
			I: e.pss.Share().I,
		},
		Ei: ei,
		Fi: fi,
	}, nil
}

// Process verifies an incoming Re-encryption reply from another node.
func (e *ThesholdDealer) Process(from pss.Node, r pre.ReencryptReply) error {
	// verify
	pkPoint := e.pk.Point()
	fi := r.Proof()
	ei := r.Challenge()
	ui := r.Share()

	ufi := e.suite.Point().Mul(fi, e.suite.Point().Add(e.u, pkPoint)) // ufi = (u*pk)^fi
	uiei := e.suite.Point().Mul(e.suite.Scalar().Neg(ei), ui.V)       // uiei = ui^-ei
	uiHat := e.suite.Point().Add(ufi, uiei)                           // uihat = uiei ^ ufi

	gfi := e.suite.Point().Mul(fi, nil)                        // gfi = g^fi
	gxi := e.pss.PublicPoly().Eval(from.Index()).V             // gxi = f(i)
	hiei := e.suite.Point().Mul(e.suite.Scalar().Neg(ei), gxi) // hiei = gxi^-ei
	hiHat := e.suite.Point().Add(gfi, hiei)                    // hiHat = gfi + hiei

	// locally produce challenge hash (random oracle)
	hash := sha256.New()
	ui.V.MarshalTo(hash)
	uiHat.MarshalTo(hash)
	hiHat.MarshalTo(hash)
	challenge := e.suite.Scalar().SetBytes(hash.Sum(nil))

	// verify local challenge
	// H(ui + uiHat + hiHat) == r.Ei
	if challenge.Equal(ei) {
		e.uis[from.Index()] = &ui
		// atomic add
		// check if threshold met - signal
		return nil
	}

	return fmt.Errorf("failed verification")
}

func (e *ThesholdDealer) Recover() (kyber.Point, error) {
	// verify length
	if len(e.uis) < e.pss.Threshold() {
		// we're not ready to recover yet
		return nil, nil
	}

	return share.RecoverCommit(e.suite, e.uis, e.pss.Threshold(), e.pss.Num())
}

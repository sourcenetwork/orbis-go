package pre

import (
	"crypto/sha256"
	"fmt"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/pss"
)

type ReencryptReply struct {
	Ui share.PubShare // nodes re-encrypted secret share
	Ei kyber.Scalar   // random oracle challenge
	Fi kyber.Scalar   // nizk proof of re-encryption
}

// ElGamalDealer is a semi-trusted dealer implementations
// of the elgamal re-encryption algorithm.
//
// In this instance, the semi-trusted dealer is responsible
// only for aggregating the respective node shares
// but is not able to recover the underlying secret. This
// helps in computation complexity and communication complexity
// between the server and the requesting client.
type ElGamalDealer struct {
	pss pss.Service

	group  kyber.Group
	random kyber.Random

	u              kyber.Point       // Encrypted Secret
	pk             crypto.PublicKey  // Receiver Public Key
	uis            []*share.PubShare // reencrypted secret shares
	verifiedShares *int64

	finished bool
}

// pk: reciever public key
// u: key share commit
func (e *ElGamalDealer) Reencrypt(pk crypto.PublicKey, u kyber.Point) (ReencryptReply, error) {

	p := e.group.Point()
	s := e.group.Scalar()

	ski := e.pss.Share().V
	uski := p.Mul(ski, u)
	pkski := p.Mul(ski, pk)

	ui := uski.Add(uski, pkski)           // ui    = (u ^ ski) * (pk ^ ski)
	si := s.Pick(e.random.RandomStream()) // si    = random scalar
	uiHat := p.Mul(si, p.Add(u, pk))      // uiHat = (u * pk) ^ si
	hiHat := p.Mul(si, nil)               // hiHat = g^si (nil implies default base g)

	hash := sha256.New()
	ui.MarshalTo(hash)
	uiHat.MarshalTo(hash)
	hiHat.MarshalTo(hash)

	ei := s.SetBytes(hash.Sum(nil)) // ei = H(ui + uiHat + hiHat)
	fi := s.Add(si, s.Mul(ei, ski)) // fi = si + (ski * ei)

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
func (e *ElGamalDealer) Process(idx int, r ReencryptReply) error {

	p := e.group.Point()
	s := e.group.Scalar()

	// verify
	ufi := p.Mul(r.Fi, p.Add(e.u, e.pk)) // ufi   = (u*pk)^fi
	uiei := p.Mul(s.Neg(r.Ei), r.Ui.V)   // uiei  = ui^-ei
	uiHat := p.Add(ufi, uiei)            // uihat = uiei ^ ufi

	gfi := p.Mul(r.Fi, nil)               // gfi   = g^fi
	gxi := e.pss.PublicPoly().Eval(idx).V // gxi   = f(i)
	hiei := p.Mul(s.Neg(r.Ei), gxi)       // hiei  = gxi^-ei
	hiHat := p.Add(gfi, hiei)             // hiHat = gfi + hiei

	// locally produce challenge hash (random oracle)
	hash := sha256.New()
	r.Ui.V.MarshalTo(hash)
	uiHat.MarshalTo(hash)
	hiHat.MarshalTo(hash)
	c := s.SetBytes(hash.Sum(nil))

	// verify
	// H(ui + uiHat + hiHat) == r.Ei
	if c.Equal(r.Ei) {
		e.uis[idx] = &r.Ui
		// atomic add
		// check if threshold met - signal
		return nil
	}

	return fmt.Errorf("failed verification")
}

func (e *ElGamalDealer) Recover() (kyber.Point, error) {
	// verify length
	if len(e.uis) < e.pss.Threshold() {
		// we're not ready to recover yet
		return nil, nil
	}

	return share.RecoverCommit(e.group, e.uis, e.pss.Threshold(), e.pss.Num())
}

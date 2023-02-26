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
	Ui kyber.Point
	Ei kyber.Scalar
	Fi kyber.Scalar
}

type ElGamalDealer struct {
	// Embedded PSS Service
	pss pss.Service

	u    kyber.Point      // Encrypted Secret
	pk   crypto.PublicKey // Receiver Public Key
	poly share.PubPoly    //
	uis  []kyber.Point    // reencrypted secret shares

	finished bool
}

// pk: reciever public key
// u: key share commit
func (e *ElGamalDealer) Reencrypt(pk crypto.PublicKey, u kyber.Point) (ReencryptReply, error) {
	// ui = (u ^ ski) * (pk ^ ski)
	// --

	// si = random scalar
	// uiHat = (u * pk) ^ si
	// hiHat = g^si
	// ei = H(ui + uiHat + hiHat)
	// fi = si + (ski * ei)

	return ReencryptReply{}, nil
}

func (e *ElGamalDealer) Process(from pss.Node, r ReencryptReply) error {
	// verify
	ufi := e.pss.Suite().Point().Mul(r.Fi, e.pss.Suite().Point().Add(e.u, e.pk))
	uiei := e.pss.Suite().Point().Mul(e.pss.Suite().Scalar().Neg(r.Ei), r.Ui)
	uiHat := e.pss.Suite().Point().Add(ufi, uiei)

	gfi := e.pss.Suite().Point().Mul(r.Fi, nil)
	gxi := e.poly.Eval(from.Index()).V
	hiei := e.pss.Suite().Point().Mul(e.pss.Suite().Scalar().Neg(r.Ei), gxi)
	hiHat := e.pss.Suite().Point().Add(gfi, hiei)

	// locally produce challenge hash (random oracle)
	hash := sha256.New()
	r.Ui.MarshalTo(hash)
	uiHat.MarshalTo(hash)
	hiHat.MarshalTo(hash)
	c := e.pss.Suite().Scalar().SetBytes(hash.Sum(nil))

	// proof verify
	if c.Equal(r.Ei) {
		e.uis[from.Index()] = r.Ui
		return nil
	}

	return fmt.Errorf("failed verification")
}

func (e *ElGamalDealer) Recover() (kyber.Point, error) {
	// verify length
	//
	// share.Recover(suite,)
}

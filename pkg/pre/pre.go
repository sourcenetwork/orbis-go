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
	Ui share.PubShare
	Ei kyber.Scalar
	Fi kyber.Scalar
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
	// Embedded PSS Service
	pss pss.Service

	u              kyber.Point       // Encrypted Secret
	pk             crypto.PublicKey  // Receiver Public Key
	poly           share.PubPoly     //
	uis            []*share.PubShare // reencrypted secret shares
	verifiedShares *int64

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

// Process verifies an incoming Re-encryption reply from another node.
func (e *ElGamalDealer) Process(from pss.Node, r ReencryptReply) error {
	// verify
	// ufi = (u*pk)^fi
	ufi := e.pss.Suite().Point().Mul(r.Fi, e.pss.Suite().Point().Add(e.u, e.pk))
	// uiei = ui^-ei
	uiei := e.pss.Suite().Point().Mul(e.pss.Suite().Scalar().Neg(r.Ei), r.Ui.V)
	// uihat = uiei ^ ufi
	uiHat := e.pss.Suite().Point().Add(ufi, uiei)

	//
	gfi := e.pss.Suite().Point().Mul(r.Fi, nil)
	//
	gxi := e.poly.Eval(from.Index()).V
	//
	hiei := e.pss.Suite().Point().Mul(e.pss.Suite().Scalar().Neg(r.Ei), gxi)
	//
	hiHat := e.pss.Suite().Point().Add(gfi, hiei)

	// locally produce challenge hash (random oracle)
	hash := sha256.New()
	r.Ui.V.MarshalTo(hash)
	uiHat.MarshalTo(hash)
	hiHat.MarshalTo(hash)
	c := e.pss.Suite().Scalar().SetBytes(hash.Sum(nil))

	// verify
	// H(ui + uiHat + hiHat) == r.Ei
	if c.Equal(r.Ei) {
		e.uis[from.Index()] = &r.Ui
		// atomic add
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

	return share.RecoverCommit(e.pss.Suite(), e.uis, e.pss.Threshold(), e.pss.Num())
}

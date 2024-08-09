package rabin

import (
	"fmt"

	ic "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/share"
	rabindkg "go.dedis.ch/kyber/v3/share/dkg/rabin"
	rabinvss "go.dedis.ch/kyber/v3/share/vss/rabin"
	"go.dedis.ch/kyber/v3/suites"

	rabinv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/rabin/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/crypto/suites/secp256k1"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	orbisdkg "github.com/sourcenetwork/orbis-go/pkg/dkg"
	p2ptransport "github.com/sourcenetwork/orbis-go/pkg/transport/p2p"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

const (
	RABIN_CUSTOM_STATE_MASK orbisdkg.State = 0b01000000

	RECIEVING orbisdkg.State = orbisdkg.CUSTOM_STATE_MASK | RABIN_CUSTOM_STATE_MASK | iota + 1 // 0b10000000
	// Processed all the deals, waiting for responses
	PROCESSED_DEALS     // 0b10000001
	PROCESSED_RESPONSES // 0b10000010
	PROCESSED_COMMITS   // 0b10000011

	// PROCESSING = PROCESSED_DEALS | PROCESSED_RESPONSES
)

var (
	ErrMissingRepoKeys = fmt.Errorf("missing repo keys")

	stateToString = map[orbisdkg.State]string{
		orbisdkg.UNSPECIFIED: orbisdkg.UNSPECIFIED.String(),
		orbisdkg.STARTED:     orbisdkg.STARTED.String(),
		orbisdkg.INITIALIZED: orbisdkg.INITIALIZED.String(),
		orbisdkg.CERTIFIED:   orbisdkg.CERTIFIED.String(),
		RECIEVING:            "Receiving Deals",
		PROCESSED_DEALS:      "Processed Deals",
		PROCESSED_RESPONSES:  "Processed Reponses",
		PROCESSED_COMMITS:    "Processed Commits",
	}
)

func IsRabinState(state orbisdkg.State) bool {
	return (state & RABIN_CUSTOM_STATE_MASK) > 0
}

type Deal = rabinv1alpha1.Deal

// /ringID/nodeID/dealIndex
func dealPkFunc(kb db.KeyBuilder, d *Deal) []byte {
	return kb.AddStringField(d.RingId).
		AddStringField(d.NodeId).
		AddInt32Field(int32(d.Index)).
		AddInt32Field(int32(d.TargetIndex)).
		Bytes()
}

type Response = rabinv1alpha1.Response

func responsePkFunc(kb db.KeyBuilder, d *Response) []byte {
	return kb.AddStringField(d.RingId).
		AddStringField(d.NodeId).
		AddInt32Field(int32(d.Index)).
		AddInt32Field(int32(d.TargetIndex)).
		Bytes()
}

type SecretCommits = rabinv1alpha1.SecretCommits

func secretCommitsPkFunc(kb db.KeyBuilder, d *SecretCommits) []byte {
	return kb.AddStringField(d.RingId).
		AddStringField(d.NodeId).
		AddInt32Field(int32(d.Index)).
		// AddInt32Field(int32(d.TargetIndex)).
		Bytes()
}

func dkgPkFunc(kb db.KeyBuilder, d *rabinv1alpha1.DKG) []byte {
	return kb.AddStringField(d.RingId).Bytes()
}

var (
	// full protocol example: /orbis/0x123/dkg/rabin/send_deal/0.0.1
	ProtocolDeal          protocol.ID = orbisdkg.ProtocolName + "/rabin/deal/0.0.1"
	ProtocolResponse      protocol.ID = orbisdkg.ProtocolName + "/rabin/response/0.0.1"
	ProtocolSecretCommits protocol.ID = orbisdkg.ProtocolName + "/rabin/secretcommits/0.0.1"

	ErrDealNotCertified = fmt.Errorf("dkg: can't give SecretCommits if deal not certified")
	ErrCouldntGetRepo   = fmt.Errorf("dkg: can't get repo")

	DealNamespace          string = "deal"
	ResponseNamespace      string = "response"
	SecretCommitsNamespace string = "secretcommits"
)

func (d *dkg) dealToProto(deal *rabindkg.Deal) (*Deal, error) {
	return dealToProto(deal)
}

func dealToProto(deal *rabindkg.Deal) (*Deal, error) {
	dkheyBytes, err := deal.Deal.DHKey.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal dhkey: %w", err)
	}

	return &Deal{
		TargetIndex: int32(deal.Target),
		Index:       deal.Index,
		Deal: &rabinv1alpha1.EncryptedDeal{
			Dhkey:     dkheyBytes,
			Signature: deal.Deal.Signature,
			Nonce:     deal.Deal.Nonce,
			Cipher:    deal.Deal.Cipher,
		},
	}, nil
}

func (d *dkg) dealFromProto(deal *Deal) (*rabindkg.Deal, error) {
	return dealFromProto(d.suite, deal)
}

func dealFromProto(suite suites.Suite, deal *Deal) (*rabindkg.Deal, error) {
	dhpoint := suite.Point()
	err := dhpoint.UnmarshalBinary(deal.Deal.Dhkey)
	if err != nil {
		return nil, fmt.Errorf("unmarshal dhkey: %w", err)
	}

	return &rabindkg.Deal{
		Index:  deal.Index,
		Target: uint32(deal.TargetIndex),
		Deal: &rabinvss.EncryptedDeal{
			DHKey:     dhpoint,
			Signature: deal.Deal.Signature,
			Nonce:     deal.Deal.Nonce,
			Cipher:    deal.Deal.Cipher,
		},
	}, nil
}

func (d *dkg) responseToProto(response *rabindkg.Response) *Response {
	return &Response{
		Index:       response.Index,
		TargetIndex: response.Target,
		Response: &rabinv1alpha1.VerifiableResponse{
			SessionId: response.Response.SessionID,
			Index:     response.Response.Index,
			Approved:  response.Response.Approved,
			Signature: response.Response.Signature,
		},
	}
}

func (d *dkg) responseFromProto(response *Response) *rabindkg.Response {
	return &rabindkg.Response{
		Index:  response.Index,
		Target: response.TargetIndex,
		Response: &rabinvss.Response{
			SessionID: response.Response.SessionId,
			Index:     response.Response.Index,
			Approved:  response.Response.Approved,
			Signature: response.Response.Signature,
		},
	}
}

func secretCommitsToProto(sc *rabindkg.SecretCommits) (*SecretCommits, error) {
	// convert kyber points
	points := make([][]byte, len(sc.Commitments))
	for i, c := range sc.Commitments {
		cBytes, err := c.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("marshal commitment: %w", err)
		}
		points[i] = cBytes
	}

	return &SecretCommits{
		Index: sc.Index,
		// TargetIndex: sc.TargetIndex,
		Commitments: points,
		SessionId:   sc.SessionID,
		Signature:   sc.Signature,
	}, nil
}

func secretCommitsFromProto(suite suites.Suite, sc *SecretCommits) (*rabindkg.SecretCommits, error) {
	// convert kyber points
	points := make([]kyber.Point, len(sc.Commitments))
	for i, c := range sc.Commitments {
		commitPoint := suite.Point()
		err := commitPoint.UnmarshalBinary(c)
		if err != nil {
			return nil, fmt.Errorf("unmarshal commitment: %w", err)
		}
		points[i] = commitPoint
	}

	return &rabindkg.SecretCommits{
		Index: sc.Index,
		// TargetIndex: sc.TargetIndex,
		Commitments: points,
		SessionID:   sc.SessionId,
		Signature:   sc.Signature,
	}, nil
}

func dkgToProto(d *dkg) (*rabinv1alpha1.DKG, error) {
	var suiteType rabinv1alpha1.SuiteType
	if d.suite != nil {
		switch d.suite.String() {
		case "Ed25519":
			suiteType = rabinv1alpha1.SuiteType_Ed25519
		case "Secp256k1":
			suiteType = rabinv1alpha1.SuiteType_Secp256k1
		default:
			return nil, fmt.Errorf("invalid suite type: %v", d.suite.String())
		}
	}

	var state rabinv1alpha1.State
	switch d.state {
	case orbisdkg.UNSPECIFIED:
		state = rabinv1alpha1.State_STATE_UNSPECIFIED
	case orbisdkg.INITIALIZED:
		state = rabinv1alpha1.State_STATE_INITIALIZED
	case orbisdkg.STARTED:
		state = rabinv1alpha1.State_STATE_STARTED
	case orbisdkg.CERTIFIED:
		state = rabinv1alpha1.State_STATE_CERTIFIED
	case RECIEVING:
		state = rabinv1alpha1.State_STATE_RECEVING
	case PROCESSED_DEALS:
		state = rabinv1alpha1.State_STATE_PROCESSED_DEALS
	case PROCESSED_RESPONSES:
		state = rabinv1alpha1.State_STATE_PROCESSED_RESPONSES
	case PROCESSED_COMMITS:
		state = rabinv1alpha1.State_STATE_PROCESSED_COMMITS
	default:
		return nil, fmt.Errorf("invalid state: %v, 0x%0x", d.state, d.state)
	}

	var nodes []*rabinv1alpha1.Node
	if d.participants != nil {
		nodes = make([]*rabinv1alpha1.Node, len(d.participants))
		for i, p := range d.participants {
			pk, err := ic.PublicKeyToProto(crypto.ToLibP2PPublicKey(p.PublicKey()))
			if err != nil {
				return nil, fmt.Errorf("couldnt convert public key to proto: %w", err)
			}
			nodes[i] = &rabinv1alpha1.Node{
				Id:        p.ID(),
				Address:   p.Address().String(),
				PublicKey: pk,
			}
		}
	}

	var distPubKey []byte
	var err error
	if d.distPubKey != nil {
		distPubKey, err = d.distPubKey.MarshalBinary()
		if err != nil {

			return nil, fmt.Errorf("couldn't marshal pubkey: %w", err)
		}
	}

	var prishare *rabinv1alpha1.PriShare
	share := d.distKeyShare.PriShare
	if share != nil {
		sbuf, err := share.V.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("couldn't marshal private share: %w", err)
		}
		prishare = &rabinv1alpha1.PriShare{
			Index: int32(share.I),
			V:     sbuf,
		}
	}

	var shareCommits [][]byte
	commits := d.distKeyShare.Commits
	if commits != nil {
		shareCommits = make([][]byte, len(commits))
		for i, cmt := range commits {
			buf, err := cmt.MarshalBinary()
			if err != nil {
				return nil, fmt.Errorf("couldn't marshal share commitment: %w", err)
			}
			shareCommits[i] = buf
		}
	}

	// polynomials and secrets
	var fPoly *rabinv1alpha1.PriPoly
	var gPoly *rabinv1alpha1.PriPoly
	var polySecret []byte
	if d.rdkg != nil {
		fPoly, err = polyToProto(d.rdkg.Dealer().FPoly())
		if err != nil {
			return nil, err
		}
		gPoly, err = polyToProto(d.rdkg.Dealer().GPoly())
		if err != nil {
			return nil, err
		}
		if secret := d.rdkg.Dealer().Secret(); secret != nil {
			polySecret, err = d.rdkg.Dealer().Secret().MarshalBinary()
			if err != nil {
				return nil, err
			}
		}
	}

	return &rabinv1alpha1.DKG{
		RingId:       string(d.ringID),
		Index:        int32(d.index),
		Num:          d.num,
		Threshold:    d.threshold,
		Suite:        suiteType,
		State:        state,
		Nodes:        nodes,
		Pubkey:       distPubKey,
		PriShare:     prishare,
		ShareCommits: shareCommits,
		F:            fPoly,
		G:            gPoly,
		PolySecret:   polySecret,
	}, nil
}

func polyToProto(p *share.PriPoly) (*rabinv1alpha1.PriPoly, error) {
	poly := &rabinv1alpha1.PriPoly{
		Coeffs: make([][]byte, p.Threshold()),
	}
	var err error
	for i, coeffs := range p.Coefficients() {
		poly.Coeffs[i], err = coeffs.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("couldn't marshal poly coeffecient: %w", err)
		}
	}
	return poly, nil
}

func polyFromProto(suite suites.Suite, p *rabinv1alpha1.PriPoly) (*share.PriPoly, error) {
	scalars := make([]kyber.Scalar, len(p.Coeffs))
	for i, coeff := range p.Coeffs {
		scalarCoeff := suite.Scalar()
		err := scalarCoeff.UnmarshalBinary(coeff)
		if err != nil {
			return nil, fmt.Errorf("couldn't unmarshal poly coeffecient")
		}
		scalars[i] = scalarCoeff
	}
	return share.CoefficientsToPriPoly(suite, scalars), nil
}

func dkgFromProto(d *rabinv1alpha1.DKG) (dkg, error) {
	var suite suites.Suite
	switch d.Suite {
	case rabinv1alpha1.SuiteType_Ed25519:
		suite = edwards25519.NewBlakeSHA256Ed25519()
	case rabinv1alpha1.SuiteType_Secp256k1:
		suite = secp256k1.NewBlakeKeccackSecp256k1()
	default:
		return dkg{}, fmt.Errorf("bad key type: %v", d.Suite.String())
	}

	var state orbisdkg.State
	switch d.State {
	case rabinv1alpha1.State_STATE_UNSPECIFIED:
		state = orbisdkg.UNSPECIFIED
	case rabinv1alpha1.State_STATE_INITIALIZED:
		state = orbisdkg.INITIALIZED
	case rabinv1alpha1.State_STATE_STARTED:
		state = orbisdkg.STARTED
	case rabinv1alpha1.State_STATE_CERTIFIED:
		state = orbisdkg.CERTIFIED
	case rabinv1alpha1.State_STATE_RECEVING:
		state = RECIEVING
	case rabinv1alpha1.State_STATE_PROCESSED_DEALS:
		state = PROCESSED_DEALS
	case rabinv1alpha1.State_STATE_PROCESSED_RESPONSES:
		state = PROCESSED_RESPONSES
	case rabinv1alpha1.State_STATE_PROCESSED_COMMITS:
		state = PROCESSED_COMMITS
	}

	participants := make([]orbisdkg.Node, len(d.Nodes))
	for i, n := range d.Nodes {
		pk, err := crypto.PublicKeyFromProto(n.PublicKey)
		if err != nil {
			return dkg{}, fmt.Errorf("couldnt convert proto to public key: %w", err)
		}
		addr, err := ma.NewMultiaddr(n.Address)
		if err != nil {
			return dkg{}, fmt.Errorf("invalid address: %w", err)
		}
		participants[i] = p2ptransport.NewNode(n.Id, pk, addr)
	}

	// pubkey
	var distPubKey kyber.Point
	if d.Pubkey != nil {
		distPubKey = suite.Point()
		log.Debug("Unmarshalling distPubKey:", d.Pubkey)
		err := distPubKey.UnmarshalBinary(d.Pubkey)
		if err != nil {
			return dkg{}, fmt.Errorf("unmarshaling pubkey: %w", err)
		}
	}

	// share
	var distKeyShare crypto.DistKeyShare
	if d.PriShare != nil {
		distKeyShare.PriShare = new(share.PriShare)
		distKeyShare.PriShare.I = int(d.PriShare.Index)
		distKeyShare.PriShare.V = suite.Scalar()
		err := distKeyShare.PriShare.V.UnmarshalBinary(d.PriShare.V)
		if err != nil {
			return dkg{}, fmt.Errorf("unmarshaling prishare: %w", err)
		}
	}

	if d.ShareCommits != nil {
		distKeyShare.Commits = make([]kyber.Point, len(d.ShareCommits))
		for i, cmt := range d.ShareCommits {
			distKeyShare.Commits[i] = suite.Point()
			err := distKeyShare.Commits[i].UnmarshalBinary(cmt)
			if err != nil {
				return dkg{}, fmt.Errorf("unmarshalling commits: %w", err)
			}
		}
	}

	// f and g polys
	fPoly, err := polyFromProto(suite, d.F)
	if err != nil {
		return dkg{}, err
	}
	gPoly, err := polyFromProto(suite, d.G)
	if err != nil {
		return dkg{}, err
	}

	// polynomial secret
	var secret kyber.Scalar
	if d.PolySecret != nil {
		secret = suite.Scalar()
		err := secret.UnmarshalBinary(d.PolySecret)
		if err != nil {
			return dkg{}, fmt.Errorf("unmarshaling secret: %w", err)
		}
	}

	return dkg{
		ringID:       types.RingID(d.RingId),
		index:        int(d.Index),
		num:          d.Num,
		threshold:    d.Threshold,
		suite:        suite,
		state:        state,
		participants: participants,
		distPubKey:   distPubKey,
		distKeyShare: distKeyShare,
		fPoly:        fPoly,
		gPoly:        gPoly,
		secret:       secret,
	}, nil
}

type dealDispatch struct {
	err  chan error
	deal *rabindkg.Deal
}

type responseDispatch struct {
	err     chan error
	respone *rabindkg.Response
}

type secretCommitsDispatch struct {
	err           chan error
	secretCommits *rabindkg.SecretCommits
}

package rabin

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/protocol"
	"go.dedis.ch/kyber/v3"
	rabindkg "go.dedis.ch/kyber/v3/share/dkg/rabin"
	rabinvss "go.dedis.ch/kyber/v3/share/vss/rabin"
	"go.dedis.ch/kyber/v3/suites"

	rabinv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/rabin/v1alpha1"
	orbisdkg "github.com/sourcenetwork/orbis-go/pkg/dkg"
)

type Deal = rabinv1alpha1.Deal
type Response = rabinv1alpha1.Response
type SecretCommits = rabinv1alpha1.SecretCommits

var (
	// full protocol example: /orbis/0x123/dkg/rabin/send_deal/0.0.1
	ProtocolDeal          protocol.ID = orbisdkg.ProtocolName + "/rabin/deal/0.0.1"
	ProtocolResponse      protocol.ID = orbisdkg.ProtocolName + "/rabin/response/0.0.1"
	ProtocolSecretCommits protocol.ID = orbisdkg.ProtocolName + "/rabin/secretcommits/0.0.1"

	ErrDealNotCertified = fmt.Errorf("dkg: can't give SecretCommits if deal not certified")
)

func (d *dkg) dealToProto(deal *rabindkg.Deal) (*Deal, error) {
	return dealToProto(deal)
}

func dealToProto(deal *rabindkg.Deal) (*Deal, error) {
	dkheyBytes, err := deal.Deal.DHKey.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return &Deal{
		Index: deal.Index,
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
	if err := dhpoint.UnmarshalBinary(deal.Deal.Dhkey); err != nil {
		return nil, err
	}

	return &rabindkg.Deal{
		Index: deal.Index,
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
		Index: response.Index,
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
		Index: response.Index,
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
			return nil, err
		}
		points[i] = cBytes
	}

	return &SecretCommits{
		Index:       sc.Index,
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
		if err := commitPoint.UnmarshalBinary(c); err != nil {
			return nil, err
		}
		points[i] = commitPoint
	}

	return &rabindkg.SecretCommits{
		Index:       sc.Index,
		Commitments: points,
		SessionID:   sc.SessionId,
		Signature:   sc.Signature,
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

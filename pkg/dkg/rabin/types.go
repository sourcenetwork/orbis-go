package rabin

import (
	"github.com/libp2p/go-libp2p/core/protocol"
	rabindkg "go.dedis.ch/kyber/v3/share/dkg/rabin"
	rabinvss "go.dedis.ch/kyber/v3/share/vss/rabin"

	rabinv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/rabin/v1alpha1"
	orbisdkg "github.com/sourcenetwork/orbis-go/pkg/dkg"
)

type Deal = rabinv1alpha1.Deal
type Response = rabinv1alpha1.Response

var (
	// full protocol example: /orbis/0x123/dkg/rabin/send_deal/0.0.1
	ProtocolDeal     protocol.ID = orbisdkg.ProtocolName + "/rabin/deal/0.0.1"
	ProtocolResponse protocol.ID = orbisdkg.ProtocolName + "/rabin/response/0.0.1"
)

func (d *dkg) dealToProto(deal *rabindkg.Deal) (*Deal, error) {
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
	dhpoint := d.suite.Point()
	if err := dhpoint.UnmarshalBinary(deal.Deal.Dhkey); err != nil {
		return nil, err
	}

	return &rabindkg.Deal{
		Index: deal.Index,
		Deal: &rabinvss.EncryptedDeal{
			DHKey:     dhpoint,
			Signature: deal.Deal.Signature,
			Nonce:     deal.Deal.Nonce,
			Cipher:    deal.Deal.Signature,
		},
	}, nil
}

func (d *dkg) responseToProto(response *rabindkg.Response) *Response {
	return &Response{
		Index: response.Index,
		Response: &rabinv1alpha1.VerifiableResponse{
			SessionID: response.Response.SessionID,
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
			SessionID: response.Response.SessionID,
			Index:     response.Response.Index,
			Approved:  response.Response.Approved,
			Signature: response.Response.Signature,
		},
	}
}

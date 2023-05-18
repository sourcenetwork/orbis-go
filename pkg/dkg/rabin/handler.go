package rabin

import (
	"context"
	"fmt"

	rabindkg "go.dedis.ch/kyber/v3/share/dkg/rabin"
	"go.dedis.ch/protobuf"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	orbisdkg "github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

func (d *dkg) setupHandlers() {
	// deal
	d.transport.AddHandler(ProtocolDeal, d.ProcessMessage)

	// response
	d.transport.AddHandler(ProtocolResponse, d.ProcessMessage)
}

func (d *dkg) processDeal(deal *rabindkg.Deal, nodes []transport.Node) error {
	response, err := d.rdkg.ProcessDeal(deal)
	if err != nil {
		return fmt.Errorf("process deal: %w", err)
	}

	for _, node := range nodes {
		if d.isMe(node) {
			continue // skip ourselves
		}

		buf, err := protobuf.Encode(response)
		if err != nil {
			return fmt.Errorf("encode response: %w", err)
		}

		// todo: context
		if err := d.send(context.TODO(), string(ProtocolResponse), buf, node); err != nil {
			return fmt.Errorf("send response: %w", err)
		}
	}

	_ = response

	return nil
}

func (d *dkg) processResponse(resp *rabindkg.Response) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// we cant process the response unless we
	// have processed the cooresponding deal
	//
	// theres a change that we missed it from the p2p
	// network or bulletin board.
	//
	// For now, lets just design it assuming all is well (temp)
	_, err := d.rdkg.ProcessResponse(resp)
	if err != nil {
		return fmt.Errorf("process response: %w", err)
	}

	if d.rdkg.Certified() {
		// interpolate shared public key
		distkey, err := d.rdkg.DistKeyShare()
		if err != nil {
			return fmt.Errorf("dist key share: %w", err)
		}

		d.share = crypto.PriShare{
			PriShare: distkey.PriShare(),
		}

		d.pubKey = distkey.Public()
		d.state = orbisdkg.CERTIFIED
	}

	return nil
}

func (d *dkg) isMe(node transport.Node) bool {
	return d.transport.Host().ID() == node.ID()
}

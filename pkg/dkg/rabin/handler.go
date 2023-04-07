package rabin

import (
	"context"

	rabindkg "go.dedis.ch/kyber/v3/share/dkg/rabin"
	"go.dedis.ch/protobuf"

	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

func (d *dkg) setupHandlers() {
	// deal
	d.transport.AddHandler(ProtocolRabinDeal, d.ProcessMessage)

	// response
	d.transport.AddHandler(ProtocolRabinResponse, d.ProcessMessage)
}

func (d *dkg) handleDeal(deal *rabindkg.Deal, nodes []transport.Node) error {
	response, err := d.rdkg.ProcessDeal(deal)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		if d.isMe(node) {
			continue // skip ourselves
		}

		buf, err := protobuf.Encode(deal)
		if err != nil {
			return err
		}

		// todo: context
		if err := d.send(context.TODO(), string(ProtocolRabinResponse), buf, node); err != nil {
			return err
		}
	}

	return nil
}

func (d *dkg) isMe(node transport.Node) bool {
	return d.transport.Host().ID() == node.ID()
}

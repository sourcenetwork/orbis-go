package rabin

import (
	"context"
	"fmt"

	rabindkg "go.dedis.ch/kyber/v3/share/dkg/rabin"
	"go.dedis.ch/protobuf"
	"google.golang.org/protobuf/proto"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	orbisdkg "github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

func (d *dkg) setupHandlers() {
	// deal
	d.transport.AddHandler(ProtocolDeal, d.ProcessMessage)

	// response
	d.transport.AddHandler(ProtocolResponse, d.ProcessMessage)

	// secretcommits
	d.transport.AddHandler(ProtocolSecretCommits, d.ProcessMessage)
}

func (d *dkg) processDeal(deal *rabindkg.Deal) error {
	response, err := d.rdkg.ProcessDeal(deal)
	if err != nil {
		return fmt.Errorf("process rabin dkg deal: %w", err)
	}

	buf, err := protobuf.Encode(response)
	if err != nil {
		return fmt.Errorf("encode response: %w", err)
	}

	for _, node := range d.participants {
		if d.isMe(node) {
			log.Infof("skipping self: %s", node.ID())
			continue // skip ourselves
		}

		// todo: context
		if err := d.send(context.TODO(), string(ProtocolResponse), buf, node); err != nil {
			return fmt.Errorf("send response: %w", err)
		}
	}

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
	just, err := d.rdkg.ProcessResponse(resp)
	if err != nil {
		return fmt.Errorf("process response: %w", err)
	}

	if just != nil {
		log.Warnf("Got justification during response process for %d: %v", resp.Index, just)
		return nil
	}

	if !d.rdkg.Certified() {
		return nil
	}

	sc, err := d.rdkg.SecretCommits()
	if err != nil && err.Error() == ErrDealNotCertified.Error() {
		return nil // skip
	} else if err != nil {
		return fmt.Errorf("generate secret commit dkg response: %w", err)
	}

	protoSC, err := secretCommitsToProto(sc)
	if err != nil {
		return fmt.Errorf("secret commits to proto: %w", err)
	}

	buf, err := proto.Marshal(protoSC)
	if err != nil {
		return fmt.Errorf("encode response: %w", err)
	}

	// send SC
	for _, node := range d.participants {
		if d.isMe(node) {
			log.Infof("skipping self: %s", node.ID())
			continue // skip ourselves
		}

		// todo: context
		if err := d.send(context.TODO(), string(ProtocolSecretCommits), buf, node); err != nil {
			return fmt.Errorf("send secret commits: %w", err)
		}
	}

	return nil
}

func (d *dkg) processSecretCommits(sc *rabindkg.SecretCommits) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.rdkg.ProcessSecretCommits(sc)
	if err != nil {
		return fmt.Errorf("process rabin dkg secretcommits: %w", err)
	}

	// If we haven't collected all deals, responses, and secret commits
	// then we can't compute the dist key share
	//
	// Also, if we've already completed the dkg setup, then we
	// can also skip
	if !d.rdkg.Finished() || d.state == orbisdkg.CERTIFIED {
		return nil
	}

	// interpolate shared public key
	distkey, err := d.rdkg.DistKeyShare()
	if err != nil {
		return fmt.Errorf("rabin dkg dist key share: %w", err)
	}

	d.share = crypto.PriShare{
		PriShare: distkey.PriShare(),
	}

	d.pubKey = distkey.Public()
	d.state = orbisdkg.CERTIFIED
	log.Infof("Finished DKG Setup. Shared Public Key: %s", d.pubKey)

	return d.save(context.TODO())
}

func (d *dkg) isMe(node transport.Node) bool {
	return d.transport.Host().ID() == node.ID()
}

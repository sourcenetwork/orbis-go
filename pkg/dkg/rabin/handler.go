package rabin

import (
	"context"
	"fmt"

	rabindkg "go.dedis.ch/kyber/v3/share/dkg/rabin"
	"go.dedis.ch/protobuf"
	"google.golang.org/protobuf/proto"

	"github.com/sourcenetwork/eventbus-go"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	orbisdkg "github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

func (d *dkg) setupHandlers() error {
	bus := d.bulletin.Events()
	var err error
	d.eventsCh, err = eventbus.Subscribe[bulletin.Event](bus)
	if err != nil {
		return err
	}

	go func() {
		for evt := range d.eventsCh {
			log.Debugf("recieved eventbus on %s from %s for %s (%s)", d.transport.Host().ID(), evt.Message.NodeId, evt.Message.TargetId, evt.Message.GetType())
			if evt.Message.TargetId != d.transport.Host().ID() {
				log.Debugf("ignoring bulletin event not for us")
				continue
			}

			// process in a dedicated goroutine so we dont block
			go func(evt bulletin.Event) {
				err := d.ProcessMessage(evt.Message)
				if err != nil {
					log.Errorf("processing bulletin message %s: %v", evt.Message.GetType(), err)
				}
			}(evt)
		}
	}()

	return nil
}

func (d *dkg) processDeal(deal *rabindkg.Deal) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	ctx := context.TODO()

	log.Debugf("processing deal, index=%v nonce=%0x sig=%0x", deal.Index, deal.Deal.Nonce, deal.Deal.Signature)
	response, err := d.rdkg.ProcessDeal(deal)
	if err != nil {
		return fmt.Errorf("process rabin dkg deal: %w", err)
	} else {
		log.Debugf("succesfully processed deal %0x", deal.Deal.Signature)
	}

	buf, err := protobuf.Encode(response)
	if err != nil {
		return fmt.Errorf("encode response: %w", err)
	}

	for _, node := range d.participants {
		if d.isMe(node) {
			continue // skip ourselves
		}
		// TODO: can we skip the sender of the deal as well?

		log.Debugf("sending response on %s for %s", d.transport.Host().ID(), node.ID())
		// /ring/<ringID>/dkg/rabin/<action>/<fromID>/<toID>

		// each node generates deals [d0, d1, d2]
		// each node processes each deal [d1, d2, d2]
		// each node creates response for (dN, SELF, TARGET)

		// /ring/<ringID>/dkg/rabin/RESPONSE/JOHN/ROY/<FOR>
		forNodeID1 := d.participants[response.Index].ID()
		// forNodeID2 := d.participants[response.Target].ID()
		log.Debugf("creating identifier from %s to %s for %s", d.NodeID(), node.ID(), forNodeID1)
		msgID := fmt.Sprintf("/%s/%s/%s/%s", ResponseNamespace, d.NodeID(), node.ID(), forNodeID1)
		if err := d.post(ctx, ResponseNamespace, d.bbnamespace, msgID, buf, node); err != nil {
			return fmt.Errorf("send response: %w", err)
		}
	}

	return nil
}

func (d *dkg) processResponse(resp *rabindkg.Response) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// we can't process the response unless we
	// have processed the cooresponding deal
	//
	// theres a chance that we missed it from the p2p
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
	if err != nil {
		if err.Error() == ErrDealNotCertified.Error() {
			return nil
		}
		return fmt.Errorf("generate secret commit: %w", err)
	}

	protoSC, err := secretCommitsToProto(sc)
	if err != nil {
		return fmt.Errorf("secret commits to proto: %w", err)
	}

	buf, err := proto.Marshal(protoSC)
	if err != nil {
		return fmt.Errorf("encode response: %w", err)
	}

	for i, node := range d.participants {
		if d.isMe(node) {
			continue
		}

		// /ring/<ringID>/dkg/rabin/RESPONSE/JOHN/ROY/<FOR>
		forNodeID1 := d.participants[sc.Index].ID()
		// forNodeID2 := d.participants[response.Target].ID()
		log.Debugf("Node %d sending secret commits to pariticipant %d", d.index, i)
		msgID := fmt.Sprintf("/%s/%s/%s/%s", SecretCommitsNamespace, d.NodeID(), d.participants[i].ID(), forNodeID1)
		err := d.post(context.TODO(), SecretCommitsNamespace, d.bbnamespace, msgID, buf, node)
		if err != nil {
			return fmt.Errorf("send secret commits: %w", err)
		}
	}

	return nil
}

func (d *dkg) processSecretCommits(sc *rabindkg.SecretCommits) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	log.Debugf("Node %d processing secret commits", d.index)
	_, err := d.rdkg.ProcessSecretCommits(sc)
	if err != nil {
		return fmt.Errorf("process rabin dkg secret commits: %w", err)
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

	d.distKeyShare = crypto.DistKeyShare{
		Commits:  distkey.Commitments(),
		PriShare: distkey.PriShare(),
	}

	d.pubKey = distkey.Public()
	d.state = orbisdkg.CERTIFIED

	log.Infof("Node %d finished setup with shared publick Key: %s", d.index, d.pubKey)
	return d.save(context.TODO())
}

func (d *dkg) isMe(node transport.Node) bool {
	return d.transport.Host().ID() == node.ID()
}

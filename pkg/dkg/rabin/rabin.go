package rabin

import (
	"context"
	"fmt"
	"sync"

	logging "github.com/ipfs/go-log"
	"go.dedis.ch/kyber/v3"
	rabindkg "go.dedis.ch/kyber/v3/share/dkg/rabin"
	"go.dedis.ch/kyber/v3/suites"
	"google.golang.org/protobuf/proto"

	rabinv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/rabin/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	orbisdkg "github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

var log = logging.Logger("orbis/dkg/rabin")

const name = "rabin"

type dkg struct {
	mu sync.Mutex

	ringID types.RingID

	rdkg         *rabindkg.DistKeyGenerator
	participants []orbisdkg.Node

	index     int
	num       int32
	threshold int32

	suite   suites.Suite
	pubKey  kyber.Point
	privKey kyber.Scalar

	dealRepo db.Repository[*rabinv1alpha1.Deal]
	respRepo db.Repository[*rabinv1alpha1.Response]

	deals     chan dealDispatch
	responses chan responseDispatch
	commits   chan secretCommitsDispatch

	share crypto.PriShare

	db        *db.DB
	transport transport.Transport
	bulletin  bulletin.Bulletin

	state orbisdkg.State
}

func New(repo *db.DB, rkeys []db.RepoKey, t transport.Transport, b bulletin.Bulletin) (*dkg, error) {

	//dealsRepo, err := db.GetRepo[db.Record](repo, rkeys[0])
	//sharesRepo, err := db.GetRepo[db.Record](repo, rkeys[1])
	return &dkg{
		db:        repo,
		transport: t,
		bulletin:  b,
		index:     -1,
	}, nil
}

// Init initializes the DKG with the target nodes
func (d *dkg) Init(ctx context.Context, pk crypto.PrivateKey, nodes []orbisdkg.Node, n int32, threshold int32) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if pk == nil {
		return fmt.Errorf("missing private key")
	}

	suite, err := crypto.SuiteForType(pk.Type())
	if err != nil {
		return err
	}

	d.suite = suite
	d.privKey = pk.Scalar()
	d.pubKey = suite.Point().Mul(d.privKey, nil) // public point for scalar
	d.num = n
	d.threshold = threshold

	if len(nodes) != int(n) {
		return orbisdkg.ErrBadNodeSet
	}

	points := make([]kyber.Point, 0, len(nodes))
	for i, n := range nodes {
		point := n.PublicKey().Point()
		if point.Equal(d.pubKey) {
			d.index = i
		}
		points = append(points, point)
	}

	// we didn't find ourselves in the list
	if d.index == -1 {
		return orbisdkg.ErrMissingSelf
	}

	rdkg, err := rabindkg.NewDistKeyGenerator(d.suite, d.privKey, points, int(d.threshold))
	if err != nil {
		return fmt.Errorf("create DKG: %w", err)
	}

	d.participants = nodes
	d.rdkg = rdkg

	// setup stream handler for transport
	d.setupHandlers()
	d.state = orbisdkg.INITIALIZED

	d.deals = make(chan dealDispatch, d.numExpectedDeals())
	d.responses = make(chan responseDispatch, d.numExpectedResponses())
	d.commits = make(chan secretCommitsDispatch, d.numExpectedCommits())

	return nil
}

func (d *dkg) Name() string {
	return name
}

func (d *dkg) PublicKey() crypto.PublicKey {
	pk, _ := crypto.PublicKeyFromPoint(d.pubKey)
	return pk
}

func (d *dkg) Share() crypto.PriShare {
	return d.share
}

func (d *dkg) State() orbisdkg.State {
	return d.state
}

// Start the DKG setup process.
func (d *dkg) Start(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.state = orbisdkg.STARTED

	// TODO
	// if !d.initialized {
	// 	return orbisdkg.ErrNotInitialized
	// }

	deals, err := d.rdkg.Deals()
	if err != nil {
		return fmt.Errorf("generate deals: %w", err)
	}

	for i, deal := range deals {
		dealproto, err := d.dealToProto(deal)
		if err != nil {
			return fmt.Errorf("convert deal to proto: %w", err)
		}

		// TODO: save deals to db
		// if err := d.deals.Create(ctx, dealproto); err != nil {
		// 	return fmt.Errorf("create deal: %w", err)
		// }

		log.Infof("node %d sending deal to partitipants %d (%x)", d.index, i, deal.Deal.Signature)
		if i == d.index {
			// TODO: deliver to ourselves
			continue
		}

		buf, err := proto.Marshal(dealproto)
		if err != nil {
			return fmt.Errorf("marshal deal: %w", err)
		}

		err = d.send(ctx, string(ProtocolDeal), buf, d.participants[i])
		if err != nil {
			return fmt.Errorf("send deal: %w", err)
		}
	}

	go d.dispatch()
	return nil
}

func (d *dkg) send(ctx context.Context, msgType string, buf []byte, node transport.Node) error {
	cid, err := types.CidFromBytes(buf)
	if err != nil {
		return fmt.Errorf("cid from bytes: %w", err)
	}
	msg, err := d.transport.NewMessage(d.ringID, cid.String(), false, buf, msgType)
	if err != nil {
		return fmt.Errorf("new message: %w", err)
	}
	log.Infof("dkg.send() node id: %s, addr: %s", node.ID(), node.Address())
	if err := d.transport.Send(ctx, node, msg); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (d *dkg) Close(_ context.Context) error {
	panic("not implemented") // TODO: Implement
}

func (d *dkg) ProcessMessage(msg *transport.Message) error {
	// todo maybe?: validate msg.PublicKey matches payload pubkeys

	switch msg.GetType() {
	case string(ProtocolDeal):
		log.Infof("dkg.ProcessMessage() ProtocolDeal: id: %s", msg.Id)

		var protoDeal rabinv1alpha1.Deal

		err := proto.Unmarshal(msg.Payload, &protoDeal)
		if err != nil {
			return fmt.Errorf("unmarshal deal message: %w", err)
		}

		deal, err := d.dealFromProto(&protoDeal)
		if err != nil {
			return fmt.Errorf("deal message from proto: %w", err)
		}

		log.Infof("dkg.ProcessMessage() node %d process deal %d (%x)", d.index, deal.Index, deal.Deal.Signature)
		err = d.dispatchDeal(deal)
		if err != nil {
			return fmt.Errorf("process deal message: %w", err)
		}

	case string(ProtocolResponse):
		log.Infof("dkg.ProcessMessage() ProtocolResponse: id: %s", msg.Id)

		var protoResponse rabinv1alpha1.Response

		err := proto.Unmarshal(msg.Payload, &protoResponse)
		if err != nil {
			return fmt.Errorf("unmarshal response message: %w", err)
		}

		resp := d.responseFromProto(&protoResponse)

		err = d.dispatchResponse(resp)
		if err != nil {
			return fmt.Errorf("process response message: %w", err)
		}
	case string(ProtocolSecretCommits):
		log.Infof("dkg.ProcessMessage() ProtocolSecretCommits: id: %s", msg.Id)
		var protoSecretCommits rabinv1alpha1.SecretCommits
		err := proto.Unmarshal(msg.Payload, &protoSecretCommits)
		if err != nil {
			return fmt.Errorf("unmarshal secret commits: %w", err)
		}

		sc, err := secretCommitsFromProto(d.suite, &protoSecretCommits)
		if err != nil {
			return fmt.Errorf("secret commits from proto: %w", err)
		}

		err = d.dispatchSecretCommit(sc)
		if err != nil {
			return fmt.Errorf("process secret commits: %w", err)
		}

	default:
		panic("bad message type") //todo
	}

	return nil
}

// dispatch is responsible for handling all the incoming
// messages, and dispatching them to their cooresponding
// handlers, but with ordering via channels. This gurantees
// that we handle all the events at their appropriate
// time.
//
// It is designed to run in a gourinte
func (d *dkg) dispatch() {
	// processDeals
	log.Infof("handling %d expected deals", d.numExpectedDeals())
	for i := 0; i < d.numExpectedDeals(); i++ {
		dd := <-d.deals
		log.Infof("%d: handling deal %d", i, dd.deal.Index)
		dd.err <- d.processDeal(dd.deal)
	}
	close(d.deals)

	// processResponses
	log.Infof("handling %d expected responses", d.numExpectedResponses())
	for i := 0; i < d.numExpectedResponses(); i++ {
		rd := <-d.responses
		log.Infof("%d: handling response %d", i, rd.respone.Index)
		rd.err <- d.processResponse(rd.respone)
	}
	close(d.responses)

	// processSecrets
	log.Infof("handling %d expected secrets", d.numExpectedCommits())
	for i := 0; i < d.numExpectedCommits(); i++ {
		sd := <-d.commits
		log.Infof("%d: handling secret %d", i, sd.secretCommits.Index)
		sd.err <- d.processSecretCommits(sd.secretCommits)
	}
	close(d.commits)
}

func (d *dkg) dispatchDeal(deal *rabindkg.Deal) error {
	dealDispatchEvent := dealDispatch{
		err:  make(chan error),
		deal: deal,
	}
	d.deals <- dealDispatchEvent   // send
	return <-dealDispatchEvent.err // recieve
}

func (d *dkg) dispatchResponse(resp *rabindkg.Response) error {
	respDispatchEvent := responseDispatch{
		err:     make(chan error),
		respone: resp,
	}
	d.responses <- respDispatchEvent // send
	return <-respDispatchEvent.err   // recieve
}

func (d *dkg) dispatchSecretCommit(sc *rabindkg.SecretCommits) error {
	scDispatchEvent := secretCommitsDispatch{
		err:           make(chan error),
		secretCommits: sc,
	}
	d.commits <- scDispatchEvent // send
	return <-scDispatchEvent.err // recieve
}

func (d *dkg) numExpectedDeals() int {
	return len(d.participants) - 1
}

func (d *dkg) numExpectedResponses() int {
	l := len(d.participants)
	return (l - 1) * (l - 1)
}

func (d *dkg) numExpectedCommits() int {
	l := len(d.participants)
	return (l - 1) * (l - 1)
}

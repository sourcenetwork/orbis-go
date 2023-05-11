package rabin

import (
	"context"
	"sync"

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

const name = "rabin"

type dkg struct {
	mu sync.Mutex

	ringID types.RingID

	rdkg         *rabindkg.DistKeyGenerator
	participants []orbisdkg.Node

	index     int
	num       int
	threshold int

	suite   suites.Suite
	pubKey  kyber.Point
	privKey kyber.Scalar

	deals db.Repository[*rabinv1alpha1.Deal]
	resps db.Repository[*rabinv1alpha1.Response]

	share crypto.PriShare

	db        *db.DB
	transport transport.Transport
	bulletin  bulletin.Bulletin

	state       orbisdkg.State
	initialized bool
}

func New(repo *db.DB, rkeys []*db.RepoKey, t transport.Transport, b bulletin.Bulletin) (*dkg, error) {

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
func (d *dkg) Init(ctx context.Context, pk crypto.PrivateKey, nodes []orbisdkg.Node, n int, threshold int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	suite, err := crypto.SuiteForType(pk.Type())
	if err != nil {
		return err
	}

	d.suite = suite
	d.privKey = pk.Scalar()
	d.pubKey = suite.Point().Mul(d.privKey, nil) // public point for scalar

	if len(nodes) != n {
		return orbisdkg.ErrBadNodeSet
	}

	points := make([]kyber.Point, len(nodes))
	for i, n := range nodes {
		points[i] = n.PublicKey().Point()
		if points[i] == d.pubKey {
			d.index = i
		}
	}

	// we didn't find ourselves in the list
	if d.index == -1 {
		return orbisdkg.ErrMissingSelf
	}

	rdkg, err := rabindkg.NewDistKeyGenerator(d.suite, d.privKey, points, d.threshold)
	if err != nil {
		return err
	}

	d.participants = nodes
	d.rdkg = rdkg
	d.num = n
	d.threshold = threshold

	// setup stream handler for transport
	d.setupHandlers()
	d.state = orbisdkg.INITIALIZED

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

	if !d.initialized {
		return orbisdkg.ErrNotInitialized
	}

	// generate deals
	// send deals to participants (and self)
	deals, err := d.rdkg.Deals()
	if err != nil {
		return err
	}
	dealprotos := make([]*Deal, len(deals))
	i := 0
	for _, deal := range deals {
		dealprotos[i], err = d.dealToProto(deal)
		if err != nil {
			return err
		}
		if err := d.deals.Create(ctx, dealprotos[i]); err != nil {
			return err
		}
		i++
	}
	// d.deals = deals

	for _, deal := range dealprotos {
		if deal.Index == uint32(d.index) {
			// todo: deliver to ourselves
			continue
		}
		buf, err := proto.Marshal(deal)
		if err != nil {
			return err
		}

		err = d.send(ctx, string(ProtocolDeal), buf, d.participants[deal.Index])
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *dkg) send(ctx context.Context, msgType string, buf []byte, node transport.Node) error {
	cid, err := types.CidFromBytes(buf)
	if err != nil {
		return err
	}
	msg, err := d.transport.NewMessage(d.ringID, cid.String(), false, buf, msgType)
	if err != nil {
		return err
	}
	if err := d.transport.Send(ctx, node, msg); err != nil {
		return err
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
		//
	case string(ProtocolResponse):
		//
	default:
		panic("bad message type") //todo
	}

	return nil
}

package rabin

import (
	"context"
	"sync"

	"go.dedis.ch/kyber/v3"
	rabindkg "go.dedis.ch/kyber/v3/share/dkg/rabin"
	"go.dedis.ch/kyber/v3/suites"
	"go.dedis.ch/protobuf"

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

	deals map[int]*rabindkg.Deal //db.Table[rabinv1alpha1.Deals]
	resps []*rabindkg.Response   // db.Table[rabinv1alpha1.Response]

	share crypto.PriShare

	repo      *db.Repository
	transport transport.Transport
	bulletin  bulletin.Bulletin

	state       orbisdkg.State
	initialized bool
}

func New(repo db.Repository, t transport.Transport, b bulletin.Bulletin, pk crypto.PrivateKey) (*dkg, error) {
	suite, err := crypto.SuiteForType(pk.Type())
	if err != nil {
		return nil, err
	}

	scalar := pk.Scalar()
	pubPoint := suite.Point().Mul(scalar, nil) // public point for scalar

	return &dkg{
		repo:      &repo,
		transport: t,
		bulletin:  b,
		suite:     suite,
		privKey:   scalar,
		pubKey:    pubPoint,
		index:     -1,
	}, nil
}

// Init initializes the DKG with the target nodes
func (d *dkg) Init(ctx context.Context, nodes []orbisdkg.Node, n int, threshold int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

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
	d.deals = deals

	for _, deal := range deals {
		if deal.Index == uint32(d.index) {
			// todo: deliver to ourselves
			continue
		}
		buf, err := protobuf.Encode(deal)
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

package rabin

import (
	"context"
	"sync"

	"go.dedis.ch/kyber/v3"
	rabindkg "go.dedis.ch/kyber/v3/share/dkg/rabin"
	"go.dedis.ch/kyber/v3/suites"

	"github.com/sourcenetwork/orbis-go/infra/logger"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	orbisdkg "github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

type dkg struct {
	mu sync.Mutex

	log logger.Logger

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

	initialized bool
}

func New(lg logger.Logger, repo *db.Repository, t transport.Transport, b bulletin.Bulletin, pk crypto.PrivateKey) (orbisdkg.DKG, error) {
	return newDKG(lg, repo, t, b, pk)
}

func newDKG(lg logger.Logger, repo *db.Repository, t transport.Transport, b bulletin.Bulletin, pk crypto.PrivateKey) (*dkg, error) {
	suite, err := crypto.SuiteForType(pk.Type())
	if err != nil {
		return nil, err
	}

	scalar := pk.Scalar()
	pubPoint := suite.Point().Mul(scalar, nil) // public point for scalar

	return &dkg{
		log:       lg,
		repo:      repo,
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

	return nil
}

func (d *dkg) Name() string {
	return "rabin"
}

func (d *dkg) PublicKey() crypto.PublicKey {
	pk, _ := crypto.PublicKeyFromPoint(d.pubKey)
	return pk
}

func (d *dkg) Share() crypto.PriShare {
	return d.share
}

func (d *dkg) State() orbisdkg.State {
	panic("not implemented") // TODO: Implement
}

// Start the DKG setup process.
func (d *dkg) Start(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

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
			// deliver to ourselves
		}

		d.transport.NewMessage(d.dealID(deal), false, dealpb, DEAL_SEND)
		d.transport.Send()
	}

	return nil
}

func (d *dkg) Close(_ context.Context) error {
	panic("not implemented") // TODO: Implement
}

func (d *dkg) ProcessMessage(_ *transport.Message) error {
	// todo maybe?: validate msg.PublicKey matches payload pubkeys
	panic("not implemented") // TODO: Implement
}

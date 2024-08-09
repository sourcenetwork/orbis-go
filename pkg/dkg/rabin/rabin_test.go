package rabin

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"math/rand"
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	rabindkg "go.dedis.ch/kyber/v3/share/dkg/rabin"
	"go.dedis.ch/kyber/v3/suites"
	"go.dedis.ch/kyber/v3/util/random"

	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/pkg/bulletin/memmap"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	"github.com/sourcenetwork/orbis-go/pkg/host"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	p2ptransport "github.com/sourcenetwork/orbis-go/pkg/transport/p2p"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

const (
	portRangeStart = 1000
	portRangeEnd   = 8999
)

func TestDealProtoSerialization(t *testing.T) {
	rand := random.New()
	suite := edwards25519.NewBlakeSHA256Ed25519()

	s1 := suite.Scalar().Pick(rand)
	p1 := suite.Point().Mul(s1, nil)
	p2 := suite.Point().Pick(rand)
	p3 := suite.Point().Pick(rand)

	points := []kyber.Point{p1, p2, p3}
	dkg, err := rabindkg.NewDistKeyGenerator(suite, s1, points, 2)
	if err != nil {
		panic(err)
	}

	deals, err := dkg.Deals()
	if err != nil {
		panic(err)
	}

	d1, err := dealToProto(deals[1])
	if err != nil {
		panic(err)
	}

	d1r, err := dealFromProto(suite, d1)
	require.NoError(t, err)

	require.Equal(t, deals[1].Index, d1r.Index)
	require.Equal(t, deals[1].Deal.Nonce, d1r.Deal.Nonce)
	require.Equal(t, deals[1].Deal.Signature, d1r.Deal.Signature)
	require.Equal(t, deals[1].Deal.Cipher, d1r.Deal.Cipher)
	require.True(t, deals[1].Deal.DHKey.Equal(d1r.Deal.DHKey))

	// require.Equal(t, deals[1].Deal.DHKey, d1r.Deal.DHKey)
	// require.Equal(t, deals[1].Index, d1r.Index)
	// require.Equal(t, deals[1].Index, d1r.Index)

	require.Equal(t, deals[1], d1r)
}

func newTestDB(t *testing.T) *db.DB {
	tmpPath := t.TempDir()
	db, err := db.New(tmpPath)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func randomPort() int {
	return rand.Intn(portRangeEnd-portRangeStart) + portRangeStart
}

func randomNodes(num int, ste suites.Suite) []transport.Node {
	nodes := make([]transport.Node, num)
	keyType, err := crypto.KeyPairTypeFromString(ste.String())
	if err != nil {
		panic(err)
	}
	for i := 0; i < num; i++ {
		_, pub, err := crypto.GenerateKeyPair(keyType, cryptorand.Reader)
		if err != nil {
			panic(err)
		}
		nodes[i] = randomNodeFromPublicKey(pub)
	}
	return nodes
}

func randomNodeFromPublicKey(pubkey crypto.PublicKey) transport.Node {
	pid, err := peer.IDFromPublicKey(crypto.ToLibP2PPublicKey(pubkey))
	if err != nil {
		panic(err)
	}
	addr, err := ma.NewMultiaddr(fmt.Sprintf("/tcp/%d", randomPort()))
	if err != nil {
		panic(err)
	}
	return p2ptransport.NewNode(pid.String(), pubkey, addr)
}

func newP2PHost(ctx context.Context) *host.Host {
	defaultHost, err := config.Default[config.Host]()
	if err != nil {
		panic(err)
	}
	defaultHost.Crypto.Seed = 1
	h, err := host.New(ctx, defaultHost)
	if err != nil {
		panic(err)
	}
	return h
}

func newBasicDKG(t *testing.T, ctx context.Context) (*dkg, crypto.PrivateKey) {
	d := newTestDB(t)
	h := newP2PHost(ctx)
	tp, err := p2ptransport.New(ctx, h, config.Transport{})
	if err != nil {
		panic(err)
	}

	b := memmap.New()
	require.NotNil(t, b)
	require.NotNil(t, b.Events())

	rkeys := []db.RepoKey{
		db.NewRepoKey("dkg"),
	}
	dkg, err := New(d, rkeys, tp, b)
	require.NoError(t, err)

	lpriv := h.Peerstore().PrivKey(h.ID())
	require.NotNil(t, lpriv)

	cpriv, err := crypto.PrivateKeyFromLibP2P(lpriv)
	require.NoError(t, err)
	pub := cpriv.GetPublic()

	nodes := randomNodes(2, suites.MustFind("Ed25519"))
	nodes = append(nodes, randomNodeFromPublicKey(pub))

	err = dkg.Init(context.Background(), cpriv, types.RingID("0x123"), nodes, 3, 2, false)
	require.NoError(t, err)

	return dkg, cpriv
}

func assertEqualDKG(t *testing.T, dkg1, dkg2 *dkg) {
	assert.Equal(t, dkg1.ringID, dkg2.ringID)
	assert.NotEmpty(t, dkg1.ringID)
	assert.Equal(t, dkg1.index, dkg2.index)
	assert.NotEmpty(t, dkg1.index)
	assert.Equal(t, dkg1.num, dkg2.num)
	assert.NotEmpty(t, dkg1.num)
	assert.Equal(t, dkg1.threshold, dkg2.threshold)
	assert.NotEmpty(t, dkg1.threshold)
	assert.Equal(t, dkg1.suite.String(), dkg2.suite.String())
	assert.Equal(t, dkg1.state, dkg2.state)
	// TODO: Improve participant equality check
	assert.Len(t, dkg1.participants, len(dkg2.participants))
	assert.NotEmpty(t, dkg1.participants)
	require.NotEmpty(t, dkg1.rdkg.Dealer().FPoly())
	assert.Equal(t, dkg1.rdkg.Dealer().FPoly().String(), dkg2.fPoly.String())
	require.NotEmpty(t, dkg1.rdkg.Dealer().GPoly())
	assert.Equal(t, dkg1.rdkg.Dealer().GPoly().String(), dkg2.gPoly.String())
	assert.Equal(t, dkg1.rdkg.Dealer().Secret(), dkg2.secret)
	assert.NotEmpty(t, dkg2.secret)
}

func TestDKGProtoSerialization(t *testing.T) {
	ctx := context.Background()
	dkg1, _ := newBasicDKG(t, ctx)

	dkgp, err := dkgToProto(dkg1)
	require.NoError(t, err)

	dkg2, err := dkgFromProto(dkgp)
	require.NoError(t, err)

	assertEqualDKG(t, dkg1, &dkg2)
}

func TestDKGSaveAndLoad(t *testing.T) {
	ctx := context.Background()
	dkg1, priv := newBasicDKG(t, ctx)
	fmt.Println(dkg1.ringID)

	fmt.Println("save start!")
	err := dkg1.save(ctx)
	require.NoError(t, err)
	fmt.Println("save end!")

	fmt.Println("========")
	// require.NoError(t, dkg1.db.Debug())
	fmt.Println("========")

	dkg3, err := New(dkg1.db, dkg1.rkeys, dkg1.transport, dkg1.bulletin)
	require.NoError(t, err)
	err = dkg3.Init(ctx, priv, dkg1.ringID, dkg1.participants, dkg1.num, dkg1.threshold, true)
	require.NoError(t, err)

	err = dkg3.load(ctx)
	require.NoError(t, err)

	assertEqualDKG(t, dkg1, dkg3)
}

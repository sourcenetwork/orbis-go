package app

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	cid "github.com/ipfs/go-cid"
	mc "github.com/multiformats/go-multicodec"
	mh "github.com/multiformats/go-multihash"

	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

type manifest struct {
	ID        string `json:"id"`
	N         int32  `json:"n"`
	T         int32  `json:"t"`
	DKG       string `json:"dkg"`
	PSS       string `json:"pss"`
	PRE       string `json:"pre"`
	Bulletin  string `json:"bulletin"`
	Transport string `json:"transport"`
	Nodes     []node `json:"nodes"`

	Authorization  string `json:"authorization"`
	Authentication string `json:"authentication"`
}

type node struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

func ringID(r *ringv1alpha1.Manifest) types.RingID {

	m := manifest{
		N:              r.N,
		T:              r.T,
		DKG:            r.Dkg,
		PSS:            r.Pss,
		PRE:            r.Pre,
		Bulletin:       r.Bulletin,
		Transport:      r.Transport,
		Nodes:          make([]node, len(r.Nodes)),
		Authorization:  r.Authorization,
		Authentication: r.Authentication,
	}

	for i, n := range r.Nodes {
		m.Nodes[i] = node{
			ID:      n.Id,
			Address: n.Address,
		}
	}

	b, err := cbor.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("marshal manifest: %w", err))
	}

	pref := cid.Prefix{
		Version:  1,
		Codec:    uint64(mc.Cbor),
		MhType:   mh.SHA2_256,
		MhLength: -1, // default length
	}

	cid, err := pref.Sum(b)
	if err != nil {
		panic(fmt.Errorf("create cid: %w", err))
	}

	return types.RingID(cid.String())
}

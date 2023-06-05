package app

import (
	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/db"
)

func ringPkFunc(kb db.KeyBuilder, r *ringv1alpha1.Ring) []byte {
	return kb.AddStringField(r.Id).
		Bytes()
}

// func ringToProto(ring *Ring) *ringv1alpha1.Ring {
// 	nodes := make([]*ringv1alpha1.Node, len(ring.nodes))
// 	for _, n := range ring.nodes {

// 	}
// 	return &ringv1alpha1.Ring{
// 		Id: string(ring.ID),
// 		N: int32(ring.N),
// 		T: int32(ring.T),
// 		Dkg: ring.DKG.Name(),
// 		Pss: ring.PSS.Name(),
// 		Pre: ring.PRE.Name(),

// 	}
// }

// func ringFromProto(ring *ringv1alpha1.Ring) *Ring {
// 	return nil
// }

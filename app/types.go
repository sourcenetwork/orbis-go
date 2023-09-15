package app

import (
	ringv1alpha1 "github.com/sourcenetwork/orbis-go/gen/proto/orbis/ring/v1alpha1"
	"github.com/sourcenetwork/orbis-go/pkg/db"
)

func ringPkFunc(kb db.KeyBuilder, r *ringv1alpha1.Ring) []byte {
	return kb.AddStringField(r.Id).Bytes()
}

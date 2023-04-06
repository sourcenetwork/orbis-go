package types

import (
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

// SecretShare is a cryptograhic share of a
// secret.
type PrivSecretShare struct{}

// SecretID is a Secret identifier
type SecretID string

// RingID is a SecretRing identifier
type RingID string

// type Node struct{}

func CidFromBytes(b []byte) (cid.Cid, error) {
	h, err := mh.Sum(b, mh.SHA2_256, -1)
	if err != nil {
		return cid.Undef, err
	}
	return cid.NewCidV1(cid.Raw, h), nil
}

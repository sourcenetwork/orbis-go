package dkg

import "fmt"

var (
	ErrBadNodeSet     = fmt.Errorf("node set size doesn't match n")
	ErrNotInitialized = fmt.Errorf("dkg not initialized")
	ErrMissingSelf    = fmt.Errorf("missing self from node set")
)

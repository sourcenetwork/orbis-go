package authz

import "context"

const (
	READ  = "read"
	WRITE = "write"
)

type Authz interface {
	Name() string

	// Init completes any additional initialization the authz service needs
	Init(ctx context.Context) error

	// Check is the main authorization API. It takes a subject and a payload
	// and returns true, or false (and optionally error). The interface
	// currently makes no assumption about the structure of contents of the
	// subject and payload byte-arrays. This is left up to the implementation
	// so a generic byte-array was the most appropriate.
	Check(ctx context.Context, perm []byte, subject string) (bool, error)
}

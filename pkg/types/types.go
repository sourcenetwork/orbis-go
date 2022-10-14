package types

type State int64

const (
	STATE_INITIALIZED = State(iota)
	STATE_UNINITIALIZED
)

// Secret is a managed secret
type Secret interface {
	Recover() ([]byte, error)
	Marshal() ([]byte, error)
}

// SecretShare is a cryptograhic share of a
// secret.
type PrivSecretShare struct{}

//
type Proof struct{}

type PublicKey struct{}

// SecretID is a Secret identifier
type SecretID string

// RingID is a SecretRing identifier
type RingID string

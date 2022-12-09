package types

// Secret is a managed secret
type Secret interface {
	Recover() ([]byte, error)
	Marshal() ([]byte, error)
}

// SecretShare is a cryptograhic share of a
// secret.
type PrivSecretShare struct{}

// SecretID is a Secret identifier
type SecretID string

// RingID is a SecretRing identifier
type RingID string

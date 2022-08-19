package types

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

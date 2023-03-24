package types

// SecretShare is a cryptograhic share of a
// secret.
type PrivSecretShare struct{}

// SecretID is a Secret identifier
type SecretID string

// RingID is a SecretRing identifier
type RingID string

// type Node struct{}

func RingFromManifest(manifest []byte) (*Ring, RingID, error) {
	return nil, "", nil
}

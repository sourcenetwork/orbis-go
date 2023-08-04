//go:generate mockery --all --with-expecter
package authn

import (
	"context"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
)

// CredentialService parses and verifies the authentication
// data from the incoming
type CredentialService interface {
	GetAndVerifyRequestMetadata(ctx context.Context) (SubjectInfo, error)
}

// Resolves the parsed subject data into fully formed UserInfo.
// Examples include DID `did` resolvers.
// This isnt necessary if:
// A) The request fully embeds the public key (such as JWKs/JWTs)
// B) The downstream consumer of the credential service doesn't
// need a PublicKey
type KeyResolver interface {
	Resolve(context.Context, string) (SubjectInfo, error)
}

// SubjectInfo is the returned authentication user info from
// a given request
type SubjectInfo struct {
	Type    string
	Subject string
	PubKey  crypto.PublicKey
}

// Metadata represents the generic version of parsed metadata
// from a request. This is used by the CredentialService so we
// can customize the parsing based on the kind of requests
// being made. Eg. GRPC, GraphQL, REST, etc.. Each of which
// could store the request metadata slightly differently.
type Metadata interface {
	Append(k string, vals ...string)
	Delete(k string)
	Get(k string) []string
	Len() int
	Set(k string, vals ...string)
}

// FromIncomingContext is a helper to strongly define
type RequestMetadataParser interface {
	Parse(ctx context.Context) (Metadata, bool)
}

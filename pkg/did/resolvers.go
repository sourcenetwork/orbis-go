package did

import (
	"context"
	"fmt"

	"github.com/TBD54566975/ssi-sdk/did/resolution"

	"github.com/sourcenetwork/orbis-go/pkg/authn"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
)

type resolver struct {
	resolution.Resolver
}

// Resolve implements auth.KeyResolver
func (r resolver) Resolve(ctx context.Context, id string) (authn.SubjectInfo, error) {
	result, err := r.Resolver.Resolve(ctx, id)
	if err != nil {
		return authn.SubjectInfo{}, fmt.Errorf("did resolving: %w", err)
	}

	if len(result.VerificationMethod) == 0 {
		return authn.SubjectInfo{}, fmt.Errorf("resolver result missing verification methods")
	}

	stdPk, err := result.VerificationMethod[0].PublicKeyJWK.ToPublicKey()
	if err != nil {
		return authn.SubjectInfo{}, fmt.Errorf("jwk to public key")
	}
	pk, err := crypto.PublicKeyFromGoPublicKey(stdPk)
	if err != nil {
		return authn.SubjectInfo{}, fmt.Errorf("crypto key conversion: %w", err)
	}
	return authn.SubjectInfo{
		Subject: result.ID,
		PubKey:  pk,
		Type:    result.VerificationMethod[0].Type.String(),
	}, nil
}

func NewResolver(r resolution.Resolver) authn.KeyResolver {
	return resolver{Resolver: r}
}

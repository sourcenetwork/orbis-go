package did

import (
	"context"
	"fmt"
	"strings"

	"github.com/TBD54566975/ssi-sdk/did"
	"github.com/TBD54566975/ssi-sdk/did/resolution"

	"github.com/sourcenetwork/orbis-go/pkg/authn"
	"github.com/sourcenetwork/orbis-go/pkg/crypto"
)

type resolver struct {
	resolution.Resolver
}

// Resolve implements auth.KeyResolver
func (r resolver) Resolve(ctx context.Context, id string) (authn.SubjectInfo, error) {
	suffixParts := strings.Split(id, "#")
	var suffix string
	if len(suffixParts) == 2 {
		suffix = suffixParts[1]
		id = suffixParts[0]
	}

	result, err := r.Resolver.Resolve(ctx, id)
	if err != nil {
		return authn.SubjectInfo{}, fmt.Errorf("did resolving: %w", err)
	}

	if len(result.VerificationMethod) == 0 {
		return authn.SubjectInfo{}, fmt.Errorf("resolver result missing verification methods")
	}

	// use the first verifcation method if none is specific via a suffix fragment
	verKeyID := result.VerificationMethod[0].ID
	if suffix != "" {
		verKeyID = suffix
	}
	stdPk, err := did.GetKeyFromVerificationMethod(result.Document, verKeyID)
	if err != nil {
		return authn.SubjectInfo{}, fmt.Errorf("did document to public key: %w", err)
	}
	pk, err := crypto.PublicKeyFromStdPublicKey(stdPk)
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

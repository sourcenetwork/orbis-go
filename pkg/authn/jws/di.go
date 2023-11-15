package jws

import (
	"fmt"

	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/pkg/authn"
	"github.com/sourcenetwork/orbis-go/pkg/db"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

var (
	_                 types.Factory[authn.CredentialService] = (*selfSignedFactory)(nil)
	SelfSignedFactory                                        = selfSignedFactory{}
)

type selfSignedFactory struct{}

func (selfSignedFactory) New(inj *do.Injector, rkeys []db.RepoKey, _ config.Config) (authn.CredentialService, error) {
	resolver, err := do.Invoke[authn.KeyResolver](inj)
	if err != nil {
		return nil, fmt.Errorf("invoke key resolver: %w", err)
	}
	metadataFn, err := do.Invoke[authn.RequestMetadataParser](inj)
	if err != nil {
		return nil, fmt.Errorf("invoke metadata parser: %w", err)
	}
	return NewSelfSignedCredentialService(resolver, metadataFn), nil
}

func (selfSignedFactory) Name() string {
	return "jws-did"
}

// Repos returns empty string to indicate no dependent repos needed
// which means an empty array will be passed to `New` for `rkeys`
func (selfSignedFactory) Repos() []string {
	return []string{}
}

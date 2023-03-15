package pss

import (
	"github.com/samber/do"

	"github.com/sourcenetwork/orbis-go/pkg/bulletin"
	"github.com/sourcenetwork/orbis-go/pkg/dkg"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
	"github.com/sourcenetwork/orbis-go/pkg/types"
)

// Testing the Factory + Dependency Injection
// patterns
//
// The github.com/samber/do Dependency Injection
// library provider/invoke model doesn't
// actually create *new* instances of the provided
// type.
//
// This is where the Factory comes into focus.
// The Factory is provided via the DI framework
// so we can *then* create new instances of the
// target services.

type testPSS struct {
	PSS // embed PSS to conform to interface (noop)
}

type testPSSFactory struct{}

func (f testPSSFactory) New(
	_ types.RingID,
	_ int,
	_ int,
	_ transport.Transport,
	_ bulletin.Bulletin,
	_ []types.Node,
	_ dkg.DKG,
) (PSS, error) {
	return &testPSS{}, nil
}

func TestProvider(d *do.Injector) {
	do.ProvideNamed(d, "pss_test", func(i *do.Injector) (testPSSFactory, error) {
		return testPSSFactory{}, nil
	})
}

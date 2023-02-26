package pss

import (
	"context"

	"github.com/samber/do"
	"go.dedis.ch/kyber/v3/suites"

	"github.com/sourcenetwork/orbis-go/pkg/crypto"
	"github.com/sourcenetwork/orbis-go/pkg/pss/types"
	"github.com/sourcenetwork/orbis-go/pkg/transport"
)

type Message interface{}

type Service interface {
	Name() string
	Suite() suites.Suite

	Start()
	Shutdown() error
	ProcessMessage(context.Context, Message)

	PublicKey() crypto.PublicKey
	PublicPoly() crypto.PubPoly

	Share() crypto.PriShare

	State() types.State
}

type Node interface {
	transport.Node
	Index() int
}

type ProviderFn = func(*do.Injector) Service

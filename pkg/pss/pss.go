package pss

import (
	"context"

	"github.com/samber/do"

	"github.com/sourcenetwork/orbis-go/pkg/pss/types"
)

type Message interface{}

type Service interface {
	Name() string

	Start()
	Shutdown() error
	ProcessMessage(context.Context, Message)

	State() types.State
}

type ProviderFn func(*do.Injector) Service

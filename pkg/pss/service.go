package pss

import (
	"context"

	"github.com/sourcenetwork/orbis-go/pkg/pss/types"
)

type Message interface{}

type Service interface {
	Name() string
	ProcessMessage(context.Context, Message)
	State() types.State
}

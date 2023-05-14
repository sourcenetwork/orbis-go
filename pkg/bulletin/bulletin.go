package bulletin

import (
	"context"
	"fmt"
)

var (
	ErrDuplicateMessage = fmt.Errorf("bulletin: duplicate message")
	ErrMessageNotFound  = fmt.Errorf("bulletin: message not found")
)

type Message []byte

type Proof []byte

type Response struct {
	Data  Message
	Proof Proof
}

type Query struct{}

type Bulletin interface {
	Name() string
	// message format := /<namespace>/
	// /ring/<ringID>/pss/<epochNum>/<nodeIndex>/<action>
	// /ring/<ringID>/pre/<nodeIndex>/<action>
	// /ring/<ringID>/dkg/<nodeIndex>/<action>
	Post(context.Context, string, Message) (Response, error)
	Read(context.Context, string) (Response, error)

	// Query Search the bulletin board using a glob based
	// text search system.
	Query(context.Context, string) ([]Response, error)

	// Verify(context.Context, Proof, string, Message) bool

	// EventBus
	// Events() eventbus.Bus
}

type Config struct {
	Proof bool
}

type Option func(*Config)

func WithProof(p bool) Option {
	return func(c *Config) {
		c.Proof = p
	}
}

package bulletin

import (
	"context"
	"fmt"
)

var (
	ErrEmptyID          = fmt.Errorf("bulletin: empty ID")
	ErrEmptyNamespace   = fmt.Errorf("bulletin: empty namespace")
	ErrDuplicateMessage = fmt.Errorf("bulletin: duplicate message")
	ErrDuplicateTopic   = fmt.Errorf("bulletin: duplicate topic")
	ErrTopicNotFound    = fmt.Errorf("bulletin: topic not found")
	ErrMessageNotFound  = fmt.Errorf("bulletin: message not found")
	ErrReadTimeout      = fmt.Errorf("bulletin: read timeout")
	ErrBadResponseType  = fmt.Errorf("bulletin: bad response type")
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
	Register(ctx context.Context, namespace string) error
	// message format := /<namespace>/
	// /ring/<ringID>/pss/<epochNum>/<nodeIndex>/<action>
	// /ring/<ringID>/pre/<nodeIndex>/<action>
	// /ring/<ringID>/dkg/<nodeIndex>/<action>
	Post(context.Context, ID, Message) (Response, error)
	Read(context.Context, ID) (Response, error)

	// Query Search the bulletin board using a glob based
	// text search system.
	Query(context.Context, string) ([]Response, error)

	// Verify(context.Context, Proof, string, Message) bool

	// EventBus
	// Events() eventbus.Bus
}

// ID
//
// /<namespace>/<service>/<key>
type ID interface {
	// Returns the full ID serialized as a string
	String() string

	// Returns the `/<namspace>/<service>` pair as a string
	ServiceNamespace() string

	// Returns the `<key>` as a string
	Key() string
}

// func IDBuilder

type Config struct {
	Proof bool
}

type Option func(*Config)

func WithProof(p bool) Option {
	return func(c *Config) {
		c.Proof = p
	}
}

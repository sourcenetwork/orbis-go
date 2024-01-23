package bulletin

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/eventbus-go"

	"github.com/sourcenetwork/orbis-go/pkg/transport"
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

func ErrDuplicateMessageF(id string) error {
	return fmt.Errorf("%w: %s", ErrDuplicateMessage, id)
}

type Message []byte

type Proof []byte

type Event struct {
	Message *transport.Message
	ID      string
}

// Response
type Response struct {
	Data  *transport.Message
	ID    string
	Proof Proof
}

// QueryResponse is the response object for a `Query()` request
// which is designed to be sent over a channel
type QueryResponse struct {
	Resp Response
	Err  error
}

type Query struct{}

type Bulletin interface {
	Name() string
	Init(context.Context) error
	Register(ctx context.Context, namespace string) error
	// message format := /<namespace>/
	// /ring/<ringID>/pss/<epochNum>/<action>/<nodeIndex>
	// /ring/<ringID>/pre/<action>/<nodeIndex>
	// /ring/<ringID>/dkg/rabin/<action>/<fromIndex>/<toIndex>
	Post(ctx context.Context, namespace, id string, msg *transport.Message) (Response, error)
	Read(ctx context.Context, namespace, id string) (Response, error)
	// Has(context.Context, string) (bool, error)

	// Query Search the bulletin board using a glob based
	// text search system.
	Query(ctx context.Context, namespace string, query string) (<-chan QueryResponse, error)

	// Verify(context.Context, Proof, string, Message) bool

	// EventBus
	Events() eventbus.Bus
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

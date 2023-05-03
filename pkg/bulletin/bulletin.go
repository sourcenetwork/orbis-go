package bulletin

import (
	"context"
	"fmt"

	"github.com/samber/do"
	"github.com/sourcenetwork/orbis-go/config"
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

type Factory interface {
	New(ctx context.Context, inj *do.Injector, cfg config.Bulletin) (Bulletin, error)
	Name() string
}

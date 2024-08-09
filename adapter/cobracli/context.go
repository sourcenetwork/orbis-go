package cobracli

import (
	"context"

	"github.com/segmentio/cli"
	"github.com/sourcenetwork/orbis-go/pkg/keyring"
)

type Context struct {
	kr keyring.Keyring
	pf cli.PrintFlusher
}

func (c *Context) Keyring() keyring.Keyring {
	return c.kr
}

func (c *Context) Output() cli.PrintFlusher {
	return c.pf
}

// cli.ClientContext()

type ctxKey string

var (
	clientCtxKey ctxKey = "clientCtxKey"
)

func FromContext(ctx context.Context) (*Context, bool) {
	clientCtx, ok := ctx.Value(clientCtxKey).(*Context)
	return clientCtx, ok
}

func WithContext(ctx context.Context, kr keyring.Keyring, pf cli.PrintFlusher) context.Context {
	clientCtx := &Context{
		kr: kr,
		pf: pf,
	}
	ctx = context.WithValue(ctx, clientCtxKey, clientCtx)
	return ctx
}

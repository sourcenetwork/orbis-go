package bulletin

import (
	"context"

	"github.com/samber/do"
)

type ProviderFn func(*do.Injector) Service

type Message interface{}

type Response struct {
	Data  interface{}
	Proof []byte
}

type Query interface{}

type Service interface {
	Post(context.Context, string, Message, ...Option) (Response, error)
	Read(context.Context, string, ...Option) (Message, error)
	Query(context.Context, Query, ...Option) ([]Message, error)

	Start()
	Shutdown()
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

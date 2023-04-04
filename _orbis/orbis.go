package orbis

import (
	"context"

	"github.com/samber/do"
)

type node struct {
	injector *do.Injector
}

// NewNode returns a newly instanciated obris node
// configured with various services
func NewNode(ctx context.Context, opts ...Option) (*node, error) {
	// todo
	return nil, nil
}

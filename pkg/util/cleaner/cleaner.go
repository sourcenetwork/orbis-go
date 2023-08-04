package cleaner

import "sync"

// Cleaner register and execute (once only) clean up functions.
type Cleaner struct {
	mu       sync.Mutex
	once     sync.Once
	cleaners []func()
}

func New() *Cleaner {
	return &Cleaner{}
}

func (c *Cleaner) Regster(fn func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cleaners = append(c.cleaners, fn)
}

// CleanUp executes all registered clean up functions (once only).
func (c *Cleaner) CleanUp() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.once.Do(func() {
		for _, fn := range c.cleaners {
			fn()
		}
	})
}

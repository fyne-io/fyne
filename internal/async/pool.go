package async

import "sync"

// Implementation inspired by https://github.com/tailscale/tailscale/blob/main/syncs/pool.go.

// Pool is the generic version of sync.Pool.
type Pool[T any] struct {
	pool sync.Pool

	// New specifies a function to generate
	// a value when Get would otherwise return the zero value of T.
	New func() T
}

// Get selects an arbitrary item from the Pool, removes it from the Pool,
// and returns it to the caller.
func (p *Pool[T]) Get() T {
	x, ok := p.pool.Get().(T)
	if !ok && p.New != nil {
		return p.New()
	}

	return x
}

// Put adds x to the pool.
func (p *Pool[T]) Put(x T) {
	p.pool.Put(x)
}

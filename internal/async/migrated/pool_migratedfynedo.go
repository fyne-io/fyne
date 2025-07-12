//go:build migrated_fynedo

package migrated

// Pool (when migrated) is just a simple single threaded data buffer.
type Pool[T any] struct {
	buffer *T

	// New specifies a function to generate
	// a value when Get would otherwise return the zero value of T.
	New func() T
}

// Get get the saved item from the Pool, removes it from the Pool,
// and returns it to the caller.
func (p *Pool[T]) Get() T {
	if p.buffer == nil {
		if p.New != nil {
			return p.New()
		}
		return *new(T)
	}

	x := *p.buffer
	p.buffer = nil
	return x
}

// Put adds x to the pool.
func (p *Pool[T]) Put(x T) {
	p.buffer = &x
}

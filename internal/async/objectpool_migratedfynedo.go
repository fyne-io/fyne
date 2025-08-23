//go:build migrated_fynedo

package async

// ObjectPool represents a pool of reusable objects. Primarily intended for use with canvas objects.
type ObjectPool[T any] struct {
	queue Queue[T]
}

// Get retrieves an object from the pool.
func (p *ObjectPool[T]) Get() T {
	return p.queue.Out()
}

// Put returns a object to the pool.
func (p *ObjectPool[T]) Put(obj T) {
	p.queue.In(obj)
}

// Clear removes all objects from the pool.
func (p *ObjectPool[T]) Clear() {
	p.queue.Clear()
}

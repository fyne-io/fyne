//go:build !migrated_fynedo

package async

// ObjectPool represents a pool of reusable objects. Primarily intended for use with canvas objects.
type ObjectPool[T comparable] struct {
	Pool[T]
}

// Clear removes all objects from the pool.
func (p *ObjectPool[T]) Clear() {
	var zero T
	for p.Get() != zero {
	}
}

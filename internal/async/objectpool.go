//go:build !migrated_fynedo

package async

// ObjectPool represents a pool of reusable canvas objects.
type ObjectPool[T comparable] struct {
	Pool[T]
}

func (p *ObjectPool[T]) Clear() {
	var zero T
	for p.Get() != zero {
	}
}

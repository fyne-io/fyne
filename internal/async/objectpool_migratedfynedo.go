//go:build migrated_fynedo

package async

import "fyne.io/fyne/v2"

// ObjectPool represents a pool of reusable canvas objects.
type ObjectPool struct {
	queue CanvasObjectQueue
}

// Get retrieves a canvas object from the pool.
func (p *ObjectPool) Get() fyne.CanvasObject {
	return p.queue.Out()
}

// Put returns a canvas object to the pool.
func (p *ObjectPool) Put(obj fyne.CanvasObject) {
	p.queue.In(obj)
}

// Clear removes all canvas objects from the pool.
func (p *ObjectPool) Clear() {
	p.queue.Clear()
}

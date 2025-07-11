//go:build migrated_fynedo

package async

import "fyne.io/fyne/v2"

// CanvasObjectQueue represents a simple single threaded queue for managing canvas objects.
type CanvasObjectQueue struct {
	items []fyne.CanvasObject
}

// NewCanvasObjectQueue returns a queue for caching values.
func NewCanvasObjectQueue() *CanvasObjectQueue {
	return &CanvasObjectQueue{items: make([]fyne.CanvasObject, 0, 32)}
}

// In puts the given value at the tail of the queue.
func (q *CanvasObjectQueue) In(v fyne.CanvasObject) {
	q.items = append(q.items, v)
}

// Out removes and returns the value at the head of the queue.
// It returns nil if the queue is empty.
func (q *CanvasObjectQueue) Out() fyne.CanvasObject {
	if len(q.items) == 0 {
		return nil
	}

	head := q.items[0]
	q.items[0] = nil
	q.items = q.items[1:]
	return head
}

// Len returns the number of items in the queue.
func (q *CanvasObjectQueue) Len() uint64 {
	return uint64(len(q.items))
}

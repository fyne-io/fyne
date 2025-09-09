//go:build migrated_fynedo

package async

import "fyne.io/fyne/v2"

const defaultQueueCapacity = 64

// CanvasObjectQueue represents a single-threaded queue for managing canvas objects using a ring buffer.
type CanvasObjectQueue struct {
	buffer []fyne.CanvasObject
	head   int
	size   int
}

// NewCanvasObjectQueue returns a queue for caching values with an initial capacity.
func NewCanvasObjectQueue() *CanvasObjectQueue {
	return &CanvasObjectQueue{buffer: make([]fyne.CanvasObject, defaultQueueCapacity)}
}

// In adds the given value to the tail of the queue.
// If the queue is full, it grows the buffer dynamically.
func (q *CanvasObjectQueue) In(v fyne.CanvasObject) {
	if q.size == len(q.buffer) {
		buffer := make([]fyne.CanvasObject, len(q.buffer)*2)
		copy(buffer, q.buffer[q.head:])
		copy(buffer[len(q.buffer)-q.head:], q.buffer[:q.head])
		q.buffer = buffer
		q.head = 0
	}

	tail := (q.head + q.size) % len(q.buffer)
	q.buffer[tail] = v
	q.size++
}

// Out removes and returns the value at the head of the queue.
// It returns nil if the queue is empty.
func (q *CanvasObjectQueue) Out() fyne.CanvasObject {
	if q.size == 0 {
		return nil
	}

	first := q.buffer[q.head]
	q.buffer[q.head] = nil
	q.head = (q.head + 1) % len(q.buffer)
	q.size--

	if q.size == 0 && len(q.buffer) > 4*defaultQueueCapacity {
		q.buffer = make([]fyne.CanvasObject, defaultQueueCapacity)
		q.head = 0
	}

	return first
}

// Len returns the number of items in the queue.
func (q *CanvasObjectQueue) Len() uint64 {
	return uint64(q.size)
}

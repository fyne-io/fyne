//go:build migrated_fynedo

package async

import "fyne.io/fyne/v2"

const defaultQueueCapacity = 64

// CanvasObjectQueue represents a single-threaded queue for managing canvas objects using a ring buffer.
type CanvasObjectQueue struct {
	Queue[fyne.CanvasObject]
}

// NewCanvasObjectQueue returns a queue for caching values with an initial capacity.
func NewCanvasObjectQueue() *CanvasObjectQueue {
	return &CanvasObjectQueue{}
}

// Queue represents a single-threaded queue for managing values using a ring buffer.
type Queue[T any] struct {
	buffer []T
	head   int
	size   int
}

// In adds the given value to the tail of the queue.
// If the queue is full, it grows the buffer dynamically.
func (q *Queue[T]) In(v T) {
	if q.buffer == nil {
		q.Clear()
	}

	if q.size == len(q.buffer) {
		buffer := make([]T, len(q.buffer)*2)
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
func (q *Queue[T]) Out() T {
	var zero T
	if q.size == 0 {
		return zero
	}

	first := q.buffer[q.head]
	q.buffer[q.head] = zero
	q.head = (q.head + 1) % len(q.buffer)
	q.size--

	if q.size == 0 && len(q.buffer) > 4*defaultQueueCapacity {
		q.Clear()
	}

	return first
}

// Len returns the number of items in the queue.
func (q *Queue[T]) Len() uint64 {
	return uint64(q.size)
}

// Clear removes all items from the queue.
func (q *Queue[T]) Clear() {
	q.buffer = make([]T, defaultQueueCapacity)
	q.head = 0
	q.size = 0
}

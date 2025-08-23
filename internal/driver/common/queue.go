package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
)

type deduplicatedObjectQueue struct {
	queue *async.CanvasObjectQueue
	dedup async.Map[fyne.CanvasObject, struct{}]
}

// In adds an object to the queue if it is not already present.
func (q *deduplicatedObjectQueue) In(obj fyne.CanvasObject) {
	_, exists := q.dedup.Load(obj)
	if exists {
		return
	}

	q.queue.In(obj)
	q.dedup.Store(obj, struct{}{})
}

// Out removes and returns the next object from the queue.
// It assumes that the whole queue is drained and defers clearing
// the deduplication map until it is empty.
func (q *deduplicatedObjectQueue) Out() fyne.CanvasObject {
	if q.queue.Len() == 0 {
		q.dedup.Clear()
		return nil
	}

	return q.queue.Out()
}

// Len returns the number of elements in the queue.
func (q *deduplicatedObjectQueue) Len() uint64 {
	return q.queue.Len()
}

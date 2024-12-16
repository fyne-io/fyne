package common

import (
	"fyne.io/fyne/v2/internal/async"
)

// Window defines common functionality for windows.
type Window struct {
	eventQueue *async.UnboundedFuncChan
}

// DestroyEventQueue destroys the event queue.
func (w *Window) DestroyEventQueue() {
	w.eventQueue.Close()
}

// InitEventQueue initializes the event queue.
func (w *Window) InitEventQueue() {
	// This channel should be closed when the window is closed.
	w.eventQueue = async.NewUnboundedFuncChan()
}

// QueueEvent uses this method to queue up a callback that handles an event. This ensures
// user interaction events for a given window are processed in order.
func (w *Window) QueueEvent(fn func()) {
	w.eventQueue.In() <- fn
}

// ProcessEventQueue runs all the items in the event queue, returning once it is empty again.
func (w *Window) ProcessEventQueue() {
	for {
		select {
		case fn := <-w.eventQueue.Out():
			fn()
		default:
			return
		}
	}
}

// WaitForEvents wait for all the events.
func (w *Window) WaitForEvents() {
	done := DonePool.Get()
	defer DonePool.Put(done)

	w.eventQueue.In() <- func() { done <- struct{}{} }
	<-done
}

var DonePool = async.Pool[chan struct{}]{
	New: func() chan struct{} {
		return make(chan struct{})
	},
}

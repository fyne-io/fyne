package common

import (
	"sync"

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

// RunEventQueue runs the event queue. This should called inside a go routine.
// This function blocks.
func (w *Window) RunEventQueue() {
	for fn := range w.eventQueue.Out() {
		fn()
	}
}

// WaitForEvents wait for all the events.
func (w *Window) WaitForEvents() {
	done := donePool.Get().(chan struct{})
	defer donePool.Put(done)

	w.eventQueue.In() <- func() { done <- struct{}{} }
	<-done
}

var donePool = sync.Pool{
	New: func() interface{} {
		return make(chan struct{})
	},
}

package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
)

// Window defines common functionality for windows.
type Window struct {
	eventQueue *async.UnboundedEventFuncChan
}

// DestroyEventQueue destroys the event queue.
func (w *Window) DestroyEventQueue() {
	w.eventQueue.Close()
}

// InitEventQueue initializes the event queue.
func (w *Window) InitEventQueue() {
	// This channel should be closed when the window is closed.
	w.eventQueue = async.NewUnboundedEventFuncChan()
}

// QueueEvent uses this method to queue up a callback that handles an event. This ensures
// user interaction events for a given window are processed in order.
func (w *Window) QueueEvent(fn fyne.EventFunc) {
	w.eventQueue.In() <- fn
}

// RunEventQueue runs the event queue. This should called inside a go routine.
// This function blocks.
func (w *Window) RunEventQueue() {
	for evfn := range w.eventQueue.Out() {
		if dragfn, ok := evfn.(*fyne.DragEventFunc); ok {
			evfn = nil

		L:
			for {
				select {
				case nevfn := <-w.eventQueue.Out():
					ndragfn, ok := nevfn.(*fyne.DragEventFunc)
					if !ok {
						evfn = nevfn
						break L
					}
					dragfn = ndragfn
				default:
					break L
				}
			}

			dragfn.Execute()
			if evfn != nil {
				evfn.Execute()
			}
			continue
		}

		evfn.Execute()
	}
}

// WaitForEvents wait for all the events.
func (w *Window) WaitForEvents() {
	done := DonePool.Get()
	defer DonePool.Put(done)

	w.eventQueue.In() <- fyne.SimpleEventFunc(func() { done <- struct{}{} })
	<-done
}

var DonePool = async.Pool[chan struct{}]{
	New: func() chan struct{} {
		return make(chan struct{})
	},
}

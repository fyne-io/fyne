package common

import (
	"sync"

	"fyne.io/fyne/v2"
)

// Window defines common functionality for windows.
type Window struct {
	eventLock  sync.RWMutex
	eventQueue chan func()
	eventWait  sync.WaitGroup
}

// DestroyEventQueue destroys the event queue.
func (w *Window) DestroyEventQueue() {
	w.eventLock.RLock()
	queue := w.eventQueue
	w.eventLock.RUnlock()

	// finish serial event queue and nil it so we don't panic if window.closed() is called twice.
	if queue != nil {
		w.WaitForEvents()

		w.eventLock.Lock()
		close(w.eventQueue)
		w.eventQueue = nil
		w.eventLock.Unlock()
	}
}

// InitEventQueue initializes the event queue.
func (w *Window) InitEventQueue() {
	// This channel should be closed when the window is closed.
	w.eventQueue = make(chan func(), 1024)
}

// QueueEvent uses this method to queue up a callback that handles an event. This ensures
// user interaction events for a given window are processed in order.
func (w *Window) QueueEvent(fn func()) {
	w.eventWait.Add(1)
	select {
	case w.eventQueue <- fn:
	default:
		fyne.LogError("EventQueue full, perhaps a callback blocked the event handler", nil)
	}
}

// RunEventQueue runs the event queue. This should called inside a go routine.
// This function blocks.
func (w *Window) RunEventQueue() {
	w.eventLock.Lock()
	queue := w.eventQueue
	w.eventLock.Unlock()

	for fn := range queue {
		fn()
		w.eventWait.Done()
	}
}

// WaitForEvents wait for all the events.
func (w *Window) WaitForEvents() {
	w.eventWait.Wait()
}

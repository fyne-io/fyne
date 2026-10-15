package glfw

import (
	"sync"
	"sync/atomic"
	"time"
)

type mainLoopEventKind uint8

const (
	mainLoopQuit mainLoopEventKind = iota
	mainLoopFunction
	mainLoopFrame
)

type mainLoopEvent struct {
	kind      mainLoopEventKind
	function  funcData
	frameTime time.Time
}

// nextMainLoopEvent waits for the next event. With a nil wake it blocks on the
// channels; otherwise it processes native events until work arrives or a frame is due.
func nextMainLoopEvent(done <-chan struct{}, functions <-chan funcData, frames <-chan time.Time,
	nextFrame time.Time, wake *mainLoopWake,
) mainLoopEvent {
	for wake != nil {
		wake.rearm()
		select {
		case <-done:
			return mainLoopEvent{kind: mainLoopQuit}
		case f := <-functions:
			return mainLoopEvent{kind: mainLoopFunction, function: f}
		case tick := <-frames:
			return mainLoopEvent{kind: mainLoopFrame, frameTime: tick}
		default:
		}
		timeout := time.Until(nextFrame)
		if timeout <= 0 {
			break // The frame tick is due, so block for it rather than spin.
		}
		wake.wait(timeout)
	}
	select {
	case <-done:
		return mainLoopEvent{kind: mainLoopQuit}
	case f := <-functions:
		return mainLoopEvent{kind: mainLoopFunction, function: f}
	case tick := <-frames:
		return mainLoopEvent{kind: mainLoopFrame, frameTime: tick}
	}
}

// mainLoopWake interrupts a native event wait when work is queued from other goroutines.
// Native wakeups must finish before GLFW is terminated.
type mainLoopWake struct {
	wait    func(time.Duration)
	mu      sync.Mutex
	post    func()
	pending atomic.Bool
}

func (w *mainLoopWake) start(post func()) {
	w.mu.Lock()
	w.post = post
	w.pending.Store(false)
	w.mu.Unlock()
}

func (w *mainLoopWake) stop() {
	w.mu.Lock()
	w.post = nil
	w.mu.Unlock()
}

// notify posts at most one native wakeup until the main loop calls rearm.
func (w *mainLoopWake) notify() {
	if w.pending.Swap(true) {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.post != nil {
		w.post()
	}
}

// rearm must be called before checking for work, so that work queued after the
// check always posts a wakeup.
func (w *mainLoopWake) rearm() {
	w.pending.Store(false)
}

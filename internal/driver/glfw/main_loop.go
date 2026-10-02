package glfw

import (
	"sync"
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

func nextMainLoopEvent(done <-chan struct{}, functions <-chan funcData, frames <-chan time.Time,
	nextFrame time.Time, waitEvents func(time.Duration),
) mainLoopEvent {
	for {
		if waitEvents == nil {
			select {
			case <-done:
				return mainLoopEvent{kind: mainLoopQuit}
			case f := <-functions:
				return mainLoopEvent{kind: mainLoopFunction, function: f}
			case tick := <-frames:
				return mainLoopEvent{kind: mainLoopFrame, frameTime: tick}
			}
		}
		select {
		case <-done:
			return mainLoopEvent{kind: mainLoopQuit}
		case f := <-functions:
			return mainLoopEvent{kind: mainLoopFunction, function: f}
		case tick := <-frames:
			return mainLoopEvent{kind: mainLoopFrame, frameTime: tick}
		default:
			waitEvents(time.Until(nextFrame))
		}
	}
}

// Native wakeups must finish before GLFW is terminated.
type mainLoopWake struct {
	mu   sync.Mutex
	post func()
}

func (w *mainLoopWake) start(post func()) {
	w.mu.Lock()
	w.post = post
	w.mu.Unlock()
}

func (w *mainLoopWake) stop() {
	w.mu.Lock()
	w.post = nil
	w.mu.Unlock()
}

func (w *mainLoopWake) notify() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.post != nil {
		w.post()
	}
}

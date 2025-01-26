package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
)

var _ fyne.Lifecycle = (*Lifecycle)(nil)

// Lifecycle represents the various phases that an app can transition through.
//
// Since: 2.1
type Lifecycle struct {
	onForeground func()
	onBackground func()
	onStarted    func()
	onStopped    func()

	onStoppedHookExecuted func()

	eventQueue *async.UnboundedChan[func()]
}

// SetOnStoppedHookExecuted is an internal function that lets Fyne schedule a clean-up after
// the user-provided stopped hook. It should only be called once during an application start-up.
func (l *Lifecycle) SetOnStoppedHookExecuted(f func()) {
	l.onStoppedHookExecuted = f
}

// SetOnEnteredForeground hooks into the app becoming foreground.
func (l *Lifecycle) SetOnEnteredForeground(f func()) {
	l.onForeground = f
}

// SetOnExitedForeground hooks into the app having moved to the background.
// Depending on the platform it may still be  visible but will not receive keyboard events.
// On some systems hover or desktop mouse move events may still occur.
func (l *Lifecycle) SetOnExitedForeground(f func()) {
	l.onBackground = f
}

// SetOnStarted hooks into an event that says the app is now running.
func (l *Lifecycle) SetOnStarted(f func()) {
	l.onStarted = f
}

// SetOnStopped hooks into an event that says the app is no longer running.
func (l *Lifecycle) SetOnStopped(f func()) {
	l.onStopped = f
}

// OnEnteredForeground returns the focus gained hook, if one is registered.
func (l *Lifecycle) OnEnteredForeground() func() {
	return l.onForeground
}

// OnExitedForeground returns the focus lost hook, if one is registered.
func (l *Lifecycle) OnExitedForeground() func() {
	return l.onBackground
}

// OnStarted returns the started hook, if one is registered.
func (l *Lifecycle) OnStarted() func() {
	return l.onStarted
}

// OnStopped returns the stopped hook, if one is registered.
func (l *Lifecycle) OnStopped() func() {
	stopped := l.onStopped
	stopHook := l.onStoppedHookExecuted
	if stopped == nil && stopHook == nil {
		return nil
	}

	if stopHook == nil {
		return stopped
	}

	if stopped == nil {
		return stopHook
	}

	// we have a stopped handle and the onStoppedHook
	return func() {
		stopped()
		stopHook()
	}
}

// DestroyEventQueue destroys the event queue.
func (l *Lifecycle) DestroyEventQueue() {
	l.eventQueue.Close()
}

// InitEventQueue initializes the event queue.
func (l *Lifecycle) InitEventQueue() {
	// This channel should be closed when the window is closed.
	l.eventQueue = async.NewUnboundedChan[func()]()
}

// QueueEvent uses this method to queue up a callback that handles an event. This ensures
// user interaction events for a given window are processed in order.
func (l *Lifecycle) QueueEvent(fn func()) {
	l.eventQueue.In() <- fn
}

// RunEventQueue runs the event queue. This should called inside a go routine.
// This function blocks.
func (l *Lifecycle) RunEventQueue(run func(func(), bool)) {
	for fn := range l.eventQueue.Out() {
		run(fn, true)
	}
}

// WaitForEvents wait for all the events.
func (l *Lifecycle) WaitForEvents() {
	done := make(chan struct{})

	l.eventQueue.In() <- func() { done <- struct{}{} }
	<-done
}

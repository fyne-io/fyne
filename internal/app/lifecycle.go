package app

import (
	"sync"

	"fyne.io/fyne/v2"
)

var _ fyne.Lifecycle = (*Lifecycle)(nil)

// Lifecycle represents the various phases that an app can transition through.
//
// Since: 2.1
type Lifecycle struct {
	onForeground, onBackground func()
	onStarted, onStopped       func()

	mux sync.Mutex
}

// SetOnEnteredForeground hooks into the the app becoming foreground.
func (l *Lifecycle) SetOnEnteredForeground(f func()) {
	l.mux.Lock()
	l.onForeground = f
	l.mux.Unlock()
}

// SetOnExitedForeground hooks into the app having moved to the background.
// Depending on the platform it may still be  visible but will not receive keyboard events.
// On some systems hover or desktop mouse move events may still occur.
func (l *Lifecycle) SetOnExitedForeground(f func()) {
	l.mux.Lock()
	l.onBackground = f
	l.mux.Unlock()
}

// SetOnStarted hooks into an event that says the app is now running.
func (l *Lifecycle) SetOnStarted(f func()) {
	l.mux.Lock()
	l.onStarted = f
	l.mux.Unlock()
}

// SetOnStopped hooks into an event that says the app is no longer running.
func (l *Lifecycle) SetOnStopped(f func()) {
	l.mux.Lock()
	l.onStopped = f
	l.mux.Unlock()
}

// TriggerEnteredForeground will call the focus gained hook, if one is registered.
func (l *Lifecycle) TriggerEnteredForeground() {
	l.mux.Lock()
	f := l.onForeground
	l.mux.Unlock()

	if f != nil {
		f()
	}
}

// TriggerExitedForeground will call the focus lost hook, if one is registered.
func (l *Lifecycle) TriggerExitedForeground() {
	l.mux.Lock()
	f := l.onBackground
	l.mux.Unlock()

	if f != nil {
		f()
	}
}

// TriggerStarted will call the started hook, if one is registered.
func (l *Lifecycle) TriggerStarted() {
	l.mux.Lock()
	f := l.onStarted
	l.mux.Unlock()

	if f != nil {
		f()
	}
}

// TriggerStopped will call the stopped hook, if one is registered.
func (l *Lifecycle) TriggerStopped() {
	l.mux.Lock()
	f := l.onStopped
	l.mux.Unlock()

	if f != nil {
		f()
	}
}

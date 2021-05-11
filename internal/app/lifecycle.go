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
	onFocusGained, onFocusLost func()
	onStarted, onStopped       func()

	mux sync.Mutex
}

// SetOnFocusGained hooks into the the app becoming foreground.
func (l *Lifecycle) SetOnFocusGained(f func()) {
	l.mux.Lock()
	l.onFocusGained = f
	l.mux.Unlock()
}

// SetOnFocusLost hooks into the app losing foreground focus.
func (l *Lifecycle) SetOnFocusLost(f func()) {
	l.mux.Lock()
	l.onFocusLost = f
}

// SetOnStarted hooks into an event that says the app is now running.
func (l *Lifecycle) SetOnStarted(f func()) {
	l.mux.Lock()
	l.onStarted = f
}

// SetOnStopped hooks into an event that says the app is no longer running.
func (l *Lifecycle) SetOnStopped(f func()) {
	l.mux.Lock()
	l.onStopped = f
}

// TriggerFocusGained will call the focus gained hook, if one is registered.
func (l *Lifecycle) TriggerFocusGained() {
	l.mux.Lock()
	f := l.onFocusGained
	l.mux.Unlock()

	if f != nil {
		f()
	}
}

// TriggerFocusLost will call the focus lost hook, if one is registered.
func (l *Lifecycle) TriggerFocusLost() {
	l.mux.Lock()
	f := l.onFocusLost
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

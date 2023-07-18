package app

import (
	"sync/atomic"

	"fyne.io/fyne/v2"
)

var _ fyne.Lifecycle = (*Lifecycle)(nil)

// Lifecycle represents the various phases that an app can transition through.
//
// Since: 2.1
type Lifecycle struct {
	onForeground atomic.Value // func()
	onBackground atomic.Value // func()
	onStarted    atomic.Value // func()
	onStopped    atomic.Value // func()

	onStoppedHookExecuted func()
}

// SetOnStoppedHookExecuted is an internal function that lets Fyne schedule a clean-up after
// the user-provided stopped hook. It should only be called once during an application start-up.
func (l *Lifecycle) SetOnStoppedHookExecuted(f func()) {
	l.onStoppedHookExecuted = f
}

// SetOnEnteredForeground hooks into the the app becoming foreground.
func (l *Lifecycle) SetOnEnteredForeground(f func()) {
	l.onForeground.Store(f)
}

// SetOnExitedForeground hooks into the app having moved to the background.
// Depending on the platform it may still be  visible but will not receive keyboard events.
// On some systems hover or desktop mouse move events may still occur.
func (l *Lifecycle) SetOnExitedForeground(f func()) {
	l.onBackground.Store(f)
}

// SetOnStarted hooks into an event that says the app is now running.
func (l *Lifecycle) SetOnStarted(f func()) {
	l.onStarted.Store(f)
}

// SetOnStopped hooks into an event that says the app is no longer running.
func (l *Lifecycle) SetOnStopped(f func()) {
	l.onStopped.Store(f)
}

// TriggerEnteredForeground will call the focus gained hook, if one is registered.
func (l *Lifecycle) TriggerEnteredForeground() {
	f := l.onForeground.Load()
	if ff, ok := f.(func()); ok && ff != nil {
		ff()
	}
}

// TriggerExitedForeground will call the focus lost hook, if one is registered.
func (l *Lifecycle) TriggerExitedForeground() {
	f := l.onBackground.Load()
	if ff, ok := f.(func()); ok && ff != nil {
		ff()
	}
}

// TriggerStarted will call the started hook, if one is registered.
func (l *Lifecycle) TriggerStarted() {
	f := l.onStarted.Load()
	if ff, ok := f.(func()); ok && ff != nil {
		ff()
	}
}

// TriggerStopped will call the stopped hook, if one is registered,
// and an internal stopped hook after that.
func (l *Lifecycle) TriggerStopped() {
	f := l.onStopped.Load()
	if ff, ok := f.(func()); ok && ff != nil {
		ff()
	}
	if l.onStoppedHookExecuted != nil {
		l.onStoppedHookExecuted()
	}
}

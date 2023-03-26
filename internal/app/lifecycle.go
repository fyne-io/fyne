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

// OnEnteredForeground returns the focus gained hook, if one is registered.
func (l *Lifecycle) OnEnteredForeground() func() {
	f := l.onForeground.Load()
	if ff, ok := f.(func()); ok {
		return ff
	}
	return nil
}

// OnExitedForeground returns the focus lost hook, if one is registered.
func (l *Lifecycle) OnExitedForeground() func() {
	f := l.onBackground.Load()
	if ff, ok := f.(func()); ok {
		return ff
	}
	return nil
}

// OnStarted returns the started hook, if one is registered.
func (l *Lifecycle) OnStarted() func() {
	f := l.onStarted.Load()
	if ff, ok := f.(func()); ok {
		return ff
	}
	return nil
}

// OnStopped returns the stopped hook, if one is registered.
func (l *Lifecycle) OnStopped() func() {
	f := l.onStopped.Load()
	if ff, ok := f.(func()); ok && ff != nil {
		return ff
	}
	return nil
}

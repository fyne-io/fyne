//go:build ci

package fyne

import "sync/atomic"

var app atomic.Pointer[App]

// SetCurrentApp is an internal function to set the app instance currently running.
func SetCurrentApp(current App) {
	app.Store(&current)
}

// CurrentApp returns the current application, for which there is only 1 per process.
func CurrentApp() App {
	val := app.Load()
	if val == nil {
		LogError("Attempt to access current Fyne app when none is started", nil)
		return nil
	}
	return *val
}

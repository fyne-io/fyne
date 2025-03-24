//go:build !ci

package fyne

var app App

// SetCurrentApp is an internal function to set the app instance currently running.
func SetCurrentApp(current App) {
	app = current
}

// CurrentApp returns the current application, for which there is only 1 per process.
func CurrentApp() App {
	if app == nil {
		LogError("Attempt to access current Fyne app when none is started", nil)
	}
	return app
}

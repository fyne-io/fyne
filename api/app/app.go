// Package app defines Fyne APIs for managing a graphical application.
package app

import "github.com/fyne-io/fyne/api/ui"

// An App is the definition of a graphical application.
// Apps can have multiple windows, it will exit when the first windows to be
// shown is closed. You can also cause the app to exit by calling Quit().
type App interface {
	// Create a new window for the application.
	// The first window to open is considered the "master" and when closed
	// the application will exit
	NewWindow(title string) ui.Window

	// Open a URL in the default browser application
	OpenURL(url string)

	// Calling Quit on the application will cause the application to exit
	// cleanly, closing all open windows.
	Quit()
}

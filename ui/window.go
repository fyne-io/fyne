// Package ui defines the user interface APIs in fyne
package ui

// Window describes a user interface window. Depending on the platform an app
// may have many windows or just the one.
type Window interface {
	// Title returns the current window title.
	// This is typically displayed in the window decorations.
	Title() string
	// SetTitle updates the current title of the window
	SetTitle(string)
	// Show the window on screen
	Show()
	// Hide the window from the user.
	// This will not destroy the window or cause the app to exit.
	Hide()
	// Close the window.
	// If it is the only open window, or the "master" window the app will Quit.
	Close()

	// Canvas returns the canvas context to render in the window
	Canvas() Canvas
}

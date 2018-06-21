package fyne

// Window describes a user interface window. Depending on the platform an app
// may have many windows or just the one.
type Window interface {
	// Title returns the current window title.
	// This is typically displayed in the window decorations.
	Title() string
	// SetTitle updates the current title of the window
	SetTitle(string)

	// Fullscreen returns whether or not this window is currently full screen
	Fullscreen() bool
	// SetFullscreen changes the requested fullscreen property
	// true for a fullscreen window and false to unset this.
	SetFullscreen(bool)

	// Show the window on screen
	Show()
	// Hide the window from the user.
	// This will not destroy the window or cause the app to exit.
	Hide()
	// Close the window.
	// If it is the only open window, or the "master" window the app will Quit.
	Close()

	// Content returns the content of this window
	Content() CanvasObject
	// SetContent sets the content of this window
	SetContent(CanvasObject)
	// Canvas returns the canvas context to render in the window
	Canvas() Canvas
}

// GetWindow returns the window containing the passed Canvas.
// It will return nil if the Canvas is not associated with any existing windows.
func GetWindow(c Canvas) Window {
	for _, window := range GetDriver().AllWindows() {
		if window.Canvas() != nil && window.Canvas() == c {
			return window
		}
	}

	return nil
}

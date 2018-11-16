package fyne

// Window describes a user interface window. Depending on the platform an app
// may have many windows or just the one.
type Window interface {
	// Title returns the current window title.
	// This is typically displayed in the window decorations.
	Title() string
	// SetTitle updates the current title of the window
	SetTitle(string)

	// FullScreen returns whether or not this window is currently full screen
	FullScreen() bool
	// SetFullScreen changes the requested fullScreen property
	// true for a fullScreen window and false to unset this.
	SetFullScreen(bool)

	// FixedSize returns whether or not this window should disable resizing.
	FixedSize() bool
	// SetFixedSize sets a hint that states whether the window should be a fixed
	// size or allow resizing.
	SetFixedSize(bool)

	SetOnClosed(func())

	// Show the window on screen
	Show()
	// Hide the window from the user.
	// This will not destroy the window or cause the app to exit.
	Hide()
	// Close the window.
	// If it is the only open window, or the "master" window the app will Quit.
	Close()

	// ShowAndRun is a shortcut to show the window and then run the application.
	// This should be called near the end of a main() function as it will block.
	ShowAndRun()

	// Content returns the content of this window
	Content() CanvasObject
	// SetContent sets the content of this window
	SetContent(CanvasObject)
	// Canvas returns the canvas context to render in the window
	Canvas() Canvas
}

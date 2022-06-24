package fyne

// Window describes a user interface window. Depending on the platform an app
// may have many windows or just the one.
type Window interface {
	// Title returns the current window title.
	// This is typically displayed in the window decorations.
	Title() string
	// SetTitle updates the current title of the window.
	SetTitle(string)

	// FullScreen returns whether or not this window is currently full screen.
	FullScreen() bool
	// SetFullScreen changes the requested fullScreen property
	// true for a fullScreen window and false to unset this.
	SetFullScreen(bool)

	// Resize this window to the requested content size.
	// The result may not be exactly as desired due to various desktop or
	// platform constraints.
	Resize(Size)

	// RequestFocus attempts to raise and focus this window.
	// This should only be called when you are sure the user would want this window
	// to steal focus from any current focused window.
	RequestFocus()

	// FixedSize returns whether or not this window should disable resizing.
	FixedSize() bool
	// SetFixedSize sets a hint that states whether the window should be a fixed
	// size or allow resizing.
	SetFixedSize(bool)

	// CenterOnScreen places a window at the center of the monitor
	// the Window object is currently positioned on.
	CenterOnScreen()

	// Padded, normally true, states whether the window should have inner
	// padding so that components do not touch the window edge.
	Padded() bool
	// SetPadded allows applications to specify that a window should have
	// no inner padding. Useful for fullscreen or graphic based applications.
	SetPadded(bool)

	// Icon returns the window icon, this is used in various ways
	// depending on operating system.
	// Most commonly this is displayed on the window border or task switcher.
	Icon() Resource

	// SetIcon sets the icon resource used for this window.
	// If none is set should return the application icon.
	SetIcon(Resource)

	// SetMaster indicates that closing this window should exit the app
	SetMaster()

	// MainMenu gets the content of the window's top level menu.
	MainMenu() *MainMenu

	// SetMainMenu adds a top level menu to this window.
	// The way this is rendered will depend on the loaded driver.
	SetMainMenu(*MainMenu)

	// SetOnClosed sets a function that runs when the window is closed.
	SetOnClosed(func())

	// SetCloseIntercept sets a function that runs instead of closing if defined.
	// Close() should be called explicitly in the interceptor to close the window.
	//
	// Since: 1.4
	SetCloseIntercept(func())

	// Show the window on screen.
	Show()
	// Hide the window from the user.
	// This will not destroy the window or cause the app to exit.
	Hide()
	// Close the window.
	// If it is he "master" window the app will Quit.
	// If it is the only open window and no menu is set via [desktop.App]
	// SetSystemTrayMenu the app will also Quit.
	Close()

	// ShowAndRun is a shortcut to show the window and then run the application.
	// This should be called near the end of a main() function as it will block.
	ShowAndRun()

	// Content returns the content of this window.
	Content() CanvasObject
	// SetContent sets the content of this window.
	SetContent(CanvasObject)
	// Canvas returns the canvas context to render in the window.
	// This can be useful to set a key handler for the window, for example.
	Canvas() Canvas

	// Clipboard returns the system clipboard
	Clipboard() Clipboard
}

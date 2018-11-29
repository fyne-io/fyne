package fyne

// An App is the definition of a graphical application.
// Apps can have multiple windows, it will exit when the first windows to be
// shown is closed. You can also cause the app to exit by calling Quit().
// To start an application you need to call Run() somewhere in your main() function.
// Alternatively use the window.ShowAndRun() function for your main window.
type App interface {
	// Create a new window for the application.
	// The first window to open is considered the "master" and when closed
	// the application will exit
	NewWindow(title string) Window

	// Open a URL in the default browser application
	OpenURL(url string)

	// Run the application - this starts the event loop and waits until Quit()
	// is called or the last window closes.
	// This should be called near the end of a main() function as it will block.
	Run()

	// Calling Quit on the application will cause the application to exit
	// cleanly, closing all open windows.
	Quit()

	// Driver returns the driver that is rendering this application.
	// Typically not needed for day to day work, mostly internal functionality.
	Driver() Driver
}

var app App

// SetCurrentApp is an internal function to set the app instance currently running
func SetCurrentApp(current App) {
	app = current
}

// CurrentApp returns the current application, for which there is only 1 per process.
func CurrentApp() App {
	return app
}

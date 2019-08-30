package fyne

import "net/url"

// An App is the definition of a graphical application.
// Apps can have multiple windows, it will exit when the first window to be
// shown is closed. You can also cause the app to exit by calling Quit().
// To start an application you need to call Run() somewhere in your main() function.
// Alternatively use the window.ShowAndRun() function for your main window.
type App interface {
	// Create a new window for the application.
	// The first window to open is considered the "master" and when closed
	// the application will exit.
	NewWindow(title string) Window

	// Open a URL in the default browser application.
	OpenURL(url *url.URL) error

	// Icon returns the application icon, this is used in various ways
	// depending on operating system.
	// This is also the default icon for new windows.
	Icon() Resource

	// SetIcon sets the icon resource used for this application instance.
	SetIcon(Resource)

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

	// UniqueID returns the application unique identifier, if set.
	// This must be set for use of the Preferences() functions... see NewWithId(string)
	UniqueID() string

	// Settings return the globally set settings, determining theme and so on.
	Settings() Settings

	// Preferences returns the application preferences, used for storing configuration and state
	Preferences() Preferences
}

var app App

// SetCurrentApp is an internal function to set the app instance currently running.
func SetCurrentApp(current App) {
	app = current
}

// CurrentApp returns the current application, for which there is only 1 per process.
func CurrentApp() App {
	return app
}

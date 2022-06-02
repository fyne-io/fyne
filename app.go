package fyne

import (
	"net/url"
	"sync/atomic"
)

// An App is the definition of a graphical application.
// Apps can have multiple windows, by default they will exit when all windows
// have been closed. This can be modified using SetMaster() or SetCloseIntercept().
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
	// This function does no thing on a mobile device as the application lifecycle is
	// managed by the operating system.
	Quit()

	// Driver returns the driver that is rendering this application.
	// Typically not needed for day to day work, mostly internal functionality.
	Driver() Driver

	// UniqueID returns the application unique identifier, if set.
	// This must be set for use of the Preferences() functions... see NewWithId(string)
	UniqueID() string

	// SendNotification sends a system notification that will be displayed in the operating system's notification area.
	SendNotification(*Notification)

	// Settings return the globally set settings, determining theme and so on.
	Settings() Settings

	// Preferences returns the application preferences, used for storing configuration and state
	Preferences() Preferences

	// Storage returns a storage handler specific to this application.
	Storage() Storage

	// Lifecycle returns a type that allows apps to hook in to lifecycle events.
	//
	// Since: 2.1
	Lifecycle() Lifecycle

	// Metadata returns the application metadata that was set at compile time.
	//
	// Since: 2.2
	Metadata() AppMetadata
}

// app contains an App variable, but due to atomic.Value restrictions on
// interfaces we need to use an indirect type, i.e. appContainer.
var app atomic.Value // appContainer

// appContainer is a dummy container that holds an App instance. This
// struct exists to guarantee that atomic.Value can store objects with
// same type.
type appContainer struct {
	current App
}

// SetCurrentApp is an internal function to set the app instance currently running.
func SetCurrentApp(current App) {
	app.Store(appContainer{current})
}

// CurrentApp returns the current application, for which there is only 1 per process.
func CurrentApp() App {
	val := app.Load()
	if val == nil {
		LogError("Attempt to access current Fyne app when none is started", nil)
		return nil
	}
	return (val).(appContainer).current
}

// AppMetadata captures the build metadata for an application.
//
// Since: 2.2
type AppMetadata struct {
	// ID is the unique ID of this application, used by many distribution platforms.
	ID string
	// Name is the human friendly name of this app.
	Name string
	// Version represents the version of this application, normally following semantic versioning.
	Version string
	// Build is the build number of this app, some times appended to the version number.
	Build int
	// Icon contains, if present, a resource of the icon that was bundled at build time.
	Icon Resource
}

// Lifecycle represents the various phases that an app can transition through.
//
// Since: 2.1
type Lifecycle interface {
	// SetOnEnteredForeground hooks into the app becoming foreground and gaining focus.
	SetOnEnteredForeground(func())
	// SetOnExitedForeground hooks into the app losing input focus and going into the background.
	SetOnExitedForeground(func())
	// SetOnStarted hooks into an event that says the app is now running.
	SetOnStarted(func())
	// SetOnStopped hooks into an event that says the app is no longer running.
	SetOnStopped(func())
}

package fyne

import (
	"errors"
	"net/url"
	"sync/atomic"
	"time"
)

// ErrPasswordRequired is returned when a password is needed to protect data but none was provided,
// for example by [App.SecretPreferences] when the operating system has no secure storage that can be used.
// The caller should obtain a password and try again.
//
// Since: 2.9
var ErrPasswordRequired = errors.New("a password is required")

// ErrPasswordIncorrect is returned when a password was provided but the protected data could not be unlocked
// with it, for example by [App.SecretPreferences] when the stored values were saved with a different password.
//
// Since: 2.9
var ErrPasswordIncorrect = errors.New("the password is incorrect")

// An App is the definition of a graphical application.
// Apps can have multiple windows, by default they will exit when all windows
// have been closed. This can be modified using SetMaster or SetCloseIntercept.
// To start an application you need to call Run somewhere in your main function.
// Alternatively use the [fyne.io/fyne/v2.Window.ShowAndRun] function for your main window.
type App interface {
	// Create a new window for the application.
	// The first window to open is considered the "master" and when closed
	// the application will exit.
	NewWindow(title string) Window

	// Open a URL in the default browser application.
	OpenURL(*url.URL) error

	// Icon returns the application icon, this is used in various ways
	// depending on operating system.
	// This is also the default icon for new windows.
	Icon() Resource

	// SetIcon sets the icon resource used for this application instance.
	SetIcon(Resource)

	// Run the application - this starts the event loop and waits until [App.Quit]
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
	// This must be set for use of the [App.Preferences]. see [NewWithID].
	UniqueID() string

	// SendNotification sends a system notification that will be displayed in the operating system's notification area.
	SendNotification(*Notification)

	// ScheduleNotification queues a notification for delivery at the given time.
	// On platforms with a system scheduler (iOS, Android, macOS, Windows) the request is
	// passed to the OS. On other platforms the app schedules delivery in-process and
	// persists the request so that an outstanding schedule survives app restarts.
	//
	// Returns the scheduled notification, whose [ScheduledNotification.ID] can be used to
	// cancel the delivery via [App.CancelScheduledNotification]. An error is returned if
	// the request could not be queued (for example a delivery time in the past on a
	// platform that does not allow it).
	//
	// Since: 2.8
	ScheduleNotification(n *Notification, deliverAt time.Time) (*ScheduledNotification, error)

	// CancelScheduledNotification removes a previously scheduled notification by its ID.
	// It is safe to call with an unknown ID or after the notification has already fired.
	//
	// Since: 2.8
	CancelScheduledNotification(id string) error

	// Settings return the globally set settings, determining theme and so on.
	Settings() Settings

	// Preferences returns the application preferences, used for storing configuration and state
	Preferences() Preferences

	// SecretPreferences returns a preference store for sensitive values such as tokens or passwords.
	// It has the same API as [App.Preferences] but the values are kept out of the plain text
	// preferences file. Where the operating system offers secure storage it is used:
	// the Keychain on iOS and macOS (when running from an app bundle), the Android Keystore
	// and the Data Protection API on Windows.
	//
	// On other platforms, and in development builds, there is no secure storage to hold a key, so the
	// values are encrypted on disk with a key derived from the passed password.
	// The password may be nil, in which case [ErrPasswordRequired] is returned if this fallback is
	// needed. The app can then obtain a password, through secure means or by asking the user, and call this
	// again. The password is not used when operating system secure storage is available, or once the store
	// has been returned by an earlier call. It is not retained, the caller should clear its content
	// once it is no longer needed.
	//
	// An error is returned, with no store, if the stored values could not be read - [ErrPasswordIncorrect]
	// when the password does not match the one the values were saved with. The stored values are left untouched.
	//
	// As with [App.Preferences] a unique ID must be set for values to be persisted, see [NewWithID].
	//
	// Since: 2.9
	SecretPreferences(password []byte) (Preferences, error)

	// Storage returns a storage handler specific to this application.
	Storage() Storage

	// Lifecycle returns a type that allows apps to hook in to lifecycle events.
	//
	// Since: 2.1
	Lifecycle() Lifecycle

	// Metadata returns the application metadata that was set at compile time.
	// The items of metadata are available after "fyne package" or when running "go run"
	// Building with "go build" may cause this to be unavailable.
	//
	// Since: 2.2
	Metadata() AppMetadata

	// CloudProvider returns the current app cloud provider,
	// if one has been registered by the developer or chosen by the user.
	//
	// Since: 2.3
	CloudProvider() CloudProvider // get the (if any) configured provider

	// SetCloudProvider allows developers to specify how this application should integrate with cloud services.
	// See [fyne.io/cloud] package for implementation details.
	//
	// Since: 2.3
	SetCloudProvider(CloudProvider) // configure cloud for this app

	// Clipboard returns the system clipboard.
	//
	// Since: 2.6
	Clipboard() Clipboard

	// Cache returns a cache handler specific to this application.
	//
	// Since: 2.8
	Cache() Cache
}

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
	// Release if true this binary was build in release mode
	// Since: 2.3
	Release bool
	// Custom contain the custom metadata defined either in FyneApp.toml or on the compile command line
	// Since: 2.3
	Custom map[string]string
	// Migrations allows an app to opt into features before they are standard
	// Since: 2.6
	Migrations map[string]bool
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

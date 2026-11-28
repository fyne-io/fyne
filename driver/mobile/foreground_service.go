package mobile

// ForegroundService is implemented by a [fyne.io/fyne/v2.App] that can ask the
// operating system to keep its process alive while it is in the background.
// Currently only Android apps implement it, so check for it with a type
// assertion:
//
//	if fs, ok := fyne.CurrentApp().(mobile.ForegroundService); ok {
//		err := fs.StartForegroundService("Title", "Content")
//	}
//
// Since: 2.9
type ForegroundService interface {
	// StartForegroundService asks the operating system to keep this app's
	// process alive while it is in the background. Android requires the user
	// to be told, so the given title and content are shown in an ongoing
	// notification for as long as the app is elevated.
	//
	// The app must opt in by supplying its own AndroidManifest.xml, next to
	// the Go files of the main package, that declares the
	// org.golang.app.GoForegroundService service along with an
	// android:foregroundServiceType and the permissions that type requires.
	// The manifest that "fyne package" generates does not declare the service.
	//
	// Calling it again while the app is elevated updates the notification.
	// An error is returned if the system refuses to start the service, for
	// example because the manifest does not declare it or a permission its
	// type requires is missing. The reason is also written to the log.
	StartForegroundService(title, content string) error

	// StopForegroundService returns this app to a normal background state
	// and removes the notification that StartForegroundService showed.
	StopForegroundService()
}

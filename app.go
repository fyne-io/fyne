package fyne

// An App is the definition of a graphical application.
// Apps can have multiple windows, it will exit when the first windows to be
// shown is closed. You can also cause the app to exit by calling Quit().
type App interface {
	// Create a new window for the application.
	// The first window to open is considered the "master" and when closed
	// the application will exit
	NewWindow(title string) Window

	// Open a URL in the default browser application
	OpenURL(url string)

	// Calling Quit on the application will cause the application to exit
	// cleanly, closing all open windows.
	Quit()
}

type fyneApp struct {
}

func (app *fyneApp) NewWindow(title string) Window {
	return GetDriver().CreateWindow(title)
}

func (app *fyneApp) Quit() {
	GetDriver().Quit()
}

func (app *fyneApp) applyTheme(Settings) {
	for _, window := range GetDriver().AllWindows() {
		content := window.Content()

		switch themed := content.(type) {
		case ThemedObject:
			themed.ApplyTheme()
		}

		window.Canvas().Refresh(content)
	}
}

// NewAppWithDriver initialises a new Fyne application using the specified driver
// and returns a handle to that App.
// As this package has no default driver one must be provided.
// Helpers are available - see desktop.NewApp() and test.NewApp().
func NewAppWithDriver(d Driver) App {
	newApp := &fyneApp{}
	setDriver(d)

	listener := make(chan Settings)
	GetSettings().AddChangeListener(listener)
	go func() {
		for {
			settings := <-listener
			newApp.applyTheme(settings)
		}
	}()

	return newApp
}

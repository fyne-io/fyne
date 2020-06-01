// Package app provides app implementations for working with Fyne graphical interfaces.
// The fastest way to get started is to call app.New() which will normally load a new desktop application.
// If the "ci" tag is passed to go (go run -tags ci myapp.go) it will run an in-memory application.
package app // import "fyne.io/fyne/app"

import (
	"os/exec"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/internal"
	helper "fyne.io/fyne/internal/app"
)

// Declare conformity with App interface
var _ fyne.App = (*fyneApp)(nil)

type fyneApp struct {
	driver   fyne.Driver
	icon     fyne.Resource
	uniqueID string

	settings *settings
	prefs    fyne.Preferences
	running  bool
	runMutex sync.Mutex
	exec     func(name string, arg ...string) *exec.Cmd
}

func (app *fyneApp) Icon() fyne.Resource {
	return app.icon
}

func (app *fyneApp) SetIcon(icon fyne.Resource) {
	app.icon = icon
}

func (app *fyneApp) UniqueID() string {
	return app.uniqueID
}

func (app *fyneApp) NewWindow(title string) fyne.Window {
	return app.driver.CreateWindow(title)
}

func (app *fyneApp) Run() {
	app.runMutex.Lock()

	if app.running {
		app.runMutex.Unlock()
		return
	}

	app.running = true
	app.runMutex.Unlock()

	app.driver.Run()
}

func (app *fyneApp) Quit() {
	for _, window := range app.driver.AllWindows() {
		window.Close()
	}

	app.driver.Quit()
	app.settings.stopWatching()
	app.running = false
}

func (app *fyneApp) Driver() fyne.Driver {
	return app.driver
}

// Settings returns the application settings currently configured.
func (app *fyneApp) Settings() fyne.Settings {
	return app.settings
}

func (app *fyneApp) Preferences() fyne.Preferences {
	return app.prefs
}

// New returns a new application instance with the default driver and no unique ID
func New() fyne.App {
	internal.LogHint("Applications should be created with a unique ID using app.NewWithID()")
	return NewWithID("")
}

// NewAppWithDriver initialises a new Fyne application using the specified
// driver and returns a handle to that App. The id should be globally unique to this app
// Built in drivers are provided in the "driver" package.
// Deprecated: Developers should not specify a driver manually but use NewAppWithID()
func NewAppWithDriver(d fyne.Driver, id string) fyne.App {
	return newAppWithDriver(d, id)
}

func newAppWithDriver(d fyne.Driver, id string) fyne.App {
	var prefs fyne.Preferences
	if id == "" {
		prefs = newPreferences()
	} else {
		prefs = loadPreferences(id)
	}
	newApp := &fyneApp{uniqueID: id, driver: d,
		prefs: prefs, exec: exec.Command}
	fyne.SetCurrentApp(newApp)
	newApp.settings = loadSettings()

	listener := make(chan fyne.Settings)
	newApp.Settings().AddChangeListener(listener)
	go func() {
		for {
			set := <-listener
			helper.ApplySettings(set, newApp)
		}
	}()
	if !d.Device().IsMobile() {
		newApp.settings.watchSettings()
	}

	return newApp
}

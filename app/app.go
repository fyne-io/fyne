// Package app provides app implementations for working with Fyne graphical interfaces.
// The fastest way to get started is to call app.New() which will normally load a new desktop application.
// If the "ci" tag is passed to go (go run -tags ci myapp.go) it will run an in-memory application.
package app // import "fyne.io/fyne/app"

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

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
	storage  *store
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
	if app.uniqueID != "" {
		return app.uniqueID
	}

	fyne.LogError("Preferences API requires a unique ID, use app.NewWithID()", nil)
	app.uniqueID = fmt.Sprintf("missing-id-%d", time.Now().Unix()) // This is a fake unique - it just has to not be reused...
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

func (app *fyneApp) Storage() fyne.Storage {
	return app.storage
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
//
// Deprecated: Developers should not specify a driver manually but use NewAppWithID()
func NewAppWithDriver(d fyne.Driver, id string) fyne.App {
	return newAppWithDriver(d, id)
}

func newAppWithDriver(d fyne.Driver, id string) fyne.App {
	newApp := &fyneApp{uniqueID: id, driver: d, exec: exec.Command}
	fyne.SetCurrentApp(newApp)

	newApp.prefs = newPreferences(newApp)
	if pref, ok := newApp.prefs.(interface{ load(string) }); ok && id != "" {
		pref.load(id)
	}
	newApp.settings = loadSettings()
	newApp.storage = &store{a: newApp}

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

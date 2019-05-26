// Package app provides app implementations for working with Fyne graphical interfaces.
// The fastest way to get started is to call app.New() which will normally load a new desktop application.
// If the "ci" tag is passed to go (go run -tags ci myapp.go) it will run an in-memory application.
package app // import "fyne.io/fyne/app"

import (
	"os"
	"sync"

	"fyne.io/fyne"
	helper "fyne.io/fyne/internal/app"
	"fyne.io/fyne/theme"
)

type fyneApp struct {
	driver fyne.Driver
	icon   fyne.Resource

	settings fyne.Settings
	running  bool
	runMutex sync.Mutex
}

func (app *fyneApp) Icon() fyne.Resource {
	return app.icon
}

func (app *fyneApp) SetIcon(icon fyne.Resource) {
	app.icon = icon
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
	app.running = false
}

func (app *fyneApp) Driver() fyne.Driver {
	return app.driver
}

// Settings returns the application settings currently configured.
func (app *fyneApp) Settings() fyne.Settings {
	return app.settings
}

func (app *fyneApp) setupTheme() {
	env := os.Getenv("FYNE_THEME")
	if env == "light" {
		app.Settings().SetTheme(theme.LightTheme())
	} else {
		app.Settings().SetTheme(theme.DarkTheme())
	}
}

// NewAppWithDriver initialises a new Fyne application using the specified
// driver and returns a handle to that App.
// Built in drivers are provided in the "driver" package.
func NewAppWithDriver(d fyne.Driver) fyne.App {
	newApp := &fyneApp{driver: d, icon: theme.FyneLogo(), settings: loadSettings()}
	fyne.SetCurrentApp(newApp)
	newApp.setupTheme()

	listener := make(chan fyne.Settings)
	newApp.Settings().AddChangeListener(listener)
	go func() {
		for {
			<-listener
			helper.ApplyTheme(newApp)
		}
	}()

	return newApp
}

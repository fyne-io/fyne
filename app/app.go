// Package app provides app implementations for working with Fyne graphical interfaces.
// The fastest way to get started is to call app.New() which will normally load a new desktop application.
// If the "ci" tag is passed to go (go run -tags ci myapp.go) it will run an in-memory application.
package app

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/theme"
)

type fyneApp struct {
	driver fyne.Driver
	icon   fyne.Resource
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
	app.driver.Run()
}

func (app *fyneApp) Quit() {
	app.driver.Quit()
}

func (app *fyneApp) Driver() fyne.Driver {
	return app.driver
}

func (app *fyneApp) applyTheme(fyne.Settings) {
	for _, window := range app.driver.AllWindows() {
		content := window.Content()

		switch themed := content.(type) {
		case fyne.ThemedObject:
			themed.ApplyTheme()
			window.Canvas().Refresh(content)
		}
	}
}

// NewAppWithDriver initialises a new Fyne application using the specified driver
// and returns a handle to that App.
// As this package has no default driver one must be provided.
// Helpers are available - see desktop.NewApp() and test.NewApp().
func NewAppWithDriver(d fyne.Driver) fyne.App {
	newApp := &fyneApp{driver: d, icon: theme.FyneLogo()}
	fyne.SetCurrentApp(newApp)

	listener := make(chan fyne.Settings)
	fyne.GlobalSettings().AddChangeListener(listener)
	go func() {
		for {
			settings := <-listener
			newApp.applyTheme(settings)
		}
	}()

	return newApp
}

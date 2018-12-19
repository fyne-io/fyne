// Package app provides app implementations for working with Fyne graphical interfaces.
// The fastest way to get started is to call app.New() which will normally load a new desktop application.
// If the "ci" tag is passed to go (go run -tags ci myapp.go) it will run an in-memory application.
package app

import (
	"os"

	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/theme"
	"github.com/fyne-io/fyne/widget"
)

func init() {
	env := os.Getenv("FYNE_THEME")
	if env == "light" {
		fyne.GlobalSettings().SetTheme(theme.LightTheme())
	} else {
		fyne.GlobalSettings().SetTheme(theme.DarkTheme())
	}
}

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

func (app *fyneApp) applyThemeTo(content fyne.CanvasObject, canvas fyne.Canvas) {
	if themed, ok := content.(fyne.ThemedObject); ok {
		themed.ApplyTheme()
		canvas.Refresh(content)
	}
	if wid, ok := content.(fyne.Widget); ok {
		widget.Renderer(wid).ApplyTheme()
		canvas.Refresh(content)

		for _, o := range widget.Renderer(wid).Objects() {
			app.applyThemeTo(o, canvas)
		}
	}
	if c, ok := content.(*fyne.Container); ok {
		for _, o := range c.Objects {
			app.applyThemeTo(o, canvas)
		}
	}
}

func (app *fyneApp) applyTheme(fyne.Settings) {
	for _, window := range app.driver.AllWindows() {
		app.applyThemeTo(window.Content(), window.Canvas())
	}
}

// NewAppWithDriver initialises a new Fyne application using the specified
// driver and returns a handle to that App.
// Built in drivers are provided in the "driver" package.
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

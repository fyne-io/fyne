// Package app provides app implementations for working with Fyne graphical interfaces.
// The fastest way to get started is to call app.New() which will normally load a new desktop application.
// If the "ci" tag is passed to go (go run -tags ci myapp.go) it will run an in-memory application.
package app // import "fyne.io/fyne/app"

import (
	"os"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
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

func (app *fyneApp) applyThemeTo(content fyne.CanvasObject, canvas fyne.Canvas) {
	if themed, ok := content.(fyne.Themeable); ok {
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
		canvas.Refresh(c)
	}
}

func (app *fyneApp) applyTheme() {
	for _, window := range app.driver.AllWindows() {
		if window.Padded() {
			window.Content().Move(fyne.NewPos(theme.Padding(), theme.Padding()))
		}
		app.applyThemeTo(window.Content(), window.Canvas())
	}
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
			newApp.applyTheme()
		}
	}()

	return newApp
}

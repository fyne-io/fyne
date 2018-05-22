// Package fyneapp provides a full implementation of the Fyne APIs
package fyneapp

import "github.com/fyne-io/fyne/app"
import "github.com/fyne-io/fyne/ui"

import "github.com/fyne-io/fyne-app/efl"

type eflApp struct {
	driver ui.Driver
}

func (app *eflApp) NewWindow(title string) ui.Window {
	return app.driver.CreateWindow(title)
}

func (app *eflApp) Quit() {
	efl.Quit()
}

func (app *eflApp) applyTheme(app.Settings) {
	for _, window := range app.driver.AllWindows() {
		window.Canvas().Refresh(window.Canvas().Content())
	}
}

// NewApp initialises a new Fyne application returning a handle to that App
func NewApp() app.App {
	newApp := &eflApp{
		driver: efl.NewEFLDriver(),
	}

	listener := make(chan app.Settings)
	app.GetSettings().AddChangeListener(listener)
	go func() {
		for {
			settings := <-listener
			newApp.applyTheme(settings)
		}
	}()

	return newApp
}

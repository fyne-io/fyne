// +build !ci

// Package fyneapp provides a full implementation of the Fyne APIs
package desktop

import "github.com/fyne-io/fyne/api/app"
import "github.com/fyne-io/fyne/api/ui"

type fyneApp struct {
}

func (app *fyneApp) NewWindow(title string) ui.Window {
	return ui.GetDriver().CreateWindow(title)
}

func (app *fyneApp) Quit() {
	ui.GetDriver().Quit()
}

func (app *fyneApp) applyTheme(app.Settings) {
	for _, window := range ui.GetDriver().AllWindows() {
		window.Canvas().Refresh(window.Canvas().Content())
	}
}

// NewApp initialises a new Fyne application returning a handle to that App
func NewApp() app.App {
	newApp := &fyneApp{}

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

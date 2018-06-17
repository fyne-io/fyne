// +build !ci

// Package desktop provides a full implementation of the Fyne APIs
package desktop

import "github.com/fyne-io/fyne/api"
import "github.com/fyne-io/fyne/api/ui"

type fyneApp struct {
}

func (app *fyneApp) NewWindow(title string) ui.Window {
	return ui.GetDriver().CreateWindow(title)
}

func (app *fyneApp) Quit() {
	ui.GetDriver().Quit()
}

func (app *fyneApp) applyTheme(fyne.Settings) {
	for _, window := range ui.GetDriver().AllWindows() {
		window.Canvas().Refresh(window.Canvas().Content())
	}
}

// NewApp initialises a new Fyne application returning a handle to that App
func NewApp() fyne.App {
	newApp := &fyneApp{}

	listener := make(chan fyne.Settings)
	fyne.GetSettings().AddChangeListener(listener)
	go func() {
		for {
			settings := <-listener
			newApp.applyTheme(settings)
		}
	}()

	return newApp
}

// +build !ci

// Package desktop provides a full implementation of the Fyne APIs
package desktop

import "github.com/fyne-io/fyne"

type fyneApp struct {
}

func (app *fyneApp) NewWindow(title string) fyne.Window {
	return fyne.GetDriver().CreateWindow(title)
}

func (app *fyneApp) Quit() {
	fyne.GetDriver().Quit()
}

func (app *fyneApp) applyTheme(fyne.Settings) {
	for _, window := range fyne.GetDriver().AllWindows() {
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

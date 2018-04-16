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

// NewApp initialises a new Fyne application returning a handle to that App
func NewApp() app.App {
	app := &eflApp{
		driver: efl.NewEFLDriver(),
	}

	return app
}

func (app *eflApp) Quit() {
	app.driver.Quit()
}

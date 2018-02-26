package fyneapp

import "github.com/fyne-io/fyne/app"
import "github.com/fyne-io/fyne/ui"

import "github.com/fyne-io/fyne-app/driver"

type eflApp struct {
	driver *driver.EFLDriver
}

func (app *eflApp) NewWindow(title string) ui.Window {
	return app.driver.CreateWindow(title)
}

func NewApp() app.App {
	app := &eflApp{
		driver: driver.NewEFLDriver(),
	}

	return app
}

func (app *eflApp) Quit() {
	app.driver.Quit()
}

package efl

import "github.com/fyne-io/fyne/app"
import "github.com/fyne-io/fyne/ui"

import "github.com/fyne-io/fyne-efl/driver"

type eflApp struct {
	driver driver.Driver
}

func (app *eflApp) NewWindow(title string) ui.Window {
        return app.driver.CreateWindow(title)
}

func NewEFLApp() app.App {
	app := &eflApp{
		driver: new(driver.EFLDriver),
	}

	return app
}

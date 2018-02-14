package main

import "github.com/fyne-io/fyne/app"
import "github.com/fyne-io/fyne/ui"

import "github.com/fyne-io/fyne-efl/driver"

func main() {
        app.SetUIDriver(new(driver.EFLDriver))

	w := app.NewWindow("Hello")
	w.Canvas().AddObject(ui.NewText("Hello Fyne!"))

	app.Run()
}

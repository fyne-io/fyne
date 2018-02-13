package main

import "github.com/fyne-io/fyne/app"
import "github.com/fyne-io/fyne/ui"

func main() {
	w := app.NewWindow("Hello")
	w.Canvas().AddObject(ui.NewText("Hello Fyne!"))

	app.Run()
}

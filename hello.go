package main

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne-efl"

func main() {
	app := efl.NewEFLApp()

	w := app.NewWindow("Hello")
	w.Canvas().SetContent(ui.NewText("Hello Fyne!"))

	w.Show()
}

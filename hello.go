package main

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne-app"

func main() {
	app := fyneapp.NewApp()

	w := app.NewWindow("Hello")
	w.Canvas().SetContent(ui.NewText("Hello Fyne!"))

	w.Show()
}

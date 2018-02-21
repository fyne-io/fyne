package main

import "fmt"

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/widget"
import "github.com/fyne-io/fyne-app"

func main() {
	app := fyneapp.NewApp()

	w := app.NewWindow("Hello")
	quit := widget.NewButton("Quit", func() {
		app.Quit()
	})
	w.Canvas().SetContent(ui.NewContainer(
		[]ui.CanvasObject{
			ui.NewText("Hello Fyne!"),
			button,
		},
		layout.NewGridLayout(1)))

	w.Show()
}

package main

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/event"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/widget"
import "github.com/fyne-io/fyne-app"

func main() {
	app := fyneapp.NewApp()

	w := app.NewWindow("Hello")
	container := ui.NewContainer(
		ui.NewText("Hello Fyne!"),
		widget.NewButton("Quit", func(e *event.MouseEvent) {
			app.Quit()
		}))
	container.Layout = layout.NewGridLayout(1)

	w.Canvas().SetContent(container)
	w.Show()
}

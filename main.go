// Package main provides various examples of Fyne API capabilities
package main

import "github.com/fyne-io/examples/examples"

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/event"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/widget"
import "github.com/fyne-io/fyne-app"

func main() {
	app := fyneapp.NewApp()

	w := app.NewWindow("Examples")
	container := ui.NewContainer(
		widget.NewLabel("Fyne Examples!"),

		widget.NewButton("Calculator", func(e *event.MouseEvent) {
			examples.Calculator(app)
		}),
		widget.NewButton("Clock", func(e *event.MouseEvent) {
			examples.Clock(app)
		}),
		widget.NewButton("Canvas", func(e *event.MouseEvent) {
			examples.Canvas(app)
		}),

		&canvas.RectangleObject{},
		widget.NewButton("Quit", func(e *event.MouseEvent) {
			app.Quit()
		}))
	container.Layout = layout.NewGridLayout(1)

	w.Canvas().SetContent(container)
	w.Show()
}

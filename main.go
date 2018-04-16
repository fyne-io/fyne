// Package main provides various examples of Fyne API capabilities
package main

import "flag"

import "github.com/fyne-io/examples/examples"

import "github.com/fyne-io/fyne/app"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/event"
import "github.com/fyne-io/fyne/ui/widget"
import "github.com/fyne-io/fyne-app"

func calcApp(app app.App) {
	examples.Calculator(app)
}

func canvasApp(app app.App) {
	examples.Canvas(app)
}

func clockApp(app app.App) {
	examples.Clock(app)
}

func welcome(app app.App) {
	w := app.NewWindow("Examples")
	w.Canvas().SetContent(widget.NewList(
		widget.NewLabel("Fyne Examples!"),

		widget.NewButton("Calculator", func(e *event.MouseEvent) {
			calcApp(app)
		}),
		widget.NewButton("Clock", func(e *event.MouseEvent) {
			clockApp(app)
		}),
		widget.NewButton("Canvas", func(e *event.MouseEvent) {
			canvasApp(app)
		}),

		&canvas.Rectangle{},
		widget.NewButton("Quit", func(e *event.MouseEvent) {
			app.Quit()
		})))
	w.Show()
}

func main() {
	app := fyneapp.NewApp()

	var ex string
	flag.StringVar(&ex, "example", "", "Launch an app directly")
	flag.Parse()

	switch ex {
	case "calculator":
		calcApp(app)
	case "canvas":
		canvasApp(app)
	case "clock":
		clockApp(app)
	default:
		welcome(app)
	}
}

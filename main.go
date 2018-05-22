// Package main provides various examples of Fyne API capabilities
package main

import "flag"

import "github.com/fyne-io/examples/examples"

import "github.com/fyne-io/fyne/app"
import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/widget"
import "github.com/fyne-io/fyne-app"

func blogApp(app app.App) {
	examples.Blog(app)
}

func calcApp(app app.App) {
	examples.Calculator(app)
}

func canvasApp(app app.App) {
	examples.Canvas(app)
}

func clockApp(app app.App) {
	examples.Clock(app)
}

func fractalApp(app app.App) {
	examples.Fractal(app)
}

func appButton(app app.App, label string, onClick func(app.App)) *widget.Button {
	return widget.NewButton(label, func() {
		onClick(app)
	})
}

func welcome(myApp app.App) {
	w := myApp.NewWindow("Examples")
	w.Canvas().SetContent(widget.NewList(
		widget.NewLabel("Fyne Examples!"),

		appButton(myApp, "Blog", blogApp),
		appButton(myApp, "Calculator", calcApp),
		appButton(myApp, "Clock", clockApp),
		appButton(myApp, "Fractal", fractalApp),
		appButton(myApp, "Canvas", canvasApp),

		&canvas.Rectangle{},
		ui.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewButton("Dark", func() {
				app.GetSettings().SetTheme("dark")
			}),
			widget.NewButton("Light", func() {
				app.GetSettings().SetTheme("light")
			}),
		),
		widget.NewButton("Quit", func() {
			myApp.Quit()
		})))
	w.Show()
}

func main() {
	app := fyneapp.NewApp()

	var ex string
	flag.StringVar(&ex, "example", "", "Launch an app directly (blog,calculator,canvas,clock)")
	flag.Parse()

	switch ex {
	case "blog":
		blogApp(app)
	case "calculator":
		calcApp(app)
	case "canvas":
		canvasApp(app)
	case "clock":
		clockApp(app)
	case "fractal":
		fractalApp(app)
	default:
		welcome(app)
	}
}

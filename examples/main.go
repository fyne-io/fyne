// Package main provides various examples of Fyne API capabilities
package main

import "flag"
import "fmt"

import "github.com/fyne-io/fyne/examples/apps"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/layout"
import "github.com/fyne-io/fyne/theme"
import "github.com/fyne-io/fyne/widget"

func blogApp(app fyne.App) {
	apps.Blog(app)
}

func calcApp(app fyne.App) {
	apps.Calculator(app)
}

func canvasApp(app fyne.App) {
	apps.Canvas(app)
}

func clockApp(app fyne.App) {
	apps.Clock(app)
}

func fractalApp(app fyne.App) {
	apps.Fractal(app)
}

func layoutApp(app fyne.App) {
	apps.Layout(app)
}

func lifeApp(app fyne.App) {
	apps.Life(app)
}

func appButton(app fyne.App, label string, onClick func(fyne.App)) *widget.Button {
	return widget.NewButton(label, func() {
		onClick(app)
	})
}

func welcome(app fyne.App) {
	w := app.NewWindow("Examples")
	w.SetContent(widget.NewList(
		widget.NewLabel("Fyne Examples!"),

		appButton(app, "Blog", blogApp),
		appButton(app, "Calculator", calcApp),
		appButton(app, "Clock", clockApp),
		appButton(app, "Fractal", fractalApp),
		appButton(app, "Canvas", canvasApp),
		appButton(app, "Layout", layoutApp),
		appButton(app, "Life", lifeApp),

		fyne.NewSpacer(),
		widget.NewEntry(),
		widget.NewCheck("Check", func(on bool) { fmt.Println("checked", on) }),
		fyne.NewSpacer(),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewButton("Dark", func() {
				fyne.GetSettings().SetTheme("dark")
			}),
			widget.NewButton("Light", func() {
				fyne.GetSettings().SetTheme("light")
			}),
		),
		widget.NewButtonWithIcon("Quit", theme.CancelIcon(), func() {
			app.Quit()
		})))
	w.Show()
}

func main() {
	app := apps.NewApp()

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
	case "layout":
		layoutApp(app)
	case "life":
		lifeApp(app)
	default:
		welcome(app)
	}
}

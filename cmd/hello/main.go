// Package main loads a very basic Hello World graphical application
package main

import (
	"net/url"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	a := app.New()
	tourURL, _ := url.Parse("https://tour.fyne.io")

	w := a.NewWindow("Hello")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewHyperlink("tour.fyne.io", tourURL),
	))

	w.ShowAndRun()
}

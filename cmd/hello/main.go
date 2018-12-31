// +build !ci

// Package main loads a very basic Hello World graphical application
package main

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/app"
	"github.com/fyne-io/fyne/widget"
)

func main() {
	app := app.New()

	w := app.NewWindow("Hello")
	w.SetContent(&widget.Box{Horizontal: false, Children: []fyne.CanvasObject{
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	}})

	w.ShowAndRun()
}

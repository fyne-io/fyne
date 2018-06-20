// +build !ci

// Package main loads a very basic Hello World graphical application
package main

import "github.com/fyne-io/fyne/widget"
import "github.com/fyne-io/fyne/desktop"

func main() {
	app := desktop.NewApp()

	w := app.NewWindow("Hello")
	w.SetContent(widget.NewList(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))

	w.Show()
}

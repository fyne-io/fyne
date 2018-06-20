// +build !ci

// Package main loads a very basic Hello World graphical application
package main

import "github.com/fyne-io/fyne/widget"
import "github.com/fyne-io/fyne/desktop"

func main() {
	app := desktop.NewApp()

	w := app.NewWindow("Hello")
	quit := widget.NewButton("Quit", func() {
		app.Quit()
	})
	w.Canvas().SetContent(widget.NewList(
		widget.NewLabel("Hello Fyne!"),
		quit))

	w.Show()
}

// +build !ci

// Package main loads a very basic Hello World graphical application
package main

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/desktop"
import "github.com/fyne-io/fyne/widget"

func main() {
	app := desktop.NewApp()

	w := app.NewWindow("Hello")
	w.SetContent(&widget.Box{Horizontal: false, Children: []fyne.CanvasObject{
		&widget.Label{Text: "Hello Fyne!"},
		&widget.Button{Text: "Quit", OnTapped: func() {
			app.Quit()
		}},
	}})

	w.ShowAndRun()
}

// +build !ci

// Package main loads a very basic Hello World graphical application
package main

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/desktop"
import W "github.com/fyne-io/fyne/widget"

func main() {
	app := desktop.NewApp()

	w := app.NewWindow("Hello")
	w.SetContent(&W.List{Children: []fyne.CanvasObject{
		&W.Label{Text: "Hello Fyne!"},
		&W.Button{Text: "Quit", OnTapped: func() {
			app.Quit()
		}},
	}})

	w.Show()
}

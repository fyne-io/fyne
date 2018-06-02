// Package fyne is the containing package for all aspects of the Fyne UI toolkit
//
// A simple application may look like this:
//
//   package main
//
//   import "github.com/fyne-io/fyne/api/ui/widget"
//   import "github.com/fyne-io/fyne/desktop"
//
//   func main() {
//   	app := desktop.NewApp()
//
//   	w := app.NewWindow("Hello")
//   	quit := widget.NewButton("Quit", func() {
//   		app.Quit()
//   	})
//   	w.Canvas().SetContent(widget.NewList(
//   		widget.NewLabel("Hello Fyne!"),
//   		quit))
//
//   	w.Show()
//   }
package fyne

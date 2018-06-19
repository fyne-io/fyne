// Package fyne describes the objects and components available to any Fyne app.
// These can all be created, manipulated and tested without renering (for speed).
// Your main package should import a driver which will actually render your UI.
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

// Package fyne describes the objects and components available to any Fyne app.
// These can all be created, manipulated and tested without rendering (for speed).
// Your main package should use the app package to create an application with
// a default driver that will render your UI.
//
// A simple application may look like this:
//
//   package main
//
//   import "fyne.io/fyne/app"
//   import "fyne.io/fyne/widget"
//
//   func main() {
//   	a := app.New()
//
//   	w := a.NewWindow("Hello")
//   	w.SetContent(widget.NewVBox(
//   		widget.NewLabel("Hello Fyne!"),
//   		widget.NewButton("Quit", func() {
//   			a.Quit()
//   		})))
//
//   	w.ShowAndRun()
//   }
package fyne // import "fyne.io/fyne"

// Package main loads a very basic Hello World graphical application.
package main

import (
	//	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	//	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello")

	notab := widget.NewLabel("12345678901234567890")
	notab.TextStyle.Monospace = true
	withtab1 := widget.NewLabel("\ttab")
	withtab1.TextStyle.Monospace = true
	withtab2 := widget.NewLabel("\t\ttab")
	withtab2.TextStyle.Monospace = true
	withtab3 := widget.NewLabel("123\t\ttab")
	withtab3.TextStyle.Monospace = true
	withtab4 := widget.NewLabel("123\ttab")
	withtab4.TextStyle.Monospace = true
	withtab5 := widget.NewEntry()
	withtab5.Text = "a\n\tb\nc"
	withtab5.TextStyle.Monospace = true
	w.SetContent(container.NewVBox(
		notab,
		withtab1,
		withtab2,
		withtab3,
		withtab4,
		withtab5,
	))

	// e := widget.NewEntry()
	// e.SetText("1234a\n\tb\n1234\tc")
	// e.TextStyle.Monospace = true

	// e.CursorRow = 1
	// e.KeyDown(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})
	// e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	// e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	// e.KeyUp(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})

	// w.SetContent(e)
	// w.Resize(fyne.NewSize(86, 86))
	// w.Canvas().Focus(e)

	// grid := widget.NewTextGridFromString("A\n\tb")
	// grid.ShowWhitespace = true
	// grid.Resize(fyne.NewSize(56, 42)) // causes refresh
	// w.SetContent(grid)

	w.ShowAndRun()
}

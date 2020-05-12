// Package main loads a very basic Hello World graphical application
package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/cmd/fyne_demo/data"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func main() {
	a := app.New()

	hello := widget.NewLabel("Hello Fyne!")
	logo := canvas.NewImageFromResource(data.FyneScene)
	quit := widget.NewButton("Quit", func() {
		a.Quit()
	})

	logo.FillMode = canvas.ImageFillContain
	logo.ScaleMode = canvas.ImageScaleSmooth
	logo.SetMinSize(fyne.NewSize(100, 100))

	w := a.NewWindow("Hello")

	w.SetContent(fyne.NewContainerWithLayout(layout.NewBorderLayout(hello, quit, nil, nil),
		hello,
		logo,
		quit,
	))

	w.ShowAndRun()
}

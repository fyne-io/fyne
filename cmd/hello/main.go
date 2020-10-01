// Package main loads a very basic Hello World graphical application.
package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello")

	icon := widget.NewFileIcon(nil)
	button := widget.NewButton("Select file", func() {
		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil || file == nil {
				return
			}

			icon.SetURI(file.URI())
		}, w)
	})

	content := fyne.NewContainerWithLayout(layout.NewBorderLayout(button, nil, nil, nil), button, icon)
	w.SetContent(content)
	w.ShowAndRun()
}

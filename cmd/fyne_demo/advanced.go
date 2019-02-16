package main

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func scaleString(c fyne.Canvas) string {
	return fmt.Sprintf("%0.2f", c.Scale())
}

// Advanced loads a window that shows details and settings that are a bit
// more detailed than normally needed.
func Advanced(app fyne.App) {
	win := app.NewWindow("Advanced")

	scale := widget.NewLabel("")

	screen := widget.NewGroup("Screen", widget.NewForm(
		&widget.FormItem{Text: "Scale", Widget: scale},
	))

	win.SetContent(widget.NewVBox(screen,
		widget.NewButton("Custom Theme", func() {
			app.Settings().SetTheme(newCustomTheme())
		}),
		widget.NewButton("Fullscreen", func() {
			win.SetFullScreen(!win.FullScreen())
		}),
	))
	win.Show()
	scale.SetText(scaleString(win.Canvas()))
}

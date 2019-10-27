package main

import (
	"fmt"
	"fyne.io/fyne/canvas"
	_ "image/png"

	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"fyne.io/fyne"
)

// View is a window with reactive elements
type View struct {
	app  fyne.App
	data *dataModel
	w    fyne.Window
}

func newView(app fyne.App, data *dataModel, logo *canvas.Image) *View {
	data.NumWindows++
	title := fmt.Sprintf("Reactive Data View #%d", data.NumWindows)

	v := &View{
		app:  app,
		data: data,
		w:    app.NewWindow(title),
	}

	v.w.SetContent(widget.NewVBox(
		widget.NewLabelWithStyle("Reactive Data", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
		layout.NewSpacer(),

	))

	v.w.Show()
	return v
}

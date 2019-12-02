package main

import (
	"fmt"
	_ "image/png"

	"fyne.io/fyne/canvas"
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
	title := fmt.Sprintf("View #%d", data.NumWindows)

	v := &View{
		app:  app,
		data: data,
		w:    app.NewWindow(title),
	}

	myWindowID := data.NumWindows
	v.w.SetContent(widget.NewVBox(
		widget.NewLabelWithStyle("Reactive Data", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
		layout.NewSpacer(),
		widget.NewForm(
			// A label bound to a Clock dataItem (a Clock is a string that auto-mutates every second)
			widget.NewFormItem("Time", widget.NewLabelWithData(data.Clock)),
			// An entry widget bound to a String dataItem
			widget.NewFormItem("Name", widget.NewEntryWithData(data.Name)),
			// A checkbox bound to a Bool dataItem
			widget.NewFormItem("", widget.NewCheckWithData(data.IsAvailable,
				"Avail",
				func(checked bool) {
					println("clicked on the avail button", checked)
				})),

			// A radioButton bound to an Int dataItem
			widget.NewFormItem("", widget.NewRadioWithData(data.Size,
				[]string{"Small", "Medium", "Large"},
				func(value string) {
					println("Radio button changed to", value)
				})),
			// A slider widget bound to a Float dataItem
			widget.NewFormItem("Delivery", widget.NewSliderWithData(data.DeliveryTime, 0.0, 100.0)),
			// A checkbox widget bound to a String dataItem  (ie - sets it true/false)
			widget.NewFormItem("", widget.NewCheckWithData(data.OnSale,
				"On Sale",
				func(checked bool) {
					println("clicked on the on sale button", checked)
				})),
			// A label widget bound to the same String dataItem as above (ie - true/false)
			widget.NewFormItem("Is On Sale ?", widget.NewLabelWithData(data.OnSale)),
			// Need a quit handler here to correctly deregister all the handlers
			// TODO - correctly clean up, else resource leaks
			widget.NewFormItem("Quit Cleanly", widget.NewButton("Quit", func() {
				data.Clock.DeleteListener(myWindowID)
				v.w.Close()
			})),
		),
	))

	v.w.Show()
	return v
}

package main

// This whole section will be removed - its a test harness to try out databinding ideas
// Will move this to a separate repo shortly

import (
	"fmt"
	_ "image/png"
	"net/url"

	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"fyne.io/fyne"
)

// View is a window with reactive elements
type View struct {
	app  fyne.App
	data *DataModel
	w    fyne.Window
}

func newView(app fyne.App, data *DataModel, cache *ImageCache) *View {
	data.NumWindows++
	title := fmt.Sprintf("View #%d", data.NumWindows)

	v := &View{
		app:  app,
		data: data,
		w:    app.NewWindow(title),
	}

	myWindowID := data.NumWindows

	myURL, err := url.Parse(data.URL.String())
	if err != nil {
		println("URL parse error:", err.Error())
	}

	v.w.SetContent(widget.NewVBox(
		widget.NewLabelWithStyle("Reactive Data", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		NewImageWidget(cache).Bind(data.Image),
		layout.NewSpacer(),
		widget.NewForm(
			// A label bound to a Clock dataItem (a Clock is a string that auto-mutates every second)
			widget.NewFormItem("Time", widget.NewLabel("").Bind(data.Clock)),
			// An entry widget bound to a String dataItem
			widget.NewFormItem("Name", widget.NewEntry().Bind(data.Name)),
			// A checkbox bound to a Bool dataItem
			widget.NewFormItem("", widget.NewCheck(
				"Avail",
				func(checked bool) {
					println("clicked on the avail button", checked)
					if checked {
						data.Image.Set(FyneAvatarAvail, 0)
					}
				}).Bind(data.IsAvailable)),

			// A radioButton bound to an Int dataItem
			widget.NewFormItem("", widget.NewRadio(
				[]string{"Small", "Medium", "Large"},
				func(value string) {
					println("Radio button changed to", value)
					data.URL.Set("http://myurl.com/"+value, 0)
					switch value {
					case "Small":
						data.Image.Set(FyneAvatarSm, 0)
					case "Medium":
						data.Image.Set(FyneAvatarMd, 0)
					case "Large":
						data.Image.Set(FyneAvatarLg, 0)
					default:
						data.Image.Set(FyneAvatar, 0)
					}
				}).Bind(data.Size)),
			// A slider widget bound to a Float dataItem
			widget.NewFormItem("Delivery", widget.NewSlider(0.0, 100.0).Bind(data.DeliveryTime)),
			// A checkbox widget bound to a String dataItem  (ie - sets it true/false)
			widget.NewFormItem("", widget.NewCheck(
				"On Sale",
				func(checked bool) {
					println("clicked on the on sale button", checked)
					if checked {
						data.Image.Set(FyneAvatarOnSale, 0)
					}
				}).Bind(data.OnSale)),
			// A label widget bound to the same String dataItem as above (ie - true/false)
			widget.NewFormItem("Is On Sale ?", widget.NewLabel("").Bind(data.OnSale)),
			widget.NewFormItem("MyLink", widget.NewHyperlink("MyLink", myURL).Bind(data.URL)),
			widget.NewFormItem("Link Button", widget.NewButton("MyLink", func() {
				println("Button pressed -> goto", data.URL.String())
			}).Bind(data.URL)),
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

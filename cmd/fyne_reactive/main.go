package main

// This whole section will be removed - its a test harness to try out databinding ideas
// Will move this to a separate repo shortly

//go:generate fyne bundle reactive-data.png > image.go

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func welcomeScreen(a fyne.App, data *DataModel, diag *canvas.Image) fyne.CanvasObject {
	cache := newImageCache()
	return widget.NewVBox(
		widget.NewLabelWithStyle("Fyne Reactive Data Update Demo", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewHBox(
			layout.NewSpacer(),
			diag,
			layout.NewSpacer(),
		),
		widget.NewHBox(
			layout.NewSpacer(),
			NewImageWidget(cache).Bind(data.Image),
			layout.NewSpacer(),
		),
		widget.NewLabel(`This demo has a single instance of a dataapi model as described above.

Each window subscribes to the dataapi model, and is automatically updated
when the model changes.

Changes to the dataapi in the view are committed to the central dataapi model,
which in turn triggers a repaint on the other subscribed views.`),
		layout.NewSpacer(),
		widget.NewGroup("Data Model Internals",
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				widget.NewLabel("Clock").Bind(data.Clock),
				widget.NewLabel("Name (String)").Bind(data.Name),
				widget.NewLabel("Avail (Bool)").Bind(data.IsAvailable),
				widget.NewLabel("OnSale (String)").Bind(data.OnSale),
				widget.NewLabel("Size (Int)").Bind(data.Size),
				widget.NewLabel("Deliv Time (Float)").Bind(data.DeliveryTime),
				widget.NewLabel("URL").Bind(data.URL),
				widget.NewLabel("Image").Bind(data.Image),
			),
		),
		widget.NewGroup("Launch New Viewer",
			fyne.NewContainerWithLayout(layout.NewGridLayout(5),
				layout.NewSpacer(),
				layout.NewSpacer(),
				widget.NewButton("+ Viewer Window", func() {
					newView(a, data, cache)
				}),
				widget.NewButton("Fx1", func() {
					newFx(a, data)
				}),
				layout.NewSpacer(),
				widget.NewButton("Close", func() { a.Quit() }),
			),
		),
	)
}

func main() {
	//logo := canvas.NewImageFromResource(theme.FyneLogo())
	//logo.SetMinSize(fyne.NewSize(320, 128))
	diag := canvas.NewImageFromResource(resourceReactiveDataPng)
	diag.SetMinSize(fyne.NewSize(400, 350))

	data := NewDataModel()
	a := app.NewWithID("io.fyne.reactive")
	w := a.NewWindow("Fyne Reactive Demo")
	w.SetMaster()

	w.SetContent(welcomeScreen(a, data, diag))

	w.ShowAndRun()
}

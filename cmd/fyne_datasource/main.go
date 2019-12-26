package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func welcomeScreen(data *DataModel) fyne.CanvasObject {
	eCountry := widget.NewEntry().Bind(data.SelectedCountry)
	eState := widget.NewEntry().Bind(data.SelectedState)

	setCountry := func(country string) {
		data.StatesList.SetFromStringSlice(data.CountryAndState.GetStates(country))
		data.SelectedState.Set("", 0)
	}

	eCountry.OnChanged = setCountry

	return fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		widget.NewVBox(
			widget.NewLabelWithStyle("Fyne Reactive Data Select Country Demo", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
			widget.NewForm(
				// A label bound to a Clock dataItem (a Clock is a string that auto-mutates every second)
				widget.NewFormItem("Time", widget.NewLabel("").Bind(data.Clock)),
				widget.NewFormItem("Country", widget.NewSelect(nil, setCountry).
					Source(data.CountriesList).
					Bind(data.SelectedCountry)),
				widget.NewFormItem("State", widget.NewSelect(nil, nil).
					Source(data.StatesList).
					Bind(data.SelectedState)),
				widget.NewFormItem("Selected Country", widget.NewLabel("").Bind(data.SelectedCountry)),
				widget.NewFormItem("Selected State", widget.NewLabel("").Bind(data.SelectedState)),
				widget.NewFormItem("Edit Country", eCountry),
				widget.NewFormItem("Edit State", eState),
			),
		),
	)
}

func main() {
	data := NewDataModel()
	a := app.NewWithID("io.fyne.reactive_select")
	w := a.NewWindow("Fyne Reactive Select Demo")
	w.SetMaster()

	w.SetContent(welcomeScreen(data))

	w.ShowAndRun()
}

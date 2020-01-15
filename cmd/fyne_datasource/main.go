package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func welcomeScreen(data *DataModel) fyne.CanvasObject {
	eCountry := widget.NewEntry()

	setCountry := func(country string) {
		data.StatesList.SetFromStringSlice(data.CountryAndState.GetStates(country))
		data.SelectedState.SetString("")
	}

	eCountry.OnChanged = setCountry

	return fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		widget.NewVBox(
			widget.NewLabelWithStyle(
				"Fyne Reactive Data Select Country Demo",
				fyne.TextAlignCenter,
				fyne.TextStyle{Bold: true},
			),
			widget.NewForm(
				widget.NewFormItem("Time", widget.NewLabel("").
					Bind(data.Clock)),
				widget.NewFormItem("Country", widget.NewSelect(nil, setCountry).
					Source(data.CountriesList).
					Bind(data.SelectedCountry)),
				widget.NewFormItem("State", widget.NewSelect(nil, nil).
					Source(data.StatesList).
					Bind(data.SelectedState)),
				widget.NewFormItem("Selected Country", widget.NewLabel("").
					Bind(data.SelectedCountry)),
				widget.NewFormItem("Selected State", widget.NewLabel("").
					Bind(data.SelectedState)),
				widget.NewFormItem("Edit Country", eCountry.
					Bind(data.SelectedCountry)),
				widget.NewFormItem("Edit State", widget.NewEntry().
					Bind(data.SelectedState)),
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

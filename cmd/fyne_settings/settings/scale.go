package settings

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type scaleItems struct {
	scale   float32
	name    string
	preview *canvas.Text
	button  *widget.Button
}

var scales = []*scaleItems{
	{scale: 0.5, name: "Tiny"},
	{scale: 0.8, name: "Small"},
	{scale: 1, name: "Normal"},
	{scale: 1.3, name: "Large"},
	{scale: 1.8, name: "Huge"}}

func (s *Settings) appliedScale(value float32) {
	for _, scale := range scales {
		scale.preview.TextSize = theme.TextSize() * scale.scale / value
	}
}

func (s *Settings) chooseScale(value float32) {
	s.fyneSettings.Scale = value

	for _, scale := range scales {
		if scale.scale == value {
			scale.button.Importance = widget.HighImportance
		} else {
			scale.button.Importance = widget.MediumImportance
		}

		scale.button.Refresh()
	}
}

func (s *Settings) makeScaleButtons() []fyne.CanvasObject {
	var buttons = make([]fyne.CanvasObject, len(scales))
	for i, scale := range scales {
		value := scale.scale
		button := widget.NewButton(scale.name, func() {
			s.chooseScale(value)
		})
		if s.fyneSettings.Scale == scale.scale {
			button.Importance = widget.HighImportance
		}

		scale.button = button
		buttons[i] = button
	}

	return buttons
}

func (s *Settings) makeScaleGroup(scale float32) *widget.Card {
	scalePreviewBox := container.NewGridWithColumns(5, s.makeScalePreviews(scale)...)
	scaleBox := container.NewGridWithColumns(5, s.makeScaleButtons()...)

	return widget.NewCard("Scale", "", container.NewVBox(scalePreviewBox, scaleBox, newRefreshMonitor(s)))
}

func (s *Settings) makeScalePreviews(value float32) []fyne.CanvasObject {
	var previews = make([]fyne.CanvasObject, len(scales))
	for i, scale := range scales {
		text := canvas.NewText("A", theme.ForegroundColor())
		text.Alignment = fyne.TextAlignCenter
		text.TextSize = theme.TextSize() * scale.scale / value

		scale.preview = text
		previews[i] = text
	}

	return previews
}

func (s *Settings) refreshScalePreviews() {
	for _, scale := range scales {
		scale.preview.Color = theme.ForegroundColor()
	}
}

// refreshMonitor is a simple widget that updates canvas components when the UI is asked to refresh.
// Captures theme and scale changes without the settings monitoring code.
type refreshMonitor struct {
	widget.Label
	settings *Settings
}

func (r *refreshMonitor) Refresh() {
	r.settings.refreshScalePreviews()
	r.Label.Refresh()
}

func newRefreshMonitor(s *Settings) *refreshMonitor {
	r := &refreshMonitor{settings: s}
	r.Hide()
	return r
}

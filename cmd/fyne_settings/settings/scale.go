package settings

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
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

func (s *Settings) chooseScale(value float32) {
	s.fyneSettings.Scale = value

	for _, scale := range scales {
		if scale.scale == value {
			scale.button.Style = widget.PrimaryButton
		} else {
			scale.button.Style = widget.DefaultButton
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
			button.Style = widget.PrimaryButton
		}

		scale.button = button
		buttons[i] = button
	}

	return buttons
}

func (s *Settings) makeScaleGroup(scale float32) *widget.Group {
	scalePreviewBox := fyne.NewContainerWithLayout(layout.NewGridLayout(5), s.makeScalePreviews(scale)...)
	scaleBox := fyne.NewContainerWithLayout(layout.NewGridLayout(5), s.makeScaleButtons()...)

	s.monitorTheme()
	return widget.NewGroup("Scale", scalePreviewBox, scaleBox)
}

func (s *Settings) makeScalePreviews(value float32) []fyne.CanvasObject {
	var previews = make([]fyne.CanvasObject, len(scales))
	for i, scale := range scales {
		text := canvas.NewText("A", theme.TextColor())
		text.Alignment = fyne.TextAlignCenter
		text.TextSize = int(float32(theme.TextSize()) * scale.scale / value)

		scale.preview = text
		previews[i] = text
	}

	return previews
}

func (s *Settings) monitorTheme() {
	watch := make(chan fyne.Settings)
	fyne.CurrentApp().Settings().AddChangeListener(watch)
	go func() {
		for range watch {
			s.refreshScalePreviews()
		}
	}()
}

func (s *Settings) refreshScalePreviews() {
	for _, scale := range scales {
		scale.preview.Color = theme.TextColor()
	}
}

func (s *Settings) setSelectedScale(value float32) {
	for _, scale := range scales {
		scale.preview.TextSize = int(float32(theme.TextSize()) * scale.scale / value)
	}
}

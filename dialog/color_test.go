package dialog

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
)

func TestColorDialog_Theme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(500, 300))

	d := NewColorPicker("Color Picker", "Pick a Color", nil, w)
	d.Show()

	test.AssertImageMatches(t, "color/light.png", w.Canvas().Capture())

	test.ApplyTheme(t, theme.DarkTheme())
	test.AssertImageMatches(t, "color/dark.png", w.Canvas().Capture())

	w.Close()
}

func TestColorDialog_Advanced_Theme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(800, 800))

	d := NewColorPicker("Color Picker", "Pick a Color", nil, w)
	d.Show()

	d.advanced.Open(0)
	d.Resize(d.win.MinSize()) // TODO FIXME Hack

	test.AssertImageMatches(t, "color/advanced_light.png", w.Canvas().Capture())

	test.ApplyTheme(t, theme.DarkTheme())
	test.AssertImageMatches(t, "color/advanced_dark.png", w.Canvas().Capture())

	w.Close()
}

func TestColorDialog_Recents(t *testing.T) {
	a := test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	// Inject recent preferences
	a.Preferences().SetString("color_recents", "#2196f3,#4caf50,#f44336")

	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(500, 300))

	d := NewColorPicker("Color Picker", "Pick a Color", nil, w)
	d.Show()

	test.AssertImageMatches(t, "color/recents.png", w.Canvas().Capture())

	w.Close()
}

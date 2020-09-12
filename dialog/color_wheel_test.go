package dialog

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
)

func Test_colorWheel_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	wheel := newColorWheel(func(h, s, l, a float64) {
		// Do nothing
	})
	wheel.SetHSLA(0.5, 0.5, 0.5, 0.5)
	window := test.NewWindow(wheel)
	window.Resize(wheel.MinSize().Max(fyne.NewSize(100, 100)))

	test.AssertImageMatches(t, "color/wheel_layout.png", window.Canvas().Capture())

	window.Close()
}

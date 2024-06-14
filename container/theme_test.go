package container

import (
	"image"
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestThemeOverride_Icons(t *testing.T) {
	b := widget.NewButtonWithIcon("", theme.HomeIcon(), func() {})
	o := NewThemeOverride(b, test.Theme())
	w := test.NewWindow(o)
	plain := w.Canvas().Capture().(*image.NRGBA)
	test.AssertImageMatches(t, "theme/icon-test-theme.png", plain)

	o.Theme = test.NewTheme()
	o.Refresh()
	changed := w.Canvas().Capture().(*image.NRGBA)
	test.AssertImageMatches(t, "theme/icon-other-theme.png", changed)
}

func TestThemeOverride_Refresh(t *testing.T) {
	b := widget.NewButton("Test", func() {})
	o := NewThemeOverride(b, test.Theme())
	w := test.NewWindow(o)
	plain := w.Canvas().Capture().(*image.NRGBA)
	test.AssertImageMatches(t, "theme/text-test-theme.png", plain)

	o.Theme = test.NewTheme()
	o.Refresh()
	changed := w.Canvas().Capture().(*image.NRGBA)
	test.AssertImageMatches(t, "theme/text-other-theme.png", changed)
}

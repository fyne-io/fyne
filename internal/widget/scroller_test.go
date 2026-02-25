package widget_test

import (
	"image/color"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func TestScrollContainer_Theme(t *testing.T) {
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(250, 250))
	scroll := widget.NewScroll(rect)

	w := test.NewTempWindow(t, scroll)
	w.SetPadded(false)
	w.Resize(fyne.NewSize(100, 100))
	test.AssertImageMatches(t, "scroll/theme_initial.png", w.Canvas().Capture())

	test.WithTestTheme(t, func() {
		time.Sleep(100 * time.Millisecond)
		scroll.Refresh()
		test.AssertImageMatches(t, "scroll/theme_changed.png", w.Canvas().Capture())
	})
}

func TestScrollContainer_ThemeOverride(t *testing.T) {
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(250, 250))
	scroll := widget.NewScroll(rect)
	scroll.Resize(fyne.NewSize(100, 100))

	w := test.NewTempWindow(t, scroll)
	w.SetPadded(false)
	w.Resize(fyne.NewSize(100, 100))
	test.ApplyTheme(t, test.NewTheme())
	test.AssertImageMatches(t, "scroll/theme_changed.png", w.Canvas().Capture())

	normal := test.Theme()
	bg := canvas.NewRectangle(normal.Color(theme.ColorNameBackground, theme.VariantDark))
	w.SetContent(container.NewStack(bg, container.NewThemeOverride(scroll, normal)))
	w.Resize(fyne.NewSize(100, 100))
	// TODO why is this off by a 1bit RGB difference?
	// test.AssertImageMatches(t, "scroll/theme_initial.png", w.Canvas().Capture())
}

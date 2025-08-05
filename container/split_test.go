package container_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func TestSplitContainer_ThemeOverride(t *testing.T) {
	test.NewTempApp(t)

	split := container.NewHSplit(canvas.NewRectangle(color.Transparent), canvas.NewRectangle(color.Transparent))
	w := test.NewWindow(split)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(100, 100))
	c := w.Canvas()

	test.AssertRendersToImage(t, "split/default.png", c)

	test.ApplyTheme(t, test.NewTheme())
	test.AssertRendersToImage(t, "split/ugly.png", c)

	// set a BG that matches the theme, this is outside our container scope
	normal := test.Theme()
	bg := canvas.NewRectangle(normal.Color(theme.ColorNameBackground, theme.VariantDark))
	w.SetContent(container.NewStack(bg, container.NewThemeOverride(split, normal)))
	w.Resize(fyne.NewSize(100, 100))
	test.AssertRendersToImage(t, "split/default.png", c)
}

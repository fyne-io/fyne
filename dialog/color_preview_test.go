package dialog

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func Test_colorPreview_Color(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	preview := newColorPreview(color.RGBA{53, 113, 233, 255})
	preview.SetColor(color.RGBA{90, 206, 80, 180})
	window := test.NewWindow(preview)
	padding := theme.Padding() * 2
	window.Resize(fyne.NewSize(128+padding, 64+padding))

	test.AssertRendersToImage(t, "color/preview_color.png", window.Canvas())

	window.Close()
}

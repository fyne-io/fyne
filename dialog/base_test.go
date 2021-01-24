package dialog

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestShowCustom_ApplyTheme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(200, 300))

	label := widget.NewLabel("Content")
	label.Alignment = fyne.TextAlignCenter

	d := NewCustom("Title", "OK", label, w)

	d.Show()
	test.AssertImageMatches(t, "dialog-custom-default.png", w.Canvas().Capture())

	test.ApplyTheme(t, test.NewTheme())
	test.AssertImageMatches(t, "dialog-custom-ugly.png", w.Canvas().Capture())
}

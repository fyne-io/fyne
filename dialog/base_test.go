package dialog

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestShowCustom_ApplyTheme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	w := test.NewWindow(canvas.NewRectangle(color.Transparent))

	label := widget.NewLabel("Content")
	label.Alignment = fyne.TextAlignCenter

	d := NewCustom("Title", "OK", label, w)
	w.Resize(d.MinSize())

	d.Show()
	test.AssertImageMatches(t, "dialog-custom-default.png", w.Canvas().Capture())

	test.ApplyTheme(t, test.NewTheme())
	w.Resize(d.MinSize())
	test.AssertImageMatches(t, "dialog-custom-ugly.png", w.Canvas().Capture())
}

func TestShowCustom_Resize(t *testing.T) {
	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(300, 300))

	label := widget.NewLabel("Content")
	label.Alignment = fyne.TextAlignCenter
	d := NewCustom("Title", "OK", label, w)

	size := fyne.NewSize(200, 200)
	d.Resize(size)
	d.Show()
	assert.Equal(t, size, d.(*dialog).win.Content.Size().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
}

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
	shadowPad := float32(50)
	w.Resize(d.MinSize().Add(fyne.NewSize(shadowPad, shadowPad)))

	d.Show()
	test.AssertRendersToImage(t, "dialog-custom-default.png", w.Canvas())

	test.ApplyTheme(t, test.NewTheme())
	w.Resize(d.MinSize().Add(fyne.NewSize(shadowPad, shadowPad)))
	d.Resize(d.MinSize()) // TODO remove once #707 is resolved
	test.AssertRendersToImage(t, "dialog-custom-ugly.png", w.Canvas())
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

func TestCustom_ApplyThemeOnShow(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(200, 300))

	label := widget.NewLabel("Content")
	label.Alignment = fyne.TextAlignCenter
	d := NewCustom("Title", "OK", label, w)

	test.ApplyTheme(t, test.Theme())
	d.Show()
	test.AssertRendersToImage(t, "dialog-onshow-theme-default.png", w.Canvas())
	d.Hide()

	test.ApplyTheme(t, test.NewTheme())
	d.Show()
	test.AssertRendersToImage(t, "dialog-onshow-theme-changed.png", w.Canvas())
	d.Hide()

	test.ApplyTheme(t, test.Theme())
	d.Show()
	test.AssertRendersToImage(t, "dialog-onshow-theme-default.png", w.Canvas())
	d.Hide()
}

func TestCustom_ResizeOnShow(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	size := fyne.NewSize(200, 300)
	w.Resize(size)

	label := widget.NewLabel("Content")
	label.Alignment = fyne.TextAlignCenter
	d := NewCustom("Title", "OK", label, w).(*dialog)

	d.Show()
	assert.Equal(t, size, d.win.Size())
	d.Hide()

	size = fyne.NewSize(500, 500)
	w.Resize(size)
	d.Show()
	assert.Equal(t, size, d.win.Size())
	d.Hide()
}

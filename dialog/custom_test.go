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
	test.NewTempApp(t)

	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))

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
	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(300, 300))

	label := widget.NewLabel("Content")
	label.Alignment = fyne.TextAlignCenter
	d := NewCustom("Title", "OK", label, w)

	size := fyne.NewSize(200, 200)
	d.Resize(size)
	d.Show()
	assert.Equal(t, size, d.dialog.win.Content.Size().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
}

func TestCustom_ApplyThemeOnShow(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
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
	test.NewTempApp(t)
	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	size := fyne.NewSize(200, 300)
	w.Resize(size)

	label := widget.NewLabel("Content")
	label.Alignment = fyne.TextAlignCenter
	d := NewCustom("Title", "OK", label, w).dialog

	d.Show()
	assert.Equal(t, size, d.win.Size())
	d.Hide()

	size = fyne.NewSize(500, 500)
	w.Resize(size)
	d.Show()
	assert.Equal(t, size, d.win.Size())
	d.Hide()
}

func TestConfirm_SetButtons(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	size := fyne.NewSize(200, 300)
	w.Resize(size)

	label := widget.NewLabel("Custom buttons")
	d := NewCustom("Test", "Initial", label, w)
	test.ApplyTheme(t, test.Theme())

	d.SetButtons([]fyne.CanvasObject{&widget.Button{Text: "1"}, &widget.Button{Text: "2"}, &widget.Button{Text: "3"}})
	d.Show()
	test.AssertRendersToMarkup(t, "dialog-custom-custom-buttons.xml", w.Canvas())
	assert.Nil(t, d.dialog.dismiss)
	d.Hide()

	d.SetButtons(nil)
	d.Show()
	test.AssertRendersToMarkup(t, "dialog-custom-no-buttons.xml", w.Canvas())
	d.Hide()
}

func TestConfirmWithoutButtons(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	size := fyne.NewSize(200, 300)
	w.Resize(size)

	test.ApplyTheme(t, test.Theme())
	label := widget.NewLabel("No buttons")
	ShowCustomWithoutButtons("Empty", label, w)
	test.AssertRendersToImage(t, "dialog-custom-without-buttons.png", w.Canvas())
}

func TestCustomConfirm_Importance(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	size := fyne.NewSize(200, 300)
	w.Resize(size)

	label := widget.NewLabel("This is dangerous!")
	d := NewCustomConfirm("Delete me?", "Delete", "Dismiss", label, nil, w)
	d.SetConfirmImportance(widget.DangerImportance)

	test.ApplyTheme(t, test.Theme())
	d.Show()
	test.AssertRendersToImage(t, "dialog-custom-confirm-importance.png", w.Canvas())
}

func TestCustom_SetIcon(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	size := fyne.NewSize(200, 300)
	w.Resize(size)

	test.ApplyTheme(t, test.Theme())
	label := widget.NewLabel("Test was successful.")
	d := NewCustom("Test result", "Close", label, w)
	d.SetIcon(theme.ConfirmIcon())
	d.Show()

	test.AssertRendersToImage(t, "dialog-custom-seticon-success.png", w.Canvas())

	d.Hide()
	d.SetIcon(nil)
	d.Show()

	test.AssertRendersToImage(t, "dialog-custom-seticon-nil.png", w.Canvas())
}

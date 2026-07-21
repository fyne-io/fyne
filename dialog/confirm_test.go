package dialog

import (
	"image/color"
	"os"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestDialog_ConfirmDoubleCallback(t *testing.T) {
	ch := make(chan int)
	cnf := NewConfirm("Test", "Test", func(_ bool) {
		ch <- 42
	}, test.NewTempWindow(t, nil))
	cnf.SetDismissText("No")
	cnf.SetConfirmText("Yes")
	cnf.SetOnClosed(func() {
		ch <- 43
	})
	cnf.Show()

	assert.False(t, cnf.win.Hidden)
	go cnf.Dismiss()
	assert.EqualValues(t, 42, <-ch)
	assert.EqualValues(t, 43, <-ch)
	assert.True(t, cnf.win.Hidden)
}

func TestDialog_ConfirmCallbackOnlyOnClosed(t *testing.T) {
	ch := make(chan int)
	cnf := NewConfirm("Test", "Test", nil, test.NewTempWindow(t, nil))
	cnf.SetDismissText("No")
	cnf.SetConfirmText("Yes")
	cnf.SetOnClosed(func() {
		ch <- 43
	})
	cnf.Show()

	assert.False(t, cnf.win.Hidden)
	go cnf.Dismiss()
	assert.EqualValues(t, 43, <-ch)
	assert.True(t, cnf.win.Hidden)
}

func TestDialog_ConfirmCallbackOnlyOnConfirm(t *testing.T) {
	ch := make(chan int)
	cnf := NewConfirm("Test", "Test", func(ok bool) {
		if !ok {
			ch <- 0
			return
		}
		ch <- 42
	}, test.NewTempWindow(t, nil))
	cnf.SetDismissText("No")
	cnf.SetConfirmText("Yes")
	cnf.Show()

	assert.False(t, cnf.win.Hidden)
	go cnf.Dismiss()
	assert.EqualValues(t, 0, <-ch)
	assert.True(t, cnf.win.Hidden)

	cnf.Show()
	assert.False(t, cnf.win.Hidden)
	go cnf.Confirm()
	assert.EqualValues(t, 42, <-ch)
	assert.True(t, cnf.win.Hidden)
}

func TestConfirmDialog_Resize(t *testing.T) {
	window := test.NewWindow(nil)
	defer window.Close()
	window.Resize(fyne.NewSize(600, 400))
	d := NewConfirm("Test", "Test", nil, window)

	d.Show() // we cannot check window size if not shown

	// Test resize - normal size scenario
	size := fyne.NewSize(300, 180) // normal size to fit (600,400)
	d.Resize(size)
	expectedWidth := float32(300)
	assert.Equal(t, expectedWidth, d.win.Content.Size().Width)
	expectedHeight := float32(180)
	assert.Equal(t, expectedHeight, d.win.Content.Size().Height)
	// Test resize - normal size scenario again
	size = fyne.NewSize(310, 280) // normal size to fit (600,400)
	d.Resize(size)
	expectedWidth = 310
	assert.Equal(t, expectedWidth, d.win.Content.Size().Width)
	expectedHeight = 280
	assert.Equal(t, expectedHeight, d.win.Content.Size().Height)
	d.Hide()

	// Test resize - greater than max size scenario
	size = fyne.NewSize(800, 600)
	d.Resize(size)
	d.Show()
	expectedWidth = 600                                // since win width only 600
	assert.Equal(t, expectedWidth, d.win.Size().Width) // max, also work
	assert.Equal(t, expectedWidth, d.win.Content.Size().Width)
	expectedHeight = 400                                 // since win height only 400
	assert.Equal(t, expectedHeight, d.win.Size().Height) // max, also work
	assert.Equal(t, expectedHeight, d.win.Content.Size().Height)
	d.Hide()

	// Test again - tiny size
	size = fyne.NewSize(1, 1)
	d.Resize(size)
	d.Show()
	expectedWidth = d.win.Content.MinSize().Width
	assert.Equal(t, expectedWidth, d.win.Content.Size().Width)
	expectedHeight = d.win.Content.MinSize().Height
	assert.Equal(t, expectedHeight, d.win.Content.Size().Height)
	d.Hide()
}

func TestConfirm_Importance(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	size := fyne.NewSize(200, 300)
	w.Resize(size)

	d := NewConfirm("Delete me?", "This is dangerous!", nil, w)
	d.SetConfirmImportance(widget.DangerImportance)

	test.ApplyTheme(t, test.Theme())
	d.Show()
	test.AssertRendersToImage(t, "dialog-confirm-importance.png", w.Canvas())
}

func TestConfirm_Importance_Blur(t *testing.T) {
	test.NewTempApp(t)
	code := widget.NewEntry()
	data, _ := os.ReadFile("./testdata/Capitalised.txt")
	code.SetText(string(data))
	w := test.NewTempWindow(t, code)
	size := fyne.NewSize(480, 320)
	w.Resize(size)

	d := NewConfirm("Delete me?", "This is dangerous!", nil, w)
	d.SetConfirmImportance(widget.DangerImportance)

	d.Show()
	test.AssertRendersToImage(t, "dialog-confirm-blur.png", w.Canvas())
}

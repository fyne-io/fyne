package dialog

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestDialog_ConfirmDoubleCallback(t *testing.T) {
	ch := make(chan int)
	cnf := NewConfirm("Test", "Test", func(_ bool) {
		ch <- 42
	}, test.NewWindow(nil))
	cnf.SetDismissText("No")
	cnf.SetConfirmText("Yes")
	cnf.SetOnClosed(func() {
		ch <- 43
	})
	cnf.Show()

	assert.False(t, cnf.win.Hidden)
	go test.Tap(cnf.dismiss)
	assert.EqualValues(t, <-ch, 43)
	assert.EqualValues(t, <-ch, 42)
	assert.True(t, cnf.win.Hidden)
}

func TestDialog_ConfirmCallbackOnlyOnClosed(t *testing.T) {
	ch := make(chan int)
	cnf := NewConfirm("Test", "Test", nil, test.NewWindow(nil))
	cnf.SetDismissText("No")
	cnf.SetConfirmText("Yes")
	cnf.SetOnClosed(func() {
		ch <- 43
	})
	cnf.Show()

	assert.False(t, cnf.win.Hidden)
	go test.Tap(cnf.dismiss)
	assert.EqualValues(t, <-ch, 43)
	assert.True(t, cnf.win.Hidden)
}

func TestDialog_ConfirmCallbackOnlyOnConfirm(t *testing.T) {
	ch := make(chan int)
	cnf := NewConfirm("Test", "Test", func(_ bool) {
		ch <- 42
	}, test.NewWindow(nil))
	cnf.SetDismissText("No")
	cnf.SetConfirmText("Yes")
	cnf.Show()

	assert.False(t, cnf.win.Hidden)
	go test.Tap(cnf.dismiss)
	assert.EqualValues(t, <-ch, 42)
	assert.True(t, cnf.win.Hidden)
}

func TestConfirmDialog_Resize(t *testing.T) {
	window := test.NewWindow(nil)
	window.Resize(fyne.NewSize(600, 400))
	defer window.Close()
	d := NewConfirm("Test", "Test", nil, window)

	theDialog := d.dialog
	d.dialog.Show() // we cannot check window size if not shown

	//Test resize - normal size scenario
	size := fyne.NewSize(300, 180) //normal size to fit (600,400)
	theDialog.Resize(size)
	expectedWidth := float32(300)
	assert.Equal(t, expectedWidth, theDialog.win.Content.Size().Width+theme.Padding()*2)
	expectedHeight := float32(180)
	assert.Equal(t, expectedHeight, theDialog.win.Content.Size().Height+theme.Padding()*2)
	//Test resize - normal size scenario again
	size = fyne.NewSize(310, 280) //normal size to fit (600,400)
	theDialog.Resize(size)
	expectedWidth = 310
	assert.Equal(t, expectedWidth, theDialog.win.Content.Size().Width+theme.Padding()*2)
	expectedHeight = 280
	assert.Equal(t, expectedHeight, theDialog.win.Content.Size().Height+theme.Padding()*2)

	//Test resize - greater than max size scenario
	size = fyne.NewSize(800, 600)
	theDialog.Resize(size)
	expectedWidth = 600                                        //since win width only 600
	assert.Equal(t, expectedWidth, theDialog.win.Size().Width) //max, also work
	assert.Equal(t, expectedWidth, theDialog.win.Content.Size().Width+theme.Padding()*2)
	expectedHeight = 400                                         //since win heigh only 400
	assert.Equal(t, expectedHeight, theDialog.win.Size().Height) //max, also work
	assert.Equal(t, expectedHeight, theDialog.win.Content.Size().Height+theme.Padding()*2)

	//Test again - extreme small size
	size = fyne.NewSize(1, 1)
	theDialog.Resize(size)
	expectedWidth = theDialog.win.Content.MinSize().Width
	assert.Equal(t, expectedWidth, theDialog.win.Content.Size().Width)
	expectedHeight = theDialog.win.Content.MinSize().Height
	assert.Equal(t, expectedHeight, theDialog.win.Content.Size().Height)
}

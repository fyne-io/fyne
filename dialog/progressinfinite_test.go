package dialog

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/stretchr/testify/assert"
)

func TestProgressInfiniteDialog_MinSize(t *testing.T) {
	window := test.NewWindow(nil)
	defer window.Close()
	d := NewProgressInfinite("title", "message", window)

	dialogContent := d.win.Content.MinSize()
	progressBar := d.bar.MinSize()

	assert.Less(t, progressBar.Width, dialogContent.Width)
}

func TestProgressInfiniteDialog_Resize(t *testing.T) {
	window := test.NewWindow(nil)
	window.Resize(fyne.NewSize(600, 400))
	defer window.Close()
	d := NewProgressInfinite("title", "message", window)
	theDialog := d.dialog

	//Test resize - normal size scenario
	size := fyne.NewSize(300, 180) //normal size to fit (600,400)
	theDialog.Resize(size)
	expectedWidth := 300
	assert.Equal(t, expectedWidth, theDialog.win.Content.Size().Width+theme.Padding()*2)
	expectedHeight := 180
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

func TestProgressInfiniteDialog_Content(t *testing.T) {
	title := "title"
	message := "message"

	window := test.NewWindow(nil)
	defer window.Close()
	d := NewProgressInfinite(title, message, window)

	assert.Equal(t, d.title, title)
	assert.Equal(t, d.content.(*widget.Label).Text, message)
}

func TestProgressInfiniteDialog_Show(t *testing.T) {
	window := test.NewWindow(nil)
	defer window.Close()
	d := NewProgressInfinite("title", "message", window)

	d.Show()

	assert.False(t, d.win.Hidden)
	assert.True(t, d.bar.Running())

	d.Hide()

	assert.True(t, d.win.Hidden)
	assert.False(t, d.bar.Running())
}

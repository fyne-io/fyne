package dialog

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
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
	expectedWidth = d.win.Content.MinSize().Width
	assert.Equal(t, expectedWidth, d.win.Content.Size().Width)
	expectedHeight = d.win.Content.MinSize().Height
	assert.Equal(t, expectedHeight, d.win.Content.Size().Height)
	d.Hide()
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

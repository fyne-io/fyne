package dialog

import (
	"testing"

	"fyne.io/fyne/test"
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

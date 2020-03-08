package dialog

import (
	"testing"

	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestDialog_MinSize(t *testing.T) {
	window := test.NewWindow(nil)
	defer window.Close()
	d := NewInformation("Looooooooooooooong title", "message...", window)
	information := d.(*dialog)

	dialogContent := information.win.Content.MinSize()
	label := information.label.MinSize()

	assert.Less(t, label.Width, dialogContent.Width)
}

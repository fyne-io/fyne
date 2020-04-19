package dialog

import (
	"errors"
	"testing"

	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestDialog_MinSize(t *testing.T) {
	d := NewInformation("Looooooooooooooong title", "message...", test.NewWindow(nil))
	information := d.(*dialog)

	dialogContent := information.win.Content.MinSize()
	label := information.label.MinSize()

	assert.Less(t, label.Width, dialogContent.Width)
}

func TestDialog_InformationCallback(t *testing.T) {
	d := NewInformation("Information", "Hello World", test.NewWindow(nil))
	tapped := false
	d.SetOnClosed(func() { tapped = true })
	d.Show()

	information := d.(*dialog)
	assert.False(t, information.win.Hidden)
	test.Tap(information.dismiss)
	assert.True(t, tapped)
	assert.True(t, information.win.Hidden)
}

func TestDialog_ErrorCallback(t *testing.T) {
	err := errors.New("Error message")
	d := NewError(err, test.NewWindow(nil))
	tapped := false
	d.SetOnClosed(func() { tapped = true })
	d.Show()

	information := d.(*dialog)
	assert.False(t, information.win.Hidden)
	test.Tap(information.dismiss)
	assert.True(t, tapped)
	assert.True(t, information.win.Hidden)
}

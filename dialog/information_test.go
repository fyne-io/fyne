package dialog

import (
	"errors"
	"testing"
	"time"

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
	tapped := make(chan bool)
	d.SetOnClosed(func() { tapped <- true })
	d.Show()

	information := d.(*dialog)
	assert.False(t, information.win.Hidden)
	go test.Tap(information.dismiss)
	func() {
		select {
		case <-tapped:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for button tap")
		}
	}()
	assert.True(t, information.win.Hidden)
}

func TestDialog_ErrorCallback(t *testing.T) {
	err := errors.New("Error message")
	d := NewError(err, test.NewWindow(nil))
	tapped := make(chan bool)
	d.SetOnClosed(func() { tapped <- true })
	d.Show()

	information := d.(*dialog)
	assert.False(t, information.win.Hidden)
	go test.Tap(information.dismiss)
	func() {
		select {
		case <-tapped:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for button tap")
		}
	}()
	assert.True(t, information.win.Hidden)
}

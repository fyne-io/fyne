package dialog

import (
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

func TestDialog_Callback(t *testing.T) {
	tapped := make(chan bool)

	d := NewInformationWithCallback("Information", "Hello World",
		func() {
			tapped <- true
		}, test.NewWindow(nil))
	d.Show()
	information := d.(*dialog)
	assert.False(t, information.win.Hidden)

	go test.Tap(information.dismiss)
	func() {
		select {
		case <-tapped:
		case <-time.After(10 * time.Second):
			assert.Fail(t, "Timed out waiting for button tap")
		}
	}()
	assert.True(t, information.win.Hidden)
}

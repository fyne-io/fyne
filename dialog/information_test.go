package dialog

import (
	"errors"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
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

func TestDialog_Resize(t *testing.T) {
	window := test.NewWindow(nil)
	window.Resize(fyne.NewSize(600, 400))
	defer window.Close()
	d := NewInformation("Looooooooooooooong title", "message...", window)
	theDialog := d.(*dialog)

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

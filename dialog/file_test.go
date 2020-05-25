package dialog

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
)

func TestShowFileOpen(t *testing.T) {
	var chosen fyne.FileReadCloser
	var openErr error
	win := test.NewWindow(widget.NewLabel("Content"))
	NewFileOpenDialog(func(file fyne.FileReadCloser, err error) {
		chosen = file
		openErr = err
	}, win).Show()

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*fyne.Container)
	title := ui.Objects[1].(*widget.Label)
	assert.Equal(t, "Open File", title.Text)

	nameLabel := ui.Objects[2].(*fyne.Container).Objects[1].(*widget.ScrollContainer).Content.(*widget.Label)
	buttons := ui.Objects[2].(*fyne.Container).Objects[0].(*widget.Box)
	open := buttons.Children[1].(*widget.Button)

	files := ui.Objects[3].(*fyne.Container).Objects[1].(*widget.ScrollContainer).Content.(*fyne.Container)
	assert.Greater(t, len(files.Objects), 0)

	fileName := files.Objects[0].(*fileDialogItem).name
	assert.Equal(t, "(Parent)", fileName)
	assert.True(t, open.Disabled())

	var target *fileDialogItem
	for _, icon := range files.Objects {
		if icon.(*fileDialogItem).dir == false {
			target = icon.(*fileDialogItem)
		}
	}

	if target == nil {
		log.Println("Could not find a file in the default directory to tap :(")
		return
	}

	// This will only execute if we have a file in the home path.
	// Until we have a way to set the directory of an open file dialog.
	test.Tap(target)
	assert.Equal(t, filepath.Base(target.path), nameLabel.Text)
	assert.False(t, open.Disabled())

	test.Tap(open)
	assert.Nil(t, win.Canvas().Overlays().Top())
	assert.Nil(t, openErr)
	assert.Equal(t, "file://"+target.path, chosen.URI().String())

	err := chosen.Close()
	assert.Nil(t, err)
}

func TestShowFileSave(t *testing.T) {
	var chosen fyne.FileWriteCloser
	var saveErr error
	win := test.NewWindow(widget.NewLabel("Content"))
	NewFileSaveDialog(func(file fyne.FileWriteCloser, err error) {
		chosen = file
		saveErr = err
	}, win).Show()

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*fyne.Container)
	title := ui.Objects[1].(*widget.Label)
	assert.Equal(t, "Save File", title.Text)

	nameEntry := ui.Objects[2].(*fyne.Container).Objects[1].(*widget.ScrollContainer).Content.(*widget.Entry)
	buttons := ui.Objects[2].(*fyne.Container).Objects[0].(*widget.Box)
	save := buttons.Children[1].(*widget.Button)

	files := ui.Objects[3].(*fyne.Container).Objects[1].(*widget.ScrollContainer).Content.(*fyne.Container)
	assert.Greater(t, len(files.Objects), 0)

	fileName := files.Objects[0].(*fileDialogItem).name
	assert.Equal(t, "(Parent)", fileName)
	assert.True(t, save.Disabled())

	var target *fileDialogItem
	for _, icon := range files.Objects {
		if icon.(*fileDialogItem).dir == false {
			target = icon.(*fileDialogItem)
		}
	}

	if target == nil {
		log.Println("Could not find a file in the default directory to tap :(")
		return
	}

	// This will only execute if we have a file in the home path.
	// Until we have a way to set the directory of an open file dialog.
	test.Tap(target)
	assert.Equal(t, filepath.Base(target.path), nameEntry.Text)
	assert.False(t, save.Disabled())

	// we are about to overwrite, a warning will show
	test.Tap(save)
	confirmUI := win.Canvas().Overlays().Top().(*widget.PopUp)
	assert.NotEqual(t, confirmUI, popup)
	confirmUI.Hide()

	// give the file a unique name and it will callback fine
	test.Type(nameEntry, "v2_")
	test.Tap(save)
	assert.Nil(t, win.Canvas().Overlays().Top())
	assert.Nil(t, saveErr)
	expectedPath := filepath.Join(filepath.Dir(target.path), "v2_"+filepath.Base(target.path))
	assert.Equal(t, "file://"+expectedPath, chosen.URI().String())

	err := chosen.Close()
	assert.Nil(t, err)
	err = os.Remove(expectedPath)
	assert.Nil(t, err)
}

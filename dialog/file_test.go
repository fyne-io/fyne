package dialog

import (
	"log"
	"path/filepath"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/stretchr/testify/assert"
)

func TestShowFileOpen(t *testing.T) {
	chosen := ""
	win := test.NewWindow(widget.NewLabel("Content"))
	ShowFileOpen(func(file string) {
		chosen = file
	}, win)

	popup := win.Canvas().Overlays().Top().(*widget.PopUp).Content
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.(*fyne.Container).Objects[1].(*fyne.Container)
	title := ui.Objects[1].(*widget.Label)
	assert.Equal(t, "Open File", title.Text)

	nameLabel := ui.Objects[2].(*fyne.Container).Objects[1].(*widget.Label)
	buttons := ui.Objects[2].(*fyne.Container).Objects[0].(*widget.Box)
	open := buttons.Children[1].(*widget.Button)

	files := ui.Objects[3].(*fyne.Container).Objects[1].(*widget.ScrollContainer).Content.(*fyne.Container)
	assert.Greater(t, len(files.Objects), 0)

	fileName := files.Objects[0].(*fileIcon).name
	assert.Equal(t, "(Parent)", fileName)
	assert.True(t, open.Disabled())

	var target *fileIcon
	for _, icon := range files.Objects {
		if icon.(*fileIcon).icon == theme.FileIcon() {
			target = icon.(*fileIcon)
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
	assert.Equal(t, target.path, chosen)
}

func TestShowFileSave(t *testing.T) {
	chosen := ""
	win := test.NewWindow(widget.NewLabel("Content"))
	ShowFileSave(func(file string) {
		chosen = file
	}, win)

	popup := win.Canvas().Overlays().Top().(*widget.PopUp).Content
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.(*fyne.Container).Objects[1].(*fyne.Container)
	title := ui.Objects[1].(*widget.Label)
	assert.Equal(t, "Save File", title.Text)

	nameEntry := ui.Objects[2].(*fyne.Container).Objects[1].(*widget.Entry)
	buttons := ui.Objects[2].(*fyne.Container).Objects[0].(*widget.Box)
	save := buttons.Children[1].(*widget.Button)

	files := ui.Objects[3].(*fyne.Container).Objects[1].(*widget.ScrollContainer).Content.(*fyne.Container)
	assert.Greater(t, len(files.Objects), 0)

	fileName := files.Objects[0].(*fileIcon).name
	assert.Equal(t, "(Parent)", fileName)
	assert.True(t, save.Disabled())

	var target *fileIcon
	for _, icon := range files.Objects {
		if icon.(*fileIcon).icon == theme.FileIcon() {
			target = icon.(*fileIcon)
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
	assert.Equal(t, filepath.Join(filepath.Dir(target.path), "v2_"+filepath.Base(target.path)), chosen)
}

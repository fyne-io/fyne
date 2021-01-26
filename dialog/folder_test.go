package dialog

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestShowFolderOpen(t *testing.T) {
	var chosen fyne.ListableURI
	var openErr error
	win := test.NewWindow(widget.NewLabel("OpenDir"))
	d := NewFolderOpen(func(file fyne.ListableURI, err error) {
		chosen = file
		openErr = err
	}, win)
	testData, _ := filepath.Abs("testdata")
	dir, err := storage.ListerForURI(storage.NewFileURI(testData))
	if err != nil {
		t.Error("Failed to open testdata dir", err)
	}
	d.SetLocation(dir)
	d.Show()

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*fyne.Container)
	title := ui.Objects[1].(*widget.Label)
	assert.Equal(t, "Open Folder", title.Text)

	nameLabel := ui.Objects[2].(*fyne.Container).Objects[1].(*container.Scroll).Content.(*widget.Label)
	buttons := ui.Objects[2].(*fyne.Container).Objects[0].(*fyne.Container)
	open := buttons.Objects[1].(*widget.Button)

	files := ui.Objects[0].(*fyne.Container).Objects[1].(*container.Scroll).Content.(*fyne.Container)
	assert.Greater(t, len(files.Objects), 0)

	fileName := files.Objects[0].(*fileDialogItem).name
	assert.Equal(t, "(Parent)", fileName)
	assert.False(t, open.Disabled())

	var target *fileDialogItem
	for _, icon := range files.Objects {
		if icon.(*fileDialogItem).dir {
			target = icon.(*fileDialogItem)
		} else {
			t.Error("Folder dialog should not list files")
		}
	}

	assert.NotNil(t, target, "Failed to find folder in testdata")
	test.Tap(target)
	assert.Equal(t, target.location.Name(), nameLabel.Text)
	assert.False(t, open.Disabled())

	test.Tap(open)
	assert.Nil(t, win.Canvas().Overlays().Top())
	assert.Nil(t, openErr)

	assert.Equal(t, target.location.String(), chosen.String())
}

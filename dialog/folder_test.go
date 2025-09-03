package dialog

import (
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2/lang"
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
	win := test.NewTempWindow(t, widget.NewLabel("OpenDir"))
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
	title := ui.Objects[1].(*fyne.Container).Objects[1].(*widget.Label)
	assert.Equal(t, lang.L("Open")+" "+lang.L("Folder"), title.Text)

	nameLabel := ui.Objects[2].(*fyne.Container).Objects[1].(*container.Scroll).Content.(*widget.Label)
	buttons := ui.Objects[2].(*fyne.Container).Objects[0].(*fyne.Container)
	open := buttons.Objects[1].(*widget.Button)

	files := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[1].(*container.Scroll).Content.(*fyne.Container).Objects[0].(*widget.GridWrap)
	assert.NotEmpty(t, d.dialog.data)

	item := test.TempWidgetRenderer(t, files).Objects()[0].(*container.Scroll).Content.(*fyne.Container).Objects[0]
	fileName := test.TempWidgetRenderer(t, item.(fyne.Widget)).Objects()[1].(*fileDialogItem).name
	assert.Equal(t, lang.L("(Parent)"), fileName)
	assert.False(t, open.Disabled())

	var target *fyne.URI
	id := -1
	for i, uri := range d.dialog.data {
		ok, _ := storage.CanList(uri)
		if ok {
			target = &uri
			id = i
		} else {
			t.Error("Folder dialog should not list files")
		}
	}

	assert.NotNil(t, target, "Failed to find folder in testdata")
	d.dialog.files.(*widget.GridWrap).Select(id)
	assert.Equal(t, (*target).Name(), nameLabel.Text)
	assert.False(t, open.Disabled())

	test.Tap(open)
	assert.Nil(t, win.Canvas().Overlays().Top())
	assert.NoError(t, openErr)

	assert.Equal(t, (*target).String(), chosen.String())
}

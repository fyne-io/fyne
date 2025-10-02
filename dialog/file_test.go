package dialog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	intWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/lang"
	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestEffectiveStartingDir(t *testing.T) {
	homeString, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("os.Gethome() failed, cannot run this test on this system (error stat()-ing ../) error was '%s'", err)
	}
	home, err := storage.ListerForURI(storage.NewFileURI(homeString))
	if err != nil {
		t.Skipf("could not get lister for working directory: %s", err)
	}

	parentURI, err := storage.Parent(home)
	if err != nil {
		t.Skipf("Could not get parent of working directory: %s", err)
	}

	parent, err := storage.ListerForURI(parentURI)
	if err != nil {
		t.Skipf("Could not get lister for parent of working directory: %s", err)
	}

	dialog := &FileDialog{}

	// test that we get wd when running with the default struct values
	res := dialog.effectiveStartingDir()
	expect := home
	if !storage.EqualURI(res, expect) {
		t.Errorf("Expected effectiveStartingDir() to be '%s', but it was '%s'",
			expect, res)
	}

	// this should always be equivalent to the preceding test
	dialog.startingLocation = nil
	res = dialog.effectiveStartingDir()
	expect = home
	if !storage.EqualURI(res, expect) {
		t.Errorf("Expected effectiveStartingDir() to be '%s', but it was '%s'",
			expect, res)
	}

	// check using StartingDirectory with some other directory
	dialog.startingLocation = parent
	res = dialog.effectiveStartingDir()
	expect = parent
	if res != expect {
		t.Errorf("Expected effectiveStartingDir() to be '%s', but it was '%s'",
			expect, res)
	}

	// make sure we fail over if the specified directory does not exist
	dialog.startingLocation, err = storage.ListerForURI(storage.NewFileURI("/some/file/that/does/not/exist"))
	if err == nil {
		t.Errorf("Should have failed to create lister for nonexistent file")
	}
	res = dialog.effectiveStartingDir()
	expect = home
	if res.String() != expect.String() {
		t.Errorf("Expected effectiveStartingDir() to be '%s', but it was '%s'",
			expect, res)
	}
}

func TestFileDialogStartRemember(t *testing.T) {
	start, err := storage.ListerForURI(storage.NewFileURI("testdata"))
	assert.NoError(t, err)

	w := test.NewTempWindow(t, widget.NewLabel("Content"))
	d := NewFileOpen(nil, w)
	d.SetLocation(start)
	d.Show()

	assert.Equal(t, start.String(), d.dialog.dir.String())
	d.Hide()

	d2 := NewFileOpen(nil, w)
	d2.Show()
	assert.Equal(t, start.String(), d.dialog.dir.String())
	d2.Hide()
}

func TestFileDialogResize(t *testing.T) {
	win := test.NewTempWindow(t, widget.NewLabel("Content"))
	win.Resize(fyne.NewSize(600, 400))
	file := NewFileOpen(func(file fyne.URIReadCloser, err error) {}, win)
	file.SetFilter(storage.NewExtensionFileFilter([]string{".png"}))

	// Mimic the fileopen dialog
	d := &fileDialog{file: file}
	open := widget.NewButton("open", func() {})
	ui := container.NewBorder(nil, nil, nil, open)
	pad := theme.Padding()
	itemMin := d.newFileItem(storage.NewFileURI("filename.txt"), false, false).MinSize()
	originalSize := ui.MinSize().Add(itemMin.AddWidthHeight(itemMin.Width+pad*6, pad*3))
	d.win = widget.NewModalPopUp(ui, file.parent.Canvas())
	d.win.Resize(originalSize)
	file.dialog = d

	// Test resize - normal size scenario
	size := fyne.NewSize(200, 180) // normal size to fit (600,400)
	file.Resize(size)
	expectedWidth := float32(200)
	assert.Equal(t, expectedWidth, file.dialog.win.Content.Size().Width+theme.Padding()*2)
	expectedHeight := float32(180)
	assert.Equal(t, expectedHeight, file.dialog.win.Content.Size().Height+theme.Padding()*2)
	// Test resize - normal size scenario again
	size = fyne.NewSize(300, 280) // normal size to fit (600,400)
	file.Resize(size)
	expectedWidth = 300
	assert.Equal(t, expectedWidth, file.dialog.win.Content.Size().Width+theme.Padding()*2)
	expectedHeight = 280
	assert.Equal(t, expectedHeight, file.dialog.win.Content.Size().Height+theme.Padding()*2)

	// Test resize - greater than max size scenario
	size = fyne.NewSize(800, 600)
	file.Resize(size)
	expectedWidth = 600                                          // since win width only 600
	assert.Equal(t, expectedWidth, file.dialog.win.Size().Width) // max, also work
	assert.Equal(t, expectedWidth, file.dialog.win.Content.Size().Width+theme.Padding()*2)
	expectedHeight = 400                                           // since win height only 400
	assert.Equal(t, expectedHeight, file.dialog.win.Size().Height) // max, also work
	assert.Equal(t, expectedHeight, file.dialog.win.Content.Size().Height+theme.Padding()*2)

	// Test again - extreme small size
	size = fyne.NewSize(1, 1)
	file.Resize(size)
	expectedWidth = file.dialog.win.Content.MinSize().Width
	assert.Equal(t, expectedWidth, file.dialog.win.Content.Size().Width)
	expectedHeight = file.dialog.win.Content.MinSize().Height
	assert.Equal(t, expectedHeight, file.dialog.win.Content.Size().Height)
}

func TestShowFileOpen(t *testing.T) {
	var chosen fyne.URIReadCloser
	var openErr error
	win := test.NewTempWindow(t, widget.NewLabel("Content"))
	d := NewFileOpen(func(file fyne.URIReadCloser, err error) {
		chosen = file
		openErr = err
	}, win)
	testData := storage.NewFileURI("testdata")
	dir, err := storage.ListerForURI(testData)
	if err != nil {
		t.Error("Failed to open testdata dir", err)
	}
	d.SetLocation(dir)
	d.Show()

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*fyne.Container)
	// header
	title := ui.Objects[1].(*fyne.Container).Objects[0].(*widget.Label)
	assert.Equal(t, lang.L("Open")+" "+lang.L("File"), title.Text)
	// optionsbuttons
	createNewFolderButton := ui.Objects[1].(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Button)
	assert.Equal(t, "", createNewFolderButton.Text)
	assert.Equal(t, theme.FolderNewIcon().Name(), createNewFolderButton.Icon.Name())
	toggleViewButton := ui.Objects[1].(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Button)
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.ListIcon().Name(), toggleViewButton.Icon.Name())
	optionsButton := ui.Objects[1].(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Button)
	assert.Equal(t, "", optionsButton.Text)
	assert.Equal(t, theme.SettingsIcon().Name(), optionsButton.Icon.Name())
	// footer
	nameLabel := ui.Objects[2].(*fyne.Container).Objects[0].(*container.Scroll).Content.(*widget.Label)
	buttons := ui.Objects[2].(*fyne.Container).Objects[1].(*fyne.Container)
	open := buttons.Objects[1].(*widget.Button)
	// body
	breadcrumb := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[1].(*container.Scroll).Content.(*fyne.Container).Objects[0].(*fyne.Container)
	assert.NotEmpty(t, breadcrumb.Objects)

	assert.NoError(t, err)
	components := strings.Split(testData.Path(), "/")
	if components[0] == "" {
		// Splitting a unix path will give a "" at the beginning, but we actually want the path bar to show "/".
		components[0] = "/"
	}
	if assert.Equal(t, len(components), len(breadcrumb.Objects)) {
		for i, object := range breadcrumb.Objects {
			assert.Equal(t, components[i], object.(*widget.Button).Text, fmt.Sprintf("Failure for %s at index: %d", testData.Path(), i))
		}
	}

	files := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[0].(*container.Scroll).Content.(*fyne.Container).Objects[0].(*widget.GridWrap)
	objects := test.TempWidgetRenderer(t, files).Objects()[0].(*container.Scroll).Content.(*fyne.Container).Objects
	assert.NotEmpty(t, objects)

	fileName := test.TempWidgetRenderer(t, objects[0].(fyne.Widget)).Objects()[1].(*fileDialogItem).name
	assert.Equal(t, lang.L("(Parent)"), fileName)
	assert.True(t, open.Disabled())

	var target *fileDialogItem
	id := 0
	for i, icon := range objects {
		item := test.TempWidgetRenderer(t, icon.(fyne.Widget)).Objects()[1].(*fileDialogItem)
		if item.dir == false {
			target = item
			id = i
			break
		}
	}
	assert.NotNil(t, target, "Failed to find file in testdata")
	d.dialog.files.(*widget.GridWrap).Select(id)
	assert.Equal(t, target.location.Name(), nameLabel.Text)
	assert.False(t, open.Disabled())

	test.Tap(open)
	assert.Nil(t, win.Canvas().Overlays().Top())
	assert.NoError(t, openErr)

	assert.Equal(t, target.location.String(), chosen.URI().String())

	err = chosen.Close()
	assert.NoError(t, err)
}

func TestHiddenFiles(t *testing.T) {
	dir, err := storage.ListerForURI(storage.NewFileURI("testdata"))
	assert.NoError(t, err)

	// git does not preserve windows hidden flag, so we have to set it.
	// just an empty function for non windows builds
	hidden, _ := storage.Child(dir, ".hidden")
	err = hideFile(hidden.Path())
	assert.NoError(t, err)

	win := test.NewTempWindow(t, widget.NewLabel("Content"))
	d := NewFileOpen(func(file fyne.URIReadCloser, err error) {
	}, win)
	d.SetLocation(dir)
	d.Show()

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*fyne.Container)

	createNewFolderButton := ui.Objects[1].(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Button)
	assert.Equal(t, "", createNewFolderButton.Text)
	assert.Equal(t, theme.FolderNewIcon().Name(), createNewFolderButton.Icon.Name())
	toggleViewButton := ui.Objects[1].(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Button)
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.ListIcon().Name(), toggleViewButton.Icon.Name())
	optionsButton := ui.Objects[1].(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Button)
	assert.Equal(t, "", optionsButton.Text)
	assert.Equal(t, theme.SettingsIcon().Name(), optionsButton.Icon.Name())

	files := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[0].(*container.Scroll).Content.(*fyne.Container).Objects[0].(*widget.GridWrap)
	objects := test.TempWidgetRenderer(t, files).Objects()[0].(*container.Scroll).Content.(*fyne.Container).Objects
	assert.NotEmpty(t, objects)

	var target *fileDialogItem
	for _, icon := range objects {
		item := test.TempWidgetRenderer(t, icon.(fyne.Widget)).Objects()[1].(*fileDialogItem)
		if item.name == ".hidden" {
			target = item
		}
	}
	assert.Nil(t, target, "Failed, .hidden found in testdata")

	d.dialog.showHidden = true
	d.dialog.refreshDir(d.dialog.dir)

	for _, icon := range objects {
		item := test.TempWidgetRenderer(t, icon.(fyne.Widget)).Objects()[1].(*fileDialogItem)
		if item.name == ".hidden" {
			target = item
		}
	}
	assert.NotNil(t, target, "Failed, .hidden not found in testdata")
}

func TestShowFileSave(t *testing.T) {
	var chosen fyne.URIWriteCloser
	var saveErr error
	win := test.NewTempWindow(t, widget.NewLabel("Content"))
	saver := NewFileSave(func(file fyne.URIWriteCloser, err error) {
		chosen = file
		saveErr = err
	}, win)
	saver.Show()

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*fyne.Container)
	title := ui.Objects[1].(*fyne.Container).Objects[0].(*widget.Label)
	assert.Equal(t, "Save File", title.Text)

	nameEntry := ui.Objects[2].(*fyne.Container).Objects[0].(*container.Scroll).Content.(*widget.Entry)
	buttons := ui.Objects[2].(*fyne.Container).Objects[1].(*fyne.Container)
	save := buttons.Objects[1].(*widget.Button)

	files := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[0].(*container.Scroll).Content.(*fyne.Container).Objects[0].(*widget.GridWrap)
	objects := test.TempWidgetRenderer(t, files).Objects()[0].(*container.Scroll).Content.(*fyne.Container).Objects
	assert.NotEmpty(t, objects)

	item := test.TempWidgetRenderer(t, objects[0].(fyne.Widget)).Objects()[1].(*fileDialogItem)
	assert.Equal(t, lang.L("(Parent)"), item.name)
	assert.True(t, save.Disabled())

	dir, err := storage.ListerForURI(storage.NewFileURI("testdata"))
	assert.NoError(t, err)
	saver.SetLocation(dir)

	var target *fileDialogItem
	id := -1
	for i, icon := range objects {
		item := test.TempWidgetRenderer(t, icon.(fyne.Widget)).Objects()[1].(*fileDialogItem)
		if item.dir == false {
			target = item
			id = i
			break
		}
	}

	if target == nil {
		log.Println("Could not find a file in the default directory to tap :(")
		return
	}

	saver.dialog.files.(*widget.GridWrap).Select(id)
	assert.Equal(t, target.location.Name(), nameEntry.Text)
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
	assert.NoError(t, saveErr)
	targetParent, err := storage.Parent(target.location)
	if err != nil {
		t.Error(err)
	}
	expectedPath, _ := storage.Child(targetParent, "v2_"+target.location.Name())
	assert.Equal(t, expectedPath.String(), chosen.URI().String())

	err = chosen.Close()
	assert.NoError(t, err)
	pathString := expectedPath.Path()
	err = os.Remove(pathString)
	assert.NoError(t, err)
}

func TestFileFilters(t *testing.T) {
	win := test.NewTempWindow(t, widget.NewLabel("Content"))
	f := NewFileOpen(func(file fyne.URIReadCloser, err error) {
	}, win)

	f.SetFilter(storage.NewExtensionFileFilter([]string{".png"}))
	f.Show()

	workingDir, err := os.Getwd()
	if err != nil {
		fyne.LogError("Could not get current working directory", err)
		t.FailNow()
	}
	testDataDir := storage.NewFileURI(filepath.Join(workingDir, "testdata"))
	testDataLister, err := storage.ListerForURI(testDataDir)
	if err != nil {
		t.Error(err)
	}

	f.dialog.setLocation(testDataLister)

	count := 0
	for _, uri := range f.dialog.data {
		ok, _ := storage.CanList(uri)
		if !ok {
			assert.Equal(t, ".png", uri.Extension())
			count++
		}
	}

	// NOTE: This count needs to be updated when more test images are added.
	assert.Equal(t, 10, count)

	f.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/jpeg"}))

	count = 0
	for _, uri := range f.dialog.data {
		ok, _ := storage.CanList(uri)
		if !ok {
			assert.Equal(t, "image/jpeg", uri.MimeType())
			count++
		}
	}

	// NOTE: This count needs to be updated when more test images are added.
	assert.Equal(t, 1, count)

	f.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/*"}))

	count = 0
	for _, uri := range f.dialog.data {
		ok, _ := storage.CanList(uri)
		if !ok {
			mimeType := strings.Split(uri.MimeType(), "/")[0]
			assert.Equal(t, "image", mimeType)
			count++
		}
	}

	// NOTE: This count needs to be updated when more test images are added.
	assert.Equal(t, 11, count)
}

func TestFileSort(t *testing.T) {
	dir, err := storage.ListerForURI(storage.NewFileURI("testdata"))
	assert.NoError(t, err)

	win := test.NewTempWindow(t, widget.NewLabel("Content"))
	d := NewFileOpen(func(file fyne.URIReadCloser, err error) {
	}, win)
	d.SetLocation(dir)
	d.Show()

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*fyne.Container)

	files := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[0].(*container.Scroll).Content.(*fyne.Container).Objects[0].(*widget.GridWrap)
	objects := test.TempWidgetRenderer(t, files).Objects()[0].(*container.Scroll).Content.(*fyne.Container).Objects
	assert.NotEmpty(t, objects)

	binPos := -1
	capitalPos := -1
	for i, icon := range objects {
		item := test.TempWidgetRenderer(t, icon.(fyne.Widget)).Objects()[1].(*fileDialogItem)
		switch item.name {
		case "bin":
			binPos = i
		case "Capitalised":
			capitalPos = i
		}
	}

	assert.NotEqual(t, -1, binPos, "bin file not found")
	assert.NotEqual(t, -1, capitalPos, "Capitalised.txt file not found")
	assert.Less(t, binPos, capitalPos)
}

func TestView(t *testing.T) {
	win := test.NewTempWindow(t, widget.NewLabel("Content"))

	dlg := NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		assert.NoError(t, err)
		assert.Nil(t, reader)
	}, win)

	dlg.SetTitleText("File Selection")
	dlg.SetConfirmText("Yes")
	dlg.SetDismissText("Dismiss")
	dlg.Show()

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*fyne.Container)
	toggleViewButton := ui.Objects[1].(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Button)
	panel := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[0].(*container.Scroll).Content.(*fyne.Container).Objects[0]

	// view should be a grid
	_, isGrid := panel.(*widget.GridWrap)
	assert.True(t, isGrid)
	// toggleViewButton should reflect to what it will do (change to a list view).
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.ListIcon().Name(), toggleViewButton.Icon.Name())

	// toggle view
	test.Tap(toggleViewButton)
	// reload files container
	panel = ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[0].(*container.Scroll).Content.(*fyne.Container).Objects[0]

	// view should be a list
	_, isList := panel.(*widget.List)
	assert.True(t, isList)
	// toggleViewButton should reflect to what it will do (change to a grid view).
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.GridIcon().Name(), toggleViewButton.Icon.Name())

	// toggle view
	test.Tap(toggleViewButton)
	// reload files container
	panel = ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[0].(*container.Scroll).Content.(*fyne.Container).Objects[0]

	// view should be a grid again
	_, isGrid = panel.(*widget.GridWrap)
	assert.True(t, isGrid)
	// toggleViewButton should reflect to what it will do (change to a list view).
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.ListIcon().Name(), toggleViewButton.Icon.Name())

	title := ui.Objects[1].(*fyne.Container).Objects[0].(*widget.Label)
	assert.Equal(t, "File Selection", title.Text)
	confirm := ui.Objects[2].(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Button)
	assert.Equal(t, "Yes", confirm.Text)
	dismiss := ui.Objects[2].(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Button)
	assert.Equal(t, "Dismiss", dismiss.Text)
}

func TestSetView(t *testing.T) {
	win := test.NewTempWindow(t, widget.NewLabel("Content"))

	fyne.CurrentApp().Preferences().SetInt(viewLayoutKey, int(defaultView))

	dlg := NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		assert.NoError(t, err)
		assert.Nil(t, reader)
	}, win)

	dlg.SetTitleText("File Selection")
	dlg.SetConfirmText("Yes")
	dlg.SetDismissText("Dismiss")

	// set view to list
	dlg.SetView(ListView)

	dlg.Show()

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*fyne.Container)
	toggleViewButton := ui.Objects[1].(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Button)
	panel := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[0].(*container.Scroll).Content.(*fyne.Container).Objects[0]

	// view should be a list
	_, isList := panel.(*widget.List)
	assert.True(t, isList)
	// toggleViewButton should reflect to what it will do (change to a grid view).
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.GridIcon(), toggleViewButton.Icon)

	title := ui.Objects[1].(*fyne.Container).Objects[0].(*widget.Label)
	assert.Equal(t, "File Selection", title.Text)
	confirm := ui.Objects[2].(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Button)
	assert.Equal(t, "Yes", confirm.Text)
	dismiss := ui.Objects[2].(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Button)
	assert.Equal(t, "Dismiss", dismiss.Text)

	// set view to grid on already opened dialog - should be updated automatically
	dlg.SetView(GridView)

	// view should be a grid again
	panel = ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[0].(*container.Scroll).Content.(*fyne.Container).Objects[0]
	_, isGrid := panel.(*widget.GridWrap)
	assert.True(t, isGrid)
	// toggleViewButton should reflect to what it will do (change to a list view).
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.ListIcon(), toggleViewButton.Icon)
}

func TestSetViewPreferences(t *testing.T) {
	win := test.NewTempWindow(t, widget.NewLabel("Content"))

	prefs := fyne.CurrentApp().Preferences()

	// set user-saved viewLayout to GridView
	prefs.SetInt(viewLayoutKey, int(GridView))

	dlg := NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		assert.NoError(t, err)
		assert.Nil(t, reader)
	}, win)

	// set default view to be ListView
	dlg.SetView(ListView)

	dlg.Show()

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*fyne.Container)
	toggleViewButton := ui.Objects[1].(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Button)
	panel := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[0].(*container.Scroll).Content.(*fyne.Container).Objects[0]

	// check that preference setting overrules configured default view
	_, isGrid := panel.(*widget.GridWrap)
	assert.True(t, isGrid)
	// toggleViewButton should reflect to what it will do (change to a list view).
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.ListIcon(), toggleViewButton.Icon)
}

func TestViewPreferences(t *testing.T) {
	win := test.NewTempWindow(t, widget.NewLabel("Content"))

	prefs := fyne.CurrentApp().Preferences()

	// set viewLayout to an invalid value to verify that this situation is handled properly
	prefs.SetInt(viewLayoutKey, -1)

	dlg := NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		assert.NoError(t, err)
		assert.Nil(t, reader)
	}, win)

	dlg.Show()

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*fyne.Container)
	toggleViewButton := ui.Objects[1].(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Button)

	// default viewLayout preference should be 'grid'
	view := ViewLayout(prefs.Int(viewLayoutKey))
	assert.Equal(t, GridView, view)

	// toggle view
	test.Tap(toggleViewButton)

	// viewLayout preference should be 'list'
	view = ViewLayout(prefs.Int(viewLayoutKey))
	assert.Equal(t, ListView, view)

	// toggle view
	test.Tap(toggleViewButton)

	// viewLayout preference should be 'grid' again
	view = ViewLayout(prefs.Int(viewLayoutKey))
	assert.Equal(t, GridView, view)
}

func TestFileFavorites(t *testing.T) {
	_ = test.NewApp()
	win := test.NewTempWindow(t, widget.NewLabel("Content"))

	dlg := NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		assert.NoError(t, err)
		assert.Nil(t, reader)
	}, win)

	dlg.Show()

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*fyne.Container)

	dlg.dialog.loadFavorites()
	favoriteLocations, _ := getFavoriteLocations()
	places := dlg.dialog.getPlaces()
	assert.Len(t, dlg.dialog.favorites, len(favoriteLocations)+len(places))

	favoritesList := ui.Objects[0].(*container.Split).Leading.(*widget.List)
	assert.Len(t, dlg.dialog.favorites, favoritesList.Length())

	for i := 0; i < favoritesList.Length(); i++ {
		favoritesList.Select(i)

		f := dlg.dialog.favorites[i]
		loc, ok := favoriteLocations[f.locName]
		if ok {
			// favoriteItem is Home, Documents, Downloads
			assert.Equal(t, loc.String(), dlg.dialog.dir.String())
		} else {
			// favoriteItem is (on windows) C:\, D:\, etc.
			assert.NotEqual(t, "Home", f.locName)
		}

		ok, err := storage.Exists(dlg.dialog.dir)
		assert.NoError(t, err)
		assert.True(t, ok)
	}

	dlg.Dismiss()
}

func TestSetFileNameBeforeShow(t *testing.T) {
	win := test.NewTempWindow(t, widget.NewLabel("Content"))
	dSave := NewFileSave(func(fyne.URIWriteCloser, error) {}, win)
	dSave.SetFileName("testfile.zip")
	dSave.Show()

	assert.Equal(t, "testfile.zip", dSave.dialog.fileName.(*widget.Entry).Text)

	// Should have no effect on FileOpen dialog
	dOpen := NewFileOpen(func(f fyne.URIReadCloser, e error) {}, win)
	dOpen.SetFileName("testfile.zip")
	dOpen.Show()

	assert.NotEqual(t, "testfile.zip", dOpen.dialog.fileName.(*widget.Label).Text)
}

func TestSetFileNameAfterShow(t *testing.T) {
	win := test.NewTempWindow(t, widget.NewLabel("Content"))
	dSave := NewFileSave(func(fyne.URIWriteCloser, error) {}, win)
	dSave.Show()
	dSave.SetFileName("testfile.zip")

	assert.Equal(t, "testfile.zip", dSave.dialog.fileName.(*widget.Entry).Text)

	// Should have no effect on FileOpen dialog
	dOpen := NewFileOpen(func(f fyne.URIReadCloser, e error) {}, win)
	dOpen.Show()
	dOpen.SetFileName("testfile.zip")

	assert.NotEqual(t, "testfile.zip", dOpen.dialog.fileName.(*widget.Label).Text)
}

func TestTapParent_GoesUpOne(t *testing.T) {
	win := test.NewTempWindow(t, widget.NewLabel("Content"))
	d := NewFileOpen(func(fyne.URIReadCloser, error) {}, win)
	home, _ := os.UserHomeDir()
	homeURI, _ := storage.ListerForURI(storage.NewFileURI(home))
	parentURI, _ := storage.Parent(homeURI)

	d.SetView(GridView)
	d.SetLocation(homeURI)
	d.Show()

	items := test.WidgetRenderer(d.dialog.files)
	item := items.Objects()[0].(*intWidget.Scroll).Content.(*fyne.Container).Objects[0]
	parent := test.WidgetRenderer(item.(fyne.Widget)).Objects()[1].(*fileDialogItem)
	test.Tap(parent)

	assert.Equal(t, d.dialog.dir.String(), parentURI.String())
}

func TestCreateNewFolderInDir(t *testing.T) {
	win := test.NewTempWindow(t, widget.NewLabel("Content"))

	folderDialog := NewFolderOpen(func(lu fyne.ListableURI, err error) {
		assert.NoError(t, err)
	}, win)
	folderDialog.SetConfirmText("Choose")
	folderDialog.SetDismissText("Cancel")
	folderDialog.Show()

	folderDialogPopup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(folderDialogPopup)
	assert.NotNil(t, folderDialogPopup)

	folderDialogUI := folderDialogPopup.Content.(*fyne.Container)

	createNewFolderButton := folderDialogUI.Objects[1].(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Button)
	assert.Equal(t, "", createNewFolderButton.Text)
	assert.Equal(t, theme.FolderNewIcon().Name(), createNewFolderButton.Icon.Name())

	// open folder name input dialog
	test.Tap(createNewFolderButton)

	inputPopup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(inputPopup)
	assert.NotNil(t, inputPopup)

	folderNameInputUI := inputPopup.Content.(*fyne.Container)

	folderNameInputTitle := folderNameInputUI.Objects[4].(*widget.Label)
	assert.Equal(t, "New Folder", folderNameInputTitle.Text)

	folderNameInputLabel := folderNameInputUI.Objects[2].(*widget.Form).Items[0].Text
	assert.Equal(t, "Name", folderNameInputLabel)

	folderNameInputEntry := folderNameInputUI.Objects[2].(*widget.Form).Items[0].Widget.(*widget.Entry)
	assert.Equal(t, "", folderNameInputEntry.Text)

	folderNameInputCancel := folderNameInputUI.Objects[3].(*fyne.Container).Objects[0].(*widget.Button)
	assert.Equal(t, "Cancel", folderNameInputCancel.Text)
	assert.Equal(t, theme.CancelIcon(), folderNameInputCancel.Icon)

	folderNameInputCreate := folderNameInputUI.Objects[3].(*fyne.Container).Objects[1].(*widget.Button)
	assert.Equal(t, theme.ConfirmIcon(), folderNameInputCreate.Icon)
}

func TestSetOnClosedBeforeShow(t *testing.T) {
	win := test.NewTempWindow(t, widget.NewLabel("Content"))
	d := NewFileSave(func(fyne.URIWriteCloser, error) {}, win)
	onClosedCalled := false
	d.SetOnClosed(func() { onClosedCalled = true })
	d.Show()
	d.Hide()
	assert.True(t, onClosedCalled)
}

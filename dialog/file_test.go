package dialog

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// comparePaths compares if two file paths point to the same thing, and calls
// t.Fatalf() if there is an error in performing the comparison.
//
// Returns true of both paths point to the same thing.
//
// You should use this if you need to compare file paths, since it explicitly
// normalizes the paths to a stable canonical form. It also nicely
// abstracts out the requisite error handling.
//
// You should only call this function on paths that you expect to be valid.
func comparePaths(t *testing.T, u1, u2 fyne.ListableURI) bool {
	p1 := u1.String()[len(u1.Scheme())+3:]
	p2 := u2.String()[len(u2.Scheme())+3:]

	a1, err := filepath.Abs(p1)
	if err != nil {
		t.Fatalf("Failed to normalize path '%s'", p1)
	}

	a2, err := filepath.Abs(p2)
	if err != nil {
		t.Fatalf("Failed to normalize path '%s'", p2)
	}

	return a1 == a2
}

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
	if !comparePaths(t, res, expect) {
		t.Errorf("Expected effectiveStartingDir() to be '%s', but it was '%s'",
			expect, res)
	}

	// this should always be equivalent to the preceding test
	dialog.startingLocation = nil
	res = dialog.effectiveStartingDir()
	expect = home
	if !comparePaths(t, res, expect) {
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
		t.Errorf("Should have failed to create lister for nonexistant file")
	}
	res = dialog.effectiveStartingDir()
	expect = home
	if res.String() != expect.String() {
		t.Errorf("Expected effectiveStartingDir() to be '%s', but it was '%s'",
			expect, res)
	}

}

func TestFileDialogResize(t *testing.T) {
	win := test.NewWindow(widget.NewLabel("Content"))
	win.Resize(fyne.NewSize(600, 400))
	file := NewFileOpen(func(file fyne.URIReadCloser, err error) {}, win)
	file.SetFilter(storage.NewExtensionFileFilter([]string{".png"}))

	//Mimic the fileopen dialog
	d := &fileDialog{file: file}
	open := widget.NewButton("open", func() {})
	ui := container.NewBorder(nil, nil, nil, open)
	pad := theme.Padding()
	itemMin := d.newFileItem(storage.NewFileURI("filename.txt"), false, false).MinSize()
	originalSize := ui.MinSize().Add(itemMin.AddWidthHeight(itemMin.Width+pad*6, pad*3))
	d.win = widget.NewModalPopUp(ui, file.parent.Canvas())
	d.win.Resize(originalSize)
	file.dialog = d

	//Test resize - normal size scenario
	size := fyne.NewSize(200, 180) //normal size to fit (600,400)
	file.Resize(size)
	expectedWidth := float32(200)
	assert.Equal(t, expectedWidth, file.dialog.win.Content.Size().Width+theme.Padding()*2)
	expectedHeight := float32(180)
	assert.Equal(t, expectedHeight, file.dialog.win.Content.Size().Height+theme.Padding()*2)
	//Test resize - normal size scenario again
	size = fyne.NewSize(300, 280) //normal size to fit (600,400)
	file.Resize(size)
	expectedWidth = 300
	assert.Equal(t, expectedWidth, file.dialog.win.Content.Size().Width+theme.Padding()*2)
	expectedHeight = 280
	assert.Equal(t, expectedHeight, file.dialog.win.Content.Size().Height+theme.Padding()*2)

	//Test resize - greater than max size scenario
	size = fyne.NewSize(800, 600)
	file.Resize(size)
	expectedWidth = 600                                          //since win width only 600
	assert.Equal(t, expectedWidth, file.dialog.win.Size().Width) //max, also work
	assert.Equal(t, expectedWidth, file.dialog.win.Content.Size().Width+theme.Padding()*2)
	expectedHeight = 400                                           //since win heigh only 400
	assert.Equal(t, expectedHeight, file.dialog.win.Size().Height) //max, also work
	assert.Equal(t, expectedHeight, file.dialog.win.Content.Size().Height+theme.Padding()*2)

	//Test again - extreme small size
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
	win := test.NewWindow(widget.NewLabel("Content"))
	d := NewFileOpen(func(file fyne.URIReadCloser, err error) {
		chosen = file
		openErr = err
	}, win)
	testDataPath, _ := filepath.Abs("testdata")
	testData := storage.NewFileURI(testDataPath)
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
	//header
	title := ui.Objects[1].(*fyne.Container).Objects[1].(*widget.Label)
	assert.Equal(t, "Open File", title.Text)
	//optionsbuttons
	createNewFolderButton := ui.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Button)
	assert.Equal(t, "", createNewFolderButton.Text)
	assert.Equal(t, theme.FolderNewIcon().Name(), createNewFolderButton.Icon.Name())
	toggleViewButton := ui.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Button)
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.ListIcon().Name(), toggleViewButton.Icon.Name())
	optionsButton := ui.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Button)
	assert.Equal(t, "", optionsButton.Text)
	assert.Equal(t, theme.SettingsIcon().Name(), optionsButton.Icon.Name())
	//footer
	nameLabel := ui.Objects[2].(*fyne.Container).Objects[1].(*container.Scroll).Content.(*widget.Label)
	buttons := ui.Objects[2].(*fyne.Container).Objects[0].(*fyne.Container)
	open := buttons.Objects[1].(*widget.Button)
	//body
	breadcrumb := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[0].(*container.Scroll).Content.(*fyne.Container).Objects[0].(*fyne.Container)
	assert.Greater(t, len(breadcrumb.Objects), 0)

	assert.Nil(t, err)
	components := strings.Split(testData.String()[7:], "/")
	if components[0] == "" {
		// Splitting a unix path will give a "" at the beginning, but we actually want the path bar to show "/".
		components[0] = "/"
	}
	if assert.Equal(t, len(components), len(breadcrumb.Objects)) {
		for i := range components {
			assert.Equal(t, components[i], breadcrumb.Objects[i].(*widget.Button).Text)
		}
	}

	files := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[1].(*container.Scroll).Content.(*fyne.Container).Objects[0].(*widget.GridWrap)
	objects := test.WidgetRenderer(files).Objects()[0].(*container.Scroll).Content.(*fyne.Container).Objects
	assert.Greater(t, len(objects), 0)

	fileName := test.WidgetRenderer(objects[0].(fyne.Widget)).Objects()[1].(*fileDialogItem).name
	assert.Equal(t, "(Parent)", fileName)
	assert.True(t, open.Disabled())

	var target *fileDialogItem
	id := 0
	for i, icon := range objects {
		item := test.WidgetRenderer(icon.(fyne.Widget)).Objects()[1].(*fileDialogItem)
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
	assert.Nil(t, openErr)

	assert.Equal(t, target.location.String(), chosen.URI().String())

	err = chosen.Close()
	assert.Nil(t, err)
}

func TestHiddenFiles(t *testing.T) {
	testDataPath, _ := filepath.Abs("testdata")
	testData := storage.NewFileURI(testDataPath)
	dir, err := storage.ListerForURI(testData)
	if err != nil {
		t.Error("Failed to open testdata dir", err)
	}

	// git does not preserve windows hidden flag so we have to set it.
	// just an empty function for non windows builds
	if err := hideFile(filepath.Join(testDataPath, ".hidden")); err != nil {
		t.Error("Failed to hide .hidden", err)
	}

	win := test.NewWindow(widget.NewLabel("Content"))
	d := NewFileOpen(func(file fyne.URIReadCloser, err error) {
	}, win)
	d.SetLocation(dir)
	d.Show()

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*fyne.Container)

	createNewFolderButton := ui.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Button)
	assert.Equal(t, "", createNewFolderButton.Text)
	assert.Equal(t, theme.FolderNewIcon().Name(), createNewFolderButton.Icon.Name())
	toggleViewButton := ui.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Button)
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.ListIcon().Name(), toggleViewButton.Icon.Name())
	optionsButton := ui.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Button)
	assert.Equal(t, "", optionsButton.Text)
	assert.Equal(t, theme.SettingsIcon().Name(), optionsButton.Icon.Name())

	files := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[1].(*container.Scroll).Content.(*fyne.Container).Objects[0].(*widget.GridWrap)
	objects := test.WidgetRenderer(files).Objects()[0].(*container.Scroll).Content.(*fyne.Container).Objects
	assert.Greater(t, len(objects), 0)

	var target *fileDialogItem
	for _, icon := range objects {
		item := test.WidgetRenderer(icon.(fyne.Widget)).Objects()[1].(*fileDialogItem)
		if item.name == ".hidden" {
			target = item
		}
	}
	assert.Nil(t, target, "Failed, .hidden found in testdata")

	d.dialog.showHidden = true
	d.dialog.refreshDir(d.dialog.dir)

	for _, icon := range objects {
		item := test.WidgetRenderer(icon.(fyne.Widget)).Objects()[1].(*fileDialogItem)
		if item.name == ".hidden" {
			target = item
		}
	}
	assert.NotNil(t, target, "Failed,.hidden not found in testdata")
}

func TestShowFileSave(t *testing.T) {
	var chosen fyne.URIWriteCloser
	var saveErr error
	win := test.NewWindow(widget.NewLabel("Content"))
	saver := NewFileSave(func(file fyne.URIWriteCloser, err error) {
		chosen = file
		saveErr = err
	}, win)
	saver.Show()

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*fyne.Container)
	title := ui.Objects[1].(*fyne.Container).Objects[1].(*widget.Label)
	assert.Equal(t, "Save File", title.Text)

	nameEntry := ui.Objects[2].(*fyne.Container).Objects[1].(*container.Scroll).Content.(*widget.Entry)
	buttons := ui.Objects[2].(*fyne.Container).Objects[0].(*fyne.Container)
	save := buttons.Objects[1].(*widget.Button)

	files := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[1].(*container.Scroll).Content.(*fyne.Container).Objects[0].(*widget.GridWrap)
	objects := test.WidgetRenderer(files).Objects()[0].(*container.Scroll).Content.(*fyne.Container).Objects
	assert.Greater(t, len(objects), 0)

	item := test.WidgetRenderer(objects[0].(fyne.Widget)).Objects()[1].(*fileDialogItem)
	assert.Equal(t, "(Parent)", item.name)
	assert.True(t, save.Disabled())

	var target *fileDialogItem
	id := -1
	for i, icon := range objects {
		item := test.WidgetRenderer(icon.(fyne.Widget)).Objects()[1].(*fileDialogItem)
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

	abs, _ := filepath.Abs("./testdata/")
	dir, _ := storage.ListerForURI(storage.NewFileURI(abs))
	saver.SetLocation(dir)
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
	assert.Nil(t, saveErr)
	targetParent, err := storage.Parent(target.location)
	if err != nil {
		t.Error(err)
	}
	expectedPath, _ := storage.Child(targetParent, "v2_"+target.location.Name())
	assert.Equal(t, expectedPath.String(), chosen.URI().String())

	err = chosen.Close()
	assert.Nil(t, err)
	pathString := expectedPath.String()[len(expectedPath.Scheme())+3:]
	err = os.Remove(pathString)
	assert.Nil(t, err)
}

func TestFileFilters(t *testing.T) {
	win := test.NewWindow(widget.NewLabel("Content"))
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
			assert.Equal(t, uri.Extension(), ".png")
			count++
		}
	}

	// NOTE: This count needs to be updated when more test images are added.
	assert.Equal(t, 8, count)

	f.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/jpeg"}))

	count = 0
	for _, uri := range f.dialog.data {
		ok, _ := storage.CanList(uri)
		if !ok {
			assert.Equal(t, uri.MimeType(), "image/jpeg")
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
			assert.Equal(t, mimeType, "image")
			count++
		}
	}

	// NOTE: This count needs to be updated when more test images are added.
	assert.Equal(t, 9, count)
}

func TestView(t *testing.T) {
	win := test.NewWindow(widget.NewLabel("Content"))

	dlg := NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		assert.Nil(t, err)
		assert.Nil(t, reader)
	}, win)

	dlg.SetConfirmText("Yes")
	dlg.SetDismissText("Dismiss")
	dlg.Show()

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*fyne.Container)
	toggleViewButton := ui.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Button)
	panel := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[1].(*container.Scroll).Content.(*fyne.Container).Objects[0]

	// view should be a grid
	_, isGrid := panel.(*widget.GridWrap)
	assert.True(t, isGrid)
	// toggleViewButton should reflect to what it will do (change to a list view).
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.ListIcon().Name(), toggleViewButton.Icon.Name())

	// toggle view
	test.Tap(toggleViewButton)
	// reload files container
	panel = ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[1].(*container.Scroll).Content.(*fyne.Container).Objects[0]

	// view should be a list
	_, isList := panel.(*widget.List)
	assert.True(t, isList)
	// toggleViewButton should reflect to what it will do (change to a grid view).
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.GridIcon().Name(), toggleViewButton.Icon.Name())

	// toggle view
	test.Tap(toggleViewButton)
	// reload files container
	panel = ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[1].(*container.Scroll).Content.(*fyne.Container).Objects[0]

	// view should be a grid again
	_, isGrid = panel.(*widget.GridWrap)
	assert.True(t, isGrid)
	// toggleViewButton should reflect to what it will do (change to a list view).
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.ListIcon().Name(), toggleViewButton.Icon.Name())

	confirm := ui.Objects[2].(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Button)
	assert.Equal(t, "Yes", confirm.Text)
	dismiss := ui.Objects[2].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Button)
	assert.Equal(t, "Dismiss", dismiss.Text)
}

func TestViewPreferences(t *testing.T) {
	win := test.NewWindow(widget.NewLabel("Content"))

	prefs := fyne.CurrentApp().Preferences()

	// set viewLayout to an invalid value to verify that this situation is handled properly
	prefs.SetInt(viewLayoutKey, -1)

	dlg := NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		assert.Nil(t, err)
		assert.Nil(t, reader)
	}, win)

	dlg.Show()

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*fyne.Container)
	toggleViewButton := ui.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Button)

	// viewLayout preference should be 'grid'
	view := viewLayout(prefs.Int(viewLayoutKey))
	assert.Equal(t, gridView, view)

	// toggle view
	test.Tap(toggleViewButton)

	// viewLayout preference should be 'list'
	view = viewLayout(prefs.Int(viewLayoutKey))
	assert.Equal(t, listView, view)

	// toggle view
	test.Tap(toggleViewButton)

	// viewLayout preference should be 'grid' again
	view = viewLayout(prefs.Int(viewLayoutKey))
	assert.Equal(t, gridView, view)
}

func TestFileFavorites(t *testing.T) {
	win := test.NewWindow(widget.NewLabel("Content"))

	dlg := NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		assert.Nil(t, err)
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
	assert.Equal(t, favoritesList.Length(), len(dlg.dialog.favorites))

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
		assert.Nil(t, err)
		assert.True(t, ok)
	}

	test.Tap(dlg.dialog.dismiss)
}

func TestSetFileNameBeforeShow(t *testing.T) {
	win := test.NewWindow(widget.NewLabel("Content"))
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

	win := test.NewWindow(widget.NewLabel("Content"))
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

func TestCreateNewFolderInDir(t *testing.T) {
	win := test.NewWindow(widget.NewLabel("Content"))

	folderDialog := NewFolderOpen(func(lu fyne.ListableURI, err error) {
		assert.Nil(t, err)
	}, win)
	folderDialog.SetConfirmText("Choose")
	folderDialog.SetDismissText("Cancel")
	folderDialog.Show()

	folderDialogPopup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(folderDialogPopup)
	assert.NotNil(t, folderDialogPopup)

	folderDialogUI := folderDialogPopup.Content.(*fyne.Container)

	createNewFolderButton := folderDialogUI.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Button)
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
	win := test.NewWindow(widget.NewLabel("Content"))
	d := NewFileSave(func(fyne.URIWriteCloser, error) {}, win)
	onClosedCalled := false
	d.SetOnClosed(func() { onClosedCalled = true })
	d.Show()
	d.Hide()
	assert.True(t, onClosedCalled)
}

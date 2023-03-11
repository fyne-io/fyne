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
	"fyne.io/fyne/v2/layout"
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
	t.Log(parentURI)
	t.Log(parent)
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
	originalSize := ui.MinSize().Add(fyne.NewSize(fileIconCellWidth*2+theme.Padding()*4,
		(fileIconSize+fileTextSize)+theme.Padding()*4))
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
	assert.Equal(t, theme.FolderNewIcon(), createNewFolderButton.Icon)
	toggleViewButton := ui.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Button)
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.ListIcon(), toggleViewButton.Icon)
	optionsButton := ui.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Button)
	assert.Equal(t, "", optionsButton.Text)
	assert.Equal(t, theme.SettingsIcon(), optionsButton.Icon)
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

	files := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[1].(*container.Scroll).Content.(*fyne.Container).Objects[0].(*fyne.Container)
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
	assert.NotNil(t, target, "Failed to find file in testdata")
	test.Tap(target)
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
	assert.Equal(t, theme.FolderNewIcon(), createNewFolderButton.Icon)
	toggleViewButton := ui.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Button)
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.ListIcon(), toggleViewButton.Icon)
	optionsButton := ui.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Button)
	assert.Equal(t, "", optionsButton.Text)
	assert.Equal(t, theme.SettingsIcon(), optionsButton.Icon)

	files := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[1].(*container.Scroll).Content.(*fyne.Container).Objects[0].(*fyne.Container)
	assert.Greater(t, len(files.Objects), 0)

	var target *fileDialogItem
	for _, icon := range files.Objects {
		if icon.(*fileDialogItem).name == ".hidden" {
			target = icon.(*fileDialogItem)
		}
	}
	assert.Nil(t, target, "Failed, .hidden found in testdata")

	d.dialog.showHidden = true
	d.dialog.refreshDir(d.dialog.dir)

	for _, icon := range files.Objects {
		if icon.(*fileDialogItem).name == ".hidden" {
			target = icon.(*fileDialogItem)
		}
	}
	assert.NotNil(t, target, "Failed,.hidden not found in testdata")
}

func TestShowFileSave(t *testing.T) {
	var chosen fyne.URIWriteCloser
	var saveErr error
	win := test.NewWindow(widget.NewLabel("Content"))
	ShowFileSave(func(file fyne.URIWriteCloser, err error) {
		chosen = file
		saveErr = err
	}, win)

	popup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(popup)
	assert.NotNil(t, popup)

	ui := popup.Content.(*fyne.Container)
	title := ui.Objects[1].(*fyne.Container).Objects[1].(*widget.Label)
	assert.Equal(t, "Save File", title.Text)

	nameEntry := ui.Objects[2].(*fyne.Container).Objects[1].(*container.Scroll).Content.(*widget.Entry)
	buttons := ui.Objects[2].(*fyne.Container).Objects[0].(*fyne.Container)
	save := buttons.Objects[1].(*widget.Button)

	files := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[1].(*container.Scroll).Content.(*fyne.Container).Objects[0].(*fyne.Container)
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
	for _, icon := range f.dialog.files.Objects {
		if icon.(*fileDialogItem).dir == false {
			uri := icon.(*fileDialogItem).location
			assert.Equal(t, uri.Extension(), ".png")
			count++
		}
	}
	assert.Equal(t, 5, count)

	f.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/jpeg"}))

	count = 0
	for _, icon := range f.dialog.files.Objects {
		if icon.(*fileDialogItem).dir == false {
			uri := icon.(*fileDialogItem).location
			assert.Equal(t, uri.MimeType(), "image/jpeg")
			count++
		}
	}
	assert.Equal(t, 1, count)

	f.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/*"}))

	count = 0
	for _, icon := range f.dialog.files.Objects {
		if icon.(*fileDialogItem).dir == false {
			uri := icon.(*fileDialogItem).location
			mimeType := strings.Split(uri.MimeType(), "/")[0]
			assert.Equal(t, mimeType, "image")
			count++
		}
	}
	assert.Equal(t, 6, count)
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
	files := ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[1].(*container.Scroll).Content.(*fyne.Container).Objects[0].(*fyne.Container)

	listLayout := layout.NewVBoxLayout()

	// view should be a grid
	assert.NotEqual(t, listLayout, files.Layout)
	// toggleViewButton should reflect to what it will do (change to a list view).
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.ListIcon(), toggleViewButton.Icon)

	// toggle view
	test.Tap(toggleViewButton)
	// reload files container
	files = ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[1].(*container.Scroll).Content.(*fyne.Container).Objects[0].(*fyne.Container)

	// view should be a list
	assert.Equal(t, listLayout, files.Layout)
	// toggleViewButton should reflect to what it will do (change to a grid view).
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.GridIcon(), toggleViewButton.Icon)

	// toggle view
	test.Tap(toggleViewButton)
	// reload files container
	files = ui.Objects[0].(*container.Split).Trailing.(*fyne.Container).Objects[1].(*container.Scroll).Content.(*fyne.Container).Objects[0].(*fyne.Container)

	// view should be a grid again
	assert.NotEqual(t, listLayout, files.Layout)
	// toggleViewButton should reflect to what it will do (change to a list view).
	assert.Equal(t, "", toggleViewButton.Text)
	assert.Equal(t, theme.ListIcon(), toggleViewButton.Icon)

	confirm := ui.Objects[2].(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Button)
	assert.Equal(t, "Yes", confirm.Text)
	dismiss := ui.Objects[2].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Button)
	assert.Equal(t, "Dismiss", dismiss.Text)
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

	folderDialogUi := folderDialogPopup.Content.(*fyne.Container)

	createNewFolderButton := folderDialogUi.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Button)
	assert.Equal(t, "", createNewFolderButton.Text)
	assert.Equal(t, theme.FolderNewIcon(), createNewFolderButton.Icon)

	// open folder name input dialog
	test.Tap(createNewFolderButton)

	inputPopup := win.Canvas().Overlays().Top().(*widget.PopUp)
	defer win.Canvas().Overlays().Remove(inputPopup)
	assert.NotNil(t, inputPopup)

	folderNameInputUi := inputPopup.Content.(*fyne.Container)

	folderNameInputTitle := folderNameInputUi.Objects[4].(*widget.Label)
	assert.Equal(t, "New Folder", folderNameInputTitle.Text)

	folderNameInputLabel := folderNameInputUi.Objects[2].(*widget.Form).Items[0].Text
	assert.Equal(t, "Name", folderNameInputLabel)

	folderNameInputEntry := folderNameInputUi.Objects[2].(*widget.Form).Items[0].Widget.(*widget.Entry)
	assert.Equal(t, "", folderNameInputEntry.Text)

	folderNameInputCancel := folderNameInputUi.Objects[3].(*fyne.Container).Objects[1].(*widget.Button)
	assert.Equal(t, "Cancel", folderNameInputCancel.Text)
	assert.Equal(t, theme.CancelIcon(), folderNameInputCancel.Icon)

	folderNameInputCreate := folderNameInputUi.Objects[3].(*fyne.Container).Objects[2].(*widget.Button)
	assert.Equal(t, theme.ConfirmIcon(), folderNameInputCreate.Icon)
}

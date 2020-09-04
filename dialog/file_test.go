package dialog

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
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
func comparePaths(t *testing.T, p1, p2 string) bool {
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
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("os.UserHomeDir) failed, cannot run this test on this system, error was '%s'", err)
	}

	parent := filepath.Dir(wd)
	_, err = os.Stat(parent)
	if err != nil {
		t.Skipf("os.Getwd() failed, cannot run this test on this system (error stat()-ing ../) error was '%s'", err)
	}

	dialog := &FileDialog{}

	// test that we get $HOME when running with the default struct values
	res := dialog.effectiveStartingDir()
	expect := home
	if !comparePaths(t, res, expect) {
		t.Errorf("Expected effectiveStartingDir() to be '%s', but it was '%s'",
			expect, res)
	}

	// this should always be equivalent to the preceding test
	dialog.StartingDirectory = ""
	res = dialog.effectiveStartingDir()
	expect = wd
	if err != nil {
		t.Skipf("os.Getwd() failed, cannot run this test on this system, error was '%s'", err)
	}
	if !comparePaths(t, res, expect) {
		t.Errorf("Expected effectiveStartingDir() to be '%s', but it was '%s'",
			expect, res)
	}

	// check using StartingDirectory with some other directory
	dialog.StartingDirectory = parent
	res = dialog.effectiveStartingDir()
	expect = parent
	if err != nil {
		t.Skipf("os.Getwd() failed, cannot run this test on this system, error was '%s'", err)
	}
	if res != expect {
		t.Errorf("Expected effectiveStartingDir() to be '%s', but it was '%s'",
			expect, res)
	}

	// make sure we fail over if the specified directory does not exist
	dialog.StartingDirectory = "/some/file/that/does/not/exist"
	res = dialog.effectiveStartingDir()
	expect = wd
	if err != nil {
		t.Skipf("os.Getwd() failed, cannot run this test on this system, error was '%s'", err)
	}
	if res != expect {
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
	ui := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, open), open)
	originalSize := ui.MinSize().Add(fyne.NewSize(fileIconCellWidth*2+theme.Padding()*4,
		(fileIconSize+fileTextSize)+theme.Padding()*4))
	d.win = widget.NewModalPopUp(ui, file.parent.Canvas())
	d.win.Resize(originalSize)
	file.dialog = d

	//Test resize - normal size scenario
	size := fyne.NewSize(200, 180) //normal size to fit (600,400)
	file.Resize(size)
	expectedWidth := 200
	assert.Equal(t, expectedWidth, file.dialog.win.Content.Size().Width+theme.Padding()*2)
	expectedHeight := 180
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
	ShowFileOpen(func(file fyne.URIReadCloser, err error) {
		chosen = file
		openErr = err
	}, win)

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
	testDataDir := filepath.Join(workingDir, "testdata")

	f.dialog.setDirectory(testDataDir)

	count := 0
	for _, icon := range f.dialog.files.Objects {
		if icon.(*fileDialogItem).dir == false {
			uri := storage.NewURI("file://" + icon.(*fileDialogItem).path)
			assert.Equal(t, uri.Extension(), ".png")
			count++
		}
	}
	assert.Equal(t, 3, count)

	f.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/jpeg"}))

	count = 0
	for _, icon := range f.dialog.files.Objects {
		if icon.(*fileDialogItem).dir == false {
			uri := storage.NewURI("file://" + icon.(*fileDialogItem).path)
			assert.Equal(t, uri.MimeType(), "image/jpeg")
			count++
		}
	}
	assert.Equal(t, 1, count)

	f.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/*"}))

	count = 0
	for _, icon := range f.dialog.files.Objects {
		if icon.(*fileDialogItem).dir == false {
			uri := storage.NewURI("file://" + icon.(*fileDialogItem).path)
			mimeType := strings.Split(uri.MimeType(), "/")[0]
			assert.Equal(t, mimeType, "image")
			count++
		}
	}
	assert.Equal(t, 4, count)
}

package dialog

import (
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2/storage"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestFileItem_Name(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()

	item := f.newFileItem(storage.NewFileURI("/path/to/filename.txt"), false)
	assert.Equal(t, "filename", item.name)

	item = f.newFileItem(storage.NewFileURI("/path/to/MyFile.jpeg"), false)
	assert.Equal(t, "MyFile", item.name)

	item = f.newFileItem(storage.NewFileURI("/path/to/.maybeHidden.txt"), false)
	assert.Equal(t, ".maybeHidden", item.name)

	item = f.newFileItem(storage.NewFileURI("/path/to/noext"), false)
	assert.Equal(t, "noext", item.name)

	// Test that the extension remains for the list view.
	f.view = listView

	item = f.newFileItem(storage.NewFileURI("/path/to/filename.txt"), false)
	assert.Equal(t, "filename.txt", item.name)

	item = f.newFileItem(storage.NewFileURI("/path/to/MyFile.jpeg"), false)
	assert.Equal(t, "MyFile.jpeg", item.name)

	item = f.newFileItem(storage.NewFileURI("/path/to/.maybeHidden.txt"), false)
	assert.Equal(t, ".maybeHidden.txt", item.name)

	item = f.newFileItem(storage.NewFileURI("/path/to/noext"), false)
	assert.Equal(t, "noext", item.name)
}

func TestFileItem_FolderName(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()

	item := f.newFileItem(storage.NewFileURI("/path/to/foldername/"), true)
	assert.Equal(t, "foldername", item.name)

	item = f.newFileItem(storage.NewFileURI("/path/to/myapp.app/"), true)
	assert.Equal(t, "myapp", item.name)

	item = f.newFileItem(storage.NewFileURI("/path/to/.maybeHidden/"), true)
	assert.Equal(t, ".maybeHidden", item.name)

	// Test that the extension remains for the list view.
	f.view = listView
	item = f.newFileItem(storage.NewFileURI("/path/to/myapp.app/"), true)
	assert.Equal(t, "myapp.app", item.name)
}

func TestNewFileItem(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()
	item := f.newFileItem(storage.NewFileURI("/path/to/filename.txt"), false)

	assert.Equal(t, "filename", item.name)

	test.Tap(item)
	assert.True(t, item.isCurrent)
	assert.Equal(t, item, f.selected)
}

func TestNewFileItem_Folder(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()
	currentDir, _ := filepath.Abs(".")
	currentLister, err := storage.ListerForURI(storage.NewFileURI(currentDir))
	if err != nil {
		t.Error(err)
	}

	parentDir := storage.NewFileURI(filepath.Dir(currentDir))
	parentLister, err := storage.ListerForURI(parentDir)
	if err != nil {
		t.Error(err)
	}
	f.setLocation(parentLister)
	item := f.newFileItem(currentLister, true)

	assert.Equal(t, filepath.Base(currentDir), item.name)

	test.Tap(item)
	assert.False(t, item.isCurrent)
	assert.Equal(t, (*fileDialogItem)(nil), f.selected)
	assert.Equal(t, storage.NewFileURI(""+currentDir).String(), f.dir.String())
}

func TestNewFileItem_ParentFolder(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()
	currentDir, _ := filepath.Abs(".")
	currentLister, err := storage.ListerForURI(storage.NewFileURI(currentDir))
	if err != nil {
		t.Error(err)
	}
	parentDir := storage.NewFileURI(filepath.Dir(currentDir))
	f.setLocation(currentLister)

	item := &fileDialogItem{picker: f, name: "(Parent)", location: parentDir, dir: true}
	item.ExtendBaseWidget(item)

	assert.Equal(t, "(Parent)", item.name)

	test.Tap(item)
	assert.False(t, item.isCurrent)
	assert.Equal(t, (*fileDialogItem)(nil), f.selected)
	assert.Equal(t, parentDir.String(), f.dir.String())
}

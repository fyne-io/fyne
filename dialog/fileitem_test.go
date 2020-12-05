package dialog

import (
	"path/filepath"
	"testing"

	"fyne.io/fyne/storage"

	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestFileItem_Name(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()

	item := f.newFileItem(storage.NewURI("file:///path/to/filename.txt"), false)
	assert.Equal(t, "filename", item.name)

	item = f.newFileItem(storage.NewURI("file:///path/to/MyFile.jpeg"), false)
	assert.Equal(t, "MyFile", item.name)

	item = f.newFileItem(storage.NewURI("file:///path/to/.maybeHidden.txt"), false)
	assert.Equal(t, ".maybeHidden", item.name)
}

func TestFileItem_FolderName(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()

	item := f.newFileItem(storage.NewURI("file:///path/to/foldername/"), true)
	assert.Equal(t, "foldername", item.name)

	item = f.newFileItem(storage.NewURI("file:///path/to/myapp.app/"), true)
	assert.Equal(t, "myapp.app", item.name)

	item = f.newFileItem(storage.NewURI("file:///path/to/.maybeHidden/"), true)
	assert.Equal(t, ".maybeHidden", item.name)
}

func TestNewFileItem(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()
	item := f.newFileItem(storage.NewURI("file:///path/to/filename.txt"), false)

	assert.Equal(t, "filename", item.name)

	test.Tap(item)
	assert.True(t, item.isCurrent)
	assert.Equal(t, item, f.selected)
}

func TestNewFileItem_Folder(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()
	currentDir, _ := filepath.Abs(".")
	currentLister, err := storage.ListerForURI(storage.NewURI("file://" + currentDir))
	if err != nil {
		t.Error(err)
	}

	parentDir := storage.NewURI("file://" + filepath.Dir(currentDir))
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
	assert.Equal(t, storage.NewURI("file://"+currentDir).String(), f.dir.String())
}

func TestNewFileItem_ParentFolder(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()
	currentDir, _ := filepath.Abs(".")
	currentLister, err := storage.ListerForURI(storage.NewURI("file://" + currentDir))
	if err != nil {
		t.Error(err)
	}
	parentDir := storage.NewURI("file://" + filepath.Dir(currentDir))
	f.setLocation(currentLister)

	item := &fileDialogItem{picker: f, name: "(Parent)", location: parentDir, dir: true}
	item.ExtendBaseWidget(item)

	assert.Equal(t, "(Parent)", item.name)

	test.Tap(item)
	assert.False(t, item.isCurrent)
	assert.Equal(t, (*fileDialogItem)(nil), f.selected)
	assert.Equal(t, parentDir.String(), f.dir.String())
}

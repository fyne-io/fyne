package dialog

import (
	"path/filepath"
	"testing"

	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func TestFileItem_Name(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()

	item := f.newFileItem("/path/to/filename.txt", false)
	assert.Equal(t, "filename", item.name)

	item = f.newFileItem("/path/to/MyFile.jpeg", false)
	assert.Equal(t, "MyFile", item.name)

	item = f.newFileItem("/path/to/.maybeHidden.txt", false)
	assert.Equal(t, ".maybeHidden", item.name)
}

func TestFileItem_FolderName(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()

	item := f.newFileItem("/path/to/foldername/", true)
	assert.Equal(t, "foldername", item.name)

	item = f.newFileItem("/path/to/myapp.app/", true)
	assert.Equal(t, "myapp.app", item.name)

	item = f.newFileItem("/path/to/.maybeHidden/", true)
	assert.Equal(t, ".maybeHidden", item.name)
}

func TestNewFileItem(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()
	item := f.newFileItem("/path/to/filename.txt", false)

	assert.Equal(t, "filename", item.name)

	test.Tap(item)
	assert.True(t, item.isCurrent)
	assert.Equal(t, item, f.selected)
}

func TestNewFileItem_Folder(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()
	currentDir, _ := filepath.Abs(".")
	parentDir := filepath.Dir(currentDir)
	f.setDirectory(parentDir)
	item := f.newFileItem(currentDir, true)

	assert.Equal(t, filepath.Base(currentDir), item.name)

	test.Tap(item)
	assert.False(t, item.isCurrent)
	assert.Equal(t, (*fileDialogItem)(nil), f.selected)
	assert.Equal(t, currentDir, f.dir)
}

func TestNewFileItem_ParentFolder(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()
	currentDir, _ := filepath.Abs(".")
	parentDir := filepath.Dir(currentDir)
	f.setDirectory(currentDir)

	item := &fileDialogItem{picker: f, icon: canvas.NewImageFromResource(theme.FolderOpenIcon()),
		name: "(Parent)", path: parentDir, dir: true}
	item.ExtendBaseWidget(item)

	assert.Equal(t, "(Parent)", item.name)

	test.Tap(item)
	assert.False(t, item.isCurrent)
	assert.Equal(t, (*fileDialogItem)(nil), f.selected)
	assert.Equal(t, parentDir, f.dir)
}

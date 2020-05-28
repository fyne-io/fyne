package dialog

import (
	"path/filepath"
	"testing"

	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func TestNewFileItem(t *testing.T) {
	f := &fileDialog{file: &FileDialog{}}
	_ = f.makeUI()
	item := f.newFileItem("/path/to/filename.txt", false)

	assert.Equal(t, item.name, "filename")

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

	assert.Equal(t, item.name, filepath.Base(currentDir))

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

	assert.Equal(t, item.name, "(Parent)")

	test.Tap(item)
	assert.False(t, item.isCurrent)
	assert.Equal(t, (*fileDialogItem)(nil), f.selected)
	assert.Equal(t, parentDir, f.dir)
}

package dialog

import (
	"path/filepath"
	"testing"

	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func TestNewFileItem(t *testing.T) {
	f := &fileDialog{}
	_ = f.makeUI()
	item := f.newFileItem(theme.FileIcon(), "/path/to/filename.txt")

	assert.Equal(t, item.name, "filename")
	assert.Equal(t, item.ext, "txt")

	test.Tap(item)
	assert.True(t, item.isCurrent)
	assert.Equal(t, item, f.selected)
}

func TestNewFileItem_Folder(t *testing.T) {
	f := &fileDialog{}
	_ = f.makeUI()
	currentDir, _ := filepath.Abs(".")
	parentDir := filepath.Dir(currentDir)
	f.setDirectory(parentDir)
	item := f.newFileItem(theme.FolderIcon(), currentDir)

	assert.Equal(t, item.name, filepath.Base(currentDir))
	assert.Equal(t, item.ext, "")

	test.Tap(item)
	assert.False(t, item.isCurrent)
	assert.Equal(t, (*fileDialogItem)(nil), f.selected)
	assert.Equal(t, currentDir, f.dir)
}

func TestNewFileItem_ParentFolder(t *testing.T) {
	f := &fileDialog{}
	_ = f.makeUI()
	currentDir, _ := filepath.Abs(".")
	parentDir := filepath.Dir(currentDir)
	f.setDirectory(currentDir)
	item := f.newFileItem(theme.FolderOpenIcon(), parentDir)

	assert.Equal(t, item.name, "(Parent)")
	assert.Equal(t, item.ext, "")

	test.Tap(item)
	assert.False(t, item.isCurrent)
	assert.Equal(t, (*fileDialogItem)(nil), f.selected)
	assert.Equal(t, parentDir, f.dir)
}

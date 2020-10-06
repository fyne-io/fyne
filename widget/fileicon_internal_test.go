// +build !ci,!windows

package widget

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
)

// Simulate being rendered by calling CreateRenderer() to update icon
func newRenderedFileIcon(uri fyne.URI) *FileIcon {
	f := NewFileIcon(uri)
	f.CreateRenderer()
	return f
}

func TestNewFileIcon(t *testing.T) {
	item := newRenderedFileIcon(storage.NewURI("file:///path/to/filename.zip"))
	assert.Equal(t, ".zip", item.extension)
	assert.Equal(t, theme.FileApplicationIcon(), item.resource)

	item = newRenderedFileIcon(storage.NewURI("file:///path/to/filename.mp3"))
	assert.Equal(t, ".mp3", item.extension)
	assert.Equal(t, theme.FileAudioIcon(), item.resource)

	item = newRenderedFileIcon(storage.NewURI("file:///path/to/filename.png"))
	assert.Equal(t, ".png", item.extension)
	assert.Equal(t, theme.FileImageIcon(), item.resource)

	item = newRenderedFileIcon(storage.NewURI("file:///path/to/filename.txt"))
	assert.Equal(t, ".txt", item.extension)
	assert.Equal(t, theme.FileTextIcon(), item.resource)

	item = newRenderedFileIcon(storage.NewURI("file:///path/to/filename.mp4"))
	assert.Equal(t, ".mp4", item.extension)
	assert.Equal(t, theme.FileVideoIcon(), item.resource)
}

func TestNewFileIconNoExtension(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		fyne.LogError("Could not get current working directory", err)
		t.FailNow()
	}
	binFileWithNoExt := filepath.Join(workingDir, "testdata/bin")
	textFileWithNoExt := filepath.Join(workingDir, "testdata/text")

	item := newRenderedFileIcon(storage.NewURI("file://" + binFileWithNoExt))
	assert.Equal(t, "", item.extension)
	assert.Equal(t, theme.FileApplicationIcon(), item.resource)

	item = newRenderedFileIcon(storage.NewURI("file://" + textFileWithNoExt))
	assert.Equal(t, "", item.extension)
	assert.Equal(t, theme.FileTextIcon(), item.resource)
}

func TestSetURI(t *testing.T) {
	item := newRenderedFileIcon(storage.NewURI("file:///path/to/filename.zip"))
	assert.Equal(t, ".zip", item.extension)
	assert.Equal(t, theme.FileApplicationIcon(), item.resource)

	item.SetURI(storage.NewURI("file:///path/to/filename.mp3"))
	assert.Equal(t, ".mp3", item.extension)
	assert.Equal(t, theme.FileAudioIcon(), item.resource)

	item.SetURI(storage.NewURI("file:///path/to/filename.png"))
	assert.Equal(t, ".png", item.extension)
	assert.Equal(t, theme.FileImageIcon(), item.resource)

	item.SetURI(storage.NewURI("file:///path/to/filename.txt"))
	assert.Equal(t, ".txt", item.extension)
	assert.Equal(t, theme.FileTextIcon(), item.resource)

	item.SetURI(storage.NewURI("file:///path/to/filename.mp4"))
	assert.Equal(t, ".mp4", item.extension)
	assert.Equal(t, theme.FileVideoIcon(), item.resource)
}

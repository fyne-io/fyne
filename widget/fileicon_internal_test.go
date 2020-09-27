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

func TestNewFileIcon(t *testing.T) {
	item := NewFileIcon(storage.NewURI("file:///path/to/filename.zip"))

	assert.Equal(t, ".zip", item.extension)
	assert.Equal(t, "application", item.mimeType)
	assert.Equal(t, "zip", item.mimeSubType)
	assert.Equal(t, theme.FileApplicationIcon(), item.resource)

	item = NewFileIcon(storage.NewURI("file:///path/to/filename.mp3"))

	assert.Equal(t, ".mp3", item.extension)
	assert.Equal(t, "audio", item.mimeType)
	assert.Equal(t, "mpeg", item.mimeSubType)
	assert.Equal(t, theme.FileAudioIcon(), item.resource)

	item = NewFileIcon(storage.NewURI("file:///path/to/filename.png"))

	assert.Equal(t, ".png", item.extension)
	assert.Equal(t, "image", item.mimeType)
	assert.Equal(t, "png", item.mimeSubType)
	assert.Equal(t, theme.FileImageIcon(), item.resource)

	item = NewFileIcon(storage.NewURI("file:///path/to/filename.txt"))

	assert.Equal(t, ".txt", item.extension)
	assert.Equal(t, "text", item.mimeType)
	assert.Equal(t, "plain", item.mimeSubType)
	assert.Equal(t, theme.FileTextIcon(), item.resource)

	item = NewFileIcon(storage.NewURI("file:///path/to/filename.mp4"))

	assert.Equal(t, ".mp4", item.extension)
	assert.Equal(t, "video", item.mimeType)
	assert.Equal(t, "mp4", item.mimeSubType)
	assert.Equal(t, theme.FileVideoIcon(), "mp4")
}

func TestNewFileIconNoExtension(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		fyne.LogError("Could not get current working directory", err)
		t.FailNow()
	}
	binFileWithNoExt := filepath.Join(workingDir, "testdata/bin")
	textFileWithNoExt := filepath.Join(workingDir, "testdata/text")

	item := NewFileIcon(storage.NewURI("file://" + binFileWithNoExt))

	assert.Equal(t, "", item.extension)
	assert.Equal(t, "application", item.mimeType)
	assert.Equal(t, "octet-stream", item.mimeSubType)
	assert.Equal(t, theme.FileApplicationIcon(), item.resource)

	item = NewFileIcon(storage.NewURI("file://" + textFileWithNoExt))

	assert.Equal(t, ".txt", item.extension)
	assert.Equal(t, "text", item.mimeType)
	assert.Equal(t, "plain", item.mimeSubType)
	assert.Equal(t, theme.FileTextIcon(), item.resource)
}

func TestUpdateURI(t *testing.T) {
	item := NewFileIcon(storage.NewURI("file:///path/to/filename.zip"))

	// assert.Equal(t, ".zip", item.extension)
	// assert.Equal(t, "application", item.mimeType)
	// assert.Equal(t, "zip", item.mimeSubType)
	// assert.Equal(t, theme.FileApplicationIcon(), item.resource)

	item.UpdateURI(storage.NewURI("file:///path/to/filename.mp3"))

	assert.Equal(t, ".mp3", item.extension)
	assert.Equal(t, "audio", item.mimeType)
	assert.Equal(t, "mpeg", item.mimeSubType)
	assert.Equal(t, theme.FileAudioIcon(), item.resource)

	item.UpdateURI(storage.NewURI("file:///path/to/filename.png"))

	assert.Equal(t, ".png", item.extension)
	assert.Equal(t, "image", item.mimeType)
	assert.Equal(t, "png", item.mimeSubType)
	assert.Equal(t, theme.FileImageIcon(), item.resource)

	item.UpdateURI(storage.NewURI("file:///path/to/filename.txt"))

	assert.Equal(t, ".txt", item.extension)
	assert.Equal(t, "text", item.mimeType)
	assert.Equal(t, "plain", item.mimeSubType)
	assert.Equal(t, theme.FileTextIcon(), item.resource)

	item.UpdateURI(storage.NewURI("file:///path/to/filename.mp4"))

	assert.Equal(t, ".mp4", item.extension)
	assert.Equal(t, "video", item.mimeType)
	assert.Equal(t, "mp4", item.mimeSubType)
	assert.Equal(t, theme.FileVideoIcon(), "mp4")
}

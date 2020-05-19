package dialog

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/theme"
)

func TestNewFileIcon(t *testing.T) {
	f := &fileDialog{}
	_ = f.makeUI()
	item := NewFileIcon("/path/to/filename.txt")

	assert.Equal(t, item.(*fileIcon).extension, ".txt")
	assert.Equal(t, item.(*fileIcon).mimeType, "text")
	assert.Equal(t, item.(*fileIcon).mimeSubType, "plain")
	assert.Equal(t, item.(*fileIcon).resource, theme.FileTextIcon())

	item = NewFileIcon("/path/to/filename.png")

	assert.Equal(t, item.(*fileIcon).extension, ".png")
	assert.Equal(t, item.(*fileIcon).mimeType, "image")
	assert.Equal(t, item.(*fileIcon).mimeSubType, "png")
	assert.Equal(t, item.(*fileIcon).resource, theme.FileImageIcon())

	item = NewFileIcon("/path/to/filename.mp3")

	assert.Equal(t, item.(*fileIcon).extension, ".mp3")
	assert.Equal(t, item.(*fileIcon).mimeType, "audio")
	assert.Equal(t, item.(*fileIcon).mimeSubType, "mpeg")
	assert.Equal(t, item.(*fileIcon).resource, theme.FileAudioIcon())

	item = NewFileIcon("/path/to/filename.mp4")

	assert.Equal(t, item.(*fileIcon).extension, ".mp4")
	assert.Equal(t, item.(*fileIcon).mimeType, "video")
	assert.Equal(t, item.(*fileIcon).mimeSubType, "mp4")
	assert.Equal(t, item.(*fileIcon).resource, theme.FileVideoIcon())

	item = NewFileIcon("/path/to/filename.exe")

	assert.Equal(t, item.(*fileIcon).extension, ".exe")
	assert.Equal(t, item.(*fileIcon).mimeType, "application")
	assert.Equal(t, item.(*fileIcon).mimeSubType, "octet-stream")
	assert.Equal(t, item.(*fileIcon).resource, theme.FileApplicationIcon())
}

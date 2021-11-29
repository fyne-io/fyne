//go:build !windows
// +build !windows

package widget

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

// Simulate being rendered by calling CreateRenderer() to update icon
func newRenderedFileIcon(uri fyne.URI) *FileIcon {
	f := NewFileIcon(uri)
	f.CreateRenderer()
	return f
}

func TestFileIcon_NewFileIcon(t *testing.T) {
	item := newRenderedFileIcon(storage.NewFileURI("/path/to/filename.zip"))
	assert.Equal(t, ".zip", item.extension)
	assert.Equal(t, theme.FileApplicationIcon(), item.resource)

	item = newRenderedFileIcon(storage.NewFileURI("/path/to/filename.mp3"))
	assert.Equal(t, ".mp3", item.extension)
	assert.Equal(t, theme.FileAudioIcon(), item.resource)

	item = newRenderedFileIcon(storage.NewFileURI("/path/to/filename.png"))
	assert.Equal(t, ".png", item.extension)
	assert.Equal(t, theme.FileImageIcon(), item.resource)

	item = newRenderedFileIcon(storage.NewFileURI("/path/to/filename.txt"))
	assert.Equal(t, ".txt", item.extension)
	assert.Equal(t, theme.FileTextIcon(), item.resource)

	item = newRenderedFileIcon(storage.NewFileURI("/path/to/filename.mp4"))
	assert.Equal(t, ".mp4", item.extension)
	assert.Equal(t, theme.FileVideoIcon(), item.resource)
}

func TestFileIcon_NewFileIcon_NoExtension(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		fyne.LogError("Could not get current working directory", err)
		t.FailNow()
	}
	binFileWithNoExt := filepath.Join(workingDir, "testdata/bin")
	textFileWithNoExt := filepath.Join(workingDir, "testdata/text")

	item := newRenderedFileIcon(storage.NewFileURI(binFileWithNoExt))
	assert.Equal(t, "", item.extension)
	assert.Equal(t, theme.FileApplicationIcon(), item.resource)

	item = newRenderedFileIcon(storage.NewFileURI(textFileWithNoExt))
	assert.Equal(t, "", item.extension)
	assert.Equal(t, theme.FileTextIcon(), item.resource)
}

func TestFileIcon_NewURI_WithFolder(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		fyne.LogError("Could not get current working directory", err)
		t.FailNow()
	}

	dir := filepath.Join(workingDir, "testdata")
	item := newRenderedFileIcon(storage.NewFileURI(dir))
	assert.Empty(t, item.extension)
	assert.Equal(t, theme.FolderIcon(), item.resource)

	item.SetURI(storage.NewFileURI(dir))
	assert.Empty(t, item.extension)
	assert.Equal(t, theme.FolderIcon(), item.resource)
}

func TestFileIcon_NewFileIcon_Rendered(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	workingDir, err := os.Getwd()
	if err != nil {
		fyne.LogError("Could not get current working directory", err)
		t.FailNow()
	}

	icon := NewFileIcon(nil)

	w := test.NewWindow(icon)
	w.Resize(fyne.NewSize(150, 150))

	test.AssertImageMatches(t, "fileicon/fileicon_nil.png", w.Canvas().Capture())

	text := filepath.Join(workingDir, "testdata/text")
	icon2 := NewFileIcon(storage.NewFileURI(text))

	w.SetContent(icon2)
	w.Resize(fyne.NewSize(150, 150))
	test.AssertImageMatches(t, "fileicon/fileicon_text.png", w.Canvas().Capture())

	text += ".txt"
	icon3 := NewFileIcon(storage.NewFileURI(text))

	w.SetContent(icon3)
	w.Resize(fyne.NewSize(150, 150))
	test.AssertImageMatches(t, "fileicon/fileicon_text_txt.png", w.Canvas().Capture())

	bin := filepath.Join(workingDir, "testdata/bin")
	icon4 := NewFileIcon(storage.NewFileURI(bin))

	w.SetContent(icon4)
	w.Resize(fyne.NewSize(150, 150))
	test.AssertImageMatches(t, "fileicon/fileicon_bin.png", w.Canvas().Capture())

	dir := filepath.Join(workingDir, "testdata")
	icon5 := NewFileIcon(storage.NewFileURI(dir))

	w.SetContent(icon5)
	w.Resize(fyne.NewSize(150, 150))
	test.AssertImageMatches(t, "fileicon/fileicon_folder.png", w.Canvas().Capture())

	w.Close()
}

func TestFileIcon_SetURI(t *testing.T) {
	item := newRenderedFileIcon(storage.NewFileURI("/path/to/filename.zip"))
	assert.Equal(t, ".zip", item.extension)
	assert.Equal(t, theme.FileApplicationIcon(), item.resource)

	item.SetURI(storage.NewFileURI("/path/to/filename.mp3"))
	assert.Equal(t, ".mp3", item.extension)
	assert.Equal(t, theme.FileAudioIcon(), item.resource)

	item.SetURI(storage.NewFileURI("/path/to/filename.png"))
	assert.Equal(t, ".png", item.extension)
	assert.Equal(t, theme.FileImageIcon(), item.resource)

	item.SetURI(storage.NewFileURI("/path/to/filename.txt"))
	assert.Equal(t, ".txt", item.extension)
	assert.Equal(t, theme.FileTextIcon(), item.resource)

	item.SetURI(storage.NewFileURI("/path/to/filename.mp4"))
	assert.Equal(t, ".mp4", item.extension)
	assert.Equal(t, theme.FileVideoIcon(), item.resource)
}

func TestFileIcon_SetURI_WithFolder(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		fyne.LogError("Could not get current working directory", err)
		t.FailNow()
	}
	dir := filepath.Join(workingDir, "testdata")

	item := newRenderedFileIcon(nil)
	assert.Empty(t, item.extension)

	item.SetURI(storage.NewFileURI(dir))
	assert.Empty(t, item.extension)
	assert.Equal(t, theme.FolderIcon(), item.resource)

	item.SetURI(storage.NewFileURI(dir))
	assert.Empty(t, item.extension)
	assert.Equal(t, theme.FolderIcon(), item.resource)
}

func TestFileIcon_DirURIUpdated(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		fyne.LogError("Could not get current working directory", err)
		t.FailNow()
	}
	testDir := filepath.Join(workingDir, "testdata", "notCreatedYet")

	uri := storage.NewFileURI(testDir)
	item := newRenderedFileIcon(uri)

	// The directory has not been created. It can not be listed yet.
	assert.Equal(t, theme.FileTextIcon(), item.resource)

	err = os.Mkdir(testDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testDir)

	item.Refresh()
	assert.Equal(t, theme.FolderIcon(), item.resource)
}

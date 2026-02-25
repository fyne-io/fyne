package widget_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/software"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestFileIcon_NewFileIcon_Rendered(t *testing.T) {
	test.NewTempApp(t)

	workingDir, err := os.Getwd()
	if err != nil {
		fyne.LogError("Could not get current working directory", err)
		t.FailNow()
	}

	icon := widget.NewFileIcon(nil)

	w := test.NewTempWindow(t, icon)
	w.Resize(fyne.NewSize(150, 150))

	test.AssertImageMatches(t, "fileicon/fileicon_nil.png", w.Canvas().Capture())

	text := filepath.Join(workingDir, "testdata/text")
	icon2 := widget.NewFileIcon(storage.NewFileURI(text))

	w.SetContent(icon2)
	w.Resize(fyne.NewSize(150, 150))
	test.AssertImageMatches(t, "fileicon/fileicon_text.png", w.Canvas().Capture())

	text += ".txt"
	icon3 := widget.NewFileIcon(storage.NewFileURI(text))

	w.SetContent(icon3)
	w.Resize(fyne.NewSize(150, 150))
	test.AssertImageMatches(t, "fileicon/fileicon_text_txt.png", w.Canvas().Capture())

	bin := filepath.Join(workingDir, "testdata/bin")
	icon4 := widget.NewFileIcon(storage.NewFileURI(bin))

	w.SetContent(icon4)
	w.Resize(fyne.NewSize(150, 150))
	test.AssertImageMatches(t, "fileicon/fileicon_bin.png", w.Canvas().Capture())

	dir := filepath.Join(workingDir, "testdata")
	icon5 := widget.NewFileIcon(storage.NewFileURI(dir))

	w.SetContent(icon5)
	w.Resize(fyne.NewSize(150, 150))
	test.AssertImageMatches(t, "fileicon/fileicon_folder.png", w.Canvas().Capture())
}

func TestFileIcon_Icon(t *testing.T) {
	dir := storage.NewFileURI("testdata")

	icon1 := widget.NewFileIcon(dir)
	trash := &customURI{URI: dir, icon: theme.DeleteIcon()}
	icon2 := widget.NewFileIcon(trash)

	// test icon change
	icon1Img := software.Render(icon1, test.Theme())
	icon2Img := software.Render(icon2, test.Theme())
	assert.NotEqual(t, icon1Img, icon2Img)
}

type customURI struct {
	fyne.URI

	name string
	icon fyne.Resource
}

func (c *customURI) Icon() fyne.Resource {
	return c.icon
}

func (c *customURI) Name() string {
	if c.name != "" {
		return c.name
	}

	return c.URI.Name()
}

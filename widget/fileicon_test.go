package widget_test

import (
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestFileIcon_NewFileIcon_Rendered(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	workingDir, err := os.Getwd()
	if err != nil {
		fyne.LogError("Could not get current working directory", err)
		t.FailNow()
	}

	icon := widget.NewFileIcon(nil)

	w := test.NewWindow(icon)
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

	w.Close()
}

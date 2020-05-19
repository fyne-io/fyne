package dialog_test

import (
	"image/color"
	"strings"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
)

func TestProgressDialog_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		title, message string
	}{
		"label": {
			title:   "Title",
			message: "Working",
		},
		"label_long_title": {
			title:   strings.Repeat("Title", 100),
			message: "Working",
		},
		"label_long_dismiss": {
			title:   "Title",
			message: strings.Repeat("Working", 100),
		},
	} {
		t.Run(name, func(t *testing.T) {
			window := fyne.CurrentApp().NewWindow(name)
			window.SetContent(canvas.NewRectangle(color.Black))
			progress := dialog.NewProgress(tt.title, tt.message, window)
			progress.Show()
			window.Resize(fyne.NewSize(400, 300))

			test.AssertImageMatches(t, "progress/layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}

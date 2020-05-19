package dialog_test

import (
	"errors"
	"image/color"
	"strings"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
)

func TestErrorDialog_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		err error
	}{
		"error": {
			err: errors.New("TestError"),
		},
		"error_long": {
			err: errors.New(strings.Repeat("TestError", 100)),
		},
	} {
		t.Run(name, func(t *testing.T) {
			window := fyne.CurrentApp().NewWindow(name)
			window.SetContent(canvas.NewRectangle(color.Black))
			dialog.ShowError(tt.err, window)
			window.Resize(fyne.NewSize(400, 300))

			test.AssertImageMatches(t, "error/layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}

func TestInformationDialog_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		title   string
		message string
	}{
		"label": {
			title:   "Title",
			message: "Information",
		},
		"label_long_title": {
			title:   strings.Repeat("Title", 100),
			message: "Information",
		},
		"label_long_message": {
			title:   "Title",
			message: strings.Repeat("Information", 100),
		},
	} {
		t.Run(name, func(t *testing.T) {
			window := fyne.CurrentApp().NewWindow(name)
			window.SetContent(canvas.NewRectangle(color.Black))
			dialog.ShowInformation(tt.title, tt.message, window)
			window.Resize(fyne.NewSize(400, 300))

			test.AssertImageMatches(t, "information/layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}

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
	"fyne.io/fyne/widget"
)

func TestCustomDialog_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		title   string
		dismiss string
		content fyne.CanvasObject
	}{
		"label": {
			title:   "Title",
			dismiss: "Dismiss",
			content: widget.NewLabel("FooBar"),
		},
		"label_long_title": {
			title:   strings.Repeat("Title", 100),
			dismiss: "Dismiss",
			content: widget.NewLabel("FooBar"),
		},
		"label_long_dismiss": {
			title:   "Title",
			dismiss: strings.Repeat("Dismiss", 100),
			content: widget.NewLabel("FooBar"),
		},
		"label_long_content": {
			title:   "Title",
			dismiss: "Dismiss",
			content: widget.NewLabel(strings.Repeat("FooBar", 100)),
		},
	} {
		t.Run(name, func(t *testing.T) {
			window := fyne.CurrentApp().NewWindow(name)
			window.SetContent(canvas.NewRectangle(color.Black))
			dialog.ShowCustom(tt.title, tt.dismiss, tt.content, window)
			window.Resize(fyne.NewSize(400, 300))

			test.AssertImageMatches(t, "custom/layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}

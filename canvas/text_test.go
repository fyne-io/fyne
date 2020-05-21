package canvas_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func TestText_MinSize(t *testing.T) {
	text := canvas.NewText("Test", color.NRGBA{0, 0, 0, 0xff})
	min := text.MinSize()

	assert.True(t, min.Width > 0)
	assert.True(t, min.Height > 0)

	text = canvas.NewText("Test2", color.NRGBA{0, 0, 0, 0xff})
	min2 := text.MinSize()
	assert.True(t, min2.Width > min.Width)
}

func TestText_MinSize_NoMultiLine(t *testing.T) {
	text := canvas.NewText("Break", color.NRGBA{0, 0, 0, 0xff})
	min := text.MinSize()

	text = canvas.NewText("Bre\nak", color.NRGBA{0, 0, 0, 0xff})
	min2 := text.MinSize()
	assert.True(t, min2.Width > min.Width)
	assert.True(t, min2.Height == min.Height)
}

func TestText_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		text string
		size fyne.Size
	}{
		"short_small": {
			text: "abc",
			size: fyne.NewSize(1, 1),
		},
		"short_large": {
			text: "abc",
			size: fyne.NewSize(100, 100),
		},
		"long_small": {
			text: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			size: fyne.NewSize(1, 1),
		},
		"long_large": {
			text: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			size: fyne.NewSize(100, 100),
		},
	} {
		t.Run(name, func(t *testing.T) {
			text := canvas.NewText(tt.text, theme.TextColor())

			window := test.NewWindow(text)
			window.Resize(text.MinSize().Max(tt.size))

			test.AssertImageMatches(t, "text/layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}

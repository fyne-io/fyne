// +build !ci

package gl

import (
	"image/color"
	"testing"
	"time"

	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"github.com/stretchr/testify/assert"
)

func TestGlCanvas_Content(t *testing.T) {
	content := &canvas.Circle{}
	w := d.CreateWindow("Test")
	w.SetContent(content)

	assert.Equal(t, content, w.Content())
}

func TestGlCanvas_NilContent(t *testing.T) {
	w := d.CreateWindow("Test")

	assert.NotNil(t, w.Content()) // never a nil canvas so we have a sensible fallback
}

func Test_glCanvas_SetContent(t *testing.T) {
	var menuHeight int
	if hasNativeMenu() {
		menuHeight = 0
	} else {
		menuHeight = widget.NewToolbar(widget.NewToolbarAction(theme.ContentCutIcon(), func() {})).MinSize().Height
	}
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	tests := []struct {
		name               string
		padding            bool
		menu               bool
		expectedPad        int
		expectedMenuHeight int
	}{
		{"window without padding", false, false, 0, 0},
		{"window with padding", true, false, 4, 0},
		{"window with menu without padding", false, true, 0, menuHeight},
		{"window with menu and padding", true, true, 4, menuHeight},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := d.CreateWindow("Test").(*window)
			w.Canvas().SetScale(1)
			w.SetPadded(tt.padding)
			if tt.menu {
				w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("Test", fyne.NewMenuItem("Test", func() {}))))
			}
			content := canvas.NewCircle(color.Black)
			w.SetContent(content)
			w.Resize(fyne.NewSize(100, 100))
			c := w.Canvas()
			canvasSize := 100

			// wait for canvas to get its size right
			for w.canvas.Size().Width == canvasSize {
				time.Sleep(time.Millisecond * 10)
			}

			newContent := canvas.NewCircle(color.White)
			assert.Equal(t, fyne.NewPos(0, 0), newContent.Position())
			assert.Equal(t, fyne.NewSize(0, 0), newContent.Size())
			c.SetContent(newContent)
			assert.Equal(t, fyne.NewPos(tt.expectedPad, tt.expectedPad+tt.expectedMenuHeight), newContent.Position())
			assert.Equal(t, fyne.NewSize(canvasSize-2*tt.expectedPad, canvasSize-2*tt.expectedPad-tt.expectedMenuHeight), newContent.Size())
		})
	}
}

// +build !ci

package gl

import (
	"image/color"
	"sync"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"github.com/go-gl/glfw/v3.2/glfw"
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
		menuHeight = 12
	}
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
			var wg sync.WaitGroup
			wg.Add(2)
			hookIntoResizedCallback(w, &wg)
			w.SetPadded(tt.padding)
			if tt.menu {
				w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("Test", fyne.NewMenuItem("Test", func() {}))))
			}
			w.SetContent(canvas.NewCircle(color.Black))
			w.Resize(fyne.NewSize(100, 100))
			c := w.Canvas()
			wg.Wait()

			newContent := canvas.NewCircle(color.White)
			assert.Equal(t, fyne.NewPos(0, 0), newContent.Position())
			assert.Equal(t, fyne.NewSize(0, 0), newContent.Size())
			c.SetContent(newContent)
			canvasSize := 99
			assert.Equal(t, fyne.NewPos(tt.expectedPad, tt.expectedPad+tt.expectedMenuHeight), newContent.Position())
			assert.Equal(t, fyne.NewSize(canvasSize-2*tt.expectedPad, canvasSize-2*tt.expectedPad-tt.expectedMenuHeight), newContent.Size())
		})
	}
}

func hookIntoResizedCallback(w *window, wg *sync.WaitGroup) {
	var prev glfw.SizeCallback
	prev = w.viewport.SetSizeCallback(func(viewport *glfw.Window, width, height int) {
		prev(viewport, width, height)
		wg.Done()
	})
}

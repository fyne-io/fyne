//go:build !no_glfw && !mobile && flakey

package glfw

import (
	"testing"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/stretchr/testify/assert"
)

func TestWindow_TappedIgnoredWhenMovedOffOfTappable(t *testing.T) {
	w := createWindow("Test")
	tapped := 0
	b1 := widget.NewButton("Tap", func() { tapped = 1 })
	b2 := widget.NewButton("Tap", func() { tapped = 2 })
	w.SetContent(container.NewVBox(b1, b2))

	runOnMain(func() {
		w.mouseMoved(w.viewport, 17, 27)
		w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
		w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

		assert.Equal(t, 1, tapped, "Button 1 should be tapped")
		tapped = 0

		w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
		w.mouseMoved(w.viewport, 17, 59)
		w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

		assert.Equal(t, 0, tapped, "button was tapped without mouse press & release on it %d", tapped)

		w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
		w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

		assert.Equal(t, 2, tapped, "Button 2 should be tapped")
	})
}

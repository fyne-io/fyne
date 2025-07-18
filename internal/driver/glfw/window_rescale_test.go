//go:build !wasm && !test_web_driver && !mobile && !no_glfw

package glfw

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestWindow_RescaleContext_Basic(t *testing.T) {
	// Basic test to ensure RescaleContext doesn't panic
	w := createWindow("Test")
	defer w.Close()
	w.SetContent(widget.NewLabel("Test"))

	// Test that RescaleContext can be called without panic
	runOnMain(func() {
		w.RescaleContext()
	})

	// Verify window still has valid canvas
	assert.NotNil(t, w.Canvas())
	assert.NotNil(t, w.Canvas().Content())
}

func TestWindow_RescaleContext_WithContent(t *testing.T) {
	// Test RescaleContext with actual content
	w := createWindow("Test")
	defer w.Close()

	content := container.NewVBox(
		widget.NewLabel("Test App"),
		widget.NewButton("Click Me", func() {}),
	)
	w.SetContent(content)

	// Call RescaleContext
	runOnMain(func() {
		w.RescaleContext()
	})

	// Verify content is preserved
	assert.Equal(t, content, w.Canvas().Content())
	assert.Greater(t, w.Canvas().Size().Width, float32(0))
	assert.Greater(t, w.Canvas().Size().Height, float32(0))
}

func TestWindow_RescaleContext_MinSize(t *testing.T) {
	// Test that RescaleContext respects minimum size
	w := createWindow("Test")
	defer w.Close()

	label := widget.NewLabel("Test with minimum size")
	label.Resize(fyne.NewSize(300, 200))
	w.SetContent(label)

	// Call RescaleContext
	runOnMain(func() {
		w.RescaleContext()
	})

	// Canvas should respect content minimum size
	minSize := label.MinSize()
	canvasSize := w.Canvas().Size()
	assert.GreaterOrEqual(t, canvasSize.Width, minSize.Width)
	assert.GreaterOrEqual(t, canvasSize.Height, minSize.Height)
}

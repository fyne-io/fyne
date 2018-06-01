package layout

import "testing"

import "image/color"

import "github.com/fyne-io/fyne/api/ui"
import "github.com/fyne-io/fyne/api/ui/canvas"

import "github.com/stretchr/testify/assert"

func TestMaxLayout(t *testing.T) {
	size := ui.NewSize(100, 100)

	obj := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	container := &ui.Container{
		Size:    size,
		Objects: []ui.CanvasObject{obj},
	}

	NewMaxLayout().Layout(container.Objects, size)

	assert.Equal(t, obj.Size, size)
}

func TestMaxLayoutMinSize(t *testing.T) {
	text := canvas.NewText("Padding", color.RGBA{0, 0xff, 0, 0})
	minSize := text.MinSize()

	container := ui.NewContainer(text)
	layoutMin := NewMaxLayout().MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestContainerMaxLayoutMinSize(t *testing.T) {
	text := canvas.NewText("Padding", color.RGBA{0, 0xff, 0, 0})
	minSize := text.MinSize()

	container := ui.NewContainer(text)
	container.Layout = NewMaxLayout()
	layoutMin := container.MinSize()

	assert.Equal(t, minSize, layoutMin)
}

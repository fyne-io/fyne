package layout

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"github.com/stretchr/testify/assert"
)

func TestMaxLayout(t *testing.T) {
	size := fyne.NewSize(100, 100)

	obj := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj},
	}
	container.Resize(size)

	NewMaxLayout().Layout(container.Objects, size)

	assert.Equal(t, obj.Size(), size)
}

func TestMaxLayoutMinSize(t *testing.T) {
	text := canvas.NewText("Padding", color.RGBA{0, 0xff, 0, 0})
	minSize := text.MinSize()

	container := fyne.NewContainer(text)
	layoutMin := NewMaxLayout().MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestContainerMaxLayoutMinSize(t *testing.T) {
	text := canvas.NewText("Padding", color.RGBA{0, 0xff, 0, 0})
	minSize := text.MinSize()

	container := fyne.NewContainer(text)
	container.Layout = NewMaxLayout()
	layoutMin := container.MinSize()

	assert.Equal(t, minSize, layoutMin)
}

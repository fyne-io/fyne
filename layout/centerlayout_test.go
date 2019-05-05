package layout

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"github.com/stretchr/testify/assert"
)

func TestCenterLayout(t *testing.T) {
	size := fyne.NewSize(100, 100)
	min := fyne.NewSize(10, 10)

	obj := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj.SetMinSize(min)
	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj},
	}
	container.Resize(size)

	NewCenterLayout().Layout(container.Objects, size)

	assert.Equal(t, obj.Size(), min)
	assert.Equal(t, fyne.NewPos(45, 45), obj.Position())
}

func TestCenterLayoutMinSize(t *testing.T) {
	text := canvas.NewText("Padding", color.RGBA{0, 0xff, 0, 0})
	minSize := text.MinSize()

	container := fyne.NewContainer(text)
	layoutMin := NewCenterLayout().MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestContainerCenterLayoutMinSize(t *testing.T) {
	text := canvas.NewText("Padding", color.RGBA{0, 0xff, 0, 0})
	minSize := text.MinSize()

	container := fyne.NewContainer(text)
	container.Layout = NewCenterLayout()
	layoutMin := container.MinSize()

	assert.Equal(t, minSize, layoutMin)
}

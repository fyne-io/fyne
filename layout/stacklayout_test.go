package layout_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	"github.com/stretchr/testify/assert"
)

func TestStackLayout(t *testing.T) {
	size := fyne.NewSize(100, 100)

	obj := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	c := &fyne.Container{
		Objects: []fyne.CanvasObject{obj},
	}
	c.Resize(size)

	layout.NewStackLayout().Layout(c.Objects, size)

	assert.Equal(t, obj.Size(), size)
}

func TestStackLayoutMinSize(t *testing.T) {
	text := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	minSize := text.MinSize()

	c := container.NewWithoutLayout(text)
	layoutMin := layout.NewStackLayout().MinSize(c.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestContainerStackLayoutMinSize(t *testing.T) {
	text := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	minSize := text.MinSize()

	c := container.NewWithoutLayout(text)
	c.Layout = layout.NewStackLayout()
	layoutMin := c.MinSize()

	assert.Equal(t, minSize, layoutMin)
}

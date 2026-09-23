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

func TestCenterLayout(t *testing.T) {
	size := fyne.NewSize(100, 100)
	minSize := fyne.NewSize(10, 10)

	obj := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj.SetMinSize(minSize)
	c := &fyne.Container{
		Objects: []fyne.CanvasObject{obj},
	}
	c.Resize(size)

	layout.NewCenterLayout().Layout(c.Objects, size)

	assert.Equal(t, obj.Size(), minSize)
	assert.Equal(t, fyne.NewPos(45, 45), obj.Position())
}

func TestCenterLayout_MinSize(t *testing.T) {
	text := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	minSize := text.MinSize()

	c := container.NewWithoutLayout(text)
	layoutMin := layout.NewCenterLayout().MinSize(c.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestCenterLayout_MinSize_Hidden(t *testing.T) {
	text1 := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	text1.Hide()
	text2 := canvas.NewText("1\n2", color.NRGBA{0, 0xff, 0, 0})

	c := container.NewWithoutLayout(text1, text2)
	layoutMin := layout.NewCenterLayout().MinSize(c.Objects)

	assert.Equal(t, text2.MinSize(), layoutMin)
}

func TestContainerCenterLayoutMinSize(t *testing.T) {
	text := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	minSize := text.MinSize()

	c := container.NewWithoutLayout(text)
	c.Layout = layout.NewCenterLayout()
	layoutMin := c.MinSize()

	assert.Equal(t, minSize, layoutMin)
}

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
	min := fyne.NewSize(10, 10)

	obj := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj.SetMinSize(min)
	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj},
	}
	container.Resize(size)

	layout.NewCenterLayout().Layout(container.Objects, size)

	assert.Equal(t, obj.Size(), min)
	assert.Equal(t, fyne.NewPos(45, 45), obj.Position())
}

func TestCenterLayout_MinSize(t *testing.T) {
	text := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	minSize := text.MinSize()

	container := container.NewWithoutLayout(text)
	layoutMin := layout.NewCenterLayout().MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestCenterLayout_MinSize_Hidden(t *testing.T) {
	text1 := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	text1.Hide()
	text2 := canvas.NewText("1\n2", color.NRGBA{0, 0xff, 0, 0})

	container := container.NewWithoutLayout(text1, text2)
	layoutMin := layout.NewCenterLayout().MinSize(container.Objects)

	assert.Equal(t, text2.MinSize(), layoutMin)
}

func TestContainerCenterLayoutMinSize(t *testing.T) {
	text := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	minSize := text.MinSize()

	container := container.NewWithoutLayout(text)
	container.Layout = layout.NewCenterLayout()
	layoutMin := container.MinSize()

	assert.Equal(t, minSize, layoutMin)
}

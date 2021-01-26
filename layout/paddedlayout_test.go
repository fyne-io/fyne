package layout_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestPaddedLayout(t *testing.T) {
	size := fyne.NewSize(100, 100)

	obj := canvas.NewRectangle(color.Black)
	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj},
	}
	container.Resize(size)

	layout.NewPaddedLayout().Layout(container.Objects, size)

	assert.Equal(t, obj.Size().Width, size.Width-2*theme.Padding())
	assert.Equal(t, obj.Size().Height, size.Height-2*theme.Padding())
}

func TestPaddedLayout_MinSize(t *testing.T) {
	text := canvas.NewText("FooBar", color.Black)
	minSize := text.MinSize()

	container := fyne.NewContainer(text)
	layoutMin := layout.NewPaddedLayout().MinSize(container.Objects)

	assert.Equal(t, minSize.Width+2*theme.Padding(), layoutMin.Width)
	assert.Equal(t, minSize.Height+2*theme.Padding(), layoutMin.Height)
}

func TestContainer_PaddedLayout_MinSize(t *testing.T) {
	text := canvas.NewText("FooBar", color.Black)
	minSize := text.MinSize()

	container := fyne.NewContainer(text)
	container.Layout = layout.NewPaddedLayout()
	layoutMin := container.MinSize()

	assert.Equal(t, minSize.Width+2*theme.Padding(), layoutMin.Width)
	assert.Equal(t, minSize.Height+2*theme.Padding(), layoutMin.Height)
}

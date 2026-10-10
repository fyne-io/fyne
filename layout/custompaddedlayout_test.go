package layout_test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func TestCustomPaddedLayout(t *testing.T) {
	size := fyne.NewSize(100, 100)

	obj := canvas.NewRectangle(color.Black)
	c := &fyne.Container{
		Objects: []fyne.CanvasObject{obj},
	}
	c.Resize(size)

	layout.NewCustomPaddedLayout(2, 3, 4, 5).Layout(c.Objects, size)

	assert.Equal(t, obj.Size().Width, size.Width-4-5)
	assert.Equal(t, obj.Size().Height, size.Height-2-3)
}

func TestCustomPaddedLayout_MinSize(t *testing.T) {
	text := canvas.NewText("FooBar", color.Black)
	minSize := text.MinSize()

	c := container.NewWithoutLayout(text)
	layoutMin := layout.NewCustomPaddedLayout(2, 3, 4, 5).MinSize(c.Objects)

	assert.Equal(t, minSize.Width+4+5, layoutMin.Width)
	assert.Equal(t, minSize.Height+2+3, layoutMin.Height)
}

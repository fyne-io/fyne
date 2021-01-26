package widget_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestFlexible_Methods(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	minSize := fyne.NewSize(110, 150)
	rect.SetMinSize(fyne.NewSize(110, 150))

	w := widget.NewFlexible(4, rect)
	assert.Equal(t, 4, w.Flex())
	assert.Equal(t, minSize, w.MinSize())
	assert.Zero(t, w.DistanceToTextBaseline())

	w = widget.NewFlexible(7, nil)
	assert.Equal(t, 7, w.Flex())
	assert.Equal(t, canvas.NewRectangle(color.Black).MinSize(), w.MinSize())
	assert.Zero(t, w.DistanceToTextBaseline())

	w = widget.NewExpanded(nil)
	assert.Equal(t, 1, w.Flex())
	assert.Equal(t, canvas.NewRectangle(color.Black).MinSize(), w.MinSize())
	assert.Zero(t, w.DistanceToTextBaseline())
}

func TestFlexible_CreateRenderer(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(120, 150))

	w := widget.NewExpanded(rect)
	renderer := test.WidgetRenderer(w)
	assert.NotNil(t, renderer)

	assert.Equal(t, rect.MinSize(), renderer.MinSize())
	assert.Len(t, renderer.Objects(), 1)
	assert.Equal(t, rect, renderer.Objects()[0])

	resize := fyne.NewSize(200, 200)
	w.Resize(resize)
	renderer.Refresh()
	renderer.Layout(resize)
	assert.Equal(t, fyne.NewPos(0, 0), rect.Position())
	assert.Equal(t, resize, rect.Size())
}

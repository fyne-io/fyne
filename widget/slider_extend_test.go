package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/cache"
)

type extendedSlider struct {
	Slider
}

func newExtendedSlider() *extendedSlider {
	slider := &extendedSlider{}
	slider.ExtendBaseWidget(slider)
	slider.Min = 0
	slider.Max = 10
	return slider
}

func TestSlider_Extended_Value(t *testing.T) {
	slider := newExtendedSlider()
	slider.Resize(slider.MinSize().Add(fyne.NewSize(20, 0)))
	objs := cache.Renderer(slider).Objects()
	assert.Equal(t, 3, len(objs))
	thumb := objs[2]
	thumbPos := thumb.Position()

	slider.Value = 2
	slider.Refresh()
	assert.Greater(t, thumb.Position().X, thumbPos.X)
	assert.Equal(t, thumbPos.Y, thumb.Position().Y)
}

func TestSlider_Extended_Drag(t *testing.T) {
	slider := newExtendedSlider()
	objs := cache.Renderer(slider).Objects()
	assert.Equal(t, 3, len(objs))
	thumb := objs[2]
	thumbPos := thumb.Position()

	drag := &fyne.DragEvent{DraggedX: 10, DraggedY: 2}
	slider.Dragged(drag)
	assert.Greater(t, thumbPos.X, thumb.Position().X)
	assert.Equal(t, thumbPos.Y, thumb.Position().Y)
}

package widget

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestScrollContainer_Scrolled(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScroller(rect)
	scroll.Resize(fyne.NewSize(100, 100))

	assert.Equal(t, 0, scroll.Offset.Y)
	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: -10})
	assert.Equal(t, 10, scroll.Offset.Y)
}

func TestScrollContainer_Scrolled_Limit(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScroller(rect)
	scroll.Resize(fyne.NewSize(80, 80))

	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: -25})
	assert.Equal(t, 20, scroll.Offset.Y)
}

func TestScrollContainer_Scrolled_Back(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScroller(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	scroll.Offset.Y = 10

	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: 10})
	assert.Equal(t, 0, scroll.Offset.Y)
}

func TestScrollContainer_Scrolled_BackLimit(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(1000, 1000))
	scroll := NewScroller(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	scroll.Offset.Y = 10

	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: 20})
	assert.Equal(t, 0, scroll.Offset.Y)
}

func TestScrollContainer_Resize(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScroller(rect)
	scroll.Resize(fyne.NewSize(80, 80))

	scroll.Scrolled(&fyne.ScrollEvent{DeltaY: -20})
	scroll.Resize(fyne.NewSize(80, 100))
	assert.Equal(t, 0, scroll.Offset.Y)
}

func TestScrollContainerRenderer_BarSize(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(100, 100))
	scroll := NewScroller(rect)
	scroll.Resize(fyne.NewSize(100, 100))
	render := Renderer(scroll).(*scrollRenderer)

	assert.Equal(t, fyne.NewSize(theme.ScrollBarSize(), 100), render.barSizeVertical())

	// resize so content is twice our size. Bar should therefore be half again.
	scroll.Resize(fyne.NewSize(50, 50))
	assert.Equal(t, fyne.NewSize(theme.ScrollBarSize(), 25), render.barSizeVertical())
}

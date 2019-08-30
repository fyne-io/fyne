package widget

import (
	"testing"

	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"
)

func TestSlider_HorizontalLayout(t *testing.T) {
	slider := NewSlider(0, 1)
	slider.Resize(fyne.NewSize(100, 10))

	render := Renderer(slider).(*sliderRenderer)

	wSize := render.slider.Size()
	tSize := render.track.Size()
	aSize := render.active.Size()

	assert.Greater(t, wSize.Width, wSize.Height)

	assert.Equal(t, wSize.Width, tSize.Width)
	assert.Greater(t, wSize.Height, tSize.Height)

	assert.Greater(t, wSize.Width, aSize.Width)
	assert.Greater(t, wSize.Height, aSize.Height)
}

func TestSlider_VerticalLayout(t *testing.T) {
	slider := NewSlider(0, 1)
	slider.Orientation = Vertical

	slider.Resize(fyne.NewSize(10, 100))

	render := Renderer(slider).(*sliderRenderer)

	wSize := render.slider.Size()
	tSize := render.track.Size()
	aSize := render.active.Size()

	assert.Greater(t, wSize.Height, wSize.Width)

	assert.Equal(t, wSize.Height, tSize.Height)
	assert.Greater(t, wSize.Width, tSize.Width)

	assert.Greater(t, wSize.Height, aSize.Height)
	assert.Greater(t, wSize.Width, aSize.Width)
}

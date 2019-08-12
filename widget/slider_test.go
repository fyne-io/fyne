package widget

import (
	"testing"

	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"
)

func TestSlider_MinMax(t *testing.T) {
	slider := NewSlider(5, 0, 10, 0)

	assert.Equal(t, 0.0, slider.Min)
	assert.Equal(t, 10.0, slider.Max)
	assert.Equal(t, 5.0, slider.Value)

	assert.Greater(t, slider.Max, slider.Min)

	slider = NewSlider(5, 10, 0, 0)

	assert.Greater(t, slider.Max, slider.Min)
	assert.Greater(t, slider.Max, slider.Value)
	assert.Greater(t, slider.Value, slider.Min)
}

func TestSlider_PrecisionClamp(t *testing.T) {
	precision := uint8(255)
	assert.Greater(t, precision, maxSliderDecimals)

	slider := NewSlider(5, 0, 10, precision)

	assert.Equal(t, slider.Precision, maxSliderDecimals)
}

func TestSlider_HorizontalLayout(t *testing.T) {
	slider := NewSlider(5, 0, 10, 0)
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
	slider := NewSlider(5, 0, 10, 0)
	slider.Vertical = true

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

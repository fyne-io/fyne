package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestSlider_HorizontalLayout(t *testing.T) {
	slider := NewSlider(0, 1)
	slider.Resize(fyne.NewSize(100, 10))

	render := test.WidgetRenderer(slider).(*sliderRenderer)
	diameter := slider.buttonDiameter()
	wSize := render.slider.Size()
	tSize := render.track.Size()
	aSize := render.active.Size()

	assert.Greater(t, wSize.Width, wSize.Height)

	assert.Equal(t, wSize.Width-diameter-theme.Padding()*2, tSize.Width)
	assert.Equal(t, theme.Padding(), tSize.Height)

	assert.Greater(t, wSize.Width, aSize.Width)
	assert.Equal(t, theme.Padding(), aSize.Height)
}

func TestSlider_OutOfRange(t *testing.T) {
	slider := NewSlider(2, 5)
	slider.Resize(fyne.NewSize(100, 10))

	assert.Equal(t, float64(2), slider.Value)
}

func TestSlider_VerticalLayout(t *testing.T) {
	slider := NewSlider(0, 1)
	slider.Orientation = Vertical
	slider.Resize(fyne.NewSize(10, 100))

	render := test.WidgetRenderer(slider).(*sliderRenderer)
	diameter := slider.buttonDiameter()
	wSize := render.slider.Size()
	tSize := render.track.Size()
	aSize := render.active.Size()

	assert.Greater(t, wSize.Height, wSize.Width)

	assert.Equal(t, wSize.Height-diameter-theme.Padding()*2, tSize.Height)
	assert.Equal(t, theme.Padding(), tSize.Width)

	assert.Greater(t, wSize.Height, aSize.Height)
	assert.Equal(t, theme.Padding(), aSize.Width)
}

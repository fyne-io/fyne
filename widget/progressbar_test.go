package widget

import (
	"testing"

	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"
)

func TestProgressBar_SetValue(t *testing.T) {
	bar := NewProgressBar()

	assert.Equal(t, 0.0, bar.Min)
	assert.Equal(t, 1.0, bar.Max)
	assert.Equal(t, 0.0, bar.Value)

	bar.SetValue(.5)
	assert.Equal(t, .5, bar.Value)
}

func TestProgressRenderer_Layout(t *testing.T) {
	bar := NewProgressBar()
	bar.Resize(fyne.NewSize(100, 10))

	render := Renderer(bar).(*progressRenderer)
	assert.Equal(t, 0, render.bar.Size().Width)

	bar.SetValue(.5)
	assert.Equal(t, 50, render.bar.Size().Width)

	bar.SetValue(1)
	assert.Equal(t, bar.Size().Width, render.bar.Size().Width)
}

func TestProgressRenderer_Layout_Overflow(t *testing.T) {
	bar := NewProgressBar()
	bar.Resize(fyne.NewSize(100, 10))

	render := Renderer(bar).(*progressRenderer)
	bar.SetValue(1)
	assert.Equal(t, bar.Size().Width, render.bar.Size().Width)

	bar.SetValue(1.2)
	assert.Equal(t, bar.Size().Width, render.bar.Size().Width)
}

func TestProgressRenderer_ApplyTheme(t *testing.T) {
	bar := NewProgressBar()
	render := Renderer(bar).(*progressRenderer)

	textSize := render.label.TextSize
	customTextSize := textSize
	withTestTheme(func() {
		render.applyTheme()
		customTextSize = render.label.TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}

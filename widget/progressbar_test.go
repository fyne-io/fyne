package widget

import (
	"fmt"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
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

func TestProgressBar_TextFormatter(t *testing.T) {
	bar := NewProgressBar()
	formatted := false

	bar.SetValue(0.2)
	assert.Equal(t, false, formatted)

	formatter := func() string {
		formatted = true
		return fmt.Sprintf("%.2f out of %.2f", bar.Value, bar.Max)
	}
	bar.TextFormatter = formatter

	bar.SetValue(0.4)

	assert.Equal(t, true, formatted)
}

func TestProgressRenderer_Layout(t *testing.T) {
	bar := NewProgressBar()
	bar.Resize(fyne.NewSize(100, 10))

	render := test.WidgetRenderer(bar).(*progressRenderer)
	assert.Equal(t, 0, render.bar.Size().Width)

	bar.SetValue(.5)
	assert.Equal(t, 50, render.bar.Size().Width)

	bar.SetValue(1)
	assert.Equal(t, bar.Size().Width, render.bar.Size().Width)
}

func TestProgressRenderer_Layout_Overflow(t *testing.T) {
	bar := NewProgressBar()
	bar.Resize(fyne.NewSize(100, 10))

	render := test.WidgetRenderer(bar).(*progressRenderer)
	bar.SetValue(1)
	assert.Equal(t, bar.Size().Width, render.bar.Size().Width)

	bar.SetValue(1.2)
	assert.Equal(t, bar.Size().Width, render.bar.Size().Width)
}

func TestProgressRenderer_ApplyTheme(t *testing.T) {
	bar := NewProgressBar()
	render := test.WidgetRenderer(bar).(*progressRenderer)

	textSize := render.label.TextSize
	customTextSize := textSize
	test.WithTestTheme(t, func() {
		render.applyTheme()
		customTextSize = render.label.TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}

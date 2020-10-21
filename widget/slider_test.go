package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestSlider_HorizontalLayout(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

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

	w := test.NewWindow(slider)
	defer w.Close()
	w.Resize(fyne.NewSize(220, 50))

	test.AssertImageMatches(t, "slider/horizontal.png", w.Canvas().Capture())
}

func TestSlider_OutOfRange(t *testing.T) {
	slider := NewSlider(2, 5)
	slider.Resize(fyne.NewSize(100, 10))

	assert.Equal(t, float64(2), slider.Value)
}

func TestSlider_VerticalLayout(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

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

	w := test.NewWindow(slider)
	defer w.Close()
	w.Resize(fyne.NewSize(50, 220))

	test.AssertImageMatches(t, "slider/vertical.png", w.Canvas().Capture())
}

func TestSlider_OnChanged(t *testing.T) {
	slider := NewSlider(0, 1)
	assert.Empty(t, slider.OnChanged)

	changes := 0

	slider.OnChanged = func(_ float64) {
		changes++
	}

	assert.Equal(t, 0, changes)

	slider.SetValue(0.5)
	assert.Equal(t, 1, changes)

	drag := &fyne.DragEvent{DraggedX: 10, DraggedY: 2}
	slider.Dragged(drag)
	assert.Equal(t, 2, changes)
}

func TestSlider_SetValue(t *testing.T) {
	slider := NewSlider(0, 2)
	assert.Equal(t, 0.0, slider.Value)

	slider.SetValue(1.0)
	assert.Equal(t, 1.0, slider.Value)

	slider.SetValue(2.0)
	assert.Equal(t, 2.0, slider.Value)
}

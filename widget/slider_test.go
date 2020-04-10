package widget

import (
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/binding"
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

func TestSlider_BindMin(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	done := make(chan bool)
	slider := NewSlider(0, 100)
	data := &binding.Float64Binding{}
	slider.BindMin(data)
	data.AddListenerFunction(func(binding.Binding) {
		done <- true
	})
	data.Set(75.0)
	select {
	case <-done:
		time.Sleep(time.Millisecond) // Powernap in case our listener runs first
	case <-time.After(time.Second):
		assert.Fail(t, "Timeout")
	}
	assert.Equal(t, 75.0, slider.Min)
}

func TestSlider_BindMax(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	done := make(chan bool)
	slider := NewSlider(0, 100)
	data := &binding.Float64Binding{}
	slider.BindMax(data)
	data.AddListenerFunction(func(binding.Binding) {
		done <- true
	})
	data.Set(75.0)
	select {
	case <-done:
		time.Sleep(time.Millisecond) // Powernap in case our listener runs first
	case <-time.After(time.Second):
		assert.Fail(t, "Timeout")
	}
	assert.Equal(t, 75.0, slider.Max)
}

func TestSlider_BindStep(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	done := make(chan bool)
	slider := NewSlider(0, 100)
	data := &binding.Float64Binding{}
	slider.BindStep(data)
	data.AddListenerFunction(func(binding.Binding) {
		done <- true
	})
	data.Set(75.0)
	select {
	case <-done:
		time.Sleep(time.Millisecond) // Powernap in case our listener runs first
	case <-time.After(time.Second):
		assert.Fail(t, "Timeout")
	}
	assert.Equal(t, 75.0, slider.Step)
}

func TestSlider_BindValue(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	done := make(chan bool)
	slider := NewSlider(0, 100)
	data := &binding.Float64Binding{}
	slider.BindValue(data)
	data.AddListenerFunction(func(binding.Binding) {
		done <- true
	})
	data.Set(75.0)
	select {
	case <-done:
		time.Sleep(time.Millisecond) // Powernap in case our listener runs first
	case <-time.After(time.Second):
		assert.Fail(t, "Timeout")
	}
	assert.Equal(t, 75.0, slider.Value)
}

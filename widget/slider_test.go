package widget

import (
	"testing"

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
	min := 1.0
	data := binding.NewFloat64Ref(&min)
	slider.BindMin(data)
	data.AddListener(binding.NewNotifyFunction(func(binding.Binding) {
		done <- true
	}))
	timedWait(t, done)
	assert.Equal(t, 1.0, slider.Min)

	// Set directly
	min = 25.0
	data.Update()
	timedWait(t, done)
	assert.Equal(t, 25.0, slider.Min)

	// Set by binding
	data.Set(75.0)
	timedWait(t, done)
	assert.Equal(t, 75.0, slider.Min)
}

func TestSlider_BindMax(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	done := make(chan bool)
	slider := NewSlider(0, 100)
	max := 1.0
	data := binding.NewFloat64Ref(&max)
	slider.BindMax(data)
	data.AddListener(binding.NewNotifyFunction(func(binding.Binding) {
		done <- true
	}))
	timedWait(t, done)
	assert.Equal(t, 1.0, slider.Max)

	// Set directly
	max = 25.0
	data.Update()
	timedWait(t, done)
	assert.Equal(t, 25.0, slider.Max)

	// Set by binding
	data.Set(75.0)
	timedWait(t, done)
	assert.Equal(t, 75.0, slider.Max)
}

func TestSlider_BindStep(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	done := make(chan bool)
	slider := NewSlider(0, 100)
	step := 1.0
	data := binding.NewFloat64Ref(&step)
	slider.BindStep(data)
	data.AddListener(binding.NewNotifyFunction(func(binding.Binding) {
		done <- true
	}))
	timedWait(t, done)
	assert.Equal(t, 1.0, slider.Step)

	// Set directly
	step = 25.0
	data.Update()
	timedWait(t, done)
	assert.Equal(t, 25.0, slider.Step)

	// Set by binding
	data.Set(75.0)
	timedWait(t, done)
	assert.Equal(t, 75.0, slider.Step)
}

func TestSlider_BindValue(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	done := make(chan bool)
	slider := NewSlider(0, 100)
	value := 1.0
	data := binding.NewFloat64Ref(&value)
	slider.BindValue(data)
	data.AddListener(binding.NewNotifyFunction(func(binding.Binding) {
		done <- true
	}))
	timedWait(t, done)
	assert.Equal(t, 1.0, slider.Value)

	// Set directly
	value = 25.0
	data.Update()
	timedWait(t, done)
	assert.Equal(t, 25.0, slider.Value)

	// Set by binding
	data.Set(75.0)
	timedWait(t, done)
	assert.Equal(t, 75.0, slider.Value)
}

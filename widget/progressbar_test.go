package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/binding"
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

func TestProgressBar_BindMin(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	done := make(chan bool)
	progressBar := NewProgressBar()
	min := 0.5
	data := binding.NewFloat64Ref(&min)
	progressBar.BindMin(data)
	data.AddListener(binding.NewNotifyFunction(func(binding.Binding) {
		done <- true
	}))
	timedWait(t, done)
	assert.Equal(t, 0.5, progressBar.Min)

	// Set directly
	min = 0.25
	data.Update()
	timedWait(t, done)
	assert.Equal(t, 0.25, progressBar.Min)

	// Set by binding
	data.Set(0.75)
	timedWait(t, done)
	assert.Equal(t, 0.75, progressBar.Min)
}

func TestProgressBar_BindMax(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	done := make(chan bool)
	progressBar := NewProgressBar()
	max := 0.5
	data := binding.NewFloat64Ref(&max)
	progressBar.BindMax(data)
	data.AddListener(binding.NewNotifyFunction(func(binding.Binding) {
		done <- true
	}))
	timedWait(t, done)
	assert.Equal(t, 0.5, progressBar.Max)

	// Set directly
	max = 0.25
	data.Update()
	timedWait(t, done)
	assert.Equal(t, 0.25, progressBar.Max)

	// Set by binding
	data.Set(0.75)
	timedWait(t, done)
	assert.Equal(t, 0.75, progressBar.Max)
}

func TestProgressBar_BindValue(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	done := make(chan bool)
	progressBar := NewProgressBar()
	value := 0.5
	data := binding.NewFloat64Ref(&value)
	progressBar.BindValue(data)
	data.AddListener(binding.NewNotifyFunction(func(binding.Binding) {
		done <- true
	}))
	timedWait(t, done)
	assert.Equal(t, 0.5, progressBar.Value)

	// Set directly
	value = 0.25
	data.Update()
	timedWait(t, done)
	assert.Equal(t, 0.25, progressBar.Value)

	// Set by binding
	data.Set(0.75)
	timedWait(t, done)
	assert.Equal(t, 0.75, progressBar.Value)
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
	withTestTheme(func() {
		render.applyTheme()
		customTextSize = render.label.TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}

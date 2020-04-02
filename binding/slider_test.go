package binding_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/binding"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
)

func TestBindSliderChanged_Binding(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	slider := widget.NewSlider(0, 100)
	data := &binding.Float64Binding{}
	binding.BindSliderChanged(slider, data)
	data.Set(75.0)
	time.Sleep(time.Second)
	assert.Equal(t, 75.0, slider.Value)
}

func TestBindSliderChanged_Event(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	slider := widget.NewSlider(0, 100)
	slider.Resize(fyne.NewSize(100, 100))
	data := &binding.Float64Binding{}
	binding.BindSliderChanged(slider, data)
	value := 0.0
	data.AddFloat64Listener(func(f float64) {
		value = f
	})

	event := &fyne.DragEvent{
		fyne.PointEvent{
			Position: fyne.NewPos(50, 50),
		},
		0,
		0,
	}
	slider.Dragged(event)

	time.Sleep(time.Second)
	assert.Equal(t, 50.0, value)
}

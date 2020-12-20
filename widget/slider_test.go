package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/data/binding"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func TestNewSliderWithData(t *testing.T) {
	val := binding.NewFloat()
	val.Set(4)

	s := NewSliderWithData(0, 10, val)
	waitForBinding()
	assert.Equal(t, 4.0, s.Value)

	s.SetValue(2.0)
	assert.Equal(t, 2.0, val.Get())
}

func TestSlider_Binding(t *testing.T) {
	s := NewSlider(0, 10)
	s.SetValue(2)
	assert.Equal(t, 2.0, s.Value)

	val := binding.NewFloat()
	s.Bind(val)
	waitForBinding()
	assert.Equal(t, 0.0, s.Value)

	val.Set(3)
	waitForBinding()
	assert.Equal(t, 3.0, s.Value)

	s.SetValue(5)
	assert.Equal(t, 5.0, val.Get())

	s.Unbind()
	waitForBinding()
	assert.Equal(t, 5.0, s.Value)
}

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

	test.AssertRendersToMarkup(t, `
		<canvas padded size="220x50">
			<content>
				<widget pos="4,4" size="212x42" type="*widget.Slider">
					<rectangle fillColor="shadow" pos="12,21" size="188x4"/>
					<rectangle fillColor="foreground" pos="12,21" size="0x4"/>
					<circle fillColor="foreground" pos="6,15" size="16x16"/>
				</widget>
			</content>
		</canvas>
	`, w.Canvas())
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

	test.AssertRendersToMarkup(t, `
		<canvas padded size="50x220">
			<content>
				<widget pos="4,4" size="42x212" type="*widget.Slider">
					<rectangle fillColor="shadow" pos="21,12" size="4x188"/>
					<rectangle fillColor="foreground" pos="21,200" size="4x0"/>
					<circle fillColor="foreground" pos="15,194" size="16x16"/>
				</widget>
			</content>
		</canvas>
	`, w.Canvas())
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

	drag := &fyne.DragEvent{Dragged: fyne.NewDelta(10, 2)}
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

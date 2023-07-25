package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestNewSliderWithData(t *testing.T) {
	val := binding.NewFloat()
	err := val.Set(4)
	assert.Nil(t, err)

	s := NewSliderWithData(0, 10, val)
	waitForBinding()
	assert.Equal(t, 4.0, s.Value)

	s.SetValue(2.0)
	f, err := val.Get()
	assert.Nil(t, err)
	assert.Equal(t, 2.0, f)
}

func TestSlider_Binding(t *testing.T) {
	s := NewSlider(0, 10)
	s.SetValue(2)
	assert.Equal(t, 2.0, s.Value)

	val := binding.NewFloat()
	s.Bind(val)
	waitForBinding()
	assert.Equal(t, 0.0, s.Value)

	err := val.Set(3)
	assert.Nil(t, err)
	waitForBinding()
	assert.Equal(t, 3.0, s.Value)

	s.SetValue(5)
	f, err := val.Get()
	assert.Nil(t, err)
	assert.Equal(t, 5.0, f)

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
	wSize := render.slider.Size()
	tSize := render.track.Size()
	aSize := render.active.Size()

	assert.Greater(t, wSize.Width, wSize.Height)

	assert.Equal(t, wSize.Width-slider.endOffset()*2, tSize.Width)
	assert.Equal(t, theme.InputBorderSize()*2, tSize.Height)

	assert.Greater(t, wSize.Width, aSize.Width)
	assert.Equal(t, theme.InputBorderSize()*2, aSize.Height)

	w := test.NewWindow(slider)
	defer w.Close()
	w.Resize(fyne.NewSize(220, 50))

	test.AssertRendersToMarkup(t, "slider/horizontal.xml", w.Canvas())
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
	wSize := render.slider.Size()
	tSize := render.track.Size()
	aSize := render.active.Size()

	assert.Greater(t, wSize.Height, wSize.Width)

	assert.Equal(t, wSize.Height-slider.endOffset()*2, tSize.Height)
	assert.Equal(t, theme.InputBorderSize()*2, tSize.Width)

	assert.Greater(t, wSize.Height, aSize.Height)
	assert.Equal(t, theme.InputBorderSize()*2, aSize.Width)

	w := test.NewWindow(slider)
	defer w.Close()
	w.Resize(fyne.NewSize(50, 220))

	test.AssertRendersToMarkup(t, "slider/vertical.xml", w.Canvas())
}

func TestSlider_OnChanged(t *testing.T) {
	slider := NewSlider(0, 2)
	slider.Resize(slider.MinSize())
	assert.Empty(t, slider.OnChanged)

	changes := 0

	slider.OnChanged = func(_ float64) {
		changes++
	}

	assert.Equal(t, 0, changes)

	slider.SetValue(0.5)
	assert.Equal(t, 0, changes)

	drag := &fyne.DragEvent{}
	drag.PointEvent.Position = fyne.NewPos(25, 2)
	slider.Dragged(drag)
	assert.Equal(t, 1, changes)

	drag.PointEvent.Position = fyne.NewPos(25, 2)
	slider.Dragged(drag)
	assert.Equal(t, 1, changes)

	drag.PointEvent.Position = fyne.NewPos(50, 2)
	slider.Dragged(drag)
	assert.Equal(t, 2, changes)

	tap := &fyne.PointEvent{}
	tap.Position = fyne.NewPos(25, 2)
	slider.Tapped(tap)
	assert.Equal(t, 3, changes)

	tap.Position = fyne.NewPos(25, 2)
	slider.Tapped(tap)
	assert.Equal(t, 3, changes)

	tap.Position = fyne.NewPos(50, 2)
	slider.Tapped(tap)
	assert.Equal(t, 4, changes)

	slider.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
	assert.Equal(t, 5, changes)
}

func TestSlider_OnChangeEnded(t *testing.T) {
	slider := NewSlider(0, 2)
	slider.Resize(slider.MinSize())
	assert.Empty(t, slider.OnChangeEnded)

	changes := 0

	slider.OnChangeEnded = func(_ float64) {
		changes++
	}

	assert.Equal(t, 0, changes)

	slider.SetValue(0.5)
	assert.Equal(t, 0, changes)

	drag := &fyne.DragEvent{}
	drag.PointEvent.Position = fyne.NewPos(25, 2)
	slider.Dragged(drag)
	assert.Equal(t, 0, changes)

	drag.PointEvent.Position = fyne.NewPos(25, 2)
	slider.Dragged(drag)
	assert.Equal(t, 0, changes)

	drag.PointEvent.Position = fyne.NewPos(50, 2)
	slider.Dragged(drag)
	slider.DragEnd()
	assert.Equal(t, 1, changes)

	tap := &fyne.PointEvent{}
	tap.Position = fyne.NewPos(25, 2)
	slider.Tapped(tap)
	assert.Equal(t, 2, changes)

	tap.Position = fyne.NewPos(25, 2)
	slider.Tapped(tap)
	assert.Equal(t, 2, changes)

	tap.Position = fyne.NewPos(50, 2)
	slider.Tapped(tap)
	assert.Equal(t, 3, changes)

	slider.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
	assert.Equal(t, 4, changes)
}

func TestSlider_OnChanged_Float(t *testing.T) {
	slider := NewSlider(0, 1)
	slider.Step = 0.5
	slider.Resize(slider.MinSize())
	assert.Empty(t, slider.OnChanged)

	changes := 0

	slider.OnChanged = func(_ float64) {
		changes++
	}

	assert.Equal(t, 0, changes)

	slider.SetValue(0.15)
	assert.Equal(t, 0, changes)

	drag := &fyne.DragEvent{}
	drag.PointEvent.Position = fyne.NewPos(25, 2)
	slider.Dragged(drag)
	assert.Equal(t, 1, changes)

	drag.PointEvent.Position = fyne.NewPos(25, 2)
	slider.Dragged(drag)
	assert.Equal(t, 1, changes)

	drag.PointEvent.Position = fyne.NewPos(50, 2)
	slider.Dragged(drag)
	assert.Equal(t, 2, changes)

	tap := &fyne.PointEvent{}
	tap.Position = fyne.NewPos(25, 2)
	slider.Tapped(tap)
	assert.Equal(t, 3, changes)

	tap.Position = fyne.NewPos(25, 2)
	slider.Tapped(tap)
	assert.Equal(t, 3, changes)

	tap.Position = fyne.NewPos(50, 2)
	slider.Tapped(tap)
	assert.Equal(t, 4, changes)

	slider.SetValue(0.9)
	assert.Equal(t, 4, changes)
}

func TestSlider_SetValue(t *testing.T) {
	slider := NewSlider(0, 2)
	assert.Equal(t, 0.0, slider.Value)

	slider.SetValue(1.0)
	assert.Equal(t, 1.0, slider.Value)

	slider.SetValue(2.0)
	assert.Equal(t, 2.0, slider.Value)
}

func TestSlider_Focus(t *testing.T) {
	slider := NewSlider(0, 5)

	slider.FocusGained()
	assert.True(t, slider.focused)

	slider.FocusLost()
	assert.False(t, slider.focused)

	slider.MouseIn(nil)
	assert.True(t, slider.hovered)

	slider.MouseOut()
	assert.False(t, slider.hovered)

	left := &fyne.KeyEvent{Name: fyne.KeyLeft}
	right := &fyne.KeyEvent{Name: fyne.KeyRight}

	for i := 1; i <= int(slider.Max); i++ {
		slider.TypedKey(right)
		assert.Equal(t, float64(i), slider.Value)
	}

	slider.TypedKey(right)
	assert.Equal(t, slider.Max, slider.Value)

	for i := 4; i >= int(slider.Min); i-- {
		slider.TypedKey(left)
		assert.Equal(t, float64(i), slider.Value)
	}

	slider.TypedKey(left)
	assert.Equal(t, slider.Min, slider.Value)

	slider.Orientation = Vertical
	up := &fyne.KeyEvent{Name: fyne.KeyUp}
	down := &fyne.KeyEvent{Name: fyne.KeyDown}

	slider.TypedKey(up)
	assert.Equal(t, slider.Min+1, slider.Value)

	slider.TypedKey(down)
	assert.Equal(t, slider.Min, slider.Value)
}

package widget

import (
	"image/color"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

const (
	maxSliderDecimals = uint8(12)
	standardScale     = 6
	minLongSide       = 150
)

var _ fyne.Draggable = (*Slider)(nil)

// Slider if a widget that can slide between two fixed values.
type Slider struct {
	baseWidget

	Value float64
	Min   float64
	Max   float64

	opts SliderOptions
}

// SliderOptions are a set of optional features to slider.
// Vertical  - if the slider should be vertical
// Precision - how many decimal places are relevant
// Drag      - callback slot for user defined behavior
type SliderOptions struct {
	Vertical  bool
	Precision uint8
	Drag      func()
}

// Resize sets a new size for a widget.
func (s *Slider) Resize(size fyne.Size) {
	s.resize(size, s)
}

// MinSize returns the smallest size this widget can shrink to
func (s *Slider) MinSize() fyne.Size {
	return s.minSize(s)
}

// Show this widget, if it was previously hidden
func (s *Slider) Show() {
	s.show(s)
}

// Hide this widget, if it was previously visible
func (s *Slider) Hide() {
	s.hide(s)
}

// Move the widget to a new position, relative to its parent.
func (s *Slider) Move(pos fyne.Position) {
	s.move(pos, s)
}

// DragEnd function.
func (s *Slider) DragEnd() {
}

// Dragged function.
func (s *Slider) Dragged(e *fyne.DragEvent) {
	ok, pos, max := s.wasSliderEvent(&(e.PointEvent))

	if ok {
		// clamp the position for drags that go out of bounds
		if pos > max {
			pos = max
		} else if pos < 0 {
			pos = 0
		}

		s.fireTrigger(pos, max)
	}
}

func (s *Slider) wasSliderEvent(e *fyne.PointEvent) (ok bool, pos, max int) {
	render := Renderer(s).(*sliderRenderer)

	tp := render.thumb.Position()
	ts := render.thumb.Size()

	pad := theme.Padding()

	x := e.Position.X
	y := e.Position.Y

	if s.opts.Vertical {
		if x > (tp.X-pad) && x < (tp.X+ts.Width+pad) {
			// if the cursor was inside the slider area
			return true, y, render.track.Size().Height
		}
	} else {
		if y > (tp.Y-pad) && y < (tp.Y+ts.Height+pad) {
			// if the cursor was inside the slider area
			return true, x, render.track.Size().Width
		}
	}
	return false, 0, 0
}

func (s *Slider) fireTrigger(pos, max int) {
	// update value
	s.updateValue(float64(pos) / float64(max))
	Refresh(s)

	if s.opts.Drag != nil {
		s.opts.Drag()
	}
}

func (s *Slider) updateValue(ratio float64) {
	if s.opts.Vertical {
		ratio = 1 - ratio
	}
	v := s.Min + ratio*(s.Max-s.Min)
	p := math.Pow(10, float64(s.opts.Precision))
	s.Value = float64(int(v*p)) / p
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (s *Slider) CreateRenderer() fyne.WidgetRenderer {
	track := canvas.NewRectangle(theme.DisabledButtonColor())
	active := canvas.NewRectangle(theme.ButtonColor())
	thumb := &canvas.Circle{
		FillColor:   theme.ButtonColor(),
		StrokeWidth: 0}

	objects := []fyne.CanvasObject{track, active, thumb}

	return &sliderRenderer{track, active, thumb, objects, s}
}

type sliderRenderer struct {
	track  *canvas.Rectangle
	active *canvas.Rectangle
	thumb  *canvas.Circle

	objects []fyne.CanvasObject
	slider  *Slider
}

// ApplyTheme is called when the Slider may need to update its look
func (s *sliderRenderer) ApplyTheme() {
	s.track.FillColor = theme.DisabledButtonColor()
	s.thumb.FillColor = theme.ButtonColor()
	s.Refresh()
}

// Refresh is used to update the widget state for drawing
func (s *sliderRenderer) Refresh() {
	s.Layout(s.slider.Size())
	canvas.Refresh(s.slider)
}

// Layout the components of the slider widget
func (s *sliderRenderer) Layout(size fyne.Size) {
	padLen := theme.Padding()
	sideLen := padLen * standardScale
	activeOffset, thumbOffset := s.moveSlide(sideLen)

	var trackPos, activePos, thumbPos fyne.Position
	var trackSize, activeSize fyne.Size

	if s.slider.opts.Vertical {
		trackPos = fyne.NewPos(size.Width/2, 0)
		activePos = fyne.NewPos(trackPos.X, activeOffset)

		trackSize = fyne.NewSize(padLen, size.Height)
		activeSize = fyne.NewSize(padLen, trackSize.Height-activeOffset)

		thumbPos = fyne.NewPos(trackPos.X-(sideLen-trackSize.Width)/2, thumbOffset)
	} else {
		trackPos = fyne.NewPos(0, size.Height/2)
		activePos = trackPos

		trackSize = fyne.NewSize(size.Width, padLen)
		activeSize = fyne.NewSize(activeOffset, padLen)

		thumbPos = fyne.NewPos(thumbOffset, trackPos.Y-(sideLen-trackSize.Height)/2)
	}

	s.track.Move(trackPos)
	s.track.Resize(trackSize)

	s.active.Move(activePos)
	s.active.Resize(activeSize)

	s.thumb.Move(thumbPos)
	s.thumb.Resize(fyne.NewSize(sideLen, sideLen))
}

// MinSize calculates the minimum size of a slider widget
func (s *sliderRenderer) MinSize() fyne.Size {
	s1, s2 := minLongSide, theme.Padding()*standardScale
	if s.slider.opts.Vertical {
		return fyne.NewSize(s2, s1)
	}
	return fyne.NewSize(s1, s2)
}

func (s *sliderRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (s *sliderRenderer) Destroy() {
}

func (s *sliderRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s *sliderRenderer) moveSlide(diameter int) (int, int) {
	w := s.slider
	t := s.track.Size()
	ratio := (w.Value - w.Min) / (w.Max - w.Min)
	if w.opts.Vertical {
		y := float64(t.Height) - (ratio * float64(t.Height))
		return int(y), int(y) - int((1-ratio)*float64(diameter))
	}
	x := ratio * float64(t.Width)
	return int(x), int(x) - int(ratio*float64(diameter))
}

func clampPrecision(p uint8) uint8 {
	if p > maxSliderDecimals {
		return maxSliderDecimals
	}
	return p
}

func checkMinMax(val, min, max float64) (float64, float64) {
	// sort the values to ensure correct order
	if val < min {
		min = val
	}
	if val > max {
		max = val
	}
	if min == max {
		min--
		max++
	}
	if min > max {
		return max, min
	}
	return min, max
}

// NewSlider returns a basic horizontal slider with
// a default precision of zero.
func NewSlider(value, min, max float64) *Slider {
	// sanitize values
	min, max = checkMinMax(value, min, max)
	slider := &Slider{
		baseWidget{},
		value, min, max,
		SliderOptions{false, 0, nil},
	}
	Renderer(slider).Layout(slider.MinSize())
	return slider
}

// NewSliderWithOptions returns a slider with the specified options.
// Options include setting the slider layout to vertical, up to 12
// decimal places of precision, and a callback function to allow custom
// behavior when the slider is dragged.
func NewSliderWithOptions(value, min, max float64, opts SliderOptions) *Slider {
	// sanitize values
	min, max = checkMinMax(value, min, max)
	opts.Precision = clampPrecision(opts.Precision)
	slider := &Slider{baseWidget{}, value, min, max, opts}
	Renderer(slider).Layout(slider.MinSize())
	return slider
}

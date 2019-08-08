package widget

import (
	"image/color"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

const (
	SLIDER_MAX_DECIMALS = uint8(12)
)

var _ fyne.Draggable = (*Slider)(nil)

type SliderOptions struct {
	Vertical  bool
	Precision uint8
	Drag      func()
}

type Slider struct {
	baseWidget

	Value float64
	Min   float64
	Max   float64

	opts SliderOptions
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

func (s *Slider) DragEnd() {
}

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

	hp := render.handle.Position()
	hs := render.handle.Size()

	pad := theme.Padding()

	x := e.Position.X
	y := e.Position.Y

	if s.opts.Vertical {
		if x > (hp.X-pad) && x < (hp.X+hs.Width+pad) {
			// if the cursor was inside the slider area
			return true, y, render.rail.Size().Height
		}
	} else {
		if y > (hp.Y-pad) && y < (hp.Y+hs.Height+pad) {
			// if the cursor was inside the slider area
			return true, x, render.rail.Size().Width
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
	rail := canvas.NewRectangle(theme.ButtonColor())
	fill := canvas.NewRectangle(theme.PrimaryColor())
	handle := &canvas.Circle{
		StrokeColor: color.RGBA{0x00, 0x00, 0x00, 0xff},
		FillColor:   color.RGBA{0x80, 0x80, 0x80, 0xff},
		StrokeWidth: 1}

	objects := []fyne.CanvasObject{rail, fill, handle}

	return &sliderRenderer{rail, fill, handle, objects, s}
}

type sliderRenderer struct {
	rail   *canvas.Rectangle
	fill   *canvas.Rectangle
	handle *canvas.Circle

	objects []fyne.CanvasObject
	slider  *Slider
}

// ApplyTheme is called when the Slider may need to update its look
func (s *sliderRenderer) ApplyTheme() {
	s.rail.FillColor = theme.ButtonColor()
	s.handle.FillColor = theme.PrimaryColor()
	s.Refresh()
}

// Refresh is used to update the widget state for drawing
func (s *sliderRenderer) Refresh() {
	s.Layout(s.slider.Size())
	canvas.Refresh(s.slider)
}

// Layout the components of the slider widget
func (s *sliderRenderer) Layout(size fyne.Size) {
	sq := theme.Padding() * 4
	d1, d2 := s.moveSlide(sq)

	var rp, fp, hp fyne.Position
	var rs, fs fyne.Size

	if s.slider.opts.Vertical {
		rp = fyne.NewPos(size.Width/2, 0)
		fp = fyne.NewPos(rp.X, d1)
		hp = fyne.NewPos(rp.X-theme.Padding(), d2)
		rs = fyne.NewSize(sq/2, size.Height)
		fs = fyne.NewSize(sq/2, rs.Height-d1)
	} else {
		rp = fyne.NewPos(0, size.Height/2)
		fp = rp
		rs = fyne.NewSize(size.Width, sq/2)
		fs = fyne.NewSize(d1, sq/2)
		hp = fyne.NewPos(d2, rp.Y-theme.Padding())
	}

	s.rail.Move(rp)
	s.rail.Resize(rs)

	s.fill.Move(fp)
	s.fill.Resize(fs)

	s.handle.Move(hp)
	s.handle.Resize(fyne.NewSize(sq, sq))
}

// MinSize calculates the minimum size of a slider widget
func (s *sliderRenderer) MinSize() fyne.Size {
	s1, s2 := 100, theme.Padding()*6
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
	r := s.rail.Size()
	ratio := (w.Value - w.Min) / (w.Max - w.Min)
	if w.opts.Vertical {
		y := float64(r.Height) - (ratio * float64(r.Height))
		return int(y), int(y) - int((1-ratio)*float64(diameter))
	}
	x := ratio * float64(r.Width)
	return int(x), int(x) - int(ratio*float64(diameter))
}

func clampPrecision(p uint8) uint8 {
	if p > SLIDER_MAX_DECIMALS {
		return SLIDER_MAX_DECIMALS
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
		min -= 1
		max += 1
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

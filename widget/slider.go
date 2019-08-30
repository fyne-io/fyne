package widget

import (
	"image/color"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

// Orientation controls the horizontal/vertical layout of a widget
type Orientation int

// Orientation constants to control widget layout
const (
	Horizontal Orientation = 0
	Vertical   Orientation = 1
)

var _ fyne.Draggable = (*Slider)(nil)

// Slider if a widget that can slide between two fixed values.
type Slider struct {
	baseWidget

	Value float64
	Min   float64
	Max   float64
	Step  float64

	Orientation Orientation
	OnChanged   func(float64)
}

// NewSlider returns a basic slider.
func NewSlider(min, max float64) *Slider {
	slider := &Slider{
		Value:       0,
		Min:         min,
		Max:         max,
		Step:        1,
		Orientation: Horizontal,
	}
	Renderer(slider).Layout(slider.MinSize())
	return slider
}

// Resize sets a new size for a widget.
func (s *Slider) Resize(size fyne.Size) {
	s.resize(size, s)
}

// MinSize returns the smallest size this widget can be.
func (s *Slider) MinSize() fyne.Size {
	return s.minSize(s)
}

// Show this widget, if it was previously hidden.
func (s *Slider) Show() {
	s.show(s)
}

// Hide this widget, if it was previously visible.
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
	ok, ratio := s.getRatio(&(e.PointEvent))

	if !ok {
		return
	}

	s.updateValue(ratio)
	Refresh(s)

	if s.OnChanged != nil {
		s.OnChanged(s.Value)
	}
}

func (s *Slider) getRatio(e *fyne.PointEvent) (bool, float64) {
	render := Renderer(s).(*sliderRenderer)

	tp := render.thumb.Position()
	ts := render.thumb.Size()

	t := render.track.Size()

	pad := theme.Padding()

	x := e.Position.X
	y := e.Position.Y

	switch s.Orientation {
	case Vertical:
		if x > (tp.X-pad) && x < (tp.X+ts.Width+pad) {
			if y > t.Height {
				return true, 0.0
			} else if y < 0 {
				return true, 1.0
			} else {
				return true, 1 - float64(y)/float64(t.Height)
			}
		}
	case Horizontal:
		if y > (tp.Y-pad) && y < (tp.Y+ts.Height+pad) {
			if x > t.Width {
				return true, 1.0
			} else if x < 0 {
				return true, 0.0
			} else {
				return true, float64(x) / float64(t.Width)
			}
		}
	}

	return false, 0.0
}

func (s *Slider) updateValue(ratio float64) {
	v := s.Min + ratio*(s.Max-s.Min)

	i := -(math.Log10(s.Step))
	p := math.Pow(10, i)

	if v >= s.Max {
		s.Value = s.Max
	} else if v <= s.Min {
		s.Value = s.Min
	} else {
		s.Value = float64(int(v*p)) / p
	}
}

// CreateRenderer links this widget to its renderer.
func (s *Slider) CreateRenderer() fyne.WidgetRenderer {
	track := canvas.NewRectangle(theme.ButtonColor())
	active := canvas.NewRectangle(theme.TextColor())
	thumb := &canvas.Circle{
		FillColor:   theme.TextColor(),
		StrokeWidth: 0}

	objects := []fyne.CanvasObject{track, active, thumb}

	return &sliderRenderer{track, active, thumb, objects, s}
}

const (
	standardScale = 6
	minLongSide   = 50
)

type sliderRenderer struct {
	track  *canvas.Rectangle
	active *canvas.Rectangle
	thumb  *canvas.Circle

	objects []fyne.CanvasObject
	slider  *Slider
}

// ApplyTheme is called when the Slider may need to update its look.
func (s *sliderRenderer) ApplyTheme() {
	s.track.FillColor = theme.ButtonColor()
	s.thumb.FillColor = theme.TextColor()
	s.active.FillColor = theme.TextColor()
	s.Refresh()
}

// Refresh updates the widget state for drawing.
func (s *sliderRenderer) Refresh() {
	s.Layout(s.slider.Size())
	canvas.Refresh(s.slider)
}

// Layout the components of the widget.
func (s *sliderRenderer) Layout(size fyne.Size) {
	pad := theme.Padding()
	diameter := pad * standardScale
	activeOffset, thumbOffset := s.getOffsets(diameter)

	var trackPos, activePos, thumbPos fyne.Position
	var trackSize, activeSize fyne.Size

	switch s.slider.Orientation {
	case Vertical:
		trackPos = fyne.NewPos(size.Width/2, 0)
		activePos = fyne.NewPos(trackPos.X, activeOffset)

		trackSize = fyne.NewSize(pad, size.Height)
		activeSize = fyne.NewSize(pad, trackSize.Height-activeOffset)

		thumbPos = fyne.NewPos(
			trackPos.X-(diameter-trackSize.Width)/2, thumbOffset)
	case Horizontal:
		trackPos = fyne.NewPos(0, size.Height/2)
		activePos = trackPos

		trackSize = fyne.NewSize(size.Width, pad)
		activeSize = fyne.NewSize(activeOffset, pad)

		thumbPos = fyne.NewPos(
			thumbOffset, trackPos.Y-(diameter-trackSize.Height)/2)
	}

	s.track.Move(trackPos)
	s.track.Resize(trackSize)

	s.active.Move(activePos)
	s.active.Resize(activeSize)

	s.thumb.Move(thumbPos)
	s.thumb.Resize(fyne.NewSize(diameter, diameter))
}

// MinSize calculates the minimum size of a widget.
func (s *sliderRenderer) MinSize() fyne.Size {
	s1, s2 := minLongSide, theme.Padding()*standardScale

	switch s.slider.Orientation {
	case Vertical:
		return fyne.NewSize(s2, s1)
	case Horizontal:
		return fyne.NewSize(s1, s2)
	}

	return fyne.Size{Width: 0, Height: 0}
}

func (s *sliderRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (s *sliderRenderer) Destroy() {
}

func (s *sliderRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s *sliderRenderer) getOffsets(diameter int) (int, int) {
	w := s.slider
	t := s.track.Size()
	ratio := (w.Value - w.Min) / (w.Max - w.Min)

	switch w.Orientation {
	case Vertical:
		y := float64(t.Height) - (ratio * float64(t.Height))
		return int(y), int(y) - int((1-ratio)*float64(diameter))
	case Horizontal:
		x := ratio * float64(t.Width)
		return int(x), int(x) - int(ratio*float64(diameter))
	}

	return 0, 0
}

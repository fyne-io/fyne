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
	BaseWidget

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

// DragEnd function.
func (s *Slider) DragEnd() {
}

// Dragged function.
func (s *Slider) Dragged(e *fyne.DragEvent) {
	ratio := s.getRatio(&(e.PointEvent))

	s.updateValue(ratio)
	Refresh(s)

	if s.OnChanged != nil {
		s.OnChanged(s.Value)
	}
}

func (s *Slider) buttonDiameter() int {
	return theme.Padding() * standardScale
}

func (s *Slider) endOffset() int {
	return s.buttonDiameter()/2 + theme.Padding()
}

func (s *Slider) getRatio(e *fyne.PointEvent) float64 {
	pad := s.endOffset()

	x := e.Position.X
	y := e.Position.Y

	switch s.Orientation {
	case Vertical:
		if y > s.size.Height-pad {
			return 0.0
		} else if y < pad {
			return 1.0
		} else {
			return 1 - float64(y-pad)/float64(s.size.Height-pad*2)
		}
	case Horizontal:
		if x > s.size.Width-pad {
			return 1.0
		} else if x < pad {
			return 0.0
		} else {
			return float64(x-pad) / float64(s.size.Width-pad*2)
		}
	}
	return 0.0
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

// MinSize returns the size that this widget should not shrink below
func (s *Slider) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	return s.BaseWidget.MinSize()
}

// CreateRenderer links this widget to its renderer.
func (s *Slider) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	track := canvas.NewRectangle(theme.ButtonColor())
	active := canvas.NewRectangle(theme.TextColor())
	thumb := &canvas.Circle{
		FillColor:   theme.TextColor(),
		StrokeWidth: 0}

	objects := []fyne.CanvasObject{track, active, thumb}

	return &sliderRenderer{track, active, thumb, objects, s}
}

const (
	standardScale = 4
	minLongSide   = 50
)

type sliderRenderer struct {
	track  *canvas.Rectangle
	active *canvas.Rectangle
	thumb  *canvas.Circle

	objects []fyne.CanvasObject
	slider  *Slider
}

// Refresh updates the widget state for drawing.
func (s *sliderRenderer) Refresh() {
	s.track.FillColor = theme.ButtonColor()
	s.thumb.FillColor = theme.TextColor()
	s.active.FillColor = theme.TextColor()

	s.Layout(s.slider.Size())
	canvas.Refresh(s.slider)
}

// Layout the components of the widget.
func (s *sliderRenderer) Layout(size fyne.Size) {
	trackWidth := theme.Padding()
	diameter := s.slider.buttonDiameter()
	activeOffset := s.getOffset()
	endPad := s.slider.endOffset()

	var trackPos, activePos, thumbPos fyne.Position
	var trackSize, activeSize fyne.Size

	switch s.slider.Orientation {
	case Vertical:
		trackPos = fyne.NewPos(size.Width/2, endPad)
		activePos = fyne.NewPos(trackPos.X, activeOffset)

		trackSize = fyne.NewSize(trackWidth, size.Height-endPad*2)
		activeSize = fyne.NewSize(trackWidth, trackSize.Height-activeOffset-endPad)

		thumbPos = fyne.NewPos(
			trackPos.X-(diameter-trackSize.Width)/2, activeOffset-((diameter-theme.Padding())/2))
	case Horizontal:
		trackPos = fyne.NewPos(endPad, size.Height/2)
		activePos = trackPos

		trackSize = fyne.NewSize(size.Width-endPad*2, trackWidth)
		activeSize = fyne.NewSize(activeOffset-endPad, trackWidth)

		thumbPos = fyne.NewPos(
			activeOffset-((diameter-theme.Padding())/2), trackPos.Y-(diameter-trackSize.Height)/2)
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
	s1, s2 := minLongSide, s.slider.buttonDiameter()

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

func (s *sliderRenderer) getOffset() int {
	w := s.slider
	size := s.track.Size()
	ratio := (w.Value - w.Min) / (w.Max - w.Min)
	endPad := s.slider.endOffset()

	switch w.Orientation {
	case Vertical:
		y := int(float64(size.Height)-ratio*float64(size.Height)) + endPad
		return y
	case Horizontal:
		x := int(ratio*float64(size.Width)) + endPad
		return x
	}

	return endPad
}

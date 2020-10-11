package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

var _ fyne.Widget = (*Separator)(nil)

// Separator is a widget for displaying a separator with themeable color.
type Separator struct {
	Base
	colorFn func() color.Color
	minSize fyne.Size
}

// NewSeparator creates a separator with the given min size and the color delivered by the given callback.
func NewSeparator(minSize fyne.Size, colorFn func() color.Color) *Separator {
	return &Separator{minSize: minSize, colorFn: colorFn}
}

// CreateRenderer returns a new renderer for the separator.
//
// Implements: fyne.Widget
func (s *Separator) CreateRenderer() fyne.WidgetRenderer {
	bar := canvas.NewRectangle(theme.DisabledTextColor())
	objects := []fyne.CanvasObject{bar}
	return &separatorRenderer{
		BaseRenderer: NewBaseRenderer(objects),
		bar:          bar,
		s:            s,
	}
}

// Hide hides the separator.
//
// Implements: fyne.Widget
func (s *Separator) Hide() {
	HideWidget(&s.Base, s)
}

// MinSize returns the minimal size of the separator.
//
// Implements: fyne.Widget
func (s *Separator) MinSize() fyne.Size {
	return MinSizeOf(s)
}

// Move sets the position of the separator relative to its parent.
//
// Implements: fyne.Widget
func (s *Separator) Move(pos fyne.Position) {
	MoveWidget(&s.Base, s, pos)
}

// Refresh triggers a redraw of the separator.
//
// Implements: fyne.Widget
func (s *Separator) Refresh() {
	RefreshWidget(s)
}

// Resize changes the size of the separator.
//
// Implements: fyne.Widget
func (s *Separator) Resize(size fyne.Size) {
	ResizeWidget(&s.Base, s, size)
}

// Show makes the separator visible.
//
// Implements: fyne.Widget
func (s *Separator) Show() {
	ShowWidget(&s.Base, s)
}

var _ fyne.WidgetRenderer = (*separatorRenderer)(nil)

type separatorRenderer struct {
	BaseRenderer
	bar *canvas.Rectangle
	s   *Separator
}

func (r *separatorRenderer) Layout(size fyne.Size) {
	r.bar.Resize(size)
}

func (r *separatorRenderer) MinSize() fyne.Size {
	return r.s.minSize
}

func (r *separatorRenderer) Refresh() {
	r.bar.FillColor = r.s.colorFn()
	canvas.Refresh(r.s)
}

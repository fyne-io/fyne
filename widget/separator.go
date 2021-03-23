package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Widget = (*Separator)(nil)

// Separator is a widget for displaying a separator with themeable color.
//
// Since: 1.4
type Separator struct {
	widget.Base
}

// NewSeparator creates a new separator.
//
// Since: 1.4
func NewSeparator() *Separator {
	return &Separator{}
}

// CreateRenderer returns a new renderer for the separator.
//
// Implements: fyne.Widget
func (s *Separator) CreateRenderer() fyne.WidgetRenderer {
	bar := canvas.NewRectangle(theme.DisabledColor())
	objects := []fyne.CanvasObject{bar}
	return &separatorRenderer{
		BaseRenderer: widget.NewBaseRenderer(objects),
		bar:          bar,
		d:            s,
	}
}

// Hide hides the separator.
//
// Implements: fyne.Widget
func (s *Separator) Hide() {
	widget.HideWidget(&s.Base, s)
}

// MinSize returns the minimal size of the separator.
//
// Implements: fyne.Widget
func (s *Separator) MinSize() fyne.Size {
	t := theme.SeparatorThicknessSize()
	return fyne.NewSize(t, t)
}

// Move sets the position of the separator relative to its parent.
//
// Implements: fyne.Widget
func (s *Separator) Move(pos fyne.Position) {
	widget.MoveWidget(&s.Base, s, pos)
}

// Refresh triggers a redraw of the separator.
//
// Implements: fyne.Widget
func (s *Separator) Refresh() {
	widget.RefreshWidget(s)
}

// Resize changes the size of the separator.
//
// Implements: fyne.Widget
func (s *Separator) Resize(size fyne.Size) {
	widget.ResizeWidget(&s.Base, s, size)
}

// Show makes the separator visible.
//
// Implements: fyne.Widget
func (s *Separator) Show() {
	widget.ShowWidget(&s.Base, s)
}

var _ fyne.WidgetRenderer = (*separatorRenderer)(nil)

type separatorRenderer struct {
	widget.BaseRenderer
	bar *canvas.Rectangle
	d   *Separator
}

func (r *separatorRenderer) Layout(size fyne.Size) {
	r.bar.Resize(size)
}

func (r *separatorRenderer) MinSize() fyne.Size {
	t := theme.SeparatorThicknessSize()
	return fyne.NewSize(t, t)
}

func (r *separatorRenderer) Refresh() {
	r.bar.FillColor = theme.DisabledColor()
	canvas.Refresh(r.d)
}

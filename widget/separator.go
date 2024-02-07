package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Widget = (*Separator)(nil)

// Separator is a widget for displaying a separator with themeable color.
//
// Since: 1.4
type Separator struct {
	BaseWidget
}

// NewSeparator creates a new separator.
//
// Since: 1.4
func NewSeparator() *Separator {
	s := &Separator{}
	s.ExtendBaseWidget(s)
	return s
}

// CreateRenderer returns a new renderer for the separator.
//
// Implements: fyne.Widget
func (s *Separator) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	th := s.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	bar := canvas.NewRectangle(th.Color(theme.ColorNameSeparator, v))
	return &separatorRenderer{
		WidgetRenderer: NewSimpleRenderer(bar),
		bar:            bar,
		d:              s,
	}
}

// MinSize returns the minimal size of the separator.
//
// Implements: fyne.Widget
func (s *Separator) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	t := s.Theme().Size(theme.SizeNameSeparatorThickness)
	return fyne.NewSize(t, t)
}

var _ fyne.WidgetRenderer = (*separatorRenderer)(nil)

type separatorRenderer struct {
	fyne.WidgetRenderer
	bar *canvas.Rectangle
	d   *Separator
}

func (r *separatorRenderer) MinSize() fyne.Size {
	t := r.d.Theme().Size(theme.SizeNameSeparatorThickness)
	return fyne.NewSize(t, t)
}

func (r *separatorRenderer) Refresh() {
	th := r.d.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	r.bar.FillColor = th.Color(theme.ColorNameSeparator, v)
	canvas.Refresh(r.d)
}

package widget

import (
	"image/color"

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

	invert bool
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
func (s *Separator) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	th := s.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	var col color.Color
	if s.invert {
		col = th.Color(theme.ColorNameForeground, v)
	} else {
		col = th.Color(theme.ColorNameSeparator, v)
	}
	bar := canvas.NewRectangle(col)

	return &separatorRenderer{
		WidgetRenderer: NewSimpleRenderer(bar),
		bar:            bar,
		d:              s,
	}
}

// MinSize returns the minimal size of the separator.
func (s *Separator) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	return s.BaseWidget.MinSize()
}

var _ fyne.WidgetRenderer = (*separatorRenderer)(nil)

type separatorRenderer struct {
	fyne.WidgetRenderer
	bar *canvas.Rectangle
	d   *Separator
}

func (r *separatorRenderer) MinSize() fyne.Size {
	return fyne.NewSquareSize(r.d.Theme().Size(theme.SizeNameSeparatorThickness))
}

func (r *separatorRenderer) Refresh() {
	th := r.d.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	if r.d.invert {
		r.bar.FillColor = th.Color(theme.ColorNameForeground, v)
	} else {
		r.bar.FillColor = th.Color(theme.ColorNameSeparator, v)
	}
	canvas.Refresh(r.d)
}

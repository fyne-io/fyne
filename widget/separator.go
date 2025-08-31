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
//
// Implements: fyne.Widget
func (s *Separator) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	var col color.Color
	if s.invert {
		col = theme.ColorForWidget(theme.ColorNameForeground, s)
	} else {
		col = theme.ColorForWidget(theme.ColorNameSeparator, s)
	}
	bar := canvas.NewRectangle(col)

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
	if r.d.invert {
		r.bar.FillColor = theme.ColorForWidget(theme.ColorNameForeground, r.d)
	} else {
		r.bar.FillColor = theme.ColorForWidget(theme.ColorNameSeparator, r.d)
	}
	canvas.Refresh(r.d)
}

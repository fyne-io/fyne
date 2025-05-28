package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type rowWrapLayout struct {
	rowCount          int
	horizontalPadding float32
	verticalPadding   float32
}

// NewRowWrapLayout returns a layout that dynamically arranges objects of similar height
// in rows and wraps them as necessary.
// The height of the rows is determined by the tallest object and the same for all rows.
//
// Since: 2.7
func NewRowWrapLayout() fyne.Layout {
	p := theme.Padding()
	return &rowWrapLayout{
		horizontalPadding: p,
		verticalPadding:   p,
	}
}

// NewRowWrapLayoutWithCustomPadding creates a new RowWrapLayout instance
// with the specified paddings.
//
// Since: 2.7
func NewRowWrapLayoutWithCustomPadding(horizontal, vertical float32) fyne.Layout {
	return &rowWrapLayout{
		horizontalPadding: horizontal,
		verticalPadding:   vertical,
	}
}

var _ fyne.Layout = (*rowWrapLayout)(nil)

// MinSize finds the smallest size that satisfies all the child objects.
// For a RowWrapLayout this is the width of the widest child
// and the height of the tallest child multiplied by the number of children,
// with appropriate padding between them.
func (l *rowWrapLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) == 0 {
		return fyne.NewSize(0, 0)
	}
	rows := l.rowCount
	if rows == 0 {
		rows = 1
	}
	var w, h float32
	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		s := o.MinSize()
		w = fyne.Max(w, s.Width)
		h = fyne.Max(h, s.Height)
	}
	return fyne.NewSize(w, h*float32(rows)+l.verticalPadding*float32(rows-1))
}

// Layout is called to pack all child objects into a specified size.
// For RowWrapLayout this will arrange all objects into rows of equal size.
func (l *rowWrapLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) == 0 {
		return
	}
	var h float32
	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		h = fyne.Max(h, o.MinSize().Height)
	}
	var minSize fyne.Size
	pos := fyne.NewPos(0, 0)
	rows := 1
	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		size := o.MinSize()
		o.Resize(size)
		w := size.Width + l.horizontalPadding
		if pos.X+w > containerSize.Width {
			y := float32(rows) * (h + l.verticalPadding)
			pos = fyne.NewPos(0, y)
			minSize.Height = fyne.Max(minSize.Height, y)
			rows++
		}
		o.Move(pos)
		pos = pos.Add(fyne.NewPos(w, 0))
		minSize.Width = fyne.Max(minSize.Width, w)
	}
	l.rowCount = rows
}

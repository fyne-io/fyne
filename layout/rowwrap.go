package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type rowWrapLayout struct {
	horizontalPadding float32
	minSize           fyne.Size
	verticalPadding   float32
}

// NewRowWrapLayout returns a layout that dynamically arranges objects of similar height
// in rows and wraps them dynamically.
// Objects are separated with horizontal and vertical padding.
//
// Since: 2.7
func NewRowWrapLayout() fyne.Layout {
	p := theme.Padding()
	return &rowWrapLayout{
		horizontalPadding: p,
		verticalPadding:   p,
	}
}

// NewRowWrapLayoutWithCustomPadding returns a new RowWrapLayout instance
// with custom horizontal and inner padding.
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
// For a RowWrapLayout this is initially the width of the widest child
// and the height of the tallest child multiplied by the number of children,
// with appropriate padding between them.
// After Layout() has run it returns the actual min size.
func (l *rowWrapLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) == 0 {
		return fyne.NewSize(0, 0)
	}
	if !l.minSize.IsZero() {
		return l.minSize
	}
	var maxW, maxH float32
	var objCount int
	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		objCount++
		s := o.MinSize()
		maxW = fyne.Max(maxW, s.Width)
		maxH = fyne.Max(maxH, s.Height)
	}
	return fyne.NewSize(maxW, l.minHeight(maxH, objCount))
}

func (l *rowWrapLayout) minHeight(rowHeight float32, rowCount int) float32 {
	return rowHeight*float32(rowCount) + l.verticalPadding*float32(rowCount-1)
}

// Layout is called to pack all child objects into a specified size.
// For RowWrapLayout this will arrange all objects into rows of equal size
// and wrap objects into additional rows as needed.
func (l *rowWrapLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) == 0 {
		return
	}
	var maxH float32
	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		maxH = fyne.Max(maxH, o.MinSize().Height)
	}
	var minSize fyne.Size
	pos := fyne.NewPos(0, 0)
	rows := 1
	isFirst := true
	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		size := o.MinSize()
		o.Resize(size)
		if !isFirst && pos.X+size.Width+l.horizontalPadding >= containerSize.Width {
			y := float32(rows) * (maxH + l.verticalPadding)
			pos = fyne.NewPos(0, y)
			rows++
		}
		isFirst = false
		minSize.Width = fyne.Max(minSize.Width, pos.X+size.Width)
		minSize.Height = l.minHeight(maxH, rows)
		o.Move(pos)
		pos = pos.Add(fyne.NewPos(size.Width+l.horizontalPadding, 0))
	}
	l.minSize = minSize
}

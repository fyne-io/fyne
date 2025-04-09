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

// NewRowWrapLayout returns a layout that dynamically arranges objects
// with the same height in rows and wraps them as necessary.
//
// Object visibility is supported.
//
// Since: 2.7
func NewRowWrapLayout() fyne.Layout {
	return &rowWrapLayout{
		horizontalPadding: theme.Padding(),
		verticalPadding:   theme.Padding(),
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
// and the height of the first child multiplied by the number of children,
// with appropriate padding between them.
func (l *rowWrapLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) == 0 {
		return fyne.NewSize(0, 0)
	}
	rows := l.rowCount
	if rows == 0 {
		rows = 1
	}
	rowHeight := objects[0].MinSize().Height
	var width float32
	for _, o := range objects {
		size := o.MinSize()
		if size.Width > width {
			width = size.Width
		}
	}
	s := fyne.NewSize(width, rowHeight*float32(rows)+l.verticalPadding*float32(rows-1))
	return s
}

// Layout is called to pack all child objects into a specified size.
// For RowWrapLayout this will arrange all objects into rows of equal size.
func (l *rowWrapLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) == 0 {
		return
	}
	rowHeight := objects[0].MinSize().Height
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
			pos = fyne.NewPos(0, float32(rows)*(rowHeight+l.verticalPadding))
			rows++
		}
		o.Move(pos)
		pos = pos.Add(fyne.NewPos(w, 0))
	}
	l.rowCount = rows
}

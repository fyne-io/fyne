package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type rowWrapLayout struct {
	rowCount int
}

// NewRowWrapLayout returns a layout that dynamically arranges objects
// with the same height in rows and wraps them as necessary.
//
// Object visibility is supported.
//
// Since: 2.7
func NewRowWrapLayout() fyne.Layout {
	return &rowWrapLayout{}
}

var _ fyne.Layout = (*rowWrapLayout)(nil)

func (l *rowWrapLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) == 0 {
		return fyne.NewSize(0, 0)
	}
	rows := l.rowCount
	if rows == 0 {
		rows = 1
	}
	rowHeight := objects[0].MinSize().Height
	var w float32
	for _, o := range objects {
		size := o.MinSize()
		if size.Width > w {
			w = size.Width
		}
	}
	s := fyne.NewSize(w, rowHeight*float32(rows)+theme.Padding()*float32(rows-1))
	return s
}

func (l *rowWrapLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) == 0 {
		return
	}
	padding := theme.Padding()
	rowHeight := objects[0].MinSize().Height
	pos := fyne.NewPos(0, 0)
	rows := 1
	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		size := o.MinSize()
		o.Resize(size)
		w := size.Width + padding
		if pos.X+w > containerSize.Width {
			pos = fyne.NewPos(0, float32(rows)*(rowHeight+padding))
			rows++
		}
		o.Move(pos)
		pos = pos.Add(fyne.NewPos(w, 0))
	}
	l.rowCount = rows
}

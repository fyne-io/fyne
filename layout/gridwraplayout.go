package layout

import (
	"math"

	"fyne.io/fyne/v2"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*gridWrapLayout)(nil)

type gridWrapLayout struct {
	BaseLayout
	CellSize fyne.Size
	colCount int
	rowCount int
}

// NewGridWrapLayout returns a new GridWrapLayout instance
func NewGridWrapLayout(size fyne.Size, options ...LayoutOption) fyne.Layout {
	l := &gridWrapLayout{CellSize: size, colCount: 1, rowCount: 1}
	for _, option := range options {
		option(l)
	}
	return l
}

// Layout is called to pack all child objects into a specified size.
// For a GridWrapLayout this will attempt to lay all the child objects in a row
// and wrap to a new row if the size is not large enough.
func (g *gridWrapLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	topPadding, _, leftPadding, _ := g.GetPaddings()
	g.colCount = 1
	g.rowCount = 0

	if size.Width > g.CellSize.Width {
		g.colCount = int(math.Floor(float64(size.Width+leftPadding) / float64(g.CellSize.Width+leftPadding)))
	}

	i, x, y := 0, float32(0), float32(0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if i%g.colCount == 0 {
			g.rowCount++
		}

		child.Move(fyne.NewPos(x, y))
		child.Resize(g.CellSize)

		if (i+1)%g.colCount == 0 {
			x = 0
			y += g.CellSize.Height + topPadding
		} else {
			x += g.CellSize.Width + leftPadding
		}
		i++
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a GridWrapLayout this is simply the specified cellsize as a single column
// layout has no padding. The returned size does not take into account the number
// of columns as this layout re-flows dynamically.
func (g *gridWrapLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	rows := g.rowCount
	if rows < 1 {
		rows = 1
	}
	topPadding, _, _, _ := g.GetPaddings()
	return fyne.NewSize(g.CellSize.Width,
		(g.CellSize.Height*float32(rows))+(float32(rows-1)*topPadding))
}

package layout

import (
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

type fixedGridLayout struct {
	CellSize fyne.Size
}

// Layout is called to pack all child objects into a specified size.
// For a FixedGridLayout this will attempt to lay all the child objects in a row
// and wrap to a new row if the size is not large enough.
func (g *fixedGridLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	cols := 1
	if size.Width > g.CellSize.Width {
		cols = int(math.Floor(float64(size.Width+theme.Padding()) / float64(g.CellSize.Width+theme.Padding())))
	}

	x, y := 0, 0
	for i, child := range objects {
		child.Move(fyne.NewPos(x, y))
		child.Resize(g.CellSize)

		if (i+1)%cols == 0 {
			x = 0
			y += g.CellSize.Height + theme.Padding()
		} else {
			x += g.CellSize.Width + theme.Padding()
		}
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a FixedGridLayout this is simply the specified cellsize as a single column
// layout has no padding. The returned sie does not take into account the number
// of child objects as this layout re-flows dynamically.
func (g *fixedGridLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return g.CellSize
}

// NewFixedGridLayout returns a new FixedGridLayout instance
func NewFixedGridLayout(size fyne.Size) fyne.Layout {
	return &fixedGridLayout{size}
}

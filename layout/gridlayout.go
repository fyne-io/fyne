package layout

import (
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*gridLayout)(nil)

type gridLayout struct {
	Cols int
}

func (g *gridLayout) countRows(objects []fyne.CanvasObject) int {
	count := 0
	for _, child := range objects {
		if child.Visible() {
			count++
		}
	}

	return int(math.Ceil(float64(count) / float64(g.Cols)))
}

// Get the leading (top or left) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func getLeading(size float64, offset int) int {
	ret := (size + float64(theme.Padding())) * float64(offset)

	return int(math.Round(ret))
}

// Get the trailing (bottom or right) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func getTrailing(size float64, offset int) int {
	return getLeading(size, offset+1) - theme.Padding()
}

// Layout is called to pack all child objects into a specified size.
// For a GridLayout this will pack objects into a table format with the number
// of columns specified in our constructor.
func (g *gridLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	rows := g.countRows(objects)

	padWidth := (g.Cols - 1) * theme.Padding()
	padHeight := (rows - 1) * theme.Padding()

	cellWidth := float64(size.Width-padWidth) / float64(g.Cols)
	cellHeight := float64(size.Height-padHeight) / float64(rows)

	row, col := 0, 0
	i := 0
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		x1 := getLeading(cellWidth, col)
		y1 := getLeading(cellHeight, row)
		x2 := getTrailing(cellWidth, col)
		y2 := getTrailing(cellHeight, row)

		child.Move(fyne.NewPos(x1, y1))
		child.Resize(fyne.NewSize(x2-x1, y2-y1))

		if (i+1)%g.Cols == 0 {
			row++
			col = 0
		} else {
			col++
		}
		i++
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a GridLayout this is the size of the largest child object multiplied by
// the required number of columns and rows, with appropriate padding between
// children.
func (g *gridLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	rows := g.countRows(objects)
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Union(child.MinSize())
	}

	minContentSize := fyne.NewSize(minSize.Width*g.Cols, minSize.Height*rows)
	return minContentSize.Add(fyne.NewSize(theme.Padding()*(g.Cols-1), theme.Padding()*(rows-1)))
}

// NewGridLayout returns a new GridLayout instance
func NewGridLayout(cols int) fyne.Layout {
	return &gridLayout{cols}
}

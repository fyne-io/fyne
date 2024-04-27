package layout

import (
	"math"

	"fyne.io/fyne/v2"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*gridLayout)(nil)

type gridLayout struct {
	BaseLayout

	Cols            int
	vertical, adapt bool
}

// NewAdaptiveGridLayout returns a new grid layout which uses columns when horizontal but rows when vertical.
func NewAdaptiveGridLayout(rowcols int, options ...LayoutOption) fyne.Layout {
	l := &gridLayout{Cols: rowcols, adapt: true}
	for _, option := range options {
		option(l)
	}
	return l
}

// NewGridLayout returns a grid layout arranged in a specified number of columns.
// The number of rows will depend on how many children are in the container that uses this layout.
func NewGridLayout(cols int, options ...LayoutOption) fyne.Layout {
	return NewGridLayoutWithColumns(cols, options...)
}

// NewGridLayoutWithColumns returns a new grid layout that specifies a column count and wrap to new rows when needed.
func NewGridLayoutWithColumns(cols int, options ...LayoutOption) fyne.Layout {
	l := &gridLayout{Cols: cols}
	for _, option := range options {
		option(l)
	}
	return l
}

// NewGridLayoutWithRows returns a new grid layout that specifies a row count that creates new rows as required.
func NewGridLayoutWithRows(rows int, options ...LayoutOption) fyne.Layout {
	l := &gridLayout{Cols: rows, vertical: true}
	for _, option := range options {
		option(l)
	}
	return l
}

func (g *gridLayout) horizontal() bool {
	if g.adapt {
		return fyne.IsHorizontal(fyne.CurrentDevice().Orientation())
	}

	return !g.vertical
}

func (g *gridLayout) countRows(objects []fyne.CanvasObject) int {
	if g.Cols < 1 {
		g.Cols = 1
	}
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
func getLeading(size float64, offset int, padding float32) float32 {
	ret := (size + float64(padding)) * float64(offset)
	return float32(ret)
}

// Get the trailing (bottom or right) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func getTrailing(size float64, offset int, leadingPadding, padding float32) float32 {
	return getLeading(size, offset+1, leadingPadding) - padding
}

// Layout is called to pack all child objects into a specified size.
// For a GridLayout this will pack objects into a table format with the number
// of columns specified in our constructor.
func (g *gridLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	rows := g.countRows(objects)

	topPadding, bottomPadding, leftPadding, rightPadding := g.GetPaddings()

	primaryObjects := rows
	secondaryObjects := g.Cols
	if g.horizontal() {
		primaryObjects, secondaryObjects = secondaryObjects, primaryObjects
	}

	padWidth := float32(primaryObjects-1) * leftPadding
	padHeight := float32(secondaryObjects-1) * topPadding
	cellWidth := float64(size.Width-padWidth) / float64(primaryObjects)
	cellHeight := float64(size.Height-padHeight) / float64(secondaryObjects)

	row, col := 0, 0
	i := 0
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		x1 := getLeading(cellWidth, col, leftPadding)
		y1 := getLeading(cellHeight, row, topPadding)
		x2 := getTrailing(cellWidth, col, leftPadding, rightPadding)
		y2 := getTrailing(cellHeight, row, topPadding, bottomPadding)

		child.Move(fyne.NewPos(x1, y1))
		child.Resize(fyne.NewSize(x2-x1, y2-y1))

		if g.horizontal() {
			if (i+1)%g.Cols == 0 {
				row++
				col = 0
			} else {
				col++
			}
		} else {
			if (i+1)%g.Cols == 0 {
				col++
				row = 0
			} else {
				row++
			}
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

		minSize = minSize.Max(child.MinSize())
	}

	topPadding, _, leftPadding, _ := g.GetPaddings()

	primaryObjects := rows
	secondaryObjects := g.Cols
	if g.horizontal() {
		primaryObjects, secondaryObjects = secondaryObjects, primaryObjects
	}

	width := minSize.Width * float32(primaryObjects)
	height := minSize.Height * float32(secondaryObjects)
	xpad := leftPadding * fyne.Max(float32(primaryObjects-1), 0)
	ypad := topPadding * fyne.Max(float32(secondaryObjects-1), 0)

	return fyne.NewSize(width+xpad, height+ypad)
}

package layout

import "math"
import "github.com/fyne-io/fyne/api/ui"
import "github.com/fyne-io/fyne/api/ui/theme"

type gridLayout struct {
	Cols int
}

func (g *gridLayout) countRows(objects []ui.CanvasObject) int {
	return int(math.Ceil(float64(len(objects)) / float64(g.Cols)))
}

// Get the leading (top or left) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row it's on.
func getLeading(size float64, offset int) int {
	ret := (size + float64(theme.Padding())) * float64(offset)

	return int(math.Round(ret))
}

// Get theh trailing (bottom or right) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row it's on.
func getTrailing(size float64, offset int) int {
	return getLeading(size, offset+1) - theme.Padding()
}

// Layout is called to pack all child objects into a specified size.
// For a GridLayout this will pack objects into a table format with the number
// of columns specified in our constructor.
func (g *gridLayout) Layout(objects []ui.CanvasObject, size ui.Size) {
	rows := g.countRows(objects)

	padWidth := (g.Cols - 1) * theme.Padding()
	padHeight := (rows - 1) * theme.Padding()

	cellWidth := float64(size.Width-padWidth) / float64(g.Cols)
	cellHeight := float64(size.Height-padHeight) / float64(rows)

	row, col := 0, 0
	for i, child := range objects {
		x1 := getLeading(cellWidth, col)
		y1 := getLeading(cellHeight, row)
		x2 := getTrailing(cellWidth, col)
		y2 := getTrailing(cellHeight, row)

		child.Move(ui.NewPos(x1, y1))
		child.Resize(ui.NewSize(x2-x1, y2-y1))

		if (i+1)%g.Cols == 0 {
			row++
			col = 0

		} else {
			col++

		}
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a GridLayout this is the size of the largest child object multiplied by
// the required number of columns and rows, with appropriate padding between
// children.
func (g *gridLayout) MinSize(objects []ui.CanvasObject) ui.Size {
	rows := g.countRows(objects)
	minSize := ui.NewSize(0, 0)
	for _, child := range objects {
		minSize = minSize.Union(child.MinSize())
	}

	minContentSize := ui.NewSize(minSize.Width*g.Cols, minSize.Height*rows)
	return minContentSize.Add(ui.NewSize(theme.Padding()*(g.Cols-1), theme.Padding()*(rows-1)))
}

// NewGridLayout returns a new GridLayout instance
func NewGridLayout(cols int) ui.Layout {
	return &gridLayout{cols}
}

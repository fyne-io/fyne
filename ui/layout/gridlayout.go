package layout

import "math"
import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/theme"

type gridLayout struct {
	Cols int
}

func (g *gridLayout) countRows(objects []ui.CanvasObject) int {
	return int(math.Ceil(float64(len(objects)) / float64(g.Cols)))
}

// Layout is called to pack all child objects into a specified size.
// For a GridLayout this will pack objects into a table format with the number
// of columns specified in our constructor.
func (g *gridLayout) Layout(objects []ui.CanvasObject, size ui.Size) {
	rows := g.countRows(objects)

	padWidth := (g.Cols - 1) * theme.Padding()
	padHeight := (rows - 1) * theme.Padding()

	cellWidth := int((size.Width - padWidth) / g.Cols)
	cellHeight := int((size.Height - padHeight) / rows)
	cellSize := ui.NewSize(cellWidth, cellHeight)

	x, y := 0, 0
	for i, child := range objects {
		child.Move(ui.NewPos(x, y))
		child.Resize(cellSize)

		if (i+1)%g.Cols == 0 {
			x = 0
			y += cellHeight + theme.Padding()
		} else {
			x += cellWidth + theme.Padding()
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
func NewGridLayout(cols int) *gridLayout {
	return &gridLayout{cols}
}

package layout

import "math"

import "github.com/fyne-io/fyne/ui"

type gridLayout struct {
	Cols int
}

func (g *gridLayout) Layout(c *ui.Container, size ui.Size) {
	rows := int(math.Ceil(float64(len(c.Objects)) / float64(g.Cols)))

	cellWidth := int(size.Width / g.Cols)
	cellHeight := int(size.Height / rows)
	cellSize := ui.NewSize(cellWidth, cellHeight)

	x, y := 0, 0
	for i, child := range c.Objects {
		child.Move(ui.NewPos(x, y))
		child.Resize(cellSize)

		if i+1%g.Cols == 0 {
			x = 0
			y += cellHeight
		} else {
			x += cellWidth
		}
	}
}

func NewGridLayout(cols int) *gridLayout {
	return &gridLayout{cols}
}

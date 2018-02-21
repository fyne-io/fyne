package layout

import "math"
import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/theme"

type gridLayout struct {
	Cols int
}

func (g *gridLayout) Layout(objects []ui.CanvasObject, size ui.Size) {
	rows := int(math.Ceil(float64(len(objects)) / float64(g.Cols)))

	padWidth := (g.Cols + 1) * theme.Padding()
	padHeight := (rows + 1) * theme.Padding()

	cellWidth := int((size.Width - padWidth) / g.Cols)
	cellHeight := int((size.Height - padHeight) / rows)
	cellSize := ui.NewSize(cellWidth, cellHeight)

	x, y := theme.Padding(), theme.Padding()
	for i, child := range objects {
		child.Move(ui.NewPos(x, y))
		child.Resize(cellSize)

		if (i+1)%g.Cols == 0 {
			x = theme.Padding()
			y += cellHeight + theme.Padding()
		} else {
			x += cellWidth + theme.Padding()
		}
	}
}

func NewGridLayout(cols int) *gridLayout {
	return &gridLayout{cols}
}

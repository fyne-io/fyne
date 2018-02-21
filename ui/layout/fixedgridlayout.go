package layout

import "math"
import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/theme"

type fixedGridLayout struct {
	CellSize ui.Size
}

func (g *fixedGridLayout) Layout(objects []ui.CanvasObject, size ui.Size) {
	cols := int(math.Floor(float64(size.Width - theme.Padding()) / float64(g.CellSize.Width + theme.Padding())))

	x, y := theme.Padding(), theme.Padding()
	for i, child := range objects {
		child.Move(ui.NewPos(x, y))
		child.Resize(g.CellSize)

		if (i+1)%cols == 0 {
			x = theme.Padding()
			y += g.CellSize.Height + theme.Padding()
		} else {
			x += g.CellSize.Width + theme.Padding()
		}
	}
}

func NewFixedGridLayout(size ui.Size) *fixedGridLayout {
	return &fixedGridLayout{size}
}

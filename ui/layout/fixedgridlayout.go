package layout

import "math"
import "github.com/fyne-io/fyne/ui"

type fixedGridLayout struct {
	Size ui.Size
}

func (g *fixedGridLayout) Layout(c *ui.Container, size ui.Size) {
	cols := int(math.Floor(float64(size.Width) / float64(g.Size.Width)))

	x, y := 0, 0
	for i, child := range c.Objects {
		child.Move(ui.NewPos(x, y))
		child.Resize(g.Size)

		if (i+1)%cols == 0 {
			x = 0
			y += g.Size.Height
		} else {
			x += g.Size.Width
		}
	}
}

func NewFixedGridLayout(size ui.Size) *fixedGridLayout {
	return &fixedGridLayout{size}
}

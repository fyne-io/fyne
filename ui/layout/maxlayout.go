package layout

import "github.com/fyne-io/fyne/ui"

func MaxLayout(c *ui.Container, size ui.Size) {
	for _, child := range c.Objects {
		child.Resize(size)
	}
}

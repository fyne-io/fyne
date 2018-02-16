package layout

import "github.com/fyne-io/fyne/ui"

type maxLayout struct {
}

func (m *maxLayout) Layout(c *ui.Container, size ui.Size) {
	for _, child := range c.Objects {
		child.Resize(size)
	}
}

func NewMaxLayout() *maxLayout {
	return &maxLayout{}
}

package layout

import "github.com/fyne-io/fyne/ui"

type maxLayout struct {
}

func (m *maxLayout) Layout(objects []ui.CanvasObject, size ui.Size) {
	for _, child := range objects {
		child.Resize(size)
	}
}

func NewMaxLayout() *maxLayout {
	return &maxLayout{}
}

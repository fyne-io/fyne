package layout

import "github.com/fyne-io/fyne/ui"

type maxLayout struct {
}

func (m *maxLayout) Layout(objects []ui.CanvasObject, size ui.Size) {
	for _, child := range objects {
		child.Resize(size)
	}
}

func (m *maxLayout) MinSize(objects []ui.CanvasObject) ui.Size {
	minSize := ui.NewSize(0, 0)
	for _, child := range objects {
		minSize = minSize.Union(child.MinSize())
	}

	return minSize
}

func NewMaxLayout() *maxLayout {
	return &maxLayout{}
}

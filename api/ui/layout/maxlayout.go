// Package layout defines the various layouts available to Fyne apps
package layout

import "github.com/fyne-io/fyne/api/ui"

type maxLayout struct {
}

// Layout is called to pack all child objects into a specified size.
// For MaxLayout this sets all children to the full size passed.
func (m *maxLayout) Layout(objects []ui.CanvasObject, size ui.Size) {
	for _, child := range objects {
		child.Resize(size)
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For MaxLayout this is determined simply as the MinSize of the largest child.
func (m *maxLayout) MinSize(objects []ui.CanvasObject) ui.Size {
	minSize := ui.NewSize(0, 0)
	for _, child := range objects {
		minSize = minSize.Union(child.MinSize())
	}

	return minSize
}

// NewMaxLayout creates a new MaxLayout instance
func NewMaxLayout() ui.Layout {
	return &maxLayout{}
}

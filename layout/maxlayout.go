// Package layout defines the various layouts available to Fyne apps
package layout // import "fyne.io/fyne/layout"

import "fyne.io/fyne"

// Declare conformity with Layout interface
var _ fyne.Layout = (*maxLayout)(nil)

type maxLayout struct {
}

// Layout is called to pack all child objects into a specified size.
// For MaxLayout this sets all children to the full size passed.
func (m *maxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, child := range objects {
		child.Resize(size)
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For MaxLayout this is determined simply as the MinSize of the largest child.
func (m *maxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Union(child.MinSize())
	}

	return minSize
}

// NewMaxLayout creates a new MaxLayout instance
func NewMaxLayout() fyne.Layout {
	return &maxLayout{}
}

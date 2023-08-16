// Package layout defines the various layouts available to Fyne apps.
package layout // import "fyne.io/fyne/v2/layout"

import "fyne.io/fyne/v2"

// Declare conformity with Layout interface
var _ fyne.Layout = (*stackLayout)(nil)

type stackLayout struct {
}

// NewStackLayout returns a new StackLayout instance. Objects are stacked
// on top of each other with later objects on top of those before.
// Having only a single object has no impact as CanvasObjects will
// fill the available space even without a Stack.
//
// Since: 2.4
func NewStackLayout() fyne.Layout {
	return &stackLayout{}
}

// NewMaxLayout creates a new MaxLayout instance
//
// Deprecated: Use layout.NewStackLayout() instead.
func NewMaxLayout() fyne.Layout {
	return NewStackLayout()
}

// Layout is called to pack all child objects into a specified size.
// For StackLayout this sets all children to the full size passed.
func (m *stackLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	topLeft := fyne.NewPos(0, 0)
	for _, child := range objects {
		child.Resize(size)
		child.Move(topLeft)
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For StackLayout this is determined simply as the MinSize of the largest child.
func (m *stackLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Max(child.MinSize())
	}

	return minSize
}

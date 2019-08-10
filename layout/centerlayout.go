package layout

import "fyne.io/fyne"

// Declare conformity with Layout interface
var _ fyne.Layout = (*centerLayout)(nil)

type centerLayout struct {
}

// Layout is called to pack all child objects into a specified size.
// For CenterLayout this sets all children to their minimum size, centered within the space.
func (c *centerLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, child := range objects {
		childMin := child.MinSize()
		child.Resize(childMin)
		child.Move(fyne.NewPos(int(float32(size.Width-childMin.Width)/2), int(float32(size.Height-childMin.Height)/2)))
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For CenterLayout this is determined simply as the MinSize of the largest child.
func (c *centerLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Union(child.MinSize())
	}

	return minSize
}

// NewCenterLayout creates a new CenterLayout instance
func NewCenterLayout() fyne.Layout {
	return &centerLayout{}
}

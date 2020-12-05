package layout

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*paddedLayout)(nil)

type paddedLayout struct {
}

// Layout is called to pack all child objects into a specified size.
// For PaddedLayout this sets all children to the full size passed minus padding all around.
func (l *paddedLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	pos := fyne.NewPos(theme.Padding(), theme.Padding())
	siz := fyne.NewSize(size.Width-2*theme.Padding(), size.Height-2*theme.Padding())
	for _, child := range objects {
		child.Resize(siz)
		child.Move(pos)
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For PaddedLayout this is determined simply as the MinSize of the largest child plus padding all around.
func (l *paddedLayout) MinSize(objects []fyne.CanvasObject) (min fyne.Size) {
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		min = min.Union(child.MinSize())
	}
	min = min.Add(fyne.NewSize(2*theme.Padding(), 2*theme.Padding()))
	return
}

// NewPaddedLayout creates a new PaddedLayout instance
//
// Since: 1.4
func NewPaddedLayout() fyne.Layout {
	return &paddedLayout{}
}

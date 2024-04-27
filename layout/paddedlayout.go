package layout

import (
	"fyne.io/fyne/v2"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*paddedLayout)(nil)

type paddedLayout struct {
	BaseLayout
}

// Layout is called to pack all child objects into a specified size.
// For PaddedLayout this sets all children to the full size passed minus padding all around.
func (l paddedLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	topPadding, bottomPadding, leftPadding, rightPadding := l.GetPaddings()
	pos := fyne.NewPos(leftPadding, topPadding)
	siz := fyne.NewSize(
		size.Width-leftPadding-rightPadding,
		size.Height-topPadding-bottomPadding)
	for _, child := range objects {
		child.Resize(siz)
		child.Move(pos)
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For PaddedLayout this is determined simply as the MinSize of the largest child plus padding all around.
func (l paddedLayout) MinSize(objects []fyne.CanvasObject) (min fyne.Size) {
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		min = min.Max(child.MinSize())
	}
	topPadding, bottomPadding, leftPadding, rightPadding := l.GetPaddings()
	min = min.Add(fyne.NewSize(leftPadding+rightPadding, topPadding+bottomPadding))
	return
}

// NewPaddedLayout creates a new PaddedLayout instance
//
// Since: 1.4
func NewPaddedLayout(options ...LayoutOption) fyne.Layout {
	l := &paddedLayout{}
	for _, option := range options {
		option(l)
	}
	return l
}

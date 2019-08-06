package layout

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*borderLayout)(nil)

type borderLayout struct {
	top, bottom, left, right fyne.CanvasObject
}

// Layout is called to pack all child objects into a specified size.
// For BorderLayout this arranges the top, bottom, left and right widgets at
// the sides and any remaining widgets are maximised in the middle space.
func (b *borderLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	var topSize, bottomSize, leftSize, rightSize fyne.Size
	if b.top != nil && b.top.Visible() {
		b.top.Resize(fyne.NewSize(size.Width, b.top.MinSize().Height))
		b.top.Move(fyne.NewPos(0, 0))
		topSize = fyne.NewSize(size.Width, b.top.MinSize().Height+theme.Padding())
	}
	if b.bottom != nil && b.bottom.Visible() {
		b.bottom.Resize(fyne.NewSize(size.Width, b.bottom.MinSize().Height))
		b.bottom.Move(fyne.NewPos(0, size.Height-b.bottom.MinSize().Height))
		bottomSize = fyne.NewSize(size.Width, b.bottom.MinSize().Height+theme.Padding())
	}
	if b.left != nil && b.left.Visible() {
		b.left.Resize(fyne.NewSize(b.left.MinSize().Width, size.Height-topSize.Height-bottomSize.Height))
		b.left.Move(fyne.NewPos(0, topSize.Height))
		leftSize = fyne.NewSize(b.left.MinSize().Width+theme.Padding(), size.Height-topSize.Height-bottomSize.Height)
	}
	if b.right != nil && b.right.Visible() {
		b.right.Resize(fyne.NewSize(b.right.MinSize().Width, size.Height-topSize.Height-bottomSize.Height))
		b.right.Move(fyne.NewPos(size.Width-b.right.MinSize().Width, topSize.Height))
		rightSize = fyne.NewSize(b.right.MinSize().Width+theme.Padding(), size.Height-topSize.Height-bottomSize.Height)
	}

	middleSize := fyne.NewSize(size.Width-leftSize.Width-rightSize.Width, size.Height-topSize.Height-bottomSize.Height)
	middlePos := fyne.NewPos(leftSize.Width, topSize.Height)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if child != b.top && child != b.bottom && child != b.left && child != b.right {
			child.Resize(middleSize)
			child.Move(middlePos)
		}
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For BorderLayout this is determined by the MinSize height of the top and
// plus the MinSize width of the left and right, plus any padding needed.
// This is then added to the union of the MinSize for any remaining content.
func (b *borderLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if child != b.top && child != b.bottom && child != b.left && child != b.right {
			minSize = minSize.Union(child.MinSize())
		}
	}

	if b.left != nil && b.left.Visible() {
		minHeight := fyne.Max(minSize.Height, b.left.MinSize().Height)
		minSize = fyne.NewSize(minSize.Width+b.left.MinSize().Width+theme.Padding(), minHeight)
	}
	if b.right != nil && b.right.Visible() {
		minHeight := fyne.Max(minSize.Height, b.right.MinSize().Height)
		minSize = fyne.NewSize(minSize.Width+b.right.MinSize().Width+theme.Padding(), minHeight)
	}

	if b.top != nil && b.top.Visible() {
		minWidth := fyne.Max(minSize.Width, b.top.MinSize().Width)
		minSize = fyne.NewSize(minWidth, minSize.Height+b.top.MinSize().Height+theme.Padding())
	}
	if b.bottom != nil && b.bottom.Visible() {
		minWidth := fyne.Max(minSize.Width, b.bottom.MinSize().Width)
		minSize = fyne.NewSize(minWidth, minSize.Height+b.bottom.MinSize().Height+theme.Padding())
	}

	return minSize
}

// NewBorderLayout creates a new BorderLayout instance with top, left, bottom
// and right objects set. All other items in the container will fill the centre
// space
func NewBorderLayout(top, bottom, left, right fyne.CanvasObject) fyne.Layout {
	return &borderLayout{top, bottom, left, right}
}

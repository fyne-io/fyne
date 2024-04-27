package layout

import (
	"fyne.io/fyne/v2"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*borderLayout)(nil)

type borderLayout struct {
	BaseLayout

	top, bottom, left, right fyne.CanvasObject
}

// NewBorderLayout creates a new BorderLayout instance with top, bottom, left
// and right objects set. All other items in the container will fill the centre
// space
func NewBorderLayout(top, bottom, left, right fyne.CanvasObject, options ...LayoutOption) fyne.Layout {
	l := &borderLayout{top: top, bottom: bottom, left: left, right: right}
	for _, option := range options {
		option(l)
	}
	return l
}

// Layout is called to pack all child objects into a specified size.
// For BorderLayout this arranges the top, bottom, left and right widgets at
// the sides and any remaining widgets are maximised in the middle space.
func (b *borderLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	topPadding, bottomPadding, leftPadding, rightPadding := b.GetPaddings()

	var topSize, bottomSize, leftSize, rightSize fyne.Size
	if b.top != nil && b.top.Visible() {
		topHeight := b.top.MinSize().Height
		b.top.Resize(fyne.NewSize(size.Width, topHeight))
		b.top.Move(fyne.NewPos(0, 0))
		topSize = fyne.NewSize(size.Width, topHeight+topPadding)
	}
	if b.bottom != nil && b.bottom.Visible() {
		bottomHeight := b.bottom.MinSize().Height
		b.bottom.Resize(fyne.NewSize(size.Width, bottomHeight))
		b.bottom.Move(fyne.NewPos(0, size.Height-bottomHeight))
		bottomSize = fyne.NewSize(size.Width, bottomHeight+bottomPadding)
	}
	if b.left != nil && b.left.Visible() {
		leftWidth := b.left.MinSize().Width
		b.left.Resize(fyne.NewSize(leftWidth, size.Height-topSize.Height-bottomSize.Height))
		b.left.Move(fyne.NewPos(0, topSize.Height))
		leftSize = fyne.NewSize(leftWidth+leftPadding, size.Height-topSize.Height-bottomSize.Height)
	}
	if b.right != nil && b.right.Visible() {
		rightWidth := b.right.MinSize().Width
		b.right.Resize(fyne.NewSize(rightWidth, size.Height-topSize.Height-bottomSize.Height))
		b.right.Move(fyne.NewPos(size.Width-rightWidth, topSize.Height))
		rightSize = fyne.NewSize(rightWidth+rightPadding, size.Height-topSize.Height-bottomSize.Height)
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
			minSize = minSize.Max(child.MinSize())
		}
	}

	topPadding, bottomPadding, leftPadding, rightPadding := b.GetPaddings()

	if b.left != nil && b.left.Visible() {
		leftMin := b.left.MinSize()
		minHeight := fyne.Max(minSize.Height, leftMin.Height)
		minSize = fyne.NewSize(minSize.Width+leftMin.Width+leftPadding, minHeight)
	}
	if b.right != nil && b.right.Visible() {
		rightMin := b.right.MinSize()
		minHeight := fyne.Max(minSize.Height, rightMin.Height)
		minSize = fyne.NewSize(minSize.Width+rightMin.Width+rightPadding, minHeight)
	}

	if b.top != nil && b.top.Visible() {
		topMin := b.top.MinSize()
		minWidth := fyne.Max(minSize.Width, topMin.Width)
		minSize = fyne.NewSize(minWidth, minSize.Height+topMin.Height+topPadding)
	}
	if b.bottom != nil && b.bottom.Visible() {
		bottomMin := b.bottom.MinSize()
		minWidth := fyne.Max(minSize.Width, bottomMin.Width)
		minSize = fyne.NewSize(minWidth, minSize.Height+bottomMin.Height+bottomPadding)
	}

	return minSize
}

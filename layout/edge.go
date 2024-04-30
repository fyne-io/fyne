package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// NewEdge creates a new edge layout instance with top, bottom, left, right
// and center objects set.
//
// Since: 2.5
func NewEdge(top, bottom, left, right, center fyne.CanvasObject) fyne.Layout {
	return edgeLayout{top: top, bottom: bottom, left: left, right: right, center: center}
}

// Declare conformity with Layout interface
var _ fyne.Layout = (*borderLayout)(nil)

type edgeLayout struct {
	top, bottom, left, right, center fyne.CanvasObject
}

// Layout is called to pack all child objects into a specified size.
// For BorderLayout this arranges the top, bottom, left and right widgets at
// the sides and any remaining widgets are maximised in the middle space.
func (b edgeLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	padding := theme.Padding()
	var topSize, bottomSize, leftSize, rightSize fyne.Size
	if b.top != nil && b.top.Visible() {
		topHeight := b.top.MinSize().Height
		b.top.Resize(fyne.NewSize(size.Width, topHeight))
		b.top.Move(fyne.NewPos(0, 0))
		topSize = fyne.NewSize(size.Width, topHeight+padding)
	}
	if b.bottom != nil && b.bottom.Visible() {
		bottomHeight := b.bottom.MinSize().Height
		b.bottom.Resize(fyne.NewSize(size.Width, bottomHeight))
		b.bottom.Move(fyne.NewPos(0, size.Height-bottomHeight))
		bottomSize = fyne.NewSize(size.Width, bottomHeight+padding)
	}
	if b.left != nil && b.left.Visible() {
		leftWidth := b.left.MinSize().Width
		b.left.Resize(fyne.NewSize(leftWidth, size.Height-topSize.Height-bottomSize.Height))
		b.left.Move(fyne.NewPos(0, topSize.Height))
		leftSize = fyne.NewSize(leftWidth+padding, size.Height-topSize.Height-bottomSize.Height)
	}
	if b.right != nil && b.right.Visible() {
		rightWidth := b.right.MinSize().Width
		b.right.Resize(fyne.NewSize(rightWidth, size.Height-topSize.Height-bottomSize.Height))
		b.right.Move(fyne.NewPos(size.Width-rightWidth, topSize.Height))
		rightSize = fyne.NewSize(rightWidth+padding, size.Height-topSize.Height-bottomSize.Height)
	}
	if b.center != nil && b.center.Visible() {
		middleSize := fyne.NewSize(size.Width-leftSize.Width-rightSize.Width, size.Height-topSize.Height-bottomSize.Height)
		middlePos := fyne.NewPos(leftSize.Width, topSize.Height)
		b.center.Resize(middleSize)
		b.center.Move(middlePos)
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For the edge layout, this is determined by the MinSize height of the top and
// plus the MinSize width of the left and right, plus any padding needed.
// This is then added to the union of the MinSize for any remaining content.
func (b edgeLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	padding := theme.Padding()

	var minSize fyne.Size
	if b.center != nil && b.center.Visible() {
		minSize = b.center.MinSize()
	}

	if b.left != nil && b.left.Visible() {
		leftMin := b.left.MinSize()
		minHeight := fyne.Max(minSize.Height, leftMin.Height)
		minSize = fyne.NewSize(minSize.Width+leftMin.Width+padding, minHeight)
	}
	if b.right != nil && b.right.Visible() {
		rightMin := b.right.MinSize()
		minHeight := fyne.Max(minSize.Height, rightMin.Height)
		minSize = fyne.NewSize(minSize.Width+rightMin.Width+padding, minHeight)
	}

	if b.top != nil && b.top.Visible() {
		topMin := b.top.MinSize()
		minWidth := fyne.Max(minSize.Width, topMin.Width)
		minSize = fyne.NewSize(minWidth, minSize.Height+topMin.Height+padding)
	}
	if b.bottom != nil && b.bottom.Visible() {
		bottomMin := b.bottom.MinSize()
		minWidth := fyne.Max(minSize.Width, bottomMin.Width)
		minSize = fyne.NewSize(minWidth, minSize.Height+bottomMin.Height+padding)
	}

	return minSize
}

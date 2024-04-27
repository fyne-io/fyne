package layout

import (
	"fyne.io/fyne/v2"
)

// NewVBoxLayout returns a vertical box layout for stacking a number of child
// canvas objects or widgets top to bottom. The objects are always displayed
// at their vertical MinSize. Use a different layout if the objects are intended
// to be larger then their vertical MinSize.
func NewVBoxLayout(options ...LayoutOption) fyne.Layout {
	l := &vBoxLayout{}
	for _, option := range options {
		option(l)
	}
	return l
}

// NewHBoxLayout returns a horizontal box layout for stacking a number of child
// canvas objects or widgets left to right. The objects are always displayed
// at their horizontal MinSize. Use a different layout if the objects are intended
// to be larger then their horizontal MinSize.
func NewHBoxLayout(options ...LayoutOption) fyne.Layout {
	l := &hBoxLayout{}
	for _, option := range options {
		option(l)
	}
	return l
}

// Declare conformity with Layout interface
var _ fyne.Layout = (*vBoxLayout)(nil)

type vBoxLayout struct {
	BaseLayout
}

// Layout is called to pack all child objects into a specified size.
// This will pack objects into a single column where each item
// is full width but the height is the minimum required.
// Any spacers added will pad the view, sharing the space if there are two or more.
func (v vBoxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	spacers := 0
	visibleObjects := 0
	// Size taken up by visible objects
	total := float32(0)

	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if isVerticalSpacer(child) {
			spacers++
			continue
		}

		visibleObjects++
		total += child.MinSize().Height
	}

	topPadding, _, _, _ := v.GetPaddings()

	// Amount of space not taken up by visible objects and inter-object padding
	extra := size.Height - total - (topPadding * float32(visibleObjects-1))

	// Spacers split extra space equally
	spacerSize := float32(0)
	if spacers > 0 {
		spacerSize = extra / float32(spacers)
	}

	x, y := float32(0), float32(0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if isVerticalSpacer(child) {
			y += spacerSize
			continue
		}
		child.Move(fyne.NewPos(x, y))

		height := child.MinSize().Height
		y += topPadding + height
		child.Resize(fyne.NewSize(size.Width, height))
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a BoxLayout this is the width of the widest item and the height is
// the sum of of all children combined with padding between each.
func (v vBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	addPadding := false
	topPadding, _, _, _ := v.GetPaddings()
	for _, child := range objects {
		if !child.Visible() || isVerticalSpacer(child) {
			continue
		}

		childMin := child.MinSize()
		minSize.Width = fyne.Max(childMin.Width, minSize.Width)
		minSize.Height += childMin.Height
		if addPadding {
			minSize.Height += topPadding
		}
		addPadding = true
	}
	return minSize
}

// Declare conformity with Layout interface
var _ fyne.Layout = (*hBoxLayout)(nil)

type hBoxLayout struct {
	BaseLayout
}

// Layout is called to pack all child objects into a specified size.
// For a VBoxLayout this will pack objects into a single column where each item
// is full width but the height is the minimum required.
// Any spacers added will pad the view, sharing the space if there are two or more.
func (g hBoxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	spacers := 0
	visibleObjects := 0
	// Size taken up by visible objects
	total := float32(0)

	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if isHorizontalSpacer(child) {
			spacers++
			continue
		}

		visibleObjects++
		total += child.MinSize().Width
	}

	_, _, leftPadding, _ := g.GetPaddings()

	// Amount of space not taken up by visible objects and inter-object padding
	extra := size.Width - total - (leftPadding * float32(visibleObjects-1))

	// Spacers split extra space equally
	spacerSize := float32(0)
	if spacers > 0 {
		spacerSize = extra / float32(spacers)
	}

	x, y := float32(0), float32(0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if isHorizontalSpacer(child) {
			x += spacerSize
			continue
		}
		child.Move(fyne.NewPos(x, y))

		width := child.MinSize().Width
		x += leftPadding + width
		child.Resize(fyne.NewSize(width, size.Height))
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a BoxLayout this is the width of the widest item and the height is
// the sum of of all children combined with padding between each.
func (g hBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	addPadding := false
	_, _, leftPadding, _ := g.GetPaddings()
	for _, child := range objects {
		if !child.Visible() || isHorizontalSpacer(child) {
			continue
		}

		childMin := child.MinSize()
		minSize.Height = fyne.Max(childMin.Height, minSize.Height)
		minSize.Width += childMin.Width
		if addPadding {
			minSize.Width += leftPadding
		}
		addPadding = true
	}
	return minSize
}

func isVerticalSpacer(obj fyne.CanvasObject) bool {
	spacer, ok := obj.(SpacerObject)
	return ok && spacer.ExpandVertical()
}

func isHorizontalSpacer(obj fyne.CanvasObject) bool {
	spacer, ok := obj.(SpacerObject)
	return ok && spacer.ExpandHorizontal()
}

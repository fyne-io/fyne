package layout

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*boxLayout)(nil)

type boxLayout struct {
	horizontal bool
}

func isVerticalSpacer(obj fyne.CanvasObject) bool {
	if spacer, ok := obj.(SpacerObject); ok {
		return spacer.ExpandVertical()
	}

	return false
}

func isHorizontalSpacer(obj fyne.CanvasObject) bool {
	if spacer, ok := obj.(SpacerObject); ok {
		return spacer.ExpandHorizontal()
	}

	return false
}

func (g *boxLayout) isSpacer(obj fyne.CanvasObject) bool {
	// invisible spacers don't impact layout
	if !obj.Visible() {
		return false
	}

	if g.horizontal {
		return isHorizontalSpacer(obj)
	}
	return isVerticalSpacer(obj)
}

// Layout is called to pack all child objects into a specified size.
// For a VBoxLayout this will pack objects into a single column where each item
// is full width but the height is the minimum required.
// Any spacers added will pad the view, sharing the space if there are two or more.
func (g *boxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	spacers := make([]fyne.CanvasObject, 0)
	total := 0
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if g.isSpacer(child) {
			spacers = append(spacers, child)
			continue
		}
		if g.horizontal {
			total += child.MinSize().Width
		} else {
			total += child.MinSize().Height
		}
	}

	x, y := 0, 0
	var extra int
	if g.horizontal {
		extra = size.Width - total - (theme.Padding() * (len(objects) - len(spacers) - 1))
	} else {
		extra = size.Height - total - (theme.Padding() * (len(objects) - len(spacers) - 1))
	}
	extraCell := 0
	if len(spacers) > 0 {
		extraCell = int(float64(extra) / float64(len(spacers)))
	}

	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		width := child.MinSize().Width
		height := child.MinSize().Height

		if g.isSpacer(child) {
			if g.horizontal {
				x += extraCell
			} else {
				y += extraCell
			}
			continue
		}
		child.Move(fyne.NewPos(x, y))

		if g.horizontal {
			x += theme.Padding() + width
			child.Resize(fyne.NewSize(width, size.Height))
		} else {
			y += theme.Padding() + height
			child.Resize(fyne.NewSize(size.Width, height))
		}
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a BoxLayout this is the width of the widest item and the height is
// the sum of of all children combined with padding between each.
func (g *boxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	spacerCount := 0
	minSize := fyne.NewSize(0, 0)
	added := false
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if g.isSpacer(child) {
			spacerCount++
			continue
		}

		if g.horizontal {
			minSize = minSize.Add(fyne.NewSize(child.MinSize().Width, 0))
			minSize.Height = fyne.Max(child.MinSize().Height, minSize.Height)
			if added {
				minSize.Width += theme.Padding()
			}
		} else {
			minSize = minSize.Add(fyne.NewSize(0, child.MinSize().Height))
			minSize.Width = fyne.Max(child.MinSize().Width, minSize.Width)
			if added {
				minSize.Height += theme.Padding()
			}
		}
		added = true
	}

	return minSize
}

// NewHBoxLayout returns a horizontal box layout for stacking a number of child
// canvas objects or widgets left to right.
func NewHBoxLayout() fyne.Layout {
	return &boxLayout{true}
}

// NewVBoxLayout returns a vertical box layout for stacking a number of child
// canvas objects or widgets top to bottom.
func NewVBoxLayout() fyne.Layout {
	return &boxLayout{false}
}

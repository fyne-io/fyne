package layout

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

var _ fyne.Layout = (*Box)(nil)

// Box is a box layout for stacking a number of child canvas objects horizontally or vertically.
type Box struct {
	Horizontal        bool
}

// Layout is called to pack all child objects into a specified size.
// For a vertical box layout this will pack objects into a single column where each item
// is full width but the height is the minimum required.
// Any spacers added will pad the view, sharing the space if there are two or more.
func (b *Box) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	spacers := make([]fyne.CanvasObject, 0)
	total := 0
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if b.isSpacer(child) {
			spacers = append(spacers, child)
			continue
		}
		if b.Horizontal {
			total += child.MinSize().Width
		} else {
			total += child.MinSize().Height
		}
	}

	x, y := 0, 0
	var extra int
	if b.Horizontal {
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

		if b.isSpacer(child) {
			if b.Horizontal {
				x += extraCell
			} else {
				y += extraCell
			}
			continue
		}
		child.Move(fyne.NewPos(x, y))

		if b.Horizontal {
			x += theme.Padding() + width
			child.Resize(fyne.NewSize(width, size.Height))
		} else {
			y += theme.Padding() + height
			child.Resize(fyne.NewSize(size.Width, height))
		}
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a box layout this is the width of the widest item and the height is
// the sum of of all children combined with padding between each.
func (b *Box) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	addPadding := false
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if b.isSpacer(child) {
			continue
		}

		if b.Horizontal {
			minSize.Height = fyne.Max(child.MinSize().Height, minSize.Height)
			minSize.Width += child.MinSize().Width
			if addPadding {
				minSize.Width += theme.Padding()
			}
		} else {
			minSize.Width = fyne.Max(child.MinSize().Width, minSize.Width)
			minSize.Height += child.MinSize().Height
			if addPadding {
				minSize.Height += theme.Padding()
			}
		}
		addPadding = true
	}
	return minSize
}

func (b *Box) isSpacer(obj fyne.CanvasObject) bool {
	// invisible spacers don't impact layout
	if !obj.Visible() {
		return false
	}
	s, ok := obj.(spacer)
	if !ok {
		return false
	}

	if b.Horizontal {
		return s.ExpandHorizontal()
	}
	return s.ExpandVertical()
}

type spacer interface {
	ExpandHorizontal() bool
	ExpandVertical() bool
}

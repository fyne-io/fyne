package layout

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/theme"
)

type boxLayout struct {
}

func isVerticalSpacer(obj interface{}) bool {
	if spacer, ok := obj.(SpacerObject); ok {
		return spacer.ExpandVertical()
	}

	return false
}

// Layout is called to pack all child objects into a specified size.
// For an VBoxLayout this will pack objects into a single column where each item
// is full width but the it's height is the minimum required.
// Any spacers added will pad the view, sharing the space if there are two or more.
func (g *boxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	spacers := make([]fyne.CanvasObject, 0)
	totalHeight := 0
	for _, child := range objects {
		if isVerticalSpacer(child) {
			spacers = append(spacers, child)
			continue
		}
		totalHeight += child.MinSize().Height
	}

	y := 0
	extraHeight := size.Height - totalHeight - (theme.Padding() * (len(objects) - len(spacers) - 1))
	extraCellHeight := 0
	if len(spacers) > 0 {
		extraCellHeight = int(float64(extraHeight) / float64(len(spacers)))
	}

	for _, child := range objects {
		height := child.MinSize().Height
		if isVerticalSpacer(child) {
			y += extraCellHeight
			continue
		}
		child.Move(fyne.NewPos(0, y))
		child.Resize(fyne.NewSize(size.Width, height))

		y += theme.Padding() + height
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a BoxLayout this is the width of the widest item and the height is
// the sum of of all children combined with padding between each.
func (g *boxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	spacerCount := 0
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if isVerticalSpacer(child) {
			spacerCount++
			continue
		}
		minSize = minSize.Add(fyne.NewSize(0,
			child.MinSize().Height))
		minSize.Width = fyne.Max(child.MinSize().Width, minSize.Width)
	}

	return minSize.Add(fyne.NewSize(0, theme.Padding()*(len(objects)-1-spacerCount)))
}

// NewVBoxLayout returns a vertical box layout for stacking a number of child
// canvas objects or widgets top to bottom.
func NewVBoxLayout() fyne.Layout {
	return new(boxLayout)
}

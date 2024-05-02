package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

const formLayoutCols = 2

// Declare conformity with Layout interface
var _ fyne.Layout = (*formLayout)(nil)

// formLayout is two column grid where each row has a label and a widget.
type formLayout struct{}

// calculateTableSizes calculates the izes of the table.
// This includes the width of the label column (maximum width of all labels),
// the width of the content column (maximum width of all content cells and remaining space in container)
// and the total minimum height of the form.
func (f *formLayout) calculateTableSizes(objects []fyne.CanvasObject, containerWidth float32) (labelWidth float32, contentWidth float32, height float32) {
	if len(objects)%formLayoutCols != 0 {
		return 0, 0, 0
	}

	rows := 0
	innerPadding := theme.InnerPadding()
	for i := 0; i < len(objects); i += formLayoutCols {
		labelCell, contentCell := objects[i], objects[i+1]
		if !labelCell.Visible() && !contentCell.Visible() {
			continue
		}

		// Label column width is the maximum of all labels.
		labelSize := labelCell.MinSize()
		if _, ok := labelCell.(*canvas.Text); ok {
			labelSize.Width += innerPadding * 2
		}
		labelWidth = fyne.Max(labelWidth, labelSize.Width)

		// Content column width is the maximum of all content items.
		contentSize := contentCell.MinSize()
		contentWidth = fyne.Max(contentWidth, contentSize.Width)

		rowHeight := fyne.Max(labelSize.Height, contentSize.Height)
		height += rowHeight
		rows++
	}

	padding := theme.Padding()
	contentWidth = fyne.Max(contentWidth, containerWidth-labelWidth-padding)
	return labelWidth, contentWidth, height + float32(rows-1)*padding
}

// Layout is called to pack all child objects into a table format with two columns.
func (f *formLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	labelWidth, contentWidth, _ := f.calculateTableSizes(objects, size.Width)

	y := float32(0)
	padding := theme.Padding()
	innerPadding := theme.InnerPadding()

	// Calculate size and position of object. Position and size is returned (instead of calling Move() and Resize()) to make inlineable.
	objectLayout := func(obj fyne.CanvasObject, offset, width, rowHeight, itemHeight float32) (fyne.Position, fyne.Size) {
		pos := fyne.NewPos(offset, y)
		size := fyne.NewSize(width, rowHeight)
		if _, ok := obj.(*canvas.Text); ok {
			pos = pos.AddXY(innerPadding, innerPadding)
			size.Width -= innerPadding * 2
			size.Height = itemHeight
		}

		return pos, size
	}

	remainer := len(objects) % formLayoutCols
	for i := 0; i < len(objects)-remainer; i += formLayoutCols {
		labelCell, contentCell := objects[i], objects[i+1]
		if !labelCell.Visible() && !contentCell.Visible() {
			continue
		}

		labelMin := labelCell.MinSize()
		contentMin := contentCell.MinSize()
		rowHeight := fyne.Max(labelMin.Height, contentMin.Height)

		pos, size := objectLayout(labelCell, 0, labelWidth, rowHeight, labelMin.Height)
		labelCell.Move(pos)
		labelCell.Resize(size)

		pos, size = objectLayout(contentCell, labelWidth+padding, contentWidth, rowHeight, contentMin.Height)
		contentCell.Move(pos)
		contentCell.Resize(size)

		y += rowHeight + padding
	}

	// Handle remaining item in the case of uneven number of objects:
	if remainer == 1 {
		lastCell := objects[len(objects)-1]
		lastMin := lastCell.MinSize()
		objectLayout(lastCell, 0, labelWidth, lastMin.Height, lastMin.Height)
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a FormLayout this is the width of the widest label and content items and the height is
// the sum of all column children combined with padding between each.
func (f *formLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	labelWidth, contentWidth, height := f.calculateTableSizes(objects, 0)
	return fyne.Size{
		Width:  labelWidth + contentWidth + theme.Padding(),
		Height: height,
	}
}

// NewFormLayout returns a new FormLayout instance
func NewFormLayout() fyne.Layout {
	return &formLayout{}
}

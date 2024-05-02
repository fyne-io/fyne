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
type formLayout struct {
}

func (f *formLayout) countRows(objects []fyne.CanvasObject) int {
	count := 0

	for i := 0; i < len(objects); i += formLayoutCols {
		if !objects[i].Visible() && !objects[i+1].Visible() {
			continue
		}
		count++
	}

	return count
}

// calculateTableSizes calculates the size of all of the cells in the table.
// It returns the width of the left column (labels), the width of the right column (content)
// as well as a slice with the height of each row. The height of each row will be set as the max
// size of the label and content cell heights for that row. The width of the label column will
// be set as the max width of all the label cells. The width of the content column will be set as
// the max width of all the content cells or the remaining space of the bounding containerWidth,
// if it is larger.
func (f *formLayout) calculateTableSizes(objects []fyne.CanvasObject, containerWidth float32) (labelWidth float32, contentWidth float32, height float32) {
	if (len(objects))%formLayoutCols != 0 {
		return 0, 0, 0
	}

	innerPadding := theme.InnerPadding()
	lowBound := 0
	highBound := 2
	rows := f.countRows(objects)
	for row := 0; row < rows; {
		currentRow := objects[lowBound:highBound]
		lowBound = highBound
		highBound += formLayoutCols
		if !currentRow[0].Visible() && !currentRow[1].Visible() {
			continue
		}

		labelCell := currentRow[0].MinSize()
		if _, ok := currentRow[0].(*canvas.Text); ok {
			labelCell.Width += innerPadding * 2
		}
		labelWidth = fyne.Max(labelWidth, labelCell.Width)

		contentCell := currentRow[1].MinSize()
		contentWidth = fyne.Max(contentWidth, contentCell.Width)

		rowHeight := fyne.Max(labelCell.Height, contentCell.Height)
		height += rowHeight
		row++
	}

	padding := theme.Padding()
	contentWidth = fyne.Max(contentWidth, containerWidth-labelWidth-padding)
	return labelWidth, contentWidth, height + float32(rows-1)*padding
}

// Layout is called to pack all child objects into a table format with two columns.
func (f *formLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	labelWidth, contentWidth, _ := f.calculateTableSizes(objects, size.Width)

	padding := theme.Padding()
	innerPadding := theme.InnerPadding()

	row := 0
	y := float32(0)
	for i := 0; i < len(objects); i += formLayoutCols {
		if !objects[i].Visible() && (i+1 < len(objects) && !objects[i+1].Visible()) {
			continue
		}

		rowHeight := objects[i].MinSize().Height
		if i+1 < len(objects) {
			rowHeight = fyne.Max(rowHeight, objects[i+1].MinSize().Height)
		}

		pos := fyne.NewPos(0, y)
		size := fyne.NewSize(labelWidth, rowHeight)
		if _, ok := objects[i].(*canvas.Text); ok {
			pos = pos.AddXY(innerPadding, innerPadding)
			size.Width -= innerPadding * 2
			size.Height = objects[i].MinSize().Height
		}

		objects[i].Move(pos)
		objects[i].Resize(size)

		if i+1 < len(objects) {
			pos = fyne.NewPos(padding+labelWidth, y)
			size = fyne.NewSize(contentWidth, rowHeight)
			if _, ok := objects[i+1].(*canvas.Text); ok {
				pos = pos.AddXY(innerPadding, innerPadding)
				size.Width -= innerPadding * 2
				size.Height = objects[i+1].MinSize().Height
			}

			objects[i+1].Move(pos)
			objects[i+1].Resize(size)
		}

		y += rowHeight + padding
		row++
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

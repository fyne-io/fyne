package layout

import (
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

const formLayoutCols = 2

// formLayout is two column grid where each row has a label and a widget.
type formLayout struct {
}

func (f *formLayout) countRows(objects []fyne.CanvasObject) int {
	return int(math.Ceil(float64(len(objects)) / float64(formLayoutCols)))
}

// tableCellsSize defines the size for all the cells of the form table.
// The height of each row will be set as the max value between the label and content cell heights.
// The width of the label column will be set as the max width value between all the label cells.
// The width of the content column will be set as the max width value between all the content cells
// or the remaining space of the bounding containerWidth, if it is larger.
func (f *formLayout) tableCellsSize(objects []fyne.CanvasObject, containerWidth int) [][2]fyne.Size {
	rows := f.countRows(objects)
	table := make([][2]fyne.Size, rows)

	if (len(objects))%formLayoutCols != 0 {
		return table
	}

	lowBound := 0
	highBound := 2
	labelCellMaxWidth := 0
	contentCellMaxWidth := 0
	for row := 0; row < rows; row++ {
		currentRow := objects[lowBound:highBound]

		labelCell := currentRow[0].MinSize()
		labelCellMaxWidth = fyne.Max(labelCellMaxWidth, labelCell.Width)

		contentCell := currentRow[1].MinSize()
		contentCellMaxWidth = fyne.Max(contentCellMaxWidth, contentCell.Width)

		rowHeight := fyne.Max(labelCell.Height, contentCell.Height)

		labelCell.Height = rowHeight
		contentCell.Height = rowHeight

		table[row][0] = labelCell
		table[row][1] = contentCell

		lowBound = highBound
		highBound += 2
	}

	contentWidth := fyne.Max(contentCellMaxWidth, containerWidth-labelCellMaxWidth-theme.Padding())
	for row := 0; row < rows; row++ {
		table[row][0].Width = labelCellMaxWidth
		table[row][1].Width = contentWidth
	}

	return table
}

// Layout is called to pack all child objects into a table format with two columns.
func (f *formLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {

	var cellWidth float64
	var cellHeight float64

	table := f.tableCellsSize(objects, size.Width)

	row := 0
	x, y := 0, 0
	deltaY := 0

	for i, child := range objects {
		if row > 0 {
			deltaY = table[row-1][0].Height + theme.Padding()
		}

		tableRow := table[row]
		cellHeight = float64(tableRow[0].Height)
		if (i)%formLayoutCols == 0 {
			x = 0
			y += deltaY
			cellWidth = float64(tableRow[0].Width)
		} else {
			x = theme.Padding() + tableRow[0].Width
			cellWidth = float64(tableRow[1].Width)
			row++
		}

		child.Move(fyne.NewPos(x, y))
		child.Resize(fyne.NewSize(int(cellWidth), int(cellHeight)))
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a FormLayout this is the width of the widest label and content items and the height is
// the sum of all column children combined with padding between each.
func (f *formLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {

	table := f.tableCellsSize(objects, 0)

	minSize := fyne.NewSize(0, 0)

	if len(table) == 0 {
		return minSize
	}

	minSize.Width = table[0][0].Width + table[0][1].Width + theme.Padding()
	for row := 0; row < len(table); row++ {
		minSize.Height += table[row][0].Height
		if row > 0 {
			minSize.Height += theme.Padding()
		}
	}
	return minSize
}

// NewFormLayout returns a new FormLayout instance
func NewFormLayout() fyne.Layout {
	return &formLayout{}
}

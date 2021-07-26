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

// tableCellsSize defines the size for all the cells of the form table.
// The height of each row will be set as the max value between the label and content cell heights.
// The width of the label column will be set as the max width value between all the label cells.
// The width of the content column will be set as the max width value between all the content cells
// or the remaining space of the bounding containerWidth, if it is larger.
func (f *formLayout) tableCellsSize(objects []fyne.CanvasObject, containerWidth float32) [][2]fyne.Size {
	rows := f.countRows(objects)
	table := make([][2]fyne.Size, rows)

	if (len(objects))%formLayoutCols != 0 {
		return table
	}

	lowBound := 0
	highBound := 2
	labelCellMaxWidth := float32(0)
	contentCellMaxWidth := float32(0)
	for row := 0; row < rows; {
		currentRow := objects[lowBound:highBound]
		lowBound = highBound
		highBound += formLayoutCols
		if !currentRow[0].Visible() && !currentRow[1].Visible() {
			continue
		}

		labelCell := currentRow[0].MinSize()
		if _, ok := currentRow[0].(*canvas.Text); ok {
			labelCell.Width += theme.Padding() * 4
		}
		labelCellMaxWidth = fyne.Max(labelCellMaxWidth, labelCell.Width)

		contentCell := currentRow[1].MinSize()
		contentCellMaxWidth = fyne.Max(contentCellMaxWidth, contentCell.Width)

		rowHeight := fyne.Max(labelCell.Height, contentCell.Height)

		labelCell.Height = rowHeight
		contentCell.Height = rowHeight

		table[row][0] = labelCell
		table[row][1] = contentCell
		row++
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
	table := f.tableCellsSize(objects, size.Width)

	row := 0
	y := float32(0)
	for i := 0; i < len(objects); i += formLayoutCols {
		if !objects[i].Visible() && (i+1 < len(objects) && !objects[i+1].Visible()) {
			continue
		}
		if row > 0 {
			y += table[row-1][0].Height + theme.Padding()
		}

		tableRow := table[row]
		if _, ok := objects[i].(*canvas.Text); ok {
			objects[i].Move(fyne.NewPos(theme.Padding()*2, y+theme.Padding()*2))
			objects[i].Resize(fyne.NewSize(tableRow[0].Width-theme.Padding()*4, objects[i].MinSize().Height))
		} else {
			objects[i].Move(fyne.NewPos(0, y))
			objects[i].Resize(fyne.NewSize(tableRow[0].Width, tableRow[0].Height))
		}

		if i+1 < len(objects) {
			objects[i+1].Move(fyne.NewPos(theme.Padding()+tableRow[0].Width, y))
			objects[i+1].Resize(fyne.NewSize(tableRow[1].Width, tableRow[0].Height))
		}
		row++
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

	added := false
	minSize.Width = table[0][0].Width + table[0][1].Width + theme.Padding()
	for row := 0; row < len(table); row++ {
		minSize.Height += table[row][0].Height
		if added {
			minSize.Height += theme.Padding()
		}
		added = true
	}
	return minSize
}

// NewFormLayout returns a new FormLayout instance
func NewFormLayout() fyne.Layout {
	return &formLayout{}
}

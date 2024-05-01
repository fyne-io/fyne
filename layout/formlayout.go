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
func (f *formLayout) calculateTableSizes(objects []fyne.CanvasObject, containerWidth float32) (float32, float32, []float32) {
	rows := f.countRows(objects)
	heights := make([]float32, rows)

	if (len(objects))%formLayoutCols != 0 {
		return 0, 0, heights
	}

	innerPadding := theme.InnerPadding()
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
			labelCell.Width += innerPadding * 2
		}
		labelCellMaxWidth = fyne.Max(labelCellMaxWidth, labelCell.Width)

		contentCell := currentRow[1].MinSize()
		contentCellMaxWidth = fyne.Max(contentCellMaxWidth, contentCell.Width)

		rowHeight := fyne.Max(labelCell.Height, contentCell.Height)
		heights[row] = rowHeight
		row++
	}

	contentCellMaxWidth = fyne.Max(contentCellMaxWidth, containerWidth-labelCellMaxWidth-theme.Padding())
	return labelCellMaxWidth, contentCellMaxWidth, heights
}

// Layout is called to pack all child objects into a table format with two columns.
func (f *formLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	labelWidth, contentWidth, heights := f.calculateTableSizes(objects, size.Width)

	padding := theme.Padding()
	innerPadding := theme.InnerPadding()

	row := 0
	y := float32(0)
	for i := 0; i < len(objects); i += formLayoutCols {
		if !objects[i].Visible() && (i+1 < len(objects) && !objects[i+1].Visible()) {
			continue
		}
		if row > 0 {
			y += heights[row-1] + padding
		}

		pos := fyne.NewPos(0, y)
		size := fyne.NewSize(labelWidth, heights[row])
		if _, ok := objects[i].(*canvas.Text); ok {
			pos = pos.AddXY(innerPadding, innerPadding)
			size.Width -= innerPadding * 2
			size.Height = objects[i].MinSize().Height
		}

		objects[i].Move(pos)
		objects[i].Resize(size)

		if i+1 < len(objects) {
			pos = fyne.NewPos(padding+labelWidth, y)
			size = fyne.NewSize(contentWidth, heights[row])
			if _, ok := objects[i+1].(*canvas.Text); ok {
				pos = pos.AddXY(innerPadding, innerPadding)
				size.Width -= innerPadding * 2
				size.Height = objects[i+1].MinSize().Height
			}

			objects[i+1].Move(pos)
			objects[i+1].Resize(size)
		}
		row++
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a FormLayout this is the width of the widest label and content items and the height is
// the sum of all column children combined with padding between each.
func (f *formLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	labelWidth, contentWidth, heights := f.calculateTableSizes(objects, 0)
	if len(heights) == 0 {
		return fyne.Size{}
	}

	padding := theme.Padding()
	minSize := fyne.Size{
		Width:  labelWidth + contentWidth + padding,
		Height: padding * float32(len(heights)-1),
	}

	for _, height := range heights {
		minSize.Height += height
	}

	return minSize
}

// NewFormLayout returns a new FormLayout instance
func NewFormLayout() fyne.Layout {
	return &formLayout{}
}

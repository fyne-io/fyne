package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

// TextGrid is a monospaced grid of characters.
// This is designed to be used by a text editor or advanced test presentation.
type TextGrid struct {
	BaseWidget
	textHandler
}

// MinSize returns the smallest size this widget can shrink to
func (t *TextGrid) MinSize() fyne.Size {
	t.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (t *TextGrid) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	render := &textGridRender{text: t}
	t.updateRowBounds()
	render.ensureGrid()

	cell := canvas.NewText("M", color.White)
	cell.TextStyle.Monospace = true
	render.cellSize = cell.MinSize()

	return render
}

// NewTextGrid creates a new textgrid widget with the specified string content.
func NewTextGrid(content string) *TextGrid {
	handler := textHandler{buffer: []rune(content)}
	grid := &TextGrid{textHandler: handler}
	grid.ExtendBaseWidget(grid)
	return grid
}

type textGridRender struct {
	text *TextGrid

	cols, rows int

	cellSize fyne.Size
	objects  []fyne.CanvasObject
}

func newTextCell(str string) *canvas.Text {
	text := canvas.NewText(str, theme.TextColor())
	text.TextStyle.Monospace = true

	if str == "·" || str == "→" || str == "↵" {
		text.Color = theme.PlaceHolderColor()
	}

	return text
}

func (t *textGridRender) ensureGrid() {
	t.cols = t.text.maxCols
	t.rows = t.text.rows()

	for _, bound := range t.text.rowBounds {
		i := 0
		for j := bound[0]; j < bound[1]; j++ {
			t.objects = append(t.objects, newTextCell(string(t.text.buffer[j])))
			i++
		}
		for ; i < t.cols; i++ {
			t.objects = append(t.objects, newTextCell(""))
		}
	}
}

func (t *textGridRender) Layout(size fyne.Size) {
	i := 0
	cellPos := fyne.NewPos(0, 0)
	for y := 0; y < t.rows; y++ {
		for x := 0; x < t.cols; x++ {
			t.objects[i].Move(cellPos)

			cellPos.X += t.cellSize.Width
			i++
		}

		cellPos.X = 0
		cellPos.Y += t.cellSize.Height
	}
}

func (t *textGridRender) MinSize() fyne.Size {
	return fyne.NewSize(t.cellSize.Width*t.cols, t.cellSize.Height*t.rows)
}

func (t *textGridRender) Refresh() {
}

func (t *textGridRender) ApplyTheme() {
}

func (t *textGridRender) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (t *textGridRender) Objects() []fyne.CanvasObject {
	return t.objects
}

func (t *textGridRender) Destroy() {
}

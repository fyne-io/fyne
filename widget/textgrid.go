package widget

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

const (
	textAreaSpaceSymbol   = '·'
	textAreaTabSymbol     = '→'
	textAreaNewLineSymbol = '↵'
)

// TextGrid is a monospaced grid of characters.
// This is designed to be used by a text editor or advanced test presentation.
type TextGrid struct {
	BaseWidget
	Buffer [][]rune

	LineNumbers bool
	Whitespace  bool
}

// MinSize returns the smallest size this widget can shrink to
func (t *TextGrid) MinSize() fyne.Size {
	t.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

// SetText updates the buffer of this textgrid to contain the specified text.
// New lines and columns will be added as required. Lines are separated by '\n'.
func (t *TextGrid) SetText(text string) {
	rows := strings.Split(text, "\n")
	var buffer [][]rune
	for _, row := range rows {
		buffer = append(buffer, []rune(row))
	}

	t.Buffer = buffer
	t.Refresh()
}

// Text returns the contents of the buffer as a single string.
// It reconstructs the lines by joining with a `\n` character.
func (t *TextGrid) Text() string {
	ret := ""
	for i, row := range t.Buffer {
		ret += string(row)

		if i < len(t.Buffer)-1 {
			ret += "\n"
		}
	}

	return ret
}

// Row returns the []rune content of a specified row. If the index is out of bounds it returns an empty slice.
func (t *TextGrid) Row(row int) []rune {
	if row < 0 || row >= len(t.Buffer) {
		return []rune{}
	}

	return t.Buffer[row]
}

// SetRow updates the specified row of the grid's buffer using the specified content and then refreshes.
// If the row is beyond the end of the current buffer it will be expanded.
func (t *TextGrid) SetRow(row int, content []rune) {
	if row < 0 {
		return
	}
	for len(t.Buffer) <= row {
		t.Buffer = append(t.Buffer, []rune{})
	}

	t.Buffer[row] = content
	t.Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (t *TextGrid) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	render := &textGridRender{text: t}
	render.update()

	cell := canvas.NewText("M", color.White)
	cell.TextStyle.Monospace = true
	render.cellSize = cell.MinSize()

	return render
}

// NewTextGrid creates a new empty TextGrid widget.
func NewTextGrid() *TextGrid {
	grid := &TextGrid{}
	grid.ExtendBaseWidget(grid)
	return grid
}

// NewTextGridFromString creates a new TextGrid widget with the specified string content.
func NewTextGridFromString(content string) *TextGrid {
	grid := NewTextGrid()
	grid.SetText(content)
	return grid
}

type textGridRender struct {
	text *TextGrid

	cols, rows int

	cellSize fyne.Size
	objects  []fyne.CanvasObject
}

func (t *textGridRender) appendTextCell(str rune) {
	text := canvas.NewText(string(str), theme.TextColor())
	text.TextStyle.Monospace = true

	t.objects = append(t.objects, text)
}

func (t *textGridRender) setCellRune(str rune, pos int) {
	text := t.objects[pos].(*canvas.Text)
	text.Text = string(str)

	if str == textAreaSpaceSymbol || str == textAreaTabSymbol || str == textAreaNewLineSymbol {
		text.Color = theme.PlaceHolderColor()
	} else {
		text.Color = theme.TextColor()
	}
}

func (t *textGridRender) update() {
	t.ensureGrid()
	t.refreshGrid()
}

func (t *textGridRender) ensureGrid() {
	cellCount := t.cols * t.rows
	if len(t.objects) == cellCount {
		return
	}
	for i := len(t.objects); i < cellCount; i++ {
		t.appendTextCell(' ')
	}
}

func (t *textGridRender) refreshGrid() {
	line := 1
	x := 0

	for rowIndex, row := range t.text.Buffer {
		if rowIndex >= t.rows { // would be an overflow - bad
			break
		}
		i := 0
		if t.text.LineNumbers {
			lineStr := []rune(fmt.Sprintf("%d", line))
			for c := 0; c < len(lineStr); c++ {
				t.setCellRune(lineStr[c], x)
				i++
				x++
			}
			for ; i < t.lineCountWidth(); i++ {
				t.setCellRune(' ', x)
				x++
			}

			t.setCellRune(' ', x)
			i++
			x++
		}
		for _, r := range row {
			if i >= t.cols { // would be an overflow - bad
				continue
			}
			if t.text.Whitespace && r == ' ' {
				r = textAreaSpaceSymbol
			}
			t.setCellRune(r, x)
			i++
			x++
		}
		if t.text.Whitespace && i < t.cols {
			t.setCellRune(textAreaNewLineSymbol, x)
			i++
			x++
		}
		for ; i < t.cols; i++ {
			t.setCellRune(' ', x)
			x++
		}

		line++
	}
	for ; x < len(t.objects); x++ {
		t.setCellRune(' ', x)
	}
	canvas.Refresh(t.text)
}

func (t *textGridRender) lineCountWidth() int {
	return len(fmt.Sprintf("%d", t.rows+1))
}

func (t *textGridRender) updateGridSize(size fyne.Size) {
	bufRows := len(t.text.Buffer)
	bufCols := 0
	for _, row := range t.text.Buffer {
		bufCols = int(math.Max(float64(bufCols), float64(len(row))))
	}
	sizeCols := int(math.Floor(float64(size.Width) / float64(t.cellSize.Width)))
	sizeRows := int(math.Floor(float64(size.Height) / float64(t.cellSize.Height)))

	t.cols = int(math.Max(float64(sizeCols), float64(bufCols)))
	t.rows = int(math.Max(float64(sizeRows), float64(bufRows)))
}

func (t *textGridRender) Layout(size fyne.Size) {
	t.updateGridSize(size)
	t.ensureGrid()

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
	t.refreshGrid()
}

func (t *textGridRender) ApplyTheme() {
}

func (t *textGridRender) BackgroundColor() color.Color {
	return color.Transparent
}

func (t *textGridRender) Objects() []fyne.CanvasObject {
	return t.objects
}

func (t *textGridRender) Destroy() {
}

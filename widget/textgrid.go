package widget

import (
	"image/color"
	"math"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/painter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

const (
	textAreaSpaceSymbol   = '·'
	textAreaTabSymbol     = '→'
	textAreaNewLineSymbol = '↵'
)

var (
	// TextGridStyleDefault is a default style for test grid cells
	TextGridStyleDefault TextGridStyle
	// TextGridStyleWhitespace is the style used for whitespace characters, if enabled
	TextGridStyleWhitespace TextGridStyle
)

// TextGridCell represents a single cell in a text grid.
// It has a rune for the text content and a style associated with it.
type TextGridCell struct {
	Rune  rune
	Style TextGridStyle
}

// TextGridRow represents a row of cells cell in a text grid.
// It contains the cells for the row and an optional style.
type TextGridRow struct {
	Cells []TextGridCell
	Style TextGridStyle
}

// TextGridStyle defines a style that can be applied to a TextGrid cell.
type TextGridStyle interface {
	TextColor() color.Color
	BackgroundColor() color.Color
}

// CustomTextGridStyle is a utility type for those not wanting to define their own style types.
type CustomTextGridStyle struct {
	FGColor, BGColor color.Color
}

// TextColor is the color a cell should use for the text.
func (c *CustomTextGridStyle) TextColor() color.Color {
	return c.FGColor
}

// BackgroundColor is the color a cell should use for the background.
func (c *CustomTextGridStyle) BackgroundColor() color.Color {
	return c.BGColor
}

// TextGrid is a monospaced grid of characters.
// This is designed to be used by a text editor, code preview or terminal emulator.
type TextGrid struct {
	BaseWidget
	Rows []TextGridRow

	ShowLineNumbers bool
	ShowWhitespace  bool
	TabWidth        int // If set to 0 the fyne.DefaultTabWidth is used
}

// MinSize returns the smallest size this widget can shrink to
func (t *TextGrid) MinSize() fyne.Size {
	t.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

// Resize is called when this widget changes size. We should make sure that we refresh cells.
func (t *TextGrid) Resize(size fyne.Size) {
	t.BaseWidget.Resize(size)
	t.Refresh()
}

// SetText updates the buffer of this textgrid to contain the specified text.
// New lines and columns will be added as required. Lines are separated by '\n'.
// The grid will use default text style and any previous content and style will be removed.
// Tab characters are padded with spaces to the next tab stop.
func (t *TextGrid) SetText(text string) {
	lines := strings.Split(text, "\n")
	rows := make([]TextGridRow, len(lines))
	for i, line := range lines {
		cells := make([]TextGridCell, 0, len(line))
		for _, r := range line {
			cells = append(cells, TextGridCell{Rune: r})
			if r == '\t' {
				col := len(cells)
				next := nextTab(col-1, t.tabWidth())
				for i := col; i < next; i++ {
					cells = append(cells, TextGridCell{Rune: ' '})
				}
			}
		}
		rows[i] = TextGridRow{Cells: cells}
	}

	t.Rows = rows
	t.Refresh()
}

// Text returns the contents of the buffer as a single string (with no style information).
// It reconstructs the lines by joining with a `\n` character.
// Tab characters have padded spaces removed.
func (t *TextGrid) Text() string {
	count := len(t.Rows) - 1 // newlines
	for _, row := range t.Rows {
		count += len(row.Cells)
	}

	if count <= 0 {
		return ""
	}

	runes := make([]rune, 0, count)

	for i, row := range t.Rows {
		next := 0
		for col, cell := range row.Cells {
			if col < next {
				continue
			}
			runes = append(runes, cell.Rune)
			if cell.Rune == '\t' {
				next = nextTab(col, t.tabWidth())
			}
		}
		if i < len(t.Rows)-1 {
			runes = append(runes, '\n')
		}
	}

	return string(runes)
}

// Row returns a copy of the content in a specified row as a TextGridRow.
// If the index is out of bounds it returns an empty row object.
func (t *TextGrid) Row(row int) TextGridRow {
	if row < 0 || row >= len(t.Rows) {
		return TextGridRow{}
	}

	return t.Rows[row]
}

// RowText returns a string representation of the content at the row specified.
// If the index is out of bounds it returns an empty string.
func (t *TextGrid) RowText(row int) string {
	rowData := t.Row(row)
	count := len(rowData.Cells)

	if count <= 0 {
		return ""
	}

	runes := make([]rune, 0, count)

	next := 0
	for col, cell := range rowData.Cells {
		if col < next {
			continue
		}
		runes = append(runes, cell.Rune)
		if cell.Rune == '\t' {
			next = nextTab(col, t.tabWidth())
		}
	}
	return string(runes)
}

// SetRow updates the specified row of the grid's contents using the specified content and style and then refreshes.
// If the row is beyond the end of the current buffer it will be expanded.
// Tab characters are not padded with spaces.
func (t *TextGrid) SetRow(row int, content TextGridRow) {
	if row < 0 {
		return
	}
	for len(t.Rows) <= row {
		t.Rows = append(t.Rows, TextGridRow{})
	}

	t.Rows[row] = content
	for col := 0; col > len(content.Cells); col++ {
		t.refreshCell(row, col)
	}
}

// SetRowStyle sets a grid style to all the cells cell at the specified row.
// Any cells in this row with their own style will override this value when displayed.
func (t *TextGrid) SetRowStyle(row int, style TextGridStyle) {
	if row < 0 {
		return
	}
	for len(t.Rows) <= row {
		t.Rows = append(t.Rows, TextGridRow{})
	}
	t.Rows[row].Style = style
}

// SetCell sets a grid data to the cell at named row and column.
func (t *TextGrid) SetCell(row, col int, cell TextGridCell) {
	if row < 0 || col < 0 {
		return
	}
	t.ensureCells(row, col)

	t.Rows[row].Cells[col] = cell
	t.refreshCell(row, col)
}

// SetRune sets a character to the cell at named row and column.
func (t *TextGrid) SetRune(row, col int, r rune) {
	if row < 0 || col < 0 {
		return
	}
	t.ensureCells(row, col)

	t.Rows[row].Cells[col].Rune = r
	t.refreshCell(row, col)
}

// SetStyle sets a grid style to the cell at named row and column.
func (t *TextGrid) SetStyle(row, col int, style TextGridStyle) {
	if row < 0 || col < 0 {
		return
	}
	t.ensureCells(row, col)

	t.Rows[row].Cells[col].Style = style
	t.refreshCell(row, col)
}

// SetStyleRange sets a grid style to all the cells between the start row and column through to the end row and column.
func (t *TextGrid) SetStyleRange(startRow, startCol, endRow, endCol int, style TextGridStyle) {
	if startRow >= len(t.Rows) || endRow < 0 {
		return
	}
	if startRow < 0 {
		startRow = 0
		startCol = 0
	}
	if endRow >= len(t.Rows) {
		endRow = len(t.Rows) - 1
		endCol = len(t.Rows[endRow].Cells) - 1
	}

	if startRow == endRow {
		for col := startCol; col <= endCol; col++ {
			t.SetStyle(startRow, col, style)
		}
		return
	}

	// first row
	for col := startCol; col < len(t.Rows[startRow].Cells); col++ {
		t.SetStyle(startRow, col, style)
	}

	// possible middle rows
	for rowNum := startRow + 1; rowNum < endRow; rowNum++ {
		for col := 0; col < len(t.Rows[rowNum].Cells); col++ {
			t.SetStyle(rowNum, col, style)
		}
	}

	// last row
	for col := 0; col <= endCol; col++ {
		t.SetStyle(endRow, col, style)
	}
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (t *TextGrid) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	render := &textGridRenderer{text: t}
	render.updateCellSize()

	TextGridStyleDefault = &CustomTextGridStyle{}
	TextGridStyleWhitespace = &CustomTextGridStyle{FGColor: theme.DisabledColor()}

	return render
}

func (t *TextGrid) ensureCells(row, col int) {
	for len(t.Rows) <= row {
		t.Rows = append(t.Rows, TextGridRow{})
	}
	data := t.Rows[row]

	for len(data.Cells) <= col {
		data.Cells = append(data.Cells, TextGridCell{})
		t.Rows[row] = data
	}
}

func (t *TextGrid) refreshCell(row, col int) {
	r := cache.Renderer(t).(*textGridRenderer)
	r.refreshCell(row, col)
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

// nextTab finds the column of the next tab stop for the given column
func nextTab(column int, tabWidth int) int {
	tabStop, _ := math.Modf(float64(column+tabWidth) / float64(tabWidth))
	return tabWidth * int(tabStop)
}

type textGridRenderer struct {
	text *TextGrid

	cols, rows int

	cellSize fyne.Size
	objects  []fyne.CanvasObject
	current  fyne.Canvas
}

func (t *textGridRenderer) appendTextCell(str rune) {
	text := canvas.NewText(string(str), theme.ForegroundColor())
	text.TextStyle.Monospace = true

	bg := canvas.NewRectangle(color.Transparent)
	t.objects = append(t.objects, bg, text)
}

func (t *textGridRenderer) refreshCell(row, col int) {
	pos := row*t.cols + col
	if pos*2+1 >= len(t.objects) {
		return
	}

	cell := t.text.Rows[row].Cells[col]
	t.setCellRune(cell.Rune, pos, cell.Style, t.text.Rows[row].Style)
}

func (t *textGridRenderer) setCellRune(str rune, pos int, style, rowStyle TextGridStyle) {
	if str == 0 {
		str = ' '
	}

	text := t.objects[pos*2+1].(*canvas.Text)
	text.TextSize = theme.TextSize()
	fg := theme.ForegroundColor()
	if style != nil && style.TextColor() != nil {
		fg = style.TextColor()
	} else if rowStyle != nil && rowStyle.TextColor() != nil {
		fg = rowStyle.TextColor()
	}
	newStr := string(str)
	if text.Text != newStr || text.Color != fg {
		text.Text = newStr
		text.Color = fg
		t.refresh(text)
	}

	rect := t.objects[pos*2].(*canvas.Rectangle)
	bg := color.Color(color.Transparent)
	if style != nil && style.BackgroundColor() != nil {
		bg = style.BackgroundColor()
	} else if rowStyle != nil && rowStyle.BackgroundColor() != nil {
		bg = rowStyle.BackgroundColor()
	}
	if rect.FillColor != bg {
		rect.FillColor = bg
		t.refresh(rect)
	}
}

func (t *textGridRenderer) addCellsIfRequired() {
	cellCount := t.cols * t.rows
	if len(t.objects) == cellCount*2 {
		return
	}
	for i := len(t.objects); i < cellCount*2; i += 2 {
		t.appendTextCell(' ')
	}
}

func (t *textGridRenderer) refreshGrid() {
	line := 1
	x := 0

	for rowIndex, row := range t.text.Rows {
		rowStyle := row.Style
		i := 0
		if t.text.ShowLineNumbers {
			lineStr := []rune(strconv.Itoa(line))
			pad := t.lineNumberWidth() - len(lineStr)
			for ; i < pad; i++ {
				t.setCellRune(' ', x, TextGridStyleWhitespace, rowStyle) // padding space
				x++
			}
			for c := 0; c < len(lineStr); c++ {
				t.setCellRune(lineStr[c], x, TextGridStyleDefault, rowStyle) // line numbers
				i++
				x++
			}

			t.setCellRune('|', x, TextGridStyleWhitespace, rowStyle) // last space
			i++
			x++
		}
		for _, r := range row.Cells {
			if i >= t.cols { // would be an overflow - bad
				continue
			}
			if t.text.ShowWhitespace && (r.Rune == ' ' || r.Rune == '\t') {
				sym := textAreaSpaceSymbol
				if r.Rune == '\t' {
					sym = textAreaTabSymbol
				}

				if r.Style != nil && r.Style.BackgroundColor() != nil {
					whitespaceBG := &CustomTextGridStyle{FGColor: TextGridStyleWhitespace.TextColor(),
						BGColor: r.Style.BackgroundColor()}
					t.setCellRune(sym, x, whitespaceBG, rowStyle) // whitespace char
				} else {
					t.setCellRune(sym, x, TextGridStyleWhitespace, rowStyle) // whitespace char
				}
			} else {
				t.setCellRune(r.Rune, x, r.Style, rowStyle) // regular char
			}
			i++
			x++
		}
		if t.text.ShowWhitespace && i < t.cols && rowIndex < len(t.text.Rows)-1 {
			t.setCellRune(textAreaNewLineSymbol, x, TextGridStyleWhitespace, rowStyle) // newline
			i++
			x++
		}
		for ; i < t.cols; i++ {
			t.setCellRune(' ', x, TextGridStyleDefault, rowStyle) // blanks
			x++
		}

		line++
	}
	for ; x < len(t.objects)/2; x++ {
		t.setCellRune(' ', x, TextGridStyleDefault, nil) // trailing cells and blank lines
	}
}

// tabWidth either returns the set tab width or if not set the returns the DefaultTabWidth
func (t *TextGrid) tabWidth() int {
	if t.TabWidth == 0 {
		return painter.DefaultTabWidth
	}
	return t.TabWidth
}

func (t *textGridRenderer) lineNumberWidth() int {
	return len(strconv.Itoa(t.rows + 1))
}

func (t *textGridRenderer) updateGridSize(size fyne.Size) {
	bufRows := len(t.text.Rows)
	bufCols := 0
	for _, row := range t.text.Rows {
		bufCols = int(math.Max(float64(bufCols), float64(len(row.Cells))))
	}
	sizeCols := math.Floor(float64(size.Width) / float64(t.cellSize.Width))
	sizeRows := math.Floor(float64(size.Height) / float64(t.cellSize.Height))

	if t.text.ShowWhitespace {
		bufCols++
	}
	if t.text.ShowLineNumbers {
		bufCols += t.lineNumberWidth()
	}

	t.cols = int(math.Max(sizeCols, float64(bufCols)))
	t.rows = int(math.Max(sizeRows, float64(bufRows)))
	t.addCellsIfRequired()
}

func (t *textGridRenderer) Layout(size fyne.Size) {
	t.updateGridSize(size)

	i := 0
	cellPos := fyne.NewPos(0, 0)
	for y := 0; y < t.rows; y++ {
		for x := 0; x < t.cols; x++ {
			t.objects[i*2+1].Move(cellPos)

			t.objects[i*2].Resize(t.cellSize)
			t.objects[i*2].Move(cellPos)
			cellPos.X += t.cellSize.Width
			i++
		}

		cellPos.X = 0
		cellPos.Y += t.cellSize.Height
	}
}

func (t *textGridRenderer) MinSize() fyne.Size {
	longestRow := float32(0)
	for _, row := range t.text.Rows {
		longestRow = fyne.Max(longestRow, float32(len(row.Cells)))
	}
	return fyne.NewSize(t.cellSize.Width*longestRow,
		t.cellSize.Height*float32(len(t.text.Rows)))
}

func (t *textGridRenderer) Refresh() {
	// we may be on a new canvas, so just update it to be sure
	if fyne.CurrentApp() != nil && fyne.CurrentApp().Driver() != nil {
		t.current = fyne.CurrentApp().Driver().CanvasForObject(t.text)
	}

	// theme could change text size
	t.updateCellSize()

	TextGridStyleWhitespace = &CustomTextGridStyle{FGColor: theme.DisabledColor()}
	t.updateGridSize(t.text.size)
	t.refreshGrid()
}

func (t *textGridRenderer) ApplyTheme() {
}

func (t *textGridRenderer) Objects() []fyne.CanvasObject {
	return t.objects
}

func (t *textGridRenderer) Destroy() {
}

func (t *textGridRenderer) refresh(obj fyne.CanvasObject) {
	if t.current == nil {
		if fyne.CurrentApp() != nil && fyne.CurrentApp().Driver() != nil {
			// cache canvas for this widget, so we don't look it up many times for every cell/row refresh!
			t.current = fyne.CurrentApp().Driver().CanvasForObject(t.text)
		}

		if t.current == nil {
			return // not yet set up perhaps?
		}
	}

	t.current.Refresh(obj)
}

func (t *textGridRenderer) updateCellSize() {
	size := fyne.MeasureText("M", theme.TextSize(), fyne.TextStyle{Monospace: true})

	// round it for seamless background
	size.Width = float32(math.Round(float64((size.Width))))
	size.Height = float32(math.Round(float64((size.Height))))

	t.cellSize = size
}

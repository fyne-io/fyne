package widget

import (
	"image/color"
	"math"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/widget"
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
	Style() fyne.TextStyle
	TextColor() color.Color
	BackgroundColor() color.Color
}

// CustomTextGridStyle is a utility type for those not wanting to define their own style types.
type CustomTextGridStyle struct {
	// Since: 2.5
	TextStyle        fyne.TextStyle
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

// Style is the text style a cell should use.
func (c *CustomTextGridStyle) Style() fyne.TextStyle {
	return c.TextStyle
}

// TextGrid is a monospaced grid of characters.
// This is designed to be used by a text editor, code preview or terminal emulator.
type TextGrid struct {
	BaseWidget
	Rows []TextGridRow

	scroll  *widget.Scroll
	content *textGridContent

	ShowLineNumbers bool
	ShowWhitespace  bool
	TabWidth        int // If set to 0 the fyne.DefaultTabWidth is used

	// Scroll can be used to turn off the scrolling of our TextGrid.
	//
	// Since: 2.6
	Scroll fyne.ScrollDirection
}

// Append will add new lines to the end of this TextGrid.
// The first character will be at the beginning of a new line and any newline characters will split the text further.
//
// Since: 2.6
func (t *TextGrid) Append(text string) {
	rows := t.parseRows(text)

	t.Rows = append(t.Rows, rows...)
	t.Refresh()
}

// CursorLocationForPosition returns the location where a cursor would be if it was located in the cell under the
// requested position. If the grid is scrolled the position will refer to the visible offset and not the distance
// from the top left of the overall document.
//
// Since: 2.6
func (t *TextGrid) CursorLocationForPosition(p fyne.Position) (row, col int) {
	y := t.content.cellSize.Height
	x := t.content.cellSize.Width

	if t.scroll != nil && t.scroll.Visible() {
		y += t.scroll.Offset.Y
		x += t.scroll.Offset.X
	}

	row = int(p.Y / y)
	col = int(p.X / x)
	return
}

// PositionForCursorLocation returns the relative position in this TextGrid for the cell at position row, col.
// If the grid has been scrolled this will be taken into account so that the position compared to top left will
// refer to the requested location.
//
// Since: 2.6
func (t *TextGrid) PositionForCursorLocation(row, col int) fyne.Position {
	y := float32(row) * t.content.cellSize.Height
	x := float32(col) * t.content.cellSize.Width

	if t.scroll != nil && t.scroll.Visible() {
		y -= t.scroll.Offset.Y
		x -= t.scroll.Offset.X
	}

	return fyne.NewPos(x, y)
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
	rows := t.parseRows(text)

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

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (t *TextGrid) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)

	th := t.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	TextGridStyleDefault = &CustomTextGridStyle{}
	TextGridStyleWhitespace = &CustomTextGridStyle{FGColor: th.Color(theme.ColorNameDisabled, v)}

	var scroll *widget.Scroll
	content := newTextGridContent(t)
	objs := make([]fyne.CanvasObject, 1)
	if t.Scroll == widget.ScrollNone {
		scroll = widget.NewScroll(nil)
		objs[0] = content
	} else {
		scroll = widget.NewScroll(content)
		scroll.Direction = t.Scroll
		objs[0] = scroll
	}
	t.scroll = scroll
	t.content = content
	r := &textGridRenderer{text: content, scroll: scroll}
	r.SetObjects(objs)
	return r
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

func (t *TextGrid) parseRows(text string) []TextGridRow {
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

	return rows
}

func (t *TextGrid) refreshCell(row, col int) {
	r := cache.Renderer(t).(*textGridRenderer).text
	r.refreshCell(row, col)
}

// NewTextGrid creates a new empty TextGrid widget.
func NewTextGrid() *TextGrid {
	grid := &TextGrid{}
	grid.Scroll = widget.ScrollNone
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
	widget.BaseRenderer

	text   *textGridContent
	scroll *widget.Scroll
}

func (t *textGridRenderer) Layout(s fyne.Size) {
	t.Objects()[0].Resize(s)
}

func (t *textGridRenderer) MinSize() fyne.Size {
	if t.text.text.Scroll == widget.ScrollNone {
		return t.text.MinSize()
	}

	return t.scroll.MinSize()
}

func (t *textGridRenderer) Refresh() {
	content := t.text
	if t.text.text.Scroll != widget.ScrollNone {
		t.scroll.Direction = t.text.text.Scroll
	}
	if t.text.text.Scroll == widget.ScrollNone && t.scroll.Content != nil {
		t.scroll.Hide()
		t.scroll.Content = nil
		content.Resize(t.text.Size())
		t.SetObjects([]fyne.CanvasObject{t.text})
	} else if (t.text.text.Scroll != widget.ScrollNone) && t.scroll.Content == nil {
		t.scroll.Content = content
		t.scroll.Show()

		t.scroll.Resize(t.text.Size())
		content.Resize(content.MinSize())
		t.SetObjects([]fyne.CanvasObject{t.scroll})
	}

	canvas.Refresh(t.text.text.super())
	t.text.Refresh()
}

type textGridContent struct {
	BaseWidget
	text *TextGrid

	rows     int
	cellSize fyne.Size
}

func newTextGridContent(t *TextGrid) *textGridContent {
	grid := &textGridContent{text: t}
	grid.ExtendBaseWidget(grid)
	return grid
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (t *textGridContent) CreateRenderer() fyne.WidgetRenderer {
	r := &textGridContentRenderer{text: t}

	r.updateCellSize()
	t.text.scroll.OnScrolled = func(_ fyne.Position) {
		r.addRowsIfRequired()
		r.Layout(t.Size())
	}
	return r
}

func (t *textGridContent) refreshCell(row, col int) {
	rows := cache.Renderer(t).Objects()
	if row >= len(rows)-1 {
		return
	}
	wid := rows[row].(*textGridRow)
	r := cache.Renderer(wid).(*textGridRowRenderer)
	r.refreshCell(col)
}

func (t *textGridContentRenderer) updateGridSize(size fyne.Size) {
	bufRows := len(t.text.text.Rows)
	bufCols := 0
	for _, row := range t.text.text.Rows {
		bufCols = int(math.Max(float64(bufCols), float64(len(row.Cells))))
	}
	sizeRows := math.Floor(float64(size.Height) / float64(t.text.cellSize.Height))

	t.text.rows = int(math.Max(sizeRows, float64(bufRows)))
	t.addRowsIfRequired()
}

type textGridContentRenderer struct {
	text *textGridContent

	visible  []fyne.CanvasObject
	itemPool async.Pool[*textGridRow]
}

func (t *textGridContentRenderer) Destroy() {
}

func (t *textGridContentRenderer) Layout(s fyne.Size) {
	size := fyne.NewSize(s.Width, t.text.cellSize.Height)
	t.updateGridSize(s)

	for _, o := range t.visible {
		o.Move(fyne.NewPos(0, float32(o.(*textGridRow).row)*t.text.cellSize.Height))
		o.Resize(size)
	}
}

func (t *textGridContentRenderer) MinSize() fyne.Size {
	longestRow := float32(0)
	for _, row := range t.text.text.Rows {
		longestRow = fyne.Max(longestRow, float32(len(row.Cells)))
	}
	return fyne.NewSize(t.text.cellSize.Width*longestRow,
		t.text.cellSize.Height*float32(len(t.text.text.Rows)))
}

func (t *textGridContentRenderer) Objects() []fyne.CanvasObject {
	return t.visible
}

func (t *textGridContentRenderer) Refresh() {
	// theme could change text size
	t.updateCellSize()
	t.updateGridSize(t.text.text.Size())

	for _, o := range t.visible {
		o.Refresh()
	}
}

func (t *textGridContentRenderer) addRowsIfRequired() {
	start := 0
	end := t.text.rows
	if t.text.text.Scroll == widget.ScrollBoth || t.text.text.Scroll == widget.ScrollVerticalOnly {
		off := t.text.text.scroll.Offset.Y
		start = int(math.Floor(float64(off / t.text.cellSize.Height)))

		off += t.text.text.Size().Height
		end = int(math.Ceil(float64(off / t.text.cellSize.Height)))
	}

	remain := t.visible[:0]
	for _, row := range t.visible {
		if row.(*textGridRow).row < start || row.(*textGridRow).row > end {
			t.itemPool.Put(row.(*textGridRow))
			continue
		}

		remain = append(remain, row.(*textGridRow))
	}
	t.visible = remain

	var newItems []fyne.CanvasObject
	for i := start; i <= end; i++ {
		found := false
		for _, row := range t.visible {
			if i == row.(*textGridRow).row {
				found = true
				break
			}
		}

		if found {
			continue
		}

		newRow := t.itemPool.Get()
		if newRow == nil {
			newRow = newTextGridRow(t.text, i)
		} else {
			newRow.setRow(i)
		}
		newItems = append(newItems, newRow)
	}

	if len(newItems) > 0 {
		t.visible = append(t.visible, newItems...)
	}
}

func (t *textGridContentRenderer) updateCellSize() {
	th := t.text.Theme()
	size := fyne.MeasureText("M", th.Size(theme.SizeNameText), fyne.TextStyle{Monospace: true})

	// round it for seamless background
	size.Width = float32(math.Round(float64(size.Width)))
	size.Height = float32(math.Round(float64(size.Height)))

	t.text.cellSize = size
}

type textGridRow struct {
	BaseWidget
	text *textGridContent

	row int
}

func newTextGridRow(t *textGridContent, row int) *textGridRow {
	return &textGridRow{text: t, row: row}
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (t *textGridRow) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	render := &textGridRowRenderer{obj: t}

	render.Refresh() // populate
	return render
}

func (t *textGridRow) setRow(row int) {
	t.row = row
	t.Refresh()
}

type textGridRowRenderer struct {
	obj *textGridRow

	cols int

	objects []fyne.CanvasObject
}

func (t *textGridRowRenderer) appendTextCell(str rune) {
	th := t.obj.text.text.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	text := canvas.NewText(string(str), th.Color(theme.ColorNameForeground, v))
	text.TextStyle.Monospace = true

	bg := canvas.NewRectangle(color.Transparent)

	ul := canvas.NewLine(color.Transparent)

	t.objects = append(t.objects, bg, text, ul)
}

func (t *textGridRowRenderer) refreshCell(col int) {
	pos := t.cols + col
	if pos*3+1 >= len(t.objects) {
		return
	}

	cell := t.obj.text.text.Rows[t.obj.row].Cells[col]
	t.setCellRune(cell.Rune, pos, cell.Style, t.obj.text.text.Rows[t.obj.row].Style)
}

func (t *textGridRowRenderer) setCellRune(str rune, pos int, style, rowStyle TextGridStyle) {
	if str == 0 {
		str = ' '
	}
	rect := t.objects[pos*3].(*canvas.Rectangle)
	text := t.objects[pos*3+1].(*canvas.Text)
	underline := t.objects[pos*3+2].(*canvas.Line)

	th := t.obj.text.text.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	fg := th.Color(theme.ColorNameForeground, v)
	text.TextSize = th.Size(theme.SizeNameText)

	var underlineStrokeWidth float32 = 1
	var underlineStrokeColor color.Color = color.Transparent
	textStyle := fyne.TextStyle{}
	if style != nil {
		textStyle = style.Style()
	} else if rowStyle != nil {
		textStyle = rowStyle.Style()
	}
	if textStyle.Bold {
		underlineStrokeWidth = 2
	}
	if textStyle.Underline {
		underlineStrokeColor = fg
	}
	textStyle.Monospace = true

	if style != nil && style.TextColor() != nil {
		fg = style.TextColor()
	} else if rowStyle != nil && rowStyle.TextColor() != nil {
		fg = rowStyle.TextColor()
	}

	newStr := string(str)
	if text.Text != newStr || text.Color != fg || textStyle != text.TextStyle {
		text.Text = newStr
		text.Color = fg
		text.TextStyle = textStyle
		text.Refresh()
	}

	if underlineStrokeWidth != underline.StrokeWidth || underlineStrokeColor != underline.StrokeColor {
		underline.StrokeWidth, underline.StrokeColor = underlineStrokeWidth, underlineStrokeColor
		underline.Refresh()
	}

	bg := color.Color(color.Transparent)
	if style != nil && style.BackgroundColor() != nil {
		bg = style.BackgroundColor()
	} else if rowStyle != nil && rowStyle.BackgroundColor() != nil {
		bg = rowStyle.BackgroundColor()
	}
	if rect.FillColor != bg {
		rect.FillColor = bg
		rect.Refresh()
	}
}

func (t *textGridRowRenderer) addCellsIfRequired() {
	cellCount := t.cols
	if len(t.objects) == cellCount*3 {
		return
	}
	for i := len(t.objects); i < cellCount*3; i += 3 {
		t.appendTextCell(' ')
	}
}

func (t *textGridRowRenderer) refreshCells() {
	x := 0
	if t.obj.row >= len(t.obj.text.text.Rows) {
		return // we can have more rows than content rows (filling space)
	}

	row := t.obj.text.text.Rows[t.obj.row]
	rowStyle := row.Style
	i := 0
	if t.obj.text.text.ShowLineNumbers {
		lineStr := []rune(strconv.Itoa(t.obj.row + 1))
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
		if t.obj.text.text.ShowWhitespace && (r.Rune == ' ' || r.Rune == '\t') {
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
	if t.obj.text.text.ShowWhitespace && i < t.cols && t.obj.row < len(t.obj.text.text.Rows)-1 {
		t.setCellRune(textAreaNewLineSymbol, x, TextGridStyleWhitespace, rowStyle) // newline
		i++
		x++
	}
	for ; i < t.cols; i++ {
		t.setCellRune(' ', x, TextGridStyleDefault, rowStyle) // blanks
		x++
	}

	for ; x < len(t.objects)/3; x++ {
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

func (t *textGridRowRenderer) lineNumberWidth() int {
	return len(strconv.Itoa(t.obj.text.rows + 1))
}

func (t *textGridRowRenderer) updateGridSize(size fyne.Size) {
	bufCols := 0
	for _, row := range t.obj.text.text.Rows {
		bufCols = int(math.Max(float64(bufCols), float64(len(row.Cells))))
	}
	sizeCols := math.Floor(float64(size.Width) / float64(t.obj.text.cellSize.Width))

	if t.obj.text.text.ShowWhitespace {
		bufCols++
	}
	if t.obj.text.text.ShowLineNumbers {
		bufCols += t.lineNumberWidth()
	}

	t.cols = int(math.Max(sizeCols, float64(bufCols)))
	t.addCellsIfRequired()
}

func (t *textGridRowRenderer) Layout(size fyne.Size) {
	t.updateGridSize(size)

	cellPos := fyne.NewPos(0, 0)
	off := 0
	for x := 0; x < t.cols; x++ {
		// rect
		t.objects[off].Resize(t.obj.text.cellSize)
		t.objects[off].Move(cellPos)

		// text
		t.objects[off+1].Move(cellPos)

		// underline
		t.objects[off+2].Move(cellPos.Add(fyne.Position{X: 0, Y: t.obj.text.cellSize.Height}))
		t.objects[off+2].Resize(fyne.Size{Width: t.obj.text.cellSize.Width})

		cellPos.X += t.obj.text.cellSize.Width
		off += 3
	}
}

func (t *textGridRowRenderer) MinSize() fyne.Size {
	longestRow := float32(0)
	for _, row := range t.obj.text.text.Rows {
		longestRow = fyne.Max(longestRow, float32(len(row.Cells)))
	}
	return fyne.NewSize(t.obj.text.cellSize.Width*longestRow, t.obj.text.cellSize.Height)
}

func (t *textGridRowRenderer) Refresh() {
	th := t.obj.text.text.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	TextGridStyleWhitespace = &CustomTextGridStyle{FGColor: th.Color(theme.ColorNameDisabled, v)}
	t.updateGridSize(t.obj.text.text.Size())
	t.refreshCells()
}

func (t *textGridRowRenderer) ApplyTheme() {
}

func (t *textGridRowRenderer) Objects() []fyne.CanvasObject {
	return t.objects
}

func (t *textGridRowRenderer) Destroy() {
}

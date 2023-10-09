package widget

import (
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

const noCellMatch = math.MaxInt32 // TODO make this MaxInt once we move to newer Go version

// allTableCellsID represents all table cells when refreshing requested cells
var allTableCellsID = TableCellID{-1, -1}

// Declare conformity with interfaces
var _ desktop.Cursorable = (*Table)(nil)
var _ fyne.Draggable = (*Table)(nil)
var _ fyne.Focusable = (*Table)(nil)
var _ desktop.Hoverable = (*Table)(nil)
var _ fyne.Tappable = (*Table)(nil)
var _ fyne.Widget = (*Table)(nil)

// TableCellID is a type that represents a cell's position in a table based on its row and column location.
type TableCellID struct {
	Row int
	Col int
}

// Table widget is a grid of items that can be scrolled and a cell selected.
// Its performance is provided by caching cell templates created with CreateCell and re-using them with UpdateCell.
// The size of the content rows/columns is returned by the Length callback.
//
// Since: 1.4
type Table struct {
	BaseWidget

	Length       func() (rows int, cols int)                      `json:"-"`
	CreateCell   func() fyne.CanvasObject                         `json:"-"`
	UpdateCell   func(id TableCellID, template fyne.CanvasObject) `json:"-"`
	OnSelected   func(id TableCellID)                             `json:"-"`
	OnUnselected func(id TableCellID)                             `json:"-"`

	// ShowHeaderRow specifies that a row should be added to the table with header content.
	// This will default to an A-Z style content, unless overridden with `CreateHeader` and `UpdateHeader` calls.
	//
	// Since: 2.4
	ShowHeaderRow bool

	// ShowHeaderColumn specifies that a column should be added to the table with header content.
	// This will default to an 1-10 style numeric content, unless overridden with `CreateHeader` and `UpdateHeader` calls.
	//
	// Since: 2.4
	ShowHeaderColumn bool

	// CreateHeader is an optional function that allows overriding of the default header widget.
	// Developers must also override `UpdateHeader`.
	//
	// Since: 2.4
	CreateHeader func() fyne.CanvasObject `json:"-"`

	// UpdateHeader is used with `CreateHeader` to support custom header content.
	// The `id` parameter will have `-1` value to indicate a header, and `> 0` where the column or row refer to data.
	//
	// Since: 2.4
	UpdateHeader func(id TableCellID, template fyne.CanvasObject) `json:"-"`

	// StickyRowCount specifies how many data rows should not scroll when the content moves.
	// If `ShowHeaderRow` us `true` then the stuck row will appear immediately underneath.
	//
	// Since: 2.4
	StickyRowCount int

	// StickyColumnCount specifies how many data columns should not scroll when the content moves.
	// If `ShowHeaderColumn` us `true` then the stuck column will appear immediately next to the header.
	//
	// Since: 2.4
	StickyColumnCount int

	currentFocus              TableCellID
	focused                   bool
	selectedCell, hoveredCell *TableCellID
	cells                     *tableCells
	columnWidths, rowHeights  map[int]float32
	moveCallback              func()
	offset                    fyne.Position
	content                   *widget.Scroll

	cellSize, headerSize                                         fyne.Size
	stuckXOff, stuckYOff, stuckWidth, stuckHeight, dragStartSize float32
	top, left, corner, dividerLayer                              *clip
	hoverHeaderRow, hoverHeaderCol, dragCol, dragRow             int
	dragStartPos                                                 fyne.Position
}

// NewTable returns a new performant table widget defined by the passed functions.
// The first returns the data size in rows and columns, second parameter is a function that returns cell
// template objects that can be cached and the third is used to apply data at specified data location to the
// passed template CanvasObject.
//
// Since: 1.4
func NewTable(length func() (rows int, cols int), create func() fyne.CanvasObject, update func(TableCellID, fyne.CanvasObject)) *Table {
	t := &Table{Length: length, CreateCell: create, UpdateCell: update}
	t.ExtendBaseWidget(t)
	return t
}

// NewTableWithHeaders returns a new performant table widget defined by the passed functions including sticky headers.
// The first returns the data size in rows and columns, second parameter is a function that returns cell
// template objects that can be cached and the third is used to apply data at specified data location to the
// passed template CanvasObject.
// The row and column headers will stick to the leading and top edges of the table and contain "1-10" and "A-Z" formatted labels.
//
// Since: 2.4
func NewTableWithHeaders(length func() (rows int, cols int), create func() fyne.CanvasObject, update func(TableCellID, fyne.CanvasObject)) *Table {
	t := NewTable(length, create, update)
	t.ShowHeaderRow = true
	t.ShowHeaderColumn = true

	return t
}

// CreateRenderer returns a new renderer for the table.
//
// Implements: fyne.Widget
func (t *Table) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)

	t.propertyLock.Lock()
	t.headerSize = t.createHeader().MinSize()
	t.cellSize = t.templateSize()
	t.cells = newTableCells(t)
	t.content = widget.NewScroll(t.cells)
	t.top = newClip(t, &fyne.Container{})
	t.left = newClip(t, &fyne.Container{})
	t.corner = newClip(t, &fyne.Container{})
	t.dividerLayer = newClip(t, &fyne.Container{})
	t.propertyLock.Unlock()
	t.dragCol = noCellMatch
	t.dragRow = noCellMatch

	r := &tableRenderer{t: t}
	r.SetObjects([]fyne.CanvasObject{t.top, t.left, t.corner, t.dividerLayer, t.content})
	t.content.OnScrolled = func(pos fyne.Position) {
		t.offset = pos
		t.cells.Refresh()
	}

	r.Layout(t.Size())
	return r
}

func (t *Table) Cursor() desktop.Cursor {
	if t.hoverHeaderRow != noCellMatch {
		return desktop.VResizeCursor
	} else if t.hoverHeaderCol != noCellMatch {
		return desktop.HResizeCursor
	}

	return desktop.DefaultCursor
}

func (t *Table) Dragged(e *fyne.DragEvent) {
	t.propertyLock.Lock()
	min := t.cellSize
	col := t.dragCol
	row := t.dragRow
	startPos := t.dragStartPos
	startSize := t.dragStartSize
	t.propertyLock.Unlock()

	if col != noCellMatch {
		newSize := startSize + (e.Position.X - startPos.X)
		if newSize < min.Width {
			newSize = min.Width
		}
		t.SetColumnWidth(t.dragCol, newSize)
	}
	if row != noCellMatch {
		newSize := startSize + (e.Position.Y - startPos.Y)
		if newSize < min.Height {
			newSize = min.Height
		}
		t.SetRowHeight(t.dragRow, newSize)
	}
}

func (t *Table) DragEnd() {
	t.dragCol = noCellMatch
	t.dragRow = noCellMatch
}

// FocusGained is called after this table has gained focus.
//
// Implements: fyne.Focusable
func (t *Table) FocusGained() {
	t.focused = true
	t.ScrollTo(t.currentFocus)
	t.RefreshItem(t.currentFocus)
}

// FocusLost is called after this Table has lost focus.
//
// Implements: fyne.Focusable
func (t *Table) FocusLost() {
	t.focused = false
	t.Refresh() //Item(t.currentFocus)
}

func (t *Table) MouseIn(ev *desktop.MouseEvent) {
	t.hoverAt(ev.Position)
}

// MouseDown response to desktop mouse event
func (t *Table) MouseDown(e *desktop.MouseEvent) {
	t.tapped(e.Position)
}

func (t *Table) MouseMoved(ev *desktop.MouseEvent) {
	t.hoverAt(ev.Position)
}

func (t *Table) MouseOut() {
	t.hoverOut()
}

// MouseUp response to desktop mouse event
func (t *Table) MouseUp(*desktop.MouseEvent) {
}

// RefreshItem refreshes a single item, specified by the item ID passed in.
//
// Since: 2.4
func (t *Table) RefreshItem(id TableCellID) {
	if t.cells == nil {
		return
	}
	r := cache.Renderer(t.cells)
	if r == nil {
		return
	}

	r.(*tableCellsRenderer).refreshForID(id)
}

// Select will mark the specified cell as selected.
func (t *Table) Select(id TableCellID) {
	if t.Length == nil {
		return
	}

	rows, cols := t.Length()
	if id.Row >= rows || id.Col >= cols {
		return
	}

	if t.selectedCell != nil && *t.selectedCell == id {
		return
	}
	if f := t.OnUnselected; f != nil && t.selectedCell != nil {
		f(*t.selectedCell)
	}
	t.selectedCell = &id

	t.ScrollTo(id)

	if f := t.OnSelected; f != nil {
		f(id)
	}
}

// SetColumnWidth supports changing the width of the specified column. Columns normally take the width of the template
// cell returned from the CreateCell callback. The width parameter uses the same units as a fyne.Size type and refers
// to the internal content width not including the divider size.
//
// Since: 1.4.1
func (t *Table) SetColumnWidth(id int, width float32) {
	t.propertyLock.Lock()
	if t.columnWidths == nil {
		t.columnWidths = make(map[int]float32)
	}
	t.columnWidths[id] = width
	t.propertyLock.Unlock()

	t.Refresh()
}

// SetRowHeight supports changing the height of the specified row. Rows normally take the height of the template
// cell returned from the CreateCell callback. The height parameter uses the same units as a fyne.Size type and refers
// to the internal content height not including the divider size.
//
// Since: 2.3
func (t *Table) SetRowHeight(id int, height float32) {
	t.propertyLock.Lock()
	if t.rowHeights == nil {
		t.rowHeights = make(map[int]float32)
	}
	t.rowHeights[id] = height
	t.propertyLock.Unlock()

	t.Refresh()
}

// TouchDown response to mobile touch event
func (t *Table) TouchDown(e *mobile.TouchEvent) {
	t.tapped(e.Position)
}

// TouchUp response to mobile touch event
func (t *Table) TouchUp(*mobile.TouchEvent) {
}

// TouchCancel response to mobile touch event
func (t *Table) TouchCancel(*mobile.TouchEvent) {
}

// TypedKey is called if a key event happens while this Table is focused.
//
// Implements: fyne.Focusable
func (t *Table) TypedKey(event *fyne.KeyEvent) {
	switch event.Name {
	case fyne.KeySpace:
		t.Select(t.currentFocus)
	case fyne.KeyDown:
		if f := t.Length; f != nil {
			rows, _ := f()
			if t.currentFocus.Row >= rows-1 {
				return
			}
		}
		t.RefreshItem(t.currentFocus)
		t.currentFocus.Row++
		t.ScrollTo(t.currentFocus)
		t.RefreshItem(t.currentFocus)
	case fyne.KeyLeft:
		if t.currentFocus.Col <= 0 {
			return
		}
		t.RefreshItem(t.currentFocus)
		t.currentFocus.Col--
		t.ScrollTo(t.currentFocus)
		t.RefreshItem(t.currentFocus)
	case fyne.KeyRight:
		if f := t.Length; f != nil {
			_, cols := f()
			if t.currentFocus.Col >= cols-1 {
				return
			}
		}
		t.RefreshItem(t.currentFocus)
		t.currentFocus.Col++
		t.ScrollTo(t.currentFocus)
		t.RefreshItem(t.currentFocus)
	case fyne.KeyUp:
		if t.currentFocus.Row <= 0 {
			return
		}
		t.RefreshItem(t.currentFocus)
		t.currentFocus.Row--
		t.ScrollTo(t.currentFocus)
		t.RefreshItem(t.currentFocus)
	}
}

// TypedRune is called if a text event happens while this Table is focused.
//
// Implements: fyne.Focusable
func (t *Table) TypedRune(_ rune) {
	// intentionally left blank
}

// Unselect will mark the cell provided by id as unselected.
func (t *Table) Unselect(id TableCellID) {
	if t.selectedCell == nil || id != *t.selectedCell {
		return
	}
	t.selectedCell = nil

	if t.moveCallback != nil {
		t.moveCallback()
	}

	if f := t.OnUnselected; f != nil {
		f(id)
	}
}

// UnselectAll will mark all cells as unselected.
//
// Since: 2.1
func (t *Table) UnselectAll() {
	if t.selectedCell == nil {
		return
	}

	selected := *t.selectedCell
	t.selectedCell = nil

	if t.moveCallback != nil {
		t.moveCallback()
	}

	if f := t.OnUnselected; f != nil {
		f(selected)
	}
}

// ScrollTo will scroll to the given cell without changing the selection.
// Attempting to scroll beyond the limits of the table will scroll to
// the edge of the table instead.
//
// Since: 2.1
func (t *Table) ScrollTo(id TableCellID) {
	if t.Length == nil {
		return
	}

	if t.content == nil {
		return
	}

	rows, cols := t.Length()
	if id.Row >= rows {
		id.Row = rows - 1
	}

	if id.Col >= cols {
		id.Col = cols - 1
	}

	scrollPos := t.offset

	cellX, cellWidth := t.findX(id.Col)
	stickCols := t.StickyColumnCount
	if stickCols > 0 {
		cellX -= t.stuckXOff + t.stuckWidth
	}
	if t.ShowHeaderColumn {
		cellX += t.headerSize.Width
		stickCols--
	}
	if stickCols == 0 || id.Col > stickCols {
		if cellX < scrollPos.X {
			scrollPos.X = cellX
		} else if cellX+cellWidth > scrollPos.X+t.content.Size().Width {
			scrollPos.X = cellX + cellWidth - t.content.Size().Width
		}
	}

	cellY, cellHeight := t.findY(id.Row)
	stickRows := t.StickyRowCount
	if stickRows > 0 {
		cellY -= t.stuckYOff + t.stuckHeight
	}
	if t.ShowHeaderRow {
		cellY += t.headerSize.Height
		stickRows--
	}
	if stickRows == 0 || id.Row >= stickRows {
		if cellY < scrollPos.Y {
			scrollPos.Y = cellY
		} else if cellY+cellHeight > scrollPos.Y+t.content.Size().Height {
			scrollPos.Y = cellY + cellHeight - t.content.Size().Height
		}
	}

	t.offset = scrollPos
	t.content.Offset = scrollPos
	t.content.Refresh()
	t.finishScroll()
}

// ScrollToBottom scrolls to the last row in the table
//
// Since: 2.1
func (t *Table) ScrollToBottom() {
	if t.Length == nil || t.content == nil {
		return
	}

	rows, _ := t.Length()
	cellY, cellHeight := t.findY(rows - 1)
	y := cellY + cellHeight - t.content.Size().Height
	if y <= 0 {
		return
	}

	t.content.Offset.Y = y
	t.offset.Y = y
	t.finishScroll()
}

// ScrollToLeading scrolls horizontally to the leading edge of the table
//
// Since: 2.1
func (t *Table) ScrollToLeading() {
	if t.content == nil {
		return
	}

	t.content.Offset.X = 0
	t.offset.X = 0
	t.finishScroll()
}

// ScrollToTop scrolls to the first row in the table
//
// Since: 2.1
func (t *Table) ScrollToTop() {
	if t.content == nil {
		return
	}

	t.content.Offset.Y = 0
	t.offset.Y = 0
	t.finishScroll()
}

// ScrollToTrailing scrolls horizontally to the trailing edge of the table
//
// Since: 2.1
func (t *Table) ScrollToTrailing() {
	if t.content == nil || t.Length == nil {
		return
	}

	_, cols := t.Length()
	cellX, cellWidth := t.findX(cols - 1)
	scrollX := cellX + cellWidth - t.content.Size().Width
	if scrollX <= 0 {
		return
	}

	t.content.Offset.X = scrollX
	t.offset.X = scrollX
	t.finishScroll()
}

func (t *Table) Tapped(e *fyne.PointEvent) {
	if e.Position.X < 0 || e.Position.X >= t.Size().Width || e.Position.Y < 0 || e.Position.Y >= t.Size().Height {
		t.selectedCell = nil
		t.Refresh()
		return
	}

	col := t.columnAt(e.Position)
	if col == noCellMatch {
		return // out of col range
	}
	row := t.rowAt(e.Position)
	if row == noCellMatch {
		return // out of row range
	}
	t.Select(TableCellID{row, col})

	if !fyne.CurrentDevice().IsMobile() {
		t.RefreshItem(t.currentFocus)
		canvas := fyne.CurrentApp().Driver().CanvasForObject(t)
		if canvas != nil {
			canvas.Focus(t)
		}
		t.currentFocus = TableCellID{row, col}
		t.RefreshItem(t.currentFocus)
	}
}

// columnAt returns a positive integer (or 0) for the column that is found at the `pos` X position.
// If the position is between cells the method will return a negative integer representing the next column,
// i.e. -1 means the gap between 0 and 1.
func (t *Table) columnAt(pos fyne.Position) int {
	dataCols := 0
	if f := t.Length; f != nil {
		_, dataCols = t.Length()
	}

	visibleColWidths, offX, minCol, maxCol := t.visibleColumnWidths(t.cellSize.Width, dataCols)
	i := minCol
	end := maxCol
	if pos.X < t.stuckXOff+t.stuckWidth {
		offX = t.stuckXOff
		end = t.StickyColumnCount
		i = 0
	} else {
		pos.X += t.content.Offset.X
		offX += t.stuckXOff
	}
	padding := theme.Padding()
	for x := offX; i < end; x += visibleColWidths[i-1] + padding {
		if pos.X < x {
			return -i // the space between i-1 and i
		} else if pos.X < x+visibleColWidths[i] {
			return i
		}
		i++
	}
	return noCellMatch
}

func (t *Table) createHeader() fyne.CanvasObject {
	if f := t.CreateHeader; f != nil {
		return f()
	}

	l := NewLabel("00")
	l.TextStyle.Bold = true
	l.Alignment = fyne.TextAlignCenter
	return l
}

func (t *Table) findX(col int) (cellX float32, cellWidth float32) {
	cellSize := t.templateSize()
	padding := theme.Padding()
	for i := 0; i <= col; i++ {
		if cellWidth > 0 {
			cellX += cellWidth + padding
		}

		width := cellSize.Width
		if w, ok := t.columnWidths[i]; ok {
			width = w
		}
		cellWidth = width
	}
	return
}

func (t *Table) findY(row int) (cellY float32, cellHeight float32) {
	cellSize := t.templateSize()
	padding := theme.Padding()
	for i := 0; i <= row; i++ {
		if cellHeight > 0 {
			cellY += cellHeight + padding
		}

		height := cellSize.Height
		if h, ok := t.rowHeights[i]; ok {
			height = h
		}
		cellHeight = height
	}
	return
}

func (t *Table) finishScroll() {
	if t.moveCallback != nil {
		t.moveCallback()
	}
	t.cells.Refresh()
}

func (t *Table) hoverAt(pos fyne.Position) {
	col := t.columnAt(pos)
	row := t.rowAt(pos)
	t.hoveredCell = &TableCellID{row, col}
	overHeaderRow := t.ShowHeaderRow && pos.Y < t.headerSize.Height
	overHeaderCol := t.ShowHeaderColumn && pos.X < t.headerSize.Width
	if overHeaderRow && !overHeaderCol {
		if col >= 0 {
			t.hoverHeaderCol = noCellMatch
		} else {
			t.hoverHeaderCol = -col - 1
		}
	} else {
		t.hoverHeaderCol = noCellMatch
	}
	if overHeaderCol && !overHeaderRow {
		if row >= 0 {
			t.hoverHeaderRow = noCellMatch
		} else {
			t.hoverHeaderRow = -row - 1
		}
	} else {
		t.hoverHeaderRow = noCellMatch
	}

	rows, cols := 0, 0
	if f := t.Length; f != nil {
		rows, cols = t.Length()
	}
	if t.hoveredCell.Col >= cols || t.hoveredCell.Row >= rows || t.hoveredCell.Col < 0 || t.hoveredCell.Row < 0 {
		t.hoverOut()
		return
	}

	if t.moveCallback != nil {
		t.moveCallback()
	}
}

func (t *Table) hoverOut() {
	t.hoveredCell = nil

	if t.moveCallback != nil {
		t.moveCallback()
	}
}

// rowAt returns a positive integer (or 0) for the row that is found at the `pos` Y position.
// If the position is between cells the method will return a negative integer representing the next row,
// i.e. -1 means the gap between rows 0 and 1.
func (t *Table) rowAt(pos fyne.Position) int {
	dataRows := 0
	if f := t.Length; f != nil {
		dataRows, _ = t.Length()
	}

	visibleRowHeights, offY, minRow, maxRow := t.visibleRowHeights(t.cellSize.Height, dataRows)
	i := minRow
	end := maxRow
	if pos.Y < t.stuckYOff+t.stuckHeight {
		offY = t.stuckYOff
		end = t.StickyRowCount
		i = 0
	} else {
		pos.Y += t.content.Offset.Y
		offY += t.stuckYOff
	}
	padding := theme.Padding()
	for y := offY; i < end; y += visibleRowHeights[i-1] + padding {
		if pos.Y < y {
			return -i // the space between i-1 and i
		} else if pos.Y >= y && pos.Y < y+visibleRowHeights[i] {
			return i
		}
		i++
	}
	return noCellMatch
}

func (t *Table) tapped(pos fyne.Position) {
	if t.dragCol == noCellMatch && t.dragRow == noCellMatch {
		t.dragStartPos = pos
		if t.hoverHeaderRow != noCellMatch {
			t.dragCol = noCellMatch
			t.dragRow = t.hoverHeaderRow
			size, ok := t.rowHeights[t.hoverHeaderRow]
			if !ok {
				size = t.cellSize.Height
			}
			t.dragStartSize = size
		} else if t.hoverHeaderCol != noCellMatch {
			t.dragCol = t.hoverHeaderCol
			t.dragRow = noCellMatch
			size, ok := t.columnWidths[t.hoverHeaderCol]
			if !ok {
				size = t.cellSize.Width
			}
			t.dragStartSize = size
		}
	}
}

func (t *Table) templateSize() fyne.Size {
	if f := t.CreateCell; f != nil {
		template := f() // don't use cache, we need new template
		if !t.ShowHeaderRow && !t.ShowHeaderColumn {
			return template.MinSize()
		}
		return template.MinSize().Max(t.createHeader().MinSize())
	}

	fyne.LogError("Missing CreateCell callback required for Table", nil)
	return fyne.Size{}
}

func (t *Table) updateHeader(id TableCellID, o fyne.CanvasObject) {
	if f := t.UpdateHeader; f != nil {
		f(id, o)
		return
	}

	l := o.(*Label)
	if id.Row < 0 {
		ids := []rune{'A' + rune(id.Col%26)}
		pre := (id.Col - id.Col%26) / 26
		for pre > 0 {
			ids = append([]rune{'A' - 1 + rune(pre%26)}, ids...)
			pre = (pre - pre%26) / 26
		}
		l.SetText(string(ids))
	} else if id.Col < 0 {
		l.SetText(strconv.Itoa(id.Row + 1))
	} else {
		l.SetText("")
	}
}

func (t *Table) stickyColumnWidths(colWidth float32, cols int) (visible []float32) {
	if cols == 0 {
		return []float32{}
	}

	max := t.StickyColumnCount
	if max > cols {
		max = cols
	}

	visible = make([]float32, max)

	if len(t.columnWidths) == 0 {
		for i := 0; i < max; i++ {
			visible[i] = colWidth
		}
		return
	}

	for i := 0; i < max; i++ {
		height := colWidth

		if h, ok := t.columnWidths[i]; ok {
			height = h
		}

		visible[i] = height
	}
	return
}

func (t *Table) visibleColumnWidths(colWidth float32, cols int) (visible map[int]float32, offX float32, minCol, maxCol int) {
	maxCol = cols
	colOffset, headWidth := float32(0), float32(0)
	isVisible := false
	visible = make(map[int]float32)

	if t.content.Size().Width <= 0 {
		return
	}

	// theme.Padding is a slow call, so we cache it
	padding := theme.Padding()
	stick := t.StickyColumnCount

	if len(t.columnWidths) == 0 {
		paddedWidth := colWidth + padding

		offX = float32(math.Floor(float64(t.offset.X/paddedWidth))) * paddedWidth
		minCol = int(math.Floor(float64(offX / paddedWidth)))
		maxCol = int(math.Ceil(float64((t.offset.X + t.size.Width) / paddedWidth)))

		if minCol > cols-1 {
			minCol = cols - 1
		}
		if minCol < 0 {
			minCol = 0
		}

		if maxCol > cols {
			maxCol = cols
		}

		visible = make(map[int]float32, maxCol-minCol+stick)
		for i := minCol; i < maxCol; i++ {
			visible[i] = colWidth
		}
		for i := 0; i < stick; i++ {
			visible[i] = colWidth
		}
		return
	}

	for i := 0; i < cols; i++ {
		width := colWidth
		if w, ok := t.columnWidths[i]; ok {
			width = w
		}

		if colOffset <= t.offset.X-width-padding {
			// before visible content
		} else if colOffset <= headWidth || colOffset <= t.offset.X {
			minCol = i
			offX = colOffset
			isVisible = true
		}
		if colOffset < t.offset.X+t.size.Width {
			maxCol = i + 1
		} else {
			break
		}

		colOffset += width + padding
		if isVisible || i < stick {
			visible[i] = width
		}
	}
	return
}

func (t *Table) stickyRowHeights(rowHeight float32, rows int) (visible []float32) {
	if rows == 0 {
		return []float32{}
	}

	max := t.StickyRowCount
	if max > rows {
		max = rows
	}

	visible = make([]float32, max)

	if len(t.rowHeights) == 0 {
		for i := 0; i < max; i++ {
			visible[i] = rowHeight
		}
		return
	}

	for i := 0; i < max; i++ {
		height := rowHeight

		if h, ok := t.rowHeights[i]; ok {
			height = h
		}

		visible[i] = height
	}
	return
}

func (t *Table) visibleRowHeights(rowHeight float32, rows int) (visible map[int]float32, offY float32, minRow, maxRow int) {
	maxRow = rows
	rowOffset, headHeight := float32(0), float32(0)
	isVisible := false
	visible = make(map[int]float32)

	if t.content.Size().Height <= 0 {
		return
	}

	// theme.Padding is a slow call, so we cache it
	padding := theme.Padding()
	stick := t.StickyRowCount

	if len(t.rowHeights) == 0 {
		paddedHeight := rowHeight + padding

		offY = float32(math.Floor(float64(t.offset.Y/paddedHeight))) * paddedHeight
		minRow = int(math.Floor(float64(offY / paddedHeight)))
		maxRow = int(math.Ceil(float64((t.offset.Y + t.size.Height) / paddedHeight)))

		if minRow > rows-1 {
			minRow = rows - 1
		}
		if minRow < 0 {
			minRow = 0
		}

		if maxRow > rows {
			maxRow = rows
		}

		visible = make(map[int]float32, maxRow-minRow+stick)
		for i := minRow; i < maxRow; i++ {
			visible[i] = rowHeight
		}
		for i := 0; i < stick; i++ {
			visible[i] = rowHeight
		}
		return
	}

	for i := 0; i < rows; i++ {
		height := rowHeight
		if h, ok := t.rowHeights[i]; ok {
			height = h
		}

		if rowOffset <= t.offset.Y-height-padding {
			// before visible content
		} else if rowOffset <= headHeight || rowOffset <= t.offset.Y {
			minRow = i
			offY = rowOffset
			isVisible = true
		}
		if rowOffset < t.offset.Y+t.size.Height {
			maxRow = i + 1
		} else {
			break
		}

		rowOffset += height + padding
		if isVisible || i < stick {
			visible[i] = height
		}
	}
	return
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*tableRenderer)(nil)

type tableRenderer struct {
	widget.BaseRenderer
	t *Table
}

func (t *tableRenderer) Layout(s fyne.Size) {
	t.t.propertyLock.RLock()

	t.calculateHeaderSizes()
	off := fyne.NewPos(t.t.stuckWidth, t.t.stuckHeight)
	if t.t.ShowHeaderRow {
		off.Y += t.t.headerSize.Height
	}
	if t.t.ShowHeaderColumn {
		off.X += t.t.headerSize.Width
	}
	t.t.propertyLock.RUnlock()

	t.t.content.Move(off)
	t.t.content.Resize(s.SubtractWidthHeight(off.X, off.Y))

	t.t.top.Move(fyne.NewPos(off.X, 0))
	t.t.top.Resize(fyne.NewSize(s.Width-off.X, off.Y))
	t.t.left.Move(fyne.NewPos(0, off.Y))
	t.t.left.Resize(fyne.NewSize(off.X, s.Height-off.Y))
	t.t.corner.Resize(fyne.NewSize(off.X, off.Y))

	t.t.dividerLayer.Resize(s)
}

func (t *tableRenderer) MinSize() fyne.Size {
	sep := theme.Padding()
	t.t.propertyLock.RLock()
	defer t.t.propertyLock.RUnlock()

	min := t.t.content.MinSize().Max(t.t.cellSize)
	if t.t.ShowHeaderRow {
		min.Height += t.t.headerSize.Height + sep
	}
	if t.t.ShowHeaderColumn {
		min.Width += t.t.headerSize.Width + sep
	}
	if t.t.StickyRowCount > 0 {
		for i := 0; i < t.t.StickyRowCount; i++ {
			height := t.t.cellSize.Height
			if h, ok := t.t.rowHeights[i]; ok {
				height = h
			}

			min.Height += height + sep
		}
	}
	if t.t.StickyColumnCount > 0 {
		for i := 0; i < t.t.StickyColumnCount; i++ {
			width := t.t.cellSize.Width
			if w, ok := t.t.columnWidths[i]; ok {
				width = w
			}

			min.Width += width + sep
		}
	}
	return min
}

func (t *tableRenderer) Refresh() {
	t.t.propertyLock.Lock()
	t.t.headerSize = t.t.createHeader().MinSize()
	t.t.cellSize = t.t.templateSize()
	t.calculateHeaderSizes()
	t.t.propertyLock.Unlock()

	t.Layout(t.t.Size())
	t.t.cells.Refresh()
}

func (t *tableRenderer) calculateHeaderSizes() {
	t.t.stuckXOff = 0
	t.t.stuckYOff = 0

	if t.t.ShowHeaderRow {
		t.t.stuckYOff = t.t.headerSize.Height
	}
	if t.t.ShowHeaderColumn {
		t.t.stuckXOff = t.t.headerSize.Width
	}

	separatorThickness := theme.Padding()
	stickyColWidths := t.t.stickyColumnWidths(t.t.cellSize.Width, t.t.StickyColumnCount)
	stickyRowHeights := t.t.stickyRowHeights(t.t.cellSize.Height, t.t.StickyRowCount)

	var stuckHeight float32
	for _, rowHeight := range stickyRowHeights {
		stuckHeight += rowHeight + separatorThickness
	}
	t.t.stuckHeight = stuckHeight
	var stuckWidth float32
	for _, colWidth := range stickyColWidths {
		stuckWidth += colWidth + separatorThickness
	}
	t.t.stuckWidth = stuckWidth
}

// Declare conformity with Widget interface.
var _ fyne.Widget = (*tableCells)(nil)

type tableCells struct {
	BaseWidget
	t *Table
}

func newTableCells(t *Table) *tableCells {
	c := &tableCells{t: t}
	c.ExtendBaseWidget(c)
	return c
}

func (c *tableCells) CreateRenderer() fyne.WidgetRenderer {
	marker := canvas.NewRectangle(theme.SelectionColor())
	marker.CornerRadius = theme.SelectionRadiusSize()
	hover := canvas.NewRectangle(theme.HoverColor())
	hover.CornerRadius = theme.SelectionRadiusSize()

	r := &tableCellsRenderer{cells: c, pool: &syncPool{}, headerPool: &syncPool{},
		visible: make(map[TableCellID]fyne.CanvasObject), headers: make(map[TableCellID]fyne.CanvasObject),
		headRowBG: canvas.NewRectangle(theme.HeaderBackgroundColor()), headColBG: canvas.NewRectangle(theme.HeaderBackgroundColor()),
		headRowStickyBG: canvas.NewRectangle(theme.HeaderBackgroundColor()), headColStickyBG: canvas.NewRectangle(theme.HeaderBackgroundColor()),
		marker: marker, hover: hover}

	c.t.moveCallback = r.moveIndicators
	return r
}

func (c *tableCells) Resize(s fyne.Size) {
	c.BaseWidget.Resize(s)
	c.Refresh() // trigger a redraw
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*tableCellsRenderer)(nil)

type tableCellsRenderer struct {
	widget.BaseRenderer

	cells            *tableCells
	pool, headerPool pool
	visible, headers map[TableCellID]fyne.CanvasObject
	hover, marker    *canvas.Rectangle
	dividers         []fyne.CanvasObject

	headColBG, headRowBG, headRowStickyBG, headColStickyBG *canvas.Rectangle
}

func (r *tableCellsRenderer) Layout(fyne.Size) {
	r.cells.propertyLock.Lock()
	r.moveIndicators()
	r.cells.propertyLock.Unlock()
}

func (r *tableCellsRenderer) MinSize() fyne.Size {
	r.cells.propertyLock.RLock()
	defer r.cells.propertyLock.RUnlock()
	rows, cols := 0, 0
	if f := r.cells.t.Length; f != nil {
		rows, cols = r.cells.t.Length()
	} else {
		fyne.LogError("Missing Length callback required for Table", nil)
	}

	stickRows := r.cells.t.StickyRowCount
	stickCols := r.cells.t.StickyColumnCount

	width := float32(0)
	if len(r.cells.t.columnWidths) == 0 {
		width = r.cells.t.cellSize.Width * float32(cols-stickCols)
	} else {
		cellWidth := r.cells.t.cellSize.Width
		for col := stickCols; col < cols; col++ {
			colWidth, ok := r.cells.t.columnWidths[col]
			if ok {
				width += colWidth
			} else {
				width += cellWidth
			}
		}
	}

	height := float32(0)
	if len(r.cells.t.rowHeights) == 0 {
		height = r.cells.t.cellSize.Height * float32(rows-stickRows)
	} else {
		cellHeight := r.cells.t.cellSize.Height
		for row := stickRows; row < rows; row++ {
			rowHeight, ok := r.cells.t.rowHeights[row]
			if ok {
				height += rowHeight
			} else {
				height += cellHeight
			}
		}
	}

	separatorSize := theme.Padding()
	return fyne.NewSize(width+float32(cols-stickCols-1)*separatorSize, height+float32(rows-stickRows-1)*separatorSize)
}

func (r *tableCellsRenderer) Refresh() {
	r.refreshForID(allTableCellsID)
}

func (r *tableCellsRenderer) refreshForID(toDraw TableCellID) {
	r.cells.propertyLock.Lock()
	separatorThickness := theme.Padding()
	dataRows, dataCols := 0, 0
	if f := r.cells.t.Length; f != nil {
		dataRows, dataCols = r.cells.t.Length()
	}
	visibleColWidths, offX, minCol, maxCol := r.cells.t.visibleColumnWidths(r.cells.t.cellSize.Width, dataCols)
	if len(visibleColWidths) == 0 && dataCols > 0 { // we can't show anything until we have some dimensions
		r.cells.propertyLock.Unlock()
		return
	}
	visibleRowHeights, offY, minRow, maxRow := r.cells.t.visibleRowHeights(r.cells.t.cellSize.Height, dataRows)
	if len(visibleRowHeights) == 0 && dataRows > 0 { // we can't show anything until we have some dimensions
		r.cells.propertyLock.Unlock()
		return
	}

	updateCell := r.cells.t.UpdateCell
	if updateCell == nil {
		fyne.LogError("Missing UpdateCell callback required for Table", nil)
	}

	var cellXOffset, cellYOffset float32
	stickRows := r.cells.t.StickyRowCount
	if r.cells.t.ShowHeaderRow {
		cellYOffset += r.cells.t.headerSize.Height
	}
	stickCols := r.cells.t.StickyColumnCount
	if r.cells.t.ShowHeaderColumn {
		cellXOffset += r.cells.t.headerSize.Width
	}
	startRow := minRow + stickRows
	if startRow < stickRows {
		startRow = stickRows
	}
	startCol := minCol + stickCols
	if startCol < stickCols {
		startCol = stickCols
	}

	wasVisible := r.visible
	r.visible = make(map[TableCellID]fyne.CanvasObject)
	var cells []fyne.CanvasObject
	displayCol := func(row, col int, rowHeight float32, cells *[]fyne.CanvasObject) {
		id := TableCellID{row, col}
		colWidth := visibleColWidths[col]
		c, ok := wasVisible[id]
		if !ok {
			c = r.pool.Obtain()
			if f := r.cells.t.CreateCell; f != nil && c == nil {
				c = f()
			}
			if c == nil {
				return
			}
		}

		c.Move(fyne.NewPos(cellXOffset, cellYOffset))
		c.Resize(fyne.NewSize(colWidth, rowHeight))

		r.visible[id] = c
		*cells = append(*cells, c)
		cellXOffset += colWidth + separatorThickness
	}

	displayRow := func(row int, cells *[]fyne.CanvasObject) {
		rowHeight := visibleRowHeights[row]
		cellXOffset = offX

		for col := startCol; col < maxCol; col++ {
			displayCol(row, col, rowHeight, cells)
		}
		cellXOffset = r.cells.t.content.Offset.X
		stick := r.cells.t.StickyColumnCount
		if r.cells.t.ShowHeaderColumn {
			cellXOffset += r.cells.t.headerSize.Width
			stick--
		}
		cellYOffset += rowHeight + separatorThickness
	}

	cellYOffset = offY
	for row := startRow; row < maxRow; row++ {
		displayRow(row, &cells)
	}

	inline := r.refreshHeaders(visibleRowHeights, visibleColWidths, offX, offY, startRow, maxRow, startCol, maxCol, separatorThickness)
	cells = append(cells, inline...)

	offX -= r.cells.t.content.Offset.X
	cellYOffset = r.cells.t.stuckYOff
	for row := 0; row < stickRows; row++ {
		displayRow(row, &r.cells.t.top.Content.(*fyne.Container).Objects)
	}

	cellYOffset = offY - r.cells.t.content.Offset.Y
	for row := startRow; row < maxRow; row++ {
		cellXOffset = r.cells.t.stuckXOff
		rowHeight := visibleRowHeights[row]
		for col := 0; col < stickCols; col++ {
			displayCol(row, col, rowHeight, &r.cells.t.left.Content.(*fyne.Container).Objects)
		}
		cellYOffset += rowHeight + separatorThickness
	}

	cellYOffset = r.cells.t.stuckYOff
	for row := 0; row < stickRows; row++ {
		cellXOffset = r.cells.t.stuckXOff
		rowHeight := visibleRowHeights[row]
		for col := 0; col < stickCols; col++ {
			displayCol(row, col, rowHeight, &r.cells.t.corner.Content.(*fyne.Container).Objects)
		}
		cellYOffset += rowHeight + separatorThickness
	}

	for id, old := range wasVisible {
		if _, ok := r.visible[id]; !ok {
			r.pool.Release(old)
		}
	}
	visible := r.visible
	headers := r.headers

	r.cells.propertyLock.Unlock()
	r.SetObjects(cells)

	if updateCell != nil {
		for id, cell := range visible {
			if toDraw != allTableCellsID && toDraw != id {
				continue
			}

			updateCell(id, cell)
		}
	}
	for id, head := range headers {
		r.cells.t.updateHeader(id, head)
	}

	r.moveIndicators()
	r.marker.FillColor = theme.SelectionColor()
	r.marker.CornerRadius = theme.SelectionRadiusSize()
	r.marker.Refresh()
	r.hover.FillColor = theme.HoverColor()
	r.hover.CornerRadius = theme.SelectionRadiusSize()
	r.hover.Refresh()
}

func (r *tableCellsRenderer) moveIndicators() {
	rows, cols := 0, 0
	if f := r.cells.t.Length; f != nil {
		rows, cols = r.cells.t.Length()
	}
	visibleColWidths, offX, minCol, maxCol := r.cells.t.visibleColumnWidths(r.cells.t.cellSize.Width, cols)
	visibleRowHeights, offY, minRow, maxRow := r.cells.t.visibleRowHeights(r.cells.t.cellSize.Height, rows)
	separatorThickness := theme.SeparatorThicknessSize()
	padding := theme.Padding()
	dividerOff := (padding - separatorThickness) / 2

	stickRows := r.cells.t.StickyRowCount
	stickCols := r.cells.t.StickyColumnCount

	if r.cells.t.ShowHeaderColumn {
		offX += r.cells.t.headerSize.Width
	}
	if r.cells.t.ShowHeaderRow {
		offY += r.cells.t.headerSize.Height
	}
	if r.cells.t.selectedCell == nil {
		r.moveMarker(r.marker, -1, -1, offX, offY, minCol, minRow, visibleColWidths, visibleRowHeights)
	} else {
		r.moveMarker(r.marker, r.cells.t.selectedCell.Row, r.cells.t.selectedCell.Col, offX, offY, minCol, minRow, visibleColWidths, visibleRowHeights)
	}
	if r.cells.t.hoveredCell == nil && !r.cells.t.focused {
		r.moveMarker(r.hover, -1, -1, offX, offY, minCol, minRow, visibleColWidths, visibleRowHeights)
	} else if r.cells.t.focused {
		r.moveMarker(r.hover, r.cells.t.currentFocus.Row, r.cells.t.currentFocus.Col, offX, offY, minCol, minRow, visibleColWidths, visibleRowHeights)
	} else {
		r.moveMarker(r.hover, r.cells.t.hoveredCell.Row, r.cells.t.hoveredCell.Col, offX, offY, minCol, minRow, visibleColWidths, visibleRowHeights)
	}

	colDivs := stickCols + maxCol - minCol - 1
	if colDivs < 0 {
		colDivs = 0
	}
	rowDivs := stickRows + maxRow - minRow - 1
	if rowDivs < 0 {
		rowDivs = 0
	}

	if colDivs < 0 {
		colDivs = 0
	}
	if rowDivs < 0 {
		rowDivs = 0
	}

	if len(r.dividers) < colDivs+rowDivs {
		for i := len(r.dividers); i < colDivs+rowDivs; i++ {
			r.dividers = append(r.dividers, NewSeparator())
		}

		objs := []fyne.CanvasObject{r.marker, r.hover}
		r.cells.t.dividerLayer.Content.(*fyne.Container).Objects = append(objs, r.dividers...)
		r.cells.t.dividerLayer.Content.Refresh()
	}

	divs := 0
	i := 0
	if stickCols > 0 {
		for x := r.cells.t.stuckXOff + visibleColWidths[i]; i < stickCols && divs < colDivs; x += visibleColWidths[i] + padding {
			i++

			xPos := x + dividerOff
			r.dividers[divs].Resize(fyne.NewSize(separatorThickness, r.cells.t.size.Height))
			r.dividers[divs].Move(fyne.NewPos(xPos, 0))
			r.dividers[divs].Show()
			divs++
		}
	}
	i = minCol + stickCols
	for x := offX + r.cells.t.stuckWidth + visibleColWidths[i]; i < maxCol-1 && divs < colDivs; x += visibleColWidths[i] + padding {
		i++

		xPos := x - r.cells.t.content.Offset.X + dividerOff
		r.dividers[divs].Resize(fyne.NewSize(separatorThickness, r.cells.t.size.Height))
		r.dividers[divs].Move(fyne.NewPos(xPos, 0))
		r.dividers[divs].Show()
		divs++
	}

	i = 0
	if stickRows > 0 {
		for y := r.cells.t.stuckYOff + visibleRowHeights[i]; i < stickRows && divs-colDivs < rowDivs; y += visibleRowHeights[i] + padding {
			i++

			yPos := y + dividerOff
			r.dividers[divs].Resize(fyne.NewSize(r.cells.t.size.Width, separatorThickness))
			r.dividers[divs].Move(fyne.NewPos(0, yPos))
			r.dividers[divs].Show()
			divs++
		}
	}
	i = minRow + stickRows
	for y := offY + r.cells.t.stuckHeight + visibleRowHeights[i]; i < maxRow-1 && divs-colDivs < rowDivs; y += visibleRowHeights[i] + padding {
		i++

		yPos := y - r.cells.t.content.Offset.Y + dividerOff
		r.dividers[divs].Resize(fyne.NewSize(r.cells.t.size.Width, separatorThickness))
		r.dividers[divs].Move(fyne.NewPos(0, yPos))
		r.dividers[divs].Show()
		divs++
	}

	for i := divs; i < len(r.dividers); i++ {
		r.dividers[i].Hide()
	}
}

func (r *tableCellsRenderer) moveMarker(marker fyne.CanvasObject, row, col int, offX, offY float32, minCol, minRow int, widths, heights map[int]float32) {
	if col == -1 || row == -1 {
		marker.Hide()
		marker.Refresh()
		return
	}

	xPos := offX
	stickCols := r.cells.t.StickyColumnCount
	if col < stickCols {
		if r.cells.t.ShowHeaderColumn {
			xPos = r.cells.t.stuckXOff
		} else {
			xPos = 0
		}
		minCol = 0
	}

	padding := theme.Padding()

	for i := minCol; i < col; i++ {
		xPos += widths[i]
		xPos += padding
	}
	x1 := xPos
	if col >= stickCols {
		x1 -= r.cells.t.content.Offset.X
	}
	x2 := x1 + widths[col]

	yPos := offY
	stickRows := r.cells.t.StickyRowCount
	if row < stickRows {
		if r.cells.t.ShowHeaderRow {
			yPos = r.cells.t.stuckYOff
		} else {
			yPos = 0
		}
		minRow = 0
	}
	for i := minRow; i < row; i++ {
		yPos += heights[i]
		yPos += padding
	}
	y1 := yPos
	if row >= stickRows {
		y1 -= r.cells.t.content.Offset.Y
	}
	y2 := y1 + heights[row]

	if x2 < 0 || x1 > r.cells.t.size.Width || y2 < 0 || y1 > r.cells.t.size.Height {
		marker.Hide()
	} else {
		left := x1
		if col >= stickCols { // clip X
			left = fyne.Max(r.cells.t.stuckXOff+r.cells.t.stuckWidth, x1)
		}
		top := y1
		if row >= stickRows { // clip Y
			top = fyne.Max(r.cells.t.stuckYOff+r.cells.t.stuckHeight, y1)
		}
		marker.Move(fyne.NewPos(left, top))
		marker.Resize(fyne.NewSize(x2-left, y2-top))

		marker.Show()
	}
	marker.Refresh()
}

func (r *tableCellsRenderer) refreshHeaders(visibleRowHeights, visibleColWidths map[int]float32, offX, offY float32,
	startRow, maxRow, startCol, maxCol int, separatorThickness float32) []fyne.CanvasObject {
	wasVisible := r.headers
	r.headers = make(map[TableCellID]fyne.CanvasObject)
	headerMin := r.cells.t.createHeader().MinSize()
	rowHeight := headerMin.Height
	colWidth := headerMin.Width

	var cells, over []fyne.CanvasObject
	corner := []fyne.CanvasObject{r.headColStickyBG, r.headRowStickyBG}
	over = []fyne.CanvasObject{r.headRowBG}
	if r.cells.t.ShowHeaderRow {
		cellXOffset := offX - r.cells.t.content.Offset.X
		displayColHeader := func(col int, list *[]fyne.CanvasObject) {
			id := TableCellID{-1, col}
			colWidth := visibleColWidths[col]
			c, ok := wasVisible[id]
			if !ok {
				c = r.headerPool.Obtain()
				if c == nil {
					c = r.cells.t.createHeader()
				}
				if c == nil {
					return
				}
			}

			c.Move(fyne.NewPos(cellXOffset, 0))
			c.Resize(fyne.NewSize(colWidth, rowHeight))

			r.headers[id] = c
			*list = append(*list, c)
			cellXOffset += colWidth + separatorThickness
		}
		for col := startCol; col < maxCol; col++ {
			displayColHeader(col, &over)
		}

		if r.cells.t.StickyColumnCount > 0 {
			cellXOffset = 0
			if r.cells.t.ShowHeaderColumn {
				cellXOffset += r.cells.t.headerSize.Width
			}

			for col := 0; col < r.cells.t.StickyColumnCount; col++ {
				displayColHeader(col, &corner)
			}
		}
	}
	r.cells.t.top.Content.(*fyne.Container).Objects = over
	r.cells.t.top.Content.Refresh()

	over = []fyne.CanvasObject{r.headColBG}
	if r.cells.t.ShowHeaderColumn {
		cellYOffset := offY - r.cells.t.content.Offset.Y
		displayRowHeader := func(row int, list *[]fyne.CanvasObject) {
			id := TableCellID{row, -1}
			rowHeight := visibleRowHeights[row]
			c, ok := wasVisible[id]
			if !ok {
				c = r.headerPool.Obtain()
				if c == nil {
					c = r.cells.t.createHeader()
				}
				if c == nil {
					return
				}
			}

			c.Move(fyne.NewPos(0, cellYOffset))
			c.Resize(fyne.NewSize(colWidth, rowHeight))

			r.headers[id] = c
			*list = append(*list, c)
			cellYOffset += rowHeight + separatorThickness
		}
		for row := startRow; row < maxRow; row++ {
			displayRowHeader(row, &over)
		}

		if r.cells.t.StickyRowCount > 0 {
			cellYOffset = 0
			if r.cells.t.ShowHeaderRow {
				cellYOffset += r.cells.t.headerSize.Height
			}

			for row := 0; row < r.cells.t.StickyRowCount; row++ {
				displayRowHeader(row, &corner)
			}
		}
	}
	r.cells.t.left.Content.(*fyne.Container).Objects = over
	r.cells.t.left.Content.Refresh()

	r.headColBG.Hidden = !r.cells.t.ShowHeaderColumn
	r.headColBG.FillColor = theme.HeaderBackgroundColor()
	r.headColBG.Resize(fyne.NewSize(colWidth, r.cells.t.Size().Height))

	r.headColStickyBG.Hidden = !r.cells.t.ShowHeaderColumn
	r.headColStickyBG.FillColor = theme.HeaderBackgroundColor()
	r.headColStickyBG.Resize(fyne.NewSize(colWidth, r.cells.t.stuckHeight+rowHeight))
	r.headRowBG.Hidden = !r.cells.t.ShowHeaderRow
	r.headRowBG.FillColor = theme.HeaderBackgroundColor()
	r.headRowBG.Resize(fyne.NewSize(r.cells.t.Size().Width, rowHeight))
	r.headRowStickyBG.Hidden = !r.cells.t.ShowHeaderRow
	r.headRowStickyBG.FillColor = theme.HeaderBackgroundColor()
	r.headRowStickyBG.Resize(fyne.NewSize(r.cells.t.stuckWidth+colWidth, rowHeight))
	r.cells.t.corner.Content.(*fyne.Container).Objects = corner
	r.cells.t.corner.Content.Refresh()

	for id, old := range wasVisible {
		if _, ok := r.headers[id]; !ok {
			r.headerPool.Release(old)
		}
	}
	return cells
}

type clip struct {
	widget.Scroll

	t *Table
}

func newClip(t *Table, o fyne.CanvasObject) *clip {
	c := &clip{t: t}
	c.Content = o
	c.Direction = widget.ScrollNone

	return c
}

func (c *clip) DragEnd() {
	c.t.DragEnd()
	c.t.dragCol = noCellMatch
	c.t.dragRow = noCellMatch
}

func (c *clip) Dragged(e *fyne.DragEvent) {
	c.t.Dragged(e)
}

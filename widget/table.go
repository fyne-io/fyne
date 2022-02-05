package widget

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

// Declare conformity with Widget interface.
var _ fyne.Widget = (*Table)(nil)

// TableCellID is a type that represents a cell's position in a table based on it's row and column location.
type TableCellID struct {
	Row int
	Col int
}

// Table widget is a grid of items that can be scrolled and a cell selected.
// It's performance is provided by caching cell templates created with CreateCell and re-using them with UpdateCell.
// The size of the content rows/columns is returned by the Length callback.
//
// Since: 1.4
type Table struct {
	BaseWidget

	Length       func() (int, int)                                `json:"-"`
	CreateCell   func() fyne.CanvasObject                         `json:"-"`
	UpdateCell   func(id TableCellID, template fyne.CanvasObject) `json:"-"`
	OnSelected   func(id TableCellID)                             `json:"-"`
	OnUnselected func(id TableCellID)                             `json:"-"`

	selectedCell, hoveredCell *TableCellID
	cells                     *tableCells
	columnWidths              map[int]float32
	moveCallback              func()
	offset                    fyne.Position
	scroll                    *widget.Scroll
}

// NewTable returns a new performant table widget defined by the passed functions.
// The first returns the data size in rows and columns, second parameter is a function that returns cell
// template objects that can be cached and the third is used to apply data at specified data location to the
// passed template CanvasObject.
//
// Since: 1.4
func NewTable(length func() (int, int), create func() fyne.CanvasObject, update func(TableCellID, fyne.CanvasObject)) *Table {
	t := &Table{Length: length, CreateCell: create, UpdateCell: update}
	t.ExtendBaseWidget(t)
	return t
}

// CreateRenderer returns a new renderer for the table.
//
// Implements: fyne.Widget
func (t *Table) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	marker := canvas.NewRectangle(theme.SelectionColor())
	hover := canvas.NewRectangle(theme.HoverColor())

	cellSize := t.templateSize()
	t.cells = newTableCells(t, cellSize)
	t.scroll = widget.NewScroll(t.cells)

	obj := []fyne.CanvasObject{marker, hover, t.scroll}
	r := &tableRenderer{t: t, scroll: t.scroll, marker: marker, hover: hover, cellSize: cellSize}
	r.SetObjects(obj)
	t.moveCallback = r.moveIndicators
	t.scroll.OnScrolled = func(pos fyne.Position) {
		t.offset = pos
		t.cells.Refresh()
		r.moveIndicators()
	}

	r.Layout(t.Size())
	return r
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
	if t.columnWidths == nil {
		t.columnWidths = make(map[int]float32)
	}
	t.columnWidths[id] = width
	t.Refresh()
	t.scroll.Refresh()
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

	if t.scroll == nil {
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
	if cellX < scrollPos.X {
		scrollPos.X = cellX
	} else if cellX+cellWidth > scrollPos.X+t.scroll.Size().Width {
		scrollPos.X = cellX + cellWidth - t.scroll.Size().Width
	}

	cellY, cellHeight := t.findY(id.Row)
	if cellY < scrollPos.Y {
		scrollPos.Y = cellY
	} else if cellY+cellHeight > scrollPos.Y+t.scroll.Size().Height {
		scrollPos.Y = cellY + cellHeight - t.scroll.Size().Height
	}

	t.scroll.Offset = scrollPos
	t.offset = scrollPos
	t.finishScroll()
}

// ScrollToBottom scrolls to the last row in the table
//
// Since: 2.1
func (t *Table) ScrollToBottom() {
	if t.Length == nil || t.scroll == nil {
		return
	}

	rows, _ := t.Length()
	cellY, cellHeight := t.findY(rows - 1)
	y := cellY + cellHeight - t.scroll.Size().Height

	t.scroll.Offset.Y = y
	t.offset.Y = y
	t.finishScroll()
}

// ScrollToLeading scrolls horizontally to the leading edge of the table
//
// Since: 2.1
func (t *Table) ScrollToLeading() {
	if t.scroll == nil {
		return
	}

	t.scroll.Offset.X = 0
	t.offset.X = 0
	t.finishScroll()
}

// ScrollToTop scrolls to the first row in the table
//
// Since: 2.1
func (t *Table) ScrollToTop() {
	if t.scroll == nil {
		return
	}

	t.scroll.Offset.Y = 0
	t.offset.Y = 0
	t.finishScroll()
}

// ScrollToTrailing scrolls horizontally to the trailing edge of the table
//
// Since: 2.1
func (t *Table) ScrollToTrailing() {
	if t.scroll == nil || t.Length == nil {
		return
	}

	_, cols := t.Length()
	cellX, cellWidth := t.findX(cols - 1)
	scrollX := cellX + cellWidth - t.scroll.Size().Width

	t.scroll.Offset.X = scrollX
	t.offset.X = scrollX
	t.finishScroll()
}

func (t *Table) findX(col int) (cellX float32, cellWidth float32) {
	cellSize := t.templateSize()
	for i := 0; i <= col; i++ {
		if cellWidth > 0 {
			cellX += cellWidth + theme.SeparatorThicknessSize()
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
	cellHeight = cellSize.Height
	cellY = float32(row) * (cellHeight + theme.SeparatorThicknessSize())
	return
}

func (t *Table) finishScroll() {
	if t.moveCallback != nil {
		t.moveCallback()
	}
	t.scroll.Refresh()
	t.cells.Refresh()
}

func (t *Table) templateSize() fyne.Size {
	if f := t.CreateCell; f != nil {
		template := f() // don't use cache, we need new template
		return template.MinSize()
	}

	fyne.LogError("Missing CreateCell callback required for Table", nil)
	return fyne.Size{}
}

func (t *Table) visibleColumnWidths(colWidth float32, cols int) (visible map[int]float32, offX float32, minCol, maxCol int) {
	maxCol = cols
	colOffset := float32(0)
	isVisible := false
	visible = make(map[int]float32)

	if t.scroll.Size().Width <= 0 {
		return
	}

	for i := 0; i < cols; i++ {
		width := colWidth
		if w, ok := t.columnWidths[i]; ok {
			width = w
		}

		if colOffset <= t.offset.X-width-theme.SeparatorThicknessSize() {
			// before scroll
		} else if colOffset <= t.offset.X {
			minCol = i
			offX = colOffset
			isVisible = true
		}
		if colOffset < t.offset.X+t.scroll.Size().Width {
			maxCol = i + 1
		} else {
			break
		}

		colOffset += width + theme.SeparatorThicknessSize()
		if isVisible {
			visible[i] = width
		}
	}
	return
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*tableRenderer)(nil)

type tableRenderer struct {
	widget.BaseRenderer
	t *Table

	scroll        *widget.Scroll
	hover, marker *canvas.Rectangle
	dividers      []fyne.CanvasObject

	cellSize fyne.Size
}

func (t *tableRenderer) Layout(s fyne.Size) {
	t.scroll.Resize(s)
	t.moveIndicators()
}

func (t *tableRenderer) MinSize() fyne.Size {
	return t.t.scroll.MinSize().Max(t.cellSize)
}

func (t *tableRenderer) Refresh() {
	t.cellSize = t.t.templateSize()
	t.moveIndicators()

	t.marker.FillColor = theme.SelectionColor()
	t.marker.Refresh()

	t.hover.FillColor = theme.HoverColor()
	t.hover.Refresh()

	t.t.cells.Refresh()
}

func (t *tableRenderer) moveIndicators() {
	rows, cols := 0, 0
	if f := t.t.Length; f != nil {
		rows, cols = t.t.Length()
	}
	visibleColWidths, offX, minCol, maxCol := t.t.visibleColumnWidths(t.cellSize.Width, cols)
	separatorThickness := theme.SeparatorThicknessSize()

	if t.t.selectedCell == nil {
		t.moveMarker(t.marker, -1, -1, offX, minCol, visibleColWidths)
	} else {
		t.moveMarker(t.marker, t.t.selectedCell.Row, t.t.selectedCell.Col, offX, minCol, visibleColWidths)
	}
	if t.t.hoveredCell == nil {
		t.moveMarker(t.hover, -1, -1, offX, minCol, visibleColWidths)
	} else {
		t.moveMarker(t.hover, t.t.hoveredCell.Row, t.t.hoveredCell.Col, offX, minCol, visibleColWidths)
	}

	colDivs := maxCol - minCol - 1
	rowDivs := int(math.Ceil(float64(t.t.size.Height+separatorThickness) / float64(t.cellSize.Height+1)))

	if len(t.dividers) < colDivs+rowDivs {
		for i := len(t.dividers); i < colDivs+rowDivs; i++ {
			t.dividers = append(t.dividers, NewSeparator())
		}

		obj := []fyne.CanvasObject{t.marker, t.hover, t.scroll}
		t.SetObjects(append(obj, t.dividers...))
	}

	divs := 0
	i := minCol
	for x := offX + visibleColWidths[i]; i < minCol+colDivs && divs < len(t.dividers); x += visibleColWidths[i] + separatorThickness {
		i++

		t.dividers[divs].Move(fyne.NewPos(x-t.scroll.Offset.X, 0))
		t.dividers[divs].Resize(fyne.NewSize(separatorThickness, t.t.size.Height))
		t.dividers[divs].Show()
		divs++
	}

	i = 0
	count := float32(math.Mod(float64(t.scroll.Offset.Y), float64(t.cellSize.Height+separatorThickness)))
	for y := t.scroll.Offset.Y - count - separatorThickness; y < t.scroll.Offset.Y+t.t.size.Height && i < rows-1 && divs < len(t.dividers); y += t.cellSize.Height + separatorThickness {
		if y < t.scroll.Offset.Y {
			continue
		}
		i++

		t.dividers[divs].Move(fyne.NewPos(0, y-t.scroll.Offset.Y))
		t.dividers[divs].Resize(fyne.NewSize(t.t.size.Width, separatorThickness))
		t.dividers[divs].Show()
		divs++
	}

	for i := divs; i < len(t.dividers); i++ {
		t.dividers[i].Hide()
	}
	canvas.Refresh(t.t)
}

func (t *tableRenderer) moveMarker(marker fyne.CanvasObject, row, col int, offX float32, minCol int, widths map[int]float32) {
	if col == -1 || row == -1 {
		marker.Hide()
		marker.Refresh()
		return
	}

	xPos := offX
	for i := minCol; i < col; i++ {
		if width, ok := widths[i]; ok {
			xPos += width
		} else {
			xPos += t.cellSize.Width
		}
		xPos += theme.SeparatorThicknessSize()
	}
	x1 := xPos - t.scroll.Offset.X
	x2 := x1 + widths[col]

	offY := float32(row)*(t.cellSize.Height+theme.SeparatorThicknessSize()) - t.scroll.Offset.Y
	y1 := offY
	y2 := y1 + t.cellSize.Height

	if x2 < 0 || x1 > t.t.size.Width || y2 < 0 || y1 > t.t.size.Height {
		marker.Hide()
	} else {
		left := fyne.Max(0, x1)
		top := fyne.Max(0, y1)
		marker.Move(fyne.NewPos(left, top))
		marker.Resize(fyne.NewSize(fyne.Min(x2, t.t.size.Width)-left, fyne.Min(y2, t.t.size.Height)-top))

		marker.Show()
	}
	marker.Refresh()
}

// Declare conformity with Hoverable interface.
var _ desktop.Hoverable = (*tableCells)(nil)

// Declare conformity with Tappable interface.
var _ fyne.Tappable = (*tableCells)(nil)

// Declare conformity with Widget interface.
var _ fyne.Widget = (*tableCells)(nil)

type tableCells struct {
	BaseWidget
	t        *Table
	cellSize fyne.Size
}

func newTableCells(t *Table, s fyne.Size) *tableCells {
	c := &tableCells{t: t, cellSize: s}
	c.ExtendBaseWidget(c)
	return c
}

func (c *tableCells) CreateRenderer() fyne.WidgetRenderer {
	return &tableCellsRenderer{cells: c, pool: &syncPool{}, visible: make(map[TableCellID]fyne.CanvasObject)}
}

func (c *tableCells) MouseIn(ev *desktop.MouseEvent) {
	c.hoverAt(ev.Position)
}

func (c *tableCells) MouseMoved(ev *desktop.MouseEvent) {
	c.hoverAt(ev.Position)
}

func (c *tableCells) MouseOut() {
	c.hoverOut()
}

func (c *tableCells) Resize(s fyne.Size) {
	if s == c.size {
		return
	}
	c.BaseWidget.Resize(s)
	c.Refresh() // trigger a redraw
}

func (c *tableCells) Tapped(e *fyne.PointEvent) {
	if e.Position.X < 0 || e.Position.X >= c.Size().Width || e.Position.Y < 0 || e.Position.Y >= c.Size().Height {
		c.t.selectedCell = nil
		c.t.Refresh()
		return
	}

	col := c.columnAt(e.Position)
	if col == -1 {
		return // out of col range
	}
	row := int(e.Position.Y / (c.cellSize.Height + theme.SeparatorThicknessSize()))
	c.t.Select(TableCellID{row, col})
}

func (c *tableCells) columnAt(pos fyne.Position) int {
	dataCols := 0
	if f := c.t.Length; f != nil {
		_, dataCols = c.t.Length()
	}

	col := -1
	visibleColWidths, offX, minCol, _ := c.t.visibleColumnWidths(c.cellSize.Width, dataCols)
	i := minCol
	for x := offX; i < minCol+len(visibleColWidths); x += visibleColWidths[i-1] + theme.SeparatorThicknessSize() {
		if pos.X >= x && pos.X < x+visibleColWidths[i] {
			col = i
		}
		i++
	}
	return col
}

func (c *tableCells) hoverAt(pos fyne.Position) {
	if pos.X < 0 || pos.X >= c.Size().Width || pos.Y < 0 || pos.Y >= c.Size().Height {
		c.hoverOut()
		return
	}

	col := c.columnAt(pos)
	row := int(pos.Y / (c.cellSize.Height + theme.SeparatorThicknessSize()))
	c.t.hoveredCell = &TableCellID{row, col}

	rows, cols := 0, 0
	if f := c.t.Length; f != nil {
		rows, cols = c.t.Length()
	}
	if c.t.hoveredCell.Col >= cols || c.t.hoveredCell.Row >= rows {
		c.hoverOut()
		return
	}

	if c.t.moveCallback != nil {
		c.t.moveCallback()
	}
}

func (c *tableCells) hoverOut() {
	c.t.hoveredCell = nil

	if c.t.moveCallback != nil {
		c.t.moveCallback()
	}
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*tableCellsRenderer)(nil)

type tableCellsRenderer struct {
	widget.BaseRenderer

	cells   *tableCells
	pool    pool
	visible map[TableCellID]fyne.CanvasObject
}

func (r *tableCellsRenderer) Layout(_ fyne.Size) {
	// we deal with cached objects so just refresh instead
}

func (r *tableCellsRenderer) MinSize() fyne.Size {
	rows, cols := 0, 0
	if f := r.cells.t.Length; f != nil {
		rows, cols = r.cells.t.Length()
	} else {
		fyne.LogError("Missing Length callback required for Table", nil)
	}

	width := float32(0)
	cellWidth := r.cells.cellSize.Width
	for col := 0; col < cols; col++ {
		colWidth, ok := r.cells.t.columnWidths[col]
		if ok {
			width += colWidth
		} else {
			width += cellWidth
		}
	}

	separatorSize := theme.SeparatorThicknessSize()
	return fyne.NewSize(width+float32(cols-1)*separatorSize, r.cells.cellSize.Height*float32(rows)+float32(rows-1)*separatorSize)
}

func (r *tableCellsRenderer) Refresh() {
	r.cells.propertyLock.Lock()
	defer r.cells.propertyLock.Unlock()
	oldSize := r.cells.cellSize
	r.cells.cellSize = r.cells.t.templateSize()
	if oldSize != r.cells.cellSize { // theme changed probably
		r.returnAllToPool()
	}

	separatorThickness := theme.SeparatorThicknessSize()
	dataRows, dataCols := 0, 0
	if f := r.cells.t.Length; f != nil {
		dataRows, dataCols = r.cells.t.Length()
	}
	rows := r.visibleRows()
	visibleColWidths, offX, minCol, maxCol := r.cells.t.visibleColumnWidths(r.cells.cellSize.Width, dataCols)
	if len(visibleColWidths) == 0 { // we can't show anything until we have some dimensions
		return
	}
	offY := r.cells.t.offset.Y - float32(int(r.cells.t.offset.Y)%int(r.cells.cellSize.Height+separatorThickness))
	minRow := int(offY / (r.cells.cellSize.Height + separatorThickness))
	maxRow := int(fyne.Min(float32(minRow+rows), float32(dataRows)))

	updateCell := r.cells.t.UpdateCell
	if updateCell == nil {
		fyne.LogError("Missing UpdateCell callback required for Table", nil)
	}

	wasVisible := r.visible
	r.visible = make(map[TableCellID]fyne.CanvasObject)
	var cells []fyne.CanvasObject
	for row := minRow; row < maxRow; row++ {
		cellOffset := offX
		for col := minCol; col < maxCol; col++ {
			id := TableCellID{row, col}
			colWidth := visibleColWidths[col]
			c, ok := wasVisible[id]
			if !ok {
				c = r.pool.Obtain()
				if f := r.cells.t.CreateCell; f != nil && c == nil {
					c = f()
				}
				if c == nil {
					continue
				}
			}

			c.Move(fyne.NewPos(cellOffset,
				float32(row)*(r.cells.cellSize.Height+separatorThickness)))
			c.Resize(fyne.NewSize(colWidth, r.cells.cellSize.Height))

			if updateCell != nil {
				updateCell(TableCellID{row, col}, c)
			}
			r.visible[id] = c
			cells = append(cells, c)
			cellOffset += colWidth + separatorThickness
		}
	}

	for id, old := range wasVisible {
		if _, ok := r.visible[id]; !ok {
			r.pool.Release(old)
		}
	}
	r.SetObjects(cells)
}

func (r *tableCellsRenderer) returnAllToPool() {
	for _, cell := range r.BaseRenderer.Objects() {
		r.pool.Release(cell)
	}
	r.visible = make(map[TableCellID]fyne.CanvasObject)
	r.SetObjects(nil)
}

func (r *tableCellsRenderer) visibleRows() int {
	rows := math.Ceil(float64(r.cells.t.Size().Height)/float64(r.cells.cellSize.Height+theme.SeparatorThicknessSize()) + 1)

	dataRows := 0
	if f := r.cells.t.Length; f != nil {
		dataRows, _ = r.cells.t.Length()
	}
	return int(fyne.Min(float32(rows), float32(dataRows)))
}

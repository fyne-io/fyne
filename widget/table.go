package widget

import (
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

const tableDividerThickness = 1

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
type Table struct {
	BaseWidget

	Length       func() (int, int)
	CreateCell   func() fyne.CanvasObject
	UpdateCell   func(id TableCellID, template fyne.CanvasObject)
	OnSelected   func(id TableCellID)
	OnUnselected func(id TableCellID)

	selectedCell, hoveredCell *TableCellID
	cells                     *tableCells
	moveCallback              func()
	offset                    fyne.Position
	scroll                    *ScrollContainer
}

// NewTable returns a new performant table widget defined by the passed functions.
// The first returns the data size in rows and columns, second parameter is a function that returns cell
// template objects that can be cached and the third is used to apply data at specified data location to the
// passed template CanvasObject.
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
	colMarker := canvas.NewRectangle(theme.PrimaryColor())
	rowMarker := canvas.NewRectangle(theme.PrimaryColor())
	colHover := canvas.NewRectangle(theme.HoverColor())
	rowHover := canvas.NewRectangle(theme.HoverColor())

	cellSize := t.templateSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
	t.cells = newTableCells(t, cellSize)
	t.scroll = NewScrollContainer(t.cells)

	obj := []fyne.CanvasObject{colMarker, rowMarker, colHover, rowHover, t.scroll}
	r := &tableRenderer{t: t, scroll: t.scroll, rowMarker: rowMarker, colMarker: colMarker,
		rowHover: rowHover, colHover: colHover, cellSize: cellSize}
	r.SetObjects(obj)
	t.moveCallback = r.moveIndicators
	t.scroll.onOffsetChanged = func() {
		t.offset = t.scroll.Offset
		t.cells.Refresh()
		r.moveIndicators()
	}

	r.Layout(t.Size())
	return r
}

// Select will mark the specified cell as selected.
func (t *Table) Select(id TableCellID) {
	if t.selectedCell != nil && *t.selectedCell == id {
		return
	}
	if f := t.OnUnselected; f != nil && t.selectedCell != nil {
		f(*t.selectedCell)
	}
	t.selectedCell = &id

	t.scrollTo(id)
	if t.moveCallback != nil {
		t.moveCallback()
	}

	if f := t.OnSelected; f != nil {
		f(id)
	}
}

// Unselect will mark the cell provided by id as unselected.
func (t *Table) Unselect(id TableCellID) {
	if t.selectedCell == nil {
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

func (t *Table) scrollTo(id TableCellID) {
	if t.scroll == nil {
		return
	}
	scrollPos := t.offset

	cellPadded := t.templateSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
	cellX := id.Col * (cellPadded.Width + tableDividerThickness)
	if cellX < scrollPos.X {
		scrollPos.X = cellX
	} else if cellX+cellPadded.Width > scrollPos.X+t.scroll.size.Width {
		scrollPos.X = cellX + cellPadded.Width - t.scroll.size.Width
	}

	cellY := id.Row * (cellPadded.Height + tableDividerThickness)
	if cellY < scrollPos.Y {
		scrollPos.Y = cellY
	} else if cellY+cellPadded.Height > scrollPos.Y+t.scroll.size.Height {
		scrollPos.Y = cellY + cellPadded.Height - t.scroll.size.Height
	}
	t.scroll.Offset = scrollPos
	t.offset = scrollPos
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

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*tableRenderer)(nil)

type tableRenderer struct {
	widget.BaseRenderer
	t *Table

	scroll               *ScrollContainer
	rowMarker, colMarker *canvas.Rectangle
	rowHover, colHover   *canvas.Rectangle
	dividers             []fyne.CanvasObject

	cellSize fyne.Size
}

func (t *tableRenderer) Layout(s fyne.Size) {
	t.moveIndicators()

	t.scroll.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
	t.scroll.Resize(s.Subtract(fyne.NewSize(theme.Padding(), theme.Padding())))
}

func (t *tableRenderer) MinSize() fyne.Size {
	return t.cellSize
}

func (t *tableRenderer) Refresh() {
	t.cellSize = t.t.templateSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
	t.moveIndicators()

	t.colMarker.FillColor = theme.PrimaryColor()
	t.colMarker.Refresh()
	t.rowMarker.FillColor = theme.PrimaryColor()
	t.rowMarker.Refresh()

	t.colHover.FillColor = theme.HoverColor()
	t.colHover.Refresh()
	t.rowHover.FillColor = theme.HoverColor()
	t.rowHover.Refresh()

	t.t.cells.Refresh()
}

func (t *tableRenderer) moveColumnMarker(marker fyne.CanvasObject, col int) {
	if col == -1 {
		marker.Hide()
	} else {
		offX := col*(t.cellSize.Width+tableDividerThickness) - t.scroll.Offset.X
		x1 := theme.Padding() + offX
		x2 := x1 + t.cellSize.Width
		if x2 < theme.Padding() || x1 > t.t.size.Width {
			marker.Hide()
		} else {
			left := fyne.Max(theme.Padding(), x1)
			marker.Move(fyne.NewPos(left, 0))
			marker.Resize(fyne.NewSize(fyne.Min(x2, t.t.size.Width)-left, theme.Padding()))

			marker.Show()
		}
	}
	marker.Refresh()
}

func (t *tableRenderer) moveIndicators() {
	if t.t.selectedCell == nil {
		t.moveColumnMarker(t.colMarker, -1)
		t.moveRowMarker(t.rowMarker, -1)
	} else {
		t.moveColumnMarker(t.colMarker, t.t.selectedCell.Col)
		t.moveRowMarker(t.rowMarker, t.t.selectedCell.Row)
	}
	if t.t.hoveredCell == nil {
		t.moveColumnMarker(t.colHover, -1)
		t.moveRowMarker(t.rowHover, -1)
	} else {
		t.moveColumnMarker(t.colHover, t.t.hoveredCell.Col)
		t.moveRowMarker(t.rowHover, t.t.hoveredCell.Row)
	}

	colDivs := int(math.Ceil(float64(t.t.size.Width+tableDividerThickness) / float64(t.cellSize.Width+1)))
	rowDivs := int(math.Ceil(float64(t.t.size.Height+tableDividerThickness) / float64(t.cellSize.Height+1)))

	if len(t.dividers) < colDivs+rowDivs {
		for i := len(t.dividers); i < colDivs+rowDivs; i++ {
			t.dividers = append(t.dividers, NewSeparator())
		}

		obj := []fyne.CanvasObject{t.scroll, t.colMarker, t.rowMarker, t.colHover, t.rowHover}
		t.SetObjects(append(obj, t.dividers...))
	}

	divs := 0
	i := 0
	rows, cols := 0, 0
	if f := t.t.Length; f != nil {
		rows, cols = t.t.Length()
	}
	for x := theme.Padding() + t.scroll.Offset.X - (t.scroll.Offset.X % (t.cellSize.Width + tableDividerThickness)) - tableDividerThickness; x < t.scroll.Offset.X+t.t.size.Width && i < cols-1; x += t.cellSize.Width + tableDividerThickness {
		if x <= theme.Padding()+t.scroll.Offset.X {
			continue
		}
		i++

		t.dividers[divs].Move(fyne.NewPos(x-t.scroll.Offset.X, theme.Padding()))
		t.dividers[divs].Resize(fyne.NewSize(tableDividerThickness, t.t.size.Height-theme.Padding()))
		t.dividers[divs].Show()
		divs++
	}

	i = 0
	for y := theme.Padding() + t.scroll.Offset.Y - (t.scroll.Offset.Y % (t.cellSize.Height + tableDividerThickness)) - tableDividerThickness; y < t.scroll.Offset.Y+t.t.size.Height && i < rows-1; y += t.cellSize.Height + tableDividerThickness {
		if y <= theme.Padding()+t.scroll.Offset.Y {
			continue
		}
		i++

		t.dividers[divs].Move(fyne.NewPos(theme.Padding(), y-t.scroll.Offset.Y))
		t.dividers[divs].Resize(fyne.NewSize(t.t.size.Width-theme.Padding(), tableDividerThickness))
		t.dividers[divs].Show()
		divs++
	}

	for i := divs; i < len(t.dividers); i++ {
		t.dividers[divs].Hide()
	}
	canvas.Refresh(t.t)
}

func (t *tableRenderer) moveRowMarker(marker fyne.CanvasObject, row int) {
	if row == -1 {
		marker.Hide()
	} else {
		offY := row*(t.cellSize.Height+tableDividerThickness) - t.scroll.Offset.Y
		y1 := theme.Padding() + offY
		y2 := y1 + t.cellSize.Height
		if y2 < theme.Padding() || y1 > t.t.size.Height {
			marker.Hide()
		} else {
			top := fyne.Max(theme.Padding(), y1)
			marker.Move(fyne.NewPos(0, top))
			marker.Resize(fyne.NewSize(theme.Padding(), fyne.Min(y2, t.t.size.Height)-top))

			marker.Show()
		}
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

	col := e.Position.X / (c.cellSize.Width + tableDividerThickness)
	row := e.Position.Y / (c.cellSize.Height + tableDividerThickness)
	c.t.Select(TableCellID{row, col})
}

func (c *tableCells) hoverAt(pos fyne.Position) {
	if pos.X < 0 || pos.X >= c.Size().Width || pos.Y < 0 || pos.Y >= c.Size().Height {
		c.hoverOut()
		return
	}

	col := pos.X / (c.cellSize.Width + tableDividerThickness)
	row := pos.Y / (c.cellSize.Height + tableDividerThickness)
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
	return fyne.NewSize(r.cells.cellSize.Width*cols+(cols-1), r.cells.cellSize.Height*rows+(rows-1))
}

func (r *tableCellsRenderer) Refresh() {
	r.cells.propertyLock.Lock()
	defer r.cells.propertyLock.Unlock()
	oldSize := r.cells.cellSize
	r.cells.cellSize = r.cells.t.templateSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
	if oldSize != r.cells.cellSize { // theme changed probably
		r.returnAllToPool()
	}

	dataRows, dataCols := 0, 0
	if f := r.cells.t.Length; f != nil {
		dataRows, dataCols = r.cells.t.Length()
	}
	rows, cols := r.visibleCount()
	offX := r.cells.t.offset.X - (r.cells.t.offset.X % (r.cells.cellSize.Width + tableDividerThickness))
	minCol := offX / (r.cells.cellSize.Width + tableDividerThickness)
	maxCol := fyne.Min(minCol+cols, dataCols)
	offY := r.cells.t.offset.Y - (r.cells.t.offset.Y % (r.cells.cellSize.Height + tableDividerThickness))
	minRow := offY / (r.cells.cellSize.Height + tableDividerThickness)
	maxRow := fyne.Min(minRow+rows, dataRows)

	updateCell := r.cells.t.UpdateCell
	if updateCell == nil {
		fyne.LogError("Missing UpdateCell callback required for Table", nil)
	}

	wasVisible := r.visible
	r.visible = make(map[TableCellID]fyne.CanvasObject)
	var cells []fyne.CanvasObject
	for row := minRow; row < maxRow; row++ {
		for col := minCol; col < maxCol; col++ {
			id := TableCellID{row, col}
			c, ok := wasVisible[id]
			if !ok {
				c = r.pool.Obtain()
				if f := r.cells.t.CreateCell; f != nil && c == nil {
					c = f()
					c.Resize(r.cells.cellSize)
				}
				if c == nil {
					continue
				}

				c.Move(fyne.NewPos(theme.Padding()+col*r.cells.cellSize.Width+(col-1)*tableDividerThickness,
					theme.Padding()+row*r.cells.cellSize.Height+(row-1)*tableDividerThickness))
			}

			if updateCell != nil {
				updateCell(TableCellID{row, col}, c)
			}
			r.visible[id] = c
			cells = append(cells, c)
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

func (r *tableCellsRenderer) visibleCount() (int, int) {
	cols := math.Ceil(float64(r.cells.t.Size().Width)/float64(r.cells.cellSize.Width+tableDividerThickness) + 1)
	rows := math.Ceil(float64(r.cells.t.Size().Height)/float64(r.cells.cellSize.Height+tableDividerThickness) + 1)

	dataRows, dataCols := 0, 0
	if f := r.cells.t.Length; f != nil {
		dataRows, dataCols = r.cells.t.Length()
	}
	return fyne.Min(int(rows), dataRows), fyne.Min(int(cols), dataCols)
}

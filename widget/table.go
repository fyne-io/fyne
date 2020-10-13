package widget

import (
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

const tableDividerThickness = 1

// Declare conformity with Widget interface.
var _ fyne.Widget = (*Table)(nil)

// Table widget is a grid of items that can be scrolled and a cell selected.
// It's performance is provided by caching cell templates created with CreateCell and re-using them with UpdateCell.
// The size of the content rows/columns is returned by the Length callback.
type Table struct {
	BaseWidget

	Length         func() (int, int)
	CreateCell     func() fyne.CanvasObject
	UpdateCell     func(row int, col int, template fyne.CanvasObject)
	OnCellSelected func(row int, col int)

	SelectedRow, SelectedColumn int
	cells                       *tableCells
	moveCallback                func()
	offset                      fyne.Position
}

// NewTable returns a new performant table widget defined by the passed functions.
// The first returns the data size in rows and columns, second parameter is a function that returns cell
// template objects that can be cached and the third is used to apply data at specified data location to the
// passed template CanvasObject.
func NewTable(length func() (int, int), create func() fyne.CanvasObject, update func(int, int, fyne.CanvasObject)) *Table {
	t := &Table{Length: length, CreateCell: create, UpdateCell: update, SelectedRow: -1, SelectedColumn: -1}
	t.ExtendBaseWidget(t)
	return t
}

// CreateRenderer returns a new renderer for the table.
//
// Implements: fyne.Widget
func (t *Table) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	marker1 := canvas.NewRectangle(theme.PrimaryColor())
	marker2 := canvas.NewRectangle(theme.PrimaryColor())

	cellSize := t.templateSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
	t.cells = newTableCells(t, cellSize)
	scroll := NewScrollContainer(t.cells)

	obj := []fyne.CanvasObject{marker1, marker2, scroll}
	r := &tableRenderer{t: t, scroll: scroll, rowMarker: marker1, colMarker: marker2, cellSize: cellSize}
	r.SetObjects(obj)
	t.moveCallback = r.moveIndicators
	scroll.onOffsetChanged = func() {
		t.offset = scroll.Offset
		t.cells.Refresh()
		r.moveIndicators()
	}

	r.Layout(t.Size())
	return r
}

func (t *Table) templateSize() fyne.Size {
	if t.CreateCell == nil {
		fyne.LogError("Missing CreateCell callback required for Table", nil)
		return fyne.Size{}
	}

	template := t.CreateCell() // don't use cache, we need new template
	return template.MinSize()
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*tableRenderer)(nil)

type tableRenderer struct {
	widget.BaseRenderer
	t *Table

	scroll               *ScrollContainer
	rowMarker, colMarker *canvas.Rectangle
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

	t.t.cells.Refresh()
}

func (t *tableRenderer) moveIndicators() {
	if t.t.SelectedColumn == -1 {
		t.colMarker.Hide()
	} else {
		offX := t.t.SelectedColumn*(t.cellSize.Width+tableDividerThickness) - t.scroll.Offset.X
		x1 := theme.Padding() + offX
		x2 := x1 + t.cellSize.Width
		if x2 < theme.Padding() || x1 > t.t.size.Width {
			t.colMarker.Hide()
		} else {
			left := fyne.Max(theme.Padding(), x1)
			t.colMarker.Move(fyne.NewPos(left, 0))
			t.colMarker.Resize(fyne.NewSize(fyne.Min(x2, t.t.size.Width)-left, theme.Padding()))

			t.colMarker.Show()
		}
	}
	t.colMarker.Refresh()

	if t.t.SelectedRow == -1 {
		t.rowMarker.Hide()
	} else {
		offY := t.t.SelectedRow*(t.cellSize.Height+tableDividerThickness) - t.scroll.Offset.Y
		y1 := theme.Padding() + offY
		y2 := y1 + t.cellSize.Height
		if y2 < theme.Padding() || y1 > t.t.size.Height {
			t.rowMarker.Hide()
		} else {
			top := fyne.Max(theme.Padding(), y1)
			t.rowMarker.Move(fyne.NewPos(0, top))
			t.rowMarker.Resize(fyne.NewSize(theme.Padding(), fyne.Min(y2, t.t.size.Height)-top))

			t.rowMarker.Show()
		}
	}
	t.rowMarker.Refresh()

	colDivs := int(math.Ceil(float64(t.t.size.Width+tableDividerThickness) / float64(t.cellSize.Width+1)))
	rowDivs := int(math.Ceil(float64(t.t.size.Height+tableDividerThickness) / float64(t.cellSize.Height+1)))

	if len(t.dividers) < colDivs+rowDivs {
		for i := len(t.dividers); i < colDivs+rowDivs; i++ {
			t.dividers = append(t.dividers, widget.NewDivider())
		}

		obj := []fyne.CanvasObject{t.scroll, t.colMarker, t.rowMarker}
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

// Declare conformity with Widget interface.
var _ fyne.Widget = (*tableCells)(nil)

// Declare conformity with Tappable interface.
var _ fyne.Tappable = (*tableCells)(nil)

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
	return &tableCellsRenderer{cells: c, pool: &syncPool{}, visible: make(map[cellID]fyne.CanvasObject)}
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
		c.t.SelectedColumn = -1
		c.t.SelectedRow = -1
		c.t.Refresh()
		return
	}

	c.t.SelectedColumn = e.Position.X / (c.cellSize.Width + tableDividerThickness)
	c.t.SelectedRow = e.Position.Y / (c.cellSize.Height + tableDividerThickness)

	if c.t.OnCellSelected != nil {
		c.t.OnCellSelected(c.t.SelectedRow, c.t.SelectedColumn)
	}

	if c.t.moveCallback != nil {
		c.t.moveCallback()
	}
}

type cellID struct {
	row, col int
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*tableCellsRenderer)(nil)

type tableCellsRenderer struct {
	widget.BaseRenderer

	cells   *tableCells
	pool    pool
	visible map[cellID]fyne.CanvasObject
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

	wasVisible := r.visible
	r.visible = make(map[cellID]fyne.CanvasObject)
	var cells []fyne.CanvasObject
	for y := minRow; y < maxRow; y++ {
		for x := minCol; x < maxCol; x++ {
			id := cellID{y, x}
			c, ok := wasVisible[id]
			if !ok {
				c = r.pool.Obtain()
				if c == nil && r.cells.t.CreateCell != nil {
					c = r.cells.t.CreateCell()
					c.Resize(r.cells.cellSize)
				}
				if c == nil {
					continue
				}

				c.Move(fyne.NewPos(theme.Padding()+x*r.cells.cellSize.Width+(x-1)*tableDividerThickness,
					theme.Padding()+y*r.cells.cellSize.Height+(y-1)*tableDividerThickness))

				if f := r.cells.t.UpdateCell; f != nil {
					r.cells.t.UpdateCell(y, x, c)
				} else {
					fyne.LogError("Missing UpdateCell callback required for Table", nil)
				}
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
	r.visible = make(map[cellID]fyne.CanvasObject)
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

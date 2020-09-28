package widget

import (
	"image/color"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

// Table widget is a grid of items that can be scrolled and a cell selected.
// It's performance is provided by caching cell templates created with NewCell and re-using them with UpdateCell.
// The size of the content rows/columns is returned by the DataSize callback.
type Table struct {
	BaseWidget

	DataSize   func() (int, int)
	NewCell    func() fyne.CanvasObject
	UpdateCell func(int, int, fyne.CanvasObject)

	cells                       *tableCells
	scroll                      *ScrollContainer
	SelectedRow, SelectedColumn int
}

// NewTable returns a new performant table widget defined by the passed functions.
// The first returns the data size in rows and columns, second parameter is a function that returns cell
// template objects that can be cached and the third is used to apply data at specified data location to the
// passed template CanvasObject.
func NewTable(size func() (int, int), create func() fyne.CanvasObject, update func(int, int, fyne.CanvasObject)) *Table {
	return &Table{DataSize: size, NewCell: create, UpdateCell: update, SelectedRow: -1, SelectedColumn: -1}
}

// CreateRenderer returns a new renderer for the table.
//
// Implements: fyne.Widget
func (t *Table) CreateRenderer() fyne.WidgetRenderer {
	marker1 := canvas.NewRectangle(theme.PrimaryColor())
	marker2 := canvas.NewRectangle(theme.PrimaryColor())

	t.cells = newTableCells(t)
	t.scroll = NewScrollContainer(t.cells)

	obj := []fyne.CanvasObject{marker1, marker2, t.scroll}
	r := &tableRenderer{t: t, rowMarker: marker1, colMarker: marker2, objects: obj}
	t.scroll.onOffsetChanged = r.moveOverlay

	r.Layout(t.Size())
	return r
}

// Resize updates this table size and adjusts the content scroller to fit.
//
// Implements: fyne.Widget
func (t *Table) Resize(s fyne.Size) {
	t.BaseWidget.Resize(s)

	if t.scroll != nil {
		t.scroll.Resize(s.Subtract(fyne.NewSize(theme.Padding(), theme.Padding())))
	}
}

type tableRenderer struct {
	t *Table

	rowMarker, colMarker *canvas.Rectangle
	dividers             []fyne.CanvasObject

	objects []fyne.CanvasObject
}

func (t *tableRenderer) moveOverlay() {
	template := t.t.NewCell()
	cellSize := template.MinSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))

	if t.t.SelectedColumn == -1 {
		t.colMarker.Hide()
	} else {
		offX := t.t.SelectedColumn*(cellSize.Width+1) - t.t.scroll.Offset.X
		x1 := theme.Padding() + offX
		x2 := x1 + cellSize.Width
		if x2 < theme.Padding() || x1 > t.t.size.Width {
			t.colMarker.Hide()
		} else {
			t.colMarker.Show()

			left := fyne.Max(theme.Padding(), x1)
			t.colMarker.Move(fyne.NewPos(left, 0))
			t.colMarker.Resize(fyne.NewSize(fyne.Min(x2, t.t.size.Width)-left, theme.Padding()))
		}
	}

	if t.t.SelectedRow == -1 {
		t.colMarker.Hide()
	} else {
		offY := t.t.SelectedRow*(cellSize.Height+1) - t.t.scroll.Offset.Y
		y1 := theme.Padding() + offY
		y2 := y1 + cellSize.Height
		if y2 < theme.Padding() || y1 > t.t.size.Height {
			t.rowMarker.Hide()
		} else {
			t.rowMarker.Show()

			top := fyne.Max(theme.Padding(), y1)
			t.rowMarker.Move(fyne.NewPos(0, top))
			t.rowMarker.Resize(fyne.NewSize(theme.Padding(), fyne.Min(y2, t.t.size.Height)-top))
		}
	}

	colDivs := int(math.Ceil(float64(t.t.size.Width+1) / float64(cellSize.Width+1)))   // +1 for div width
	rowDivs := int(math.Ceil(float64(t.t.size.Height+1) / float64(cellSize.Height+1))) // +1 for div width

	if len(t.dividers) < colDivs+rowDivs {
		for i := len(t.dividers); i < colDivs+rowDivs; i++ {
			t.dividers = append(t.dividers, canvas.NewRectangle(theme.ShadowColor()))
		}

		obj := []fyne.CanvasObject{t.t.scroll, t.colMarker, t.rowMarker}
		t.objects = append(obj, t.dividers...)
	}

	divs := 0
	i := 0
	rows, cols := t.t.DataSize()
	for x := theme.Padding() + t.t.scroll.Offset.X - (t.t.scroll.Offset.X % (cellSize.Width + 1)) - 1; x < t.t.scroll.Offset.X+t.t.size.Width && i < cols-1; x += cellSize.Width + 1 {
		if x <= theme.Padding()+t.t.scroll.Offset.X {
			continue
		}
		i++

		t.dividers[divs].Move(fyne.NewPos(x-t.t.scroll.Offset.X, theme.Padding()))
		t.dividers[divs].Resize(fyne.NewSize(1, t.t.size.Height-theme.Padding()*2))
		t.dividers[divs].Show()
		divs++
	}

	i = 0
	for y := theme.Padding() + t.t.scroll.Offset.Y - (t.t.scroll.Offset.Y % (cellSize.Height + 1)) - 1; y < t.t.scroll.Offset.Y+t.t.size.Height && i < rows-1; y += cellSize.Height + 1 {
		if y <= theme.Padding()+t.t.scroll.Offset.Y {
			continue
		}
		i++

		t.dividers[divs].Move(fyne.NewPos(theme.Padding(), y-t.t.scroll.Offset.Y))
		t.dividers[divs].Resize(fyne.NewSize(t.t.size.Width-theme.Padding()*2, 1))
		t.dividers[divs].Show()
		divs++
	}

	for i := divs; i < len(t.dividers); i++ {
		t.dividers[divs].Hide()
	}
}

func (t *tableRenderer) Layout(s fyne.Size) {
	t.moveOverlay()

	t.t.scroll.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
	t.t.scroll.Resize(s.Subtract(fyne.NewSize(theme.Padding(), theme.Padding())))
}

func (t *tableRenderer) MinSize() fyne.Size {
	template := t.t.NewCell()
	rows, cols := t.t.DataSize()
	return fyne.NewSize(template.MinSize().Width*cols+(cols-1), template.MinSize().Height*rows+(rows-1))
}

func (t *tableRenderer) Refresh() {
	t.colMarker.FillColor = theme.PrimaryColor()
	t.colMarker.Refresh()
	t.rowMarker.FillColor = theme.PrimaryColor()
	t.rowMarker.Refresh()

	for _, div := range t.dividers {
		div.(*canvas.Rectangle).FillColor = theme.ShadowColor()
		div.Refresh()
	}
	t.t.cells.Refresh()
}

func (t *tableRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (t *tableRenderer) Objects() []fyne.CanvasObject {
	return t.objects
}

func (t *tableRenderer) Destroy() {
}

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
	return &tableCellsRenderer{t: c.t, pool: &syncPool{}, objects: []fyne.CanvasObject{}}
}

type tableCellsRenderer struct {
	t       *Table
	pool    pool
	objects []fyne.CanvasObject
}

func (r *tableCellsRenderer) Layout(s fyne.Size) {
	template := r.pool.Obtain()
	if template == nil {
		template = r.t.NewCell()
	}
	defer r.pool.Release(template)
	rows, cols := r.t.DataSize()
	cellSize := template.MinSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))

	// TODO only visible
	if len(r.objects) == 0 {
		r.ensureCells()
	}
	i := 0
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			c := r.objects[i]
			c.Move(fyne.NewPos(theme.Padding()+x*cellSize.Width+(x-1), theme.Padding()+y*cellSize.Height+(y-1)))
			c.Resize(cellSize)
			i++
		}
	}
}

func (r *tableCellsRenderer) MinSize() fyne.Size {
	template := r.pool.Obtain()
	if template == nil {
		template = r.t.NewCell()
	}
	defer r.pool.Release(template)
	rows, cols := r.t.DataSize()
	return fyne.NewSize((template.MinSize().Width+theme.Padding()*2)*cols+(cols-1), (template.MinSize().Height+theme.Padding()*2)*rows+(rows-1))
}

func (r *tableCellsRenderer) Refresh() {
	r.ensureCells() // TODO force refresh, once caching doesn't reset them all each time :)
}

func (r *tableCellsRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *tableCellsRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *tableCellsRenderer) Destroy() {
	r.clearPool()
}

func (r *tableCellsRenderer) clearPool() {
	for _, cell := range r.objects {
		r.pool.Release(cell)
	}
	r.objects = nil
}

func (r *tableCellsRenderer) ensureCells() {
	r.clearPool()

	rows, cols := r.t.DataSize()
	var cells []fyne.CanvasObject
	i := 0
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			c := r.pool.Obtain()
			if c == nil {
				c = r.t.NewCell()
			}

			r.t.UpdateCell(y, x, c)
			cells = append(cells, c)
			i++
		}
	}
	r.objects = cells
}

package widget

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
)

type selectable struct {
	BaseWidget
	cursorRow, cursorColumn int

	// selectRow and selectColumn represent the selection start location
	// The selection will span from selectRow/Column to CursorRow/Column -- note that the cursor
	// position may occur before or after the select start position in the text.
	selectRow, selectColumn int

	selecting, password bool
	style               fyne.TextStyle

	provider *RichText
	theme    fyne.Theme

	// TODO maybe render?
	selections []fyne.CanvasObject
}

func (s *selectable) CreateRenderer() fyne.WidgetRenderer {
	return &selectableRenderer{sel: s}
}

func (s *selectable) Cursor() desktop.Cursor {
	return desktop.TextCursor
}

func (s *selectable) DragEnd() {
	if s.cursorColumn == s.selectColumn && s.cursorRow == s.selectRow {
		s.selecting = false
	}

	shouldRefresh := !s.selecting
	if shouldRefresh {
		s.Refresh()
	}
}

func (s *selectable) Dragged(d *fyne.DragEvent) {
	if !s.selecting {
		startPos := d.Position.Subtract(d.Dragged)
		s.selectRow, s.selectColumn = s.getRowCol(startPos)
		s.selecting = true
	}

	s.updateMousePointer(d.Position)
	s.Refresh()
}

func (s *selectable) MouseDown(m *desktop.MouseEvent) {
	//if e.isTripleTap(time.Now().UnixMilli()) {
	//	e.selectCurrentRow()
	//	return
	//}
	if s.selecting && m.Button == desktop.MouseButtonPrimary {
		s.selecting = false
	}

	if m.Button == desktop.MouseButtonPrimary {
		s.updateMousePointer(m.Position)
	}
}

func (s *selectable) MouseUp(ev *desktop.MouseEvent) {
	if ev.Button == desktop.MouseButtonSecondary {
		c := fyne.CurrentApp().Driver().CanvasForObject(s)
		if c == nil {
			return
		}

		m := fyne.NewMenu("",
			fyne.NewMenuItem(lang.L("Copy"), func() {
				fyne.CurrentApp().Clipboard().SetContent(s.SelectedText())
			}))
		ShowPopUpMenuAtPosition(m, c, ev.AbsolutePosition)

		return
	}

	start, _ := s.selection()
	if start == -1 && s.selecting {
		s.selecting = false
	}
}

// SelectedText returns the text currently selected in this Entry.
// If there is no selection it will return the empty string.
func (s *selectable) SelectedText() string {
	if s == nil || !s.selecting {
		return ""
	}

	start, stop := s.selection()
	if start == stop {
		return ""
	}
	r := ([]rune)(s.provider.String())
	return string(r[start:stop])
}

func (s *selectable) cursorColAt(text []rune, pos fyne.Position) int {
	th := s.theme
	textSize := th.Size(theme.SizeNameText)
	innerPad := th.Size(theme.SizeNameInnerPadding)

	for i := 0; i < len(text); i++ {
		str := string(text[0:i])
		wid := fyne.MeasureText(str, textSize, fyne.TextStyle{}).Width                 // todo e.TextStyle
		charWid := fyne.MeasureText(string(text[i]), textSize, fyne.TextStyle{}).Width // todo e.TextStyle
		if pos.X < innerPad+wid+(charWid/2) {
			return i
		}
	}
	return len(text)
}

func (s *selectable) getRowCol(p fyne.Position) (int, int) {
	th := s.theme
	textSize := th.Size(theme.SizeNameText)
	innerPad := th.Size(theme.SizeNameInnerPadding)

	rowHeight := s.provider.charMinSize(false, fyne.TextStyle{}, textSize).Height // TODO (e.Password, e.TextStyle, textSize).Height
	row := int(math.Floor(float64(p.Y-innerPad+th.Size(theme.SizeNameLineSpacing)) / float64(rowHeight)))
	col := 0
	if row < 0 {
		row = 0
	} else if row >= s.provider.rows() {
		row = s.provider.rows() - 1
		col = s.provider.rowLength(row)
	} else {
		col = s.cursorColAt(s.provider.row(row), p)
	}

	return row, col
}

// selection returns the start and end text positions for the selected span of text
// Note: this functionality depends on the relationship between the selection start row/col and
// the current cursor row/column.
// eg: (whitespace for clarity, '_' denotes cursor)
//
//	"T  e  s [t  i]_n  g" == 3, 5
//	"T  e  s_[t  i] n  g" == 3, 5
//	"T  e_[s  t  i] n  g" == 2, 5
func (s *selectable) selection() (int, int) {
	noSelection := !s.selecting || (s.cursorRow == s.selectRow && s.cursorColumn == s.selectColumn)

	if noSelection {
		return -1, -1
	}

	// Find the selection start
	rowA, colA := s.cursorRow, s.cursorColumn
	rowB, colB := s.selectRow, s.selectColumn
	// Reposition if the cursors row is more than select start row, or if the row is the same and
	// the cursors col is more that the select start column
	if rowA > s.selectRow || (rowA == s.selectRow && colA > s.selectColumn) {
		rowA, colA = s.selectRow, s.selectColumn
		rowB, colB = s.cursorRow, s.cursorColumn
	}

	return textPosFromRowCol(rowA, colA, s.provider), textPosFromRowCol(rowB, colB, s.provider)
}

// Obtains textual position from a given row and col
// expects a read or write lock to be held by the caller
func textPosFromRowCol(row, col int, prov *RichText) int {
	b := prov.rowBoundary(row)
	if b == nil {
		return col
	}
	return b.begin + col
}

func (s *selectable) updateMousePointer(p fyne.Position) {
	row, col := s.getRowCol(p)
	s.cursorRow, s.cursorColumn = row, col

	if !s.selecting {
		s.selectRow = row
		s.selectColumn = col
	}
}

type selectableRenderer struct {
	sel *selectable
}

func (r *selectableRenderer) Destroy() {
}

func (r *selectableRenderer) Layout(fyne.Size) {
}

func (r *selectableRenderer) MinSize() fyne.Size {
	return fyne.Size{}
}

func (r *selectableRenderer) Objects() []fyne.CanvasObject {
	return r.sel.selections
}

func (r *selectableRenderer) Refresh() {
	r.buildSelection()

	selections := r.sel.selections
	v := fyne.CurrentApp().Settings().ThemeVariant()

	selectionColor := r.sel.theme.Color(theme.ColorNameSelection, v)
	for _, selection := range selections {
		rect := selection.(*canvas.Rectangle)
		rect.FillColor = selectionColor
	}

	canvas.Refresh(r.sel)
}

// This process builds a slice of rectangles:
// - one entry per row of text
// - ordered by row order as they occur in multiline text
// This process could be optimized in the scenario where the user is selecting upwards:
// If the upwards case instead produces an order-reversed slice then only the newest rectangle would
// require movement and resizing. The existing solution creates a new rectangle and then moves/resizes
// all rectangles to comply with the occurrence order as stated above.
func (r *selectableRenderer) buildSelection() {
	th := r.sel.theme
	v := fyne.CurrentApp().Settings().ThemeVariant()
	textSize := th.Size(theme.SizeNameText)

	cursorRow, cursorCol := r.sel.cursorRow, r.sel.cursorColumn
	selectRow, selectCol := -1, -1
	if r.sel.selecting {
		selectRow = r.sel.selectRow
		selectCol = r.sel.selectColumn
	}

	if selectRow == -1 || (cursorRow == selectRow && cursorCol == selectCol) {
		r.sel.selections = r.sel.selections[:0]

		return
	}

	provider := r.sel.provider
	innerPad := th.Size(theme.SizeNameInnerPadding)
	// Convert column, row into x,y
	getCoordinates := func(column int, row int) (float32, float32) {
		sz := provider.lineSizeToColumn(column, row, textSize, innerPad)
		return sz.Width, sz.Height*float32(row) - th.Size(theme.SizeNameInputBorder) + innerPad
	}

	lineHeight := r.sel.provider.charMinSize(r.sel.password, r.sel.style, textSize).Height

	minmax := func(a, b int) (int, int) {
		if a < b {
			return a, b
		}
		return b, a
	}

	// The remainder of the function calculates the set of boxes and add them to r.selection

	selectStartRow, selectEndRow := minmax(selectRow, cursorRow)
	selectStartCol, selectEndCol := minmax(selectCol, cursorCol)
	if selectRow < cursorRow {
		selectStartCol, selectEndCol = selectCol, cursorCol
	}
	if selectRow > cursorRow {
		selectStartCol, selectEndCol = cursorCol, selectCol
	}
	rowCount := selectEndRow - selectStartRow + 1

	// trim r.selection to remove unwanted old rectangles
	if len(r.sel.selections) > rowCount {
		r.sel.selections = r.sel.selections[:rowCount]
	}

	// build a rectangle for each row and add it to r.selection
	for i := 0; i < rowCount; i++ {
		if len(r.sel.selections) <= i {
			box := canvas.NewRectangle(th.Color(theme.ColorNameSelection, v))
			r.sel.selections = append(r.sel.selections, box)
		}

		// determine starting/ending columns for this rectangle
		row := selectStartRow + i
		startCol, endCol := selectStartCol, selectEndCol
		if selectStartRow < row {
			startCol = 0
		}
		if selectEndRow > row {
			endCol = provider.rowLength(row)
		}

		// translate columns and row into draw coordinates
		x1, y1 := getCoordinates(startCol, row)
		x2, _ := getCoordinates(endCol, row)

		// resize and reposition each rectangle
		r.sel.selections[i].Resize(fyne.NewSize(x2-x1+1, lineHeight))
		r.sel.selections[i].Move(fyne.NewPos(x1-1, y1))
	}
}

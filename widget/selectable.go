package widget

import (
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/driver/mobile"
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

	focussed, selecting, selectEnded, password bool
	sizeName                                   fyne.ThemeSizeName
	style                                      fyne.TextStyle

	provider *RichText
	theme    fyne.Theme
	focus    fyne.Focusable

	// doubleTappedAtUnixMillis stores the time the entry was last DoubleTapped
	// used for deciding whether the next MouseDown/TouchDown is a triple-tap or not
	doubleTappedAtUnixMillis int64
}

func (s *selectable) CreateRenderer() fyne.WidgetRenderer {
	return &selectableRenderer{sel: s}
}

func (s *selectable) Cursor() desktop.Cursor {
	return desktop.TextCursor
}

func (s *selectable) DoubleTapped(p *fyne.PointEvent) {
	s.doubleTappedAtUnixMillis = time.Now().UnixMilli()
	s.updateMousePointer(p.Position)
	row := s.provider.row(s.cursorRow)
	start, end := getTextWhitespaceRegion(row, s.cursorColumn, false)
	if start == -1 || end == -1 {
		return
	}

	s.selectRow = s.cursorRow
	s.selectColumn = start
	s.cursorColumn = end

	s.selecting = true
	s.grabFocus()
	s.Refresh()
}

func (s *selectable) DragEnd() {
	if s.cursorColumn == s.selectColumn && s.cursorRow == s.selectRow {
		s.selecting = false
	}

	shouldRefresh := !s.selecting
	if shouldRefresh {
		s.Refresh()
	}
	s.selectEnded = true
}

func (s *selectable) Dragged(d *fyne.DragEvent) {
	s.dragged(d, true)
}

func (s *selectable) dragged(d *fyne.DragEvent, focus bool) {
	if !s.selecting || s.selectEnded {
		s.selectEnded = false
		s.updateMousePointer(d.Position)

		startPos := d.Position.Subtract(d.Dragged)
		s.selectRow, s.selectColumn = s.getRowCol(startPos)
		s.selecting = true

		s.grabFocus()
	}

	s.updateMousePointer(d.Position)
	s.Refresh()
}

func (s *selectable) MouseDown(m *desktop.MouseEvent) {
	if isTripleTap(s.doubleTappedAtUnixMillis, time.Now().UnixMilli()) {
		s.selectCurrentRow(false)
		return
	}
	s.grabFocus()
	if s.selecting && m.Button == desktop.MouseButtonPrimary {
		s.selecting = false
	}
}

func (s *selectable) MouseUp(ev *desktop.MouseEvent) {
	if ev.Button == desktop.MouseButtonSecondary {
		return
	}

	start, _ := s.selection()
	if (start == -1 || (s.selectRow == s.cursorRow && s.selectColumn == s.cursorColumn)) && s.selecting {
		s.selecting = false
	}
	s.Refresh()
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

func (s *selectable) Tapped(*fyne.PointEvent) {
	if !fyne.CurrentDevice().IsMobile() {
		return
	}

	if s.doubleTappedAtUnixMillis != 0 {
		s.doubleTappedAtUnixMillis = 0
		return // was a triple (TappedDouble plus Tapped)
	}
	s.selecting = false
	s.Refresh()
}

func (s *selectable) TappedSecondary(ev *fyne.PointEvent) {
	c := fyne.CurrentApp().Driver().CanvasForObject(s.focus.(fyne.CanvasObject))
	if c == nil {
		return
	}

	m := fyne.NewMenu("",
		fyne.NewMenuItem(lang.L("Copy"), func() {
			fyne.CurrentApp().Clipboard().SetContent(s.SelectedText())
		}))
	ShowPopUpMenuAtPosition(m, c, ev.AbsolutePosition)
}

func (s *selectable) TouchCancel(m *mobile.TouchEvent) {
	s.TouchUp(m)
}

func (s *selectable) TouchDown(m *mobile.TouchEvent) {
	if isTripleTap(s.doubleTappedAtUnixMillis, time.Now().UnixMilli()) {
		s.selectCurrentRow(true)
		return
	}
}

func (s *selectable) TouchUp(*mobile.TouchEvent) {
}

func (s *selectable) TypedShortcut(sh fyne.Shortcut) {
	switch sh.(type) {
	case *fyne.ShortcutCopy:
		fyne.CurrentApp().Clipboard().SetContent(s.SelectedText())
	}
}

func (s *selectable) cursorColAt(text []rune, pos fyne.Position) int {
	th := s.theme
	textSize := th.Size(s.getSizeName())
	innerPad := th.Size(theme.SizeNameInnerPadding)

	for i := 0; i < len(text); i++ {
		str := string(text[0:i])
		wid := fyne.MeasureText(str, textSize, s.style).Width
		charWid := fyne.MeasureText(string(text[i]), textSize, s.style).Width
		if pos.X < innerPad+wid+(charWid/2) {
			return i
		}
	}
	return len(text)
}

func (s *selectable) getRowCol(p fyne.Position) (int, int) {
	th := s.theme
	textSize := th.Size(s.getSizeName())
	innerPad := th.Size(theme.SizeNameInnerPadding)

	rowHeight := s.provider.charMinSize(false, s.style, textSize).Height // TODO handle Password
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

// Selects the row where the cursorColumn is currently positioned
func (s *selectable) selectCurrentRow(focus bool) {
	s.grabFocus()
	provider := s.provider
	s.selectRow = s.cursorRow
	s.selectColumn = 0
	s.cursorColumn = provider.rowLength(s.cursorRow)
	s.Refresh()
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

func (s *selectable) getSizeName() fyne.ThemeSizeName {
	if s.sizeName != "" {
		return s.sizeName
	}
	return theme.SizeNameText
}

type selectableRenderer struct {
	sel *selectable

	selections []fyne.CanvasObject
}

func (r *selectableRenderer) Destroy() {
}

func (r *selectableRenderer) Layout(fyne.Size) {
}

func (r *selectableRenderer) MinSize() fyne.Size {
	return fyne.Size{}
}

func (r *selectableRenderer) Objects() []fyne.CanvasObject {
	return r.selections
}

func (r *selectableRenderer) Refresh() {
	r.buildSelection()
	selections := r.selections
	v := fyne.CurrentApp().Settings().ThemeVariant()

	selectionColor := r.sel.theme.Color(theme.ColorNameSelection, v)
	for _, selection := range selections {
		rect := selection.(*canvas.Rectangle)
		rect.FillColor = selectionColor

		if r.sel.focussed {
			rect.Show()
		} else {
			rect.Hide()
		}
	}

	canvas.Refresh(r.sel.impl)
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
	textSize := th.Size(r.sel.getSizeName())

	cursorRow, cursorCol := r.sel.cursorRow, r.sel.cursorColumn
	selectRow, selectCol := -1, -1
	if r.sel.selecting {
		selectRow = r.sel.selectRow
		selectCol = r.sel.selectColumn
	}

	if selectRow == -1 || (cursorRow == selectRow && cursorCol == selectCol) {
		r.selections = r.selections[:0]
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
	if len(r.selections) > rowCount {
		r.selections = r.selections[:rowCount]
	}

	// build a rectangle for each row and add it to r.selection
	for i := 0; i < rowCount; i++ {
		if len(r.selections) <= i {
			box := canvas.NewRectangle(th.Color(theme.ColorNameSelection, v))
			r.selections = append(r.selections, box)
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
		r.selections[i].Resize(fyne.NewSize(x2-x1+1, lineHeight))
		r.selections[i].Move(fyne.NewPos(x1-1, y1))
	}
}

func (s *selectable) grabFocus() {
	if c := fyne.CurrentApp().Driver().CanvasForObject(s.focus.(fyne.CanvasObject)); c != nil {
		c.Focus(s.focus)
	}
}

func isTripleTap(double, nowMilli int64) bool {
	return nowMilli-double <= fyne.CurrentApp().Driver().DoubleTapDelay().Milliseconds()
}

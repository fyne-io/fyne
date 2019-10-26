package widget

import (
	"image/color"
	"math"
	"strings"
	"sync"
	"unicode"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

const (
	multiLineRows            = 3
	doubleClickWordSeperator = "`~!@#$%^&*()-=+[{]}\\|;:'\",.<>/?"
)

type entryRenderer struct {
	text         *textProvider
	placeholder  *textProvider
	line, cursor *canvas.Rectangle
	selection    []fyne.CanvasObject

	objects []fyne.CanvasObject
	entry   *Entry
}

// MinSize calculates the minimum size of an entry widget.
// This is based on the contained text with a standard amount of padding added.
// If MultiLine is true then we will reserve space for at leasts 3 lines
func (e *entryRenderer) MinSize() fyne.Size {
	minSize := e.placeholder.MinSize()

	if e.text.len() > 0 {
		minSize = e.text.MinSize()
	}

	if e.entry.MultiLine == true {
		// ensure multiline height is at least charMinSize * multilineRows
		minSize.Height = fyne.Max(minSize.Height, e.text.charMinSize().Height*multiLineRows)
	}

	return minSize.Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
}

// This process builds a slice of rectangles:
// - one entry per row of text
// - ordered by row order as they occur in multiline text
// This process could be optimized in the scenario where the user is selecting upwards:
// If the upwards case instead produces an order-reversed slice then only the newest rectangle would
// require movement and resizing. The existing solution creates a new rectangle and then moves/resizes
// all rectangles to comply with the occurance order as stated above.
func (e *entryRenderer) buildSelection() {

	e.entry.RLock()
	cursorRow, cursorCol := e.entry.CursorRow, e.entry.CursorColumn
	selectRow, selectCol := -1, -1
	if e.entry.selecting {
		selectRow = e.entry.selectRow
		selectCol = e.entry.selectColumn
	}
	e.entry.RUnlock()

	if selectRow == -1 {
		e.selection = e.selection[:0]
		return
	}

	textRenderer := Renderer(e.text).(*textRenderer)

	// Convert column, row into x,y
	getCoordinates := func(column int, row int) (int, int) {
		sz := textRenderer.lineSizeToColumn(column, row)
		return sz.Width + theme.Padding()*2, sz.Height*row + theme.Padding()*2
	}

	lineHeight := e.text.charMinSize().Height

	minmax := func(a, b int) (int, int) {
		if a < b {
			return a, b
		}
		return b, a
	}

	// The remainder of the function calculates the set of boxes and add them to e.selection

	selectStartRow, selectEndRow := minmax(selectRow, cursorRow)
	selectStartCol, selectEndCol := minmax(selectCol, cursorCol)
	if selectRow < cursorRow {
		selectStartCol, selectEndCol = selectCol, cursorCol
	}
	if selectRow > cursorRow {
		selectStartCol, selectEndCol = cursorCol, selectCol
	}
	rowCount := selectEndRow - selectStartRow + 1

	// trim e.selection to remove unwanted old rectangles
	if len(e.selection) > rowCount {
		e.selection = e.selection[:rowCount]
	}

	// build a rectangle for each row and add it to e.selection
	for i := 0; i < rowCount; i++ {
		if len(e.selection) <= i {
			box := canvas.NewRectangle(theme.FocusColor())
			e.selection = append(e.selection, box)
		}

		// determine starting/ending columns for this rectangle
		row := selectStartRow + i
		startCol, endCol := selectStartCol, selectEndCol
		if selectStartRow < row {
			startCol = 0
		}
		if selectEndRow > row {
			endCol = textRenderer.provider.rowLength(row)
		}

		// translate columns and row into draw coordinates
		x1, y1 := getCoordinates(startCol, row)
		x2, _ := getCoordinates(endCol, row)

		// resize and reposition each rectangle
		e.selection[i].Resize(fyne.NewSize(x2-x1+1, lineHeight))
		e.selection[i].Move(fyne.NewPos(x1-1, y1))
	}
}

func (e *entryRenderer) moveCursor() {
	textRenderer := Renderer(e.text).(*textRenderer)
	e.entry.RLock()
	size := textRenderer.lineSizeToColumn(e.entry.CursorColumn, e.entry.CursorRow)
	xPos := size.Width
	yPos := size.Height * e.entry.CursorRow
	e.entry.RUnlock()

	// build e.selection[] if the user has made a selection
	e.buildSelection()

	lineHeight := e.text.charMinSize().Height
	e.cursor.Resize(fyne.NewSize(2, lineHeight))
	e.cursor.Move(fyne.NewPos(xPos-1+theme.Padding()*2, yPos+theme.Padding()*2))

	if e.entry.OnCursorChanged != nil {
		e.entry.OnCursorChanged()
	}
	canvas.Refresh(e.cursor)
}

// Layout the components of the entry widget.
func (e *entryRenderer) Layout(size fyne.Size) {
	e.line.Resize(fyne.NewSize(size.Width, theme.Padding()))
	e.line.Move(fyne.NewPos(0, size.Height-theme.Padding()))

	e.text.Resize(size.Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
	e.text.Move(fyne.NewPos(theme.Padding(), theme.Padding()))

	e.placeholder.Resize(size.Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
	e.placeholder.Move(fyne.NewPos(theme.Padding(), theme.Padding()))

	e.moveCursor()
}

// ApplyTheme is called when the Entry may need to update its look.
func (e *entryRenderer) ApplyTheme() {
	Renderer(e.text).ApplyTheme()
	if e.entry.focused {
		e.line.FillColor = theme.FocusColor()
	} else {
		e.line.FillColor = theme.ButtonColor()
	}

	e.cursor.FillColor = theme.FocusColor()
	for _, selection := range e.selection {
		selection.(*canvas.Rectangle).FillColor = theme.FocusColor()
	}

	e.Refresh()
}

func (e *entryRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (e *entryRenderer) Refresh() {
	if e.text.len() == 0 && e.entry.Visible() {
		e.placeholder.Show()
	} else if e.placeholder.Visible() {
		e.placeholder.Hide()
	}

	if e.entry.focused {
		e.cursor.Show()
		e.line.FillColor = theme.FocusColor()
	} else {
		e.cursor.Hide()
		e.line.FillColor = theme.ButtonColor()
	}

	for _, selection := range e.selection {
		selection.(*canvas.Rectangle).Hidden = !e.entry.focused
	}

	canvas.Refresh(e.entry)
}

func (e *entryRenderer) Objects() []fyne.CanvasObject {
	// Objects are generated dynamically force selection rectangles to appear underneath the text
	if e.entry.selecting {
		return append(e.selection, e.objects...)
	}
	return e.objects
}

func (e *entryRenderer) Destroy() {
	if e.entry.popUp != nil {
		c := fyne.CurrentApp().Driver().CanvasForObject(e.entry)
		c.SetOverlay(nil)
		Renderer(e.entry.popUp).Destroy()
		e.entry.popUp = nil
	}
}

// Declare conformity with interfaces
var _ fyne.Draggable = (*Entry)(nil)
var _ fyne.Tappable = (*Entry)(nil)
var _ desktop.Mouseable = (*Entry)(nil)
var _ desktop.Keyable = (*Entry)(nil)

// Entry widget allows simple text to be input when focused.
type Entry struct {
	baseWidget
	sync.RWMutex
	shortcut    fyne.ShortcutHandler
	Text        string
	PlaceHolder string
	OnChanged   func(string) `json:"-"`
	Password    bool
	ReadOnly    bool
	MultiLine   bool

	CursorRow, CursorColumn int
	OnCursorChanged         func() `json:"-"`

	focused bool

	// selectRow and selectColumn represent the selection start location
	// The selection will span from selectRow/Column to CursorRow/Column -- note that the cursor
	// position may occur before or after the select start position in the text.
	selectRow, selectColumn int

	// selectKeyDown indicates whether left shift or right shift is currently held down
	selectKeyDown bool

	// selecting indicates whether the cursor has moved since it was at the selection start location
	selecting bool
	popUp     *PopUp
	// TODO: Add OnSelectChanged
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (e *Entry) Resize(size fyne.Size) {
	e.resize(size, e)
}

// Move the widget to a new position, relative to its parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (e *Entry) Move(pos fyne.Position) {
	e.move(pos, e)
}

// MinSize returns the smallest size this widget can shrink to
func (e *Entry) MinSize() fyne.Size {
	return e.minSize(e)
}

// Show this widget, if it was previously hidden
func (e *Entry) Show() {
	e.show(e)
	if len(e.Text) != 0 {
		e.placeholderProvider().Hide()
	}
	if !e.Focused() {
		Renderer(e).(*entryRenderer).cursor.Hide()
	}
}

// Hide this widget, if it was previously visible
func (e *Entry) Hide() {
	if e.focused {
		fyne.CurrentApp().Driver().CanvasForObject(e).Focus(nil)
	}
	if e.popUp != nil {
		e.popUp.Hide()
	}
	e.hide(e)
}

// SetText manually sets the text of the Entry to the given text value.
func (e *Entry) SetText(text string) {
	e.textProvider().SetText(text)
	e.updateText(text)

	if text == "" {
		e.Lock()
		e.CursorColumn = 0
		e.CursorRow = 0
		e.Unlock()
		Renderer(e).(*entryRenderer).moveCursor()
	} else {
		provider := e.textProvider()
		if e.CursorRow >= provider.rows() {
			e.CursorRow = provider.rows() - 1
		}
		rowLength := provider.rowLength(e.CursorRow)
		if e.CursorColumn >= rowLength {
			e.CursorColumn = rowLength
		}
	}
}

// SetPlaceHolder sets the text that will be displayed if the entry is otherwise empty
func (e *Entry) SetPlaceHolder(text string) {
	e.PlaceHolder = text
	e.placeholderProvider().SetText(text) // refreshes
}

// SetReadOnly sets whether or not the Entry should not be editable
func (e *Entry) SetReadOnly(ro bool) {
	e.ReadOnly = ro

	Refresh(e)
}

// updateText updates the internal text to the given value
func (e *Entry) updateText(text string) {
	changed := e.Text != text
	e.Lock()
	e.Text = text
	e.Unlock()
	if changed && e.OnChanged != nil {
		e.OnChanged(text)
	}

	Refresh(e)
}

// selection returns the start and end text positions for the selected span of text
// Note: this functionality depends on the relationship between the selection start row/col and
// the current cursor row/column.
// eg: (whitespace for clarity, '_' denotes cursor)
//   "T  e  s [t  i]_n  g" == 3, 5
//   "T  e  s_[t  i] n  g" == 3, 5
//   "T  e_[s  t  i] n  g" == 2, 5
func (e *Entry) selection() (int, int) {
	e.RLock()
	defer e.RUnlock()

	if e.selecting == false {
		return -1, -1
	}
	if e.CursorRow == e.selectRow && e.CursorColumn == e.selectColumn {
		return -1, -1
	}

	// Find the selection start
	rowA, colA := e.CursorRow, e.CursorColumn
	rowB, colB := e.selectRow, e.selectColumn
	// Reposition if the cursors row is more than select start row, or if the row is the same and
	// the cursors col is more that the select start column
	if rowA > e.selectRow || (rowA == e.selectRow && colA > e.selectColumn) {
		rowA, colA = e.selectRow, e.selectColumn
		rowB, colB = e.CursorRow, e.CursorColumn
	}

	return e.textPosFromRowCol(rowA, colA), e.textPosFromRowCol(rowB, colB)
}

// selectedText returns the text currently selected in this Entry.
// If there is no selection it will return the empty string.
func (e *Entry) selectedText() string {
	if e.selecting == false {
		return ""
	}

	start, stop := e.selection()
	return string(e.textProvider().buffer[start:stop])
}

// Obtains row,col from a given textual position
// expects a read or write lock to be held by the caller
func (e *Entry) rowColFromTextPos(pos int) (int, int) {
	provider := e.textProvider()
	for i := 0; i < provider.rows(); i++ {
		rowLength := provider.rowLength(i)
		if rowLength+1 > pos {
			return i, pos
		}
		pos -= rowLength + 1 // +1 for newline
	}
	return 0, 0
}

// Obtains textual position from a given row and col
// expects a read or write lock to be held by the caller
func (e *Entry) textPosFromRowCol(row, col int) int {
	pos := 0
	provider := e.textProvider()
	for i := 0; i < row; i++ {
		rowLength := provider.rowLength(i)
		pos += rowLength + 1
	}
	pos += col
	return pos
}

func (e *Entry) cursorTextPos() int {
	pos := 0
	e.RLock()
	provider := e.textProvider()
	for i := 0; i < e.CursorRow; i++ {
		rowLength := provider.rowLength(i)
		pos += rowLength + 1
	}
	pos += e.CursorColumn
	e.RUnlock()
	return pos
}

// FocusGained is called when the Entry has been given focus.
func (e *Entry) FocusGained() {
	if e.ReadOnly {
		return
	}
	e.focused = true

	Refresh(e)
}

// FocusLost is called when the Entry has had focus removed.
func (e *Entry) FocusLost() {
	e.focused = false

	Refresh(e)
}

// Focused returns whether or not this Entry has focus.
func (e *Entry) Focused() bool {
	return e.focused
}

func (e *Entry) cursorColAt(text []rune, pos fyne.Position) int {
	for i := 0; i < len(text); i++ {
		str := string(text[0 : i+1])
		wid := textMinSize(str, theme.TextSize(), e.textStyle()).Width + theme.Padding()
		if wid > pos.X {
			return i
		}
	}
	return len(text)
}

// Tapped is called when this entry has been tapped so we should update the cursor position.
func (e *Entry) Tapped(ev *fyne.PointEvent) {
	e.updateMousePointer(ev, false)
}

// pasteFromClipboard inserts text from the clipboard content,
// starting from the cursor position.
func (e *Entry) pasteFromClipboard(clipboard fyne.Clipboard) {
	if e.selecting {
		e.eraseSelection()
	}
	text := clipboard.Content()
	if !e.MultiLine {
		// format clipboard content to be compatible with single line entry
		text = strings.Replace(text, "\n", " ", -1)
	}
	provider := e.textProvider()
	runes := []rune(text)
	provider.insertAt(e.cursorTextPos(), runes)

	newlines := strings.Count(text, "\n")
	if newlines == 0 {
		e.CursorColumn += len(runes)
	} else {
		e.CursorRow += newlines
		lastNewline := strings.LastIndex(text, "\n")
		e.CursorColumn = len(runes) - lastNewline - 1
	}
	e.updateText(provider.String())
	Renderer(e).(*entryRenderer).moveCursor()
}

// selectAll selects all text in entry
func (e *Entry) selectAll() {
	e.Lock()
	e.selectRow = 0
	e.selectColumn = 0

	lastRow := e.textProvider().rows() - 1
	e.CursorColumn = e.textProvider().rowLength(lastRow)
	e.CursorRow = lastRow
	e.selecting = true
	e.Unlock()

	Renderer(e).(*entryRenderer).moveCursor()
}

// TappedSecondary is called when right or alternative tap is invoked.
//
// Opens the PopUpMenu with `Paste` item to paste text from the clipboard.
func (e *Entry) TappedSecondary(pe *fyne.PointEvent) {
	if e.ReadOnly {
		return
	}

	c := fyne.CurrentApp().Driver().CanvasForObject(e)

	pasteItem := fyne.NewMenuItem("Paste", func() {
		clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
		e.pasteFromClipboard(clipboard)
	})
	selectAllItem := fyne.NewMenuItem("Select all", e.selectAll)
	e.popUp = NewPopUpMenu(fyne.NewMenu("", pasteItem, selectAllItem), c)

	entryPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(e)
	popUpPos := entryPos.Add(fyne.NewPos(pe.Position.X, pe.Position.Y))

	e.popUp.Move(popUpPos)
}

// MouseDown called on mouse click, this triggers a mouse click which can move the cursor,
// update the existing selection (if shift is held), or start a selection dragging operation.
func (e *Entry) MouseDown(m *desktop.MouseEvent) {
	if e.selectKeyDown {
		e.selecting = true
	}
	if e.selecting && e.selectKeyDown == false {
		e.selecting = false
	}
	e.updateMousePointer(&m.PointEvent, e.selecting == false)
}

// MouseUp called on mouse release
// If a mouse drag event has completed then check to see if it has resulted in an empty selection,
// if so, and if a text select key isn't held, then disable selecting
func (e *Entry) MouseUp(m *desktop.MouseEvent) {
	start, _ := e.selection()
	if start == -1 && e.selecting && e.selectKeyDown == false {
		e.selecting = false
	}
}

// Dragged is called when the pointer moves while a button is held down
func (e *Entry) Dragged(d *fyne.DragEvent) {
	e.selecting = true
	e.updateMousePointer(&d.PointEvent, false)
}

// DragEnd is called at end of a drag event - currently ignored
func (e *Entry) DragEnd() {
}

func (e *Entry) updateMousePointer(ev *fyne.PointEvent, startSelect bool) {

	if !e.focused {
		e.FocusGained()
	}

	rowHeight := e.textProvider().charMinSize().Height
	row := int(math.Floor(float64(ev.Position.Y-theme.Padding()) / float64(rowHeight)))
	col := 0
	if row < 0 {
		row = 0
	} else if row >= e.textProvider().rows() {
		row = e.textProvider().rows() - 1
		col = 0
	} else {
		col = e.cursorColAt(e.textProvider().row(row), ev.Position)
	}

	e.Lock()
	e.CursorRow = row
	e.CursorColumn = col
	if startSelect {
		e.selectRow = row
		e.selectColumn = col
	}
	e.Unlock()
	Renderer(e).(*entryRenderer).moveCursor()
}

// getTextWhitespaceRegion returns the start/end markers for selection highlight on starting from col
// and expanding to the start and end of the whitespace or text underneat the specified position.
func getTextWhitespaceRegion(row []rune, col int) (int, int) {

	if len(row) == 0 || col < 0 {
		return -1, -1
	}

	// If the click position exceeds the length of text then snap it to the end
	if col >= len(row) {
		col = len(row) - 1
	}

	// maps: " fi-sh 日本語本語日  \t "
	// into: " -- -- ------   "
	space := func(r rune) rune {
		if unicode.IsSpace(r) {
			return ' '
		}
		// If this rune is a typical word separator then classify it as whitespace
		if strings.ContainsRune(doubleClickWordSeperator, r) {
			return ' '
		}
		return '-'
	}
	toks := strings.Map(space, string(row))

	c := byte(' ')
	if toks[col] == ' ' {
		c = byte('-')
	}

	// LastIndexByte + 1 ensures that the position of the unwanted character 'c' is excluded
	// +1 also has the added side effect whereby if 'c' isn't found then -1 is snapped to 0
	start := strings.LastIndexByte(toks[:col], c) + 1

	// IndexByte will find the position of the next unwanted character, this is to be the end
	// marker for the selection
	end := strings.IndexByte(toks[col:], c)

	if end == -1 {
		end = len(toks) // snap end to len(toks) if it results in -1
	} else {
		end += col // otherwise include the text slice position
	}
	return start, end
}

// DoubleTapped is called when this entry has been double tapped so we should select text below the pointer
func (e *Entry) DoubleTapped(ev *fyne.PointEvent) {
	row := e.textProvider().row(e.CursorRow)

	start, end := getTextWhitespaceRegion(row, e.CursorColumn)
	if start == -1 || end == -1 {
		return
	}

	e.Lock()
	if e.selectKeyDown == false {
		e.selectRow = e.CursorRow
		e.selectColumn = start
	}
	// Always aim to maximise the selected region
	if e.selectRow > e.CursorRow || (e.selectRow == e.CursorRow && e.selectColumn > e.CursorColumn) {
		e.CursorColumn = start
	} else {
		e.CursorColumn = end
	}
	e.selecting = true
	e.Unlock()
	Renderer(e).(*entryRenderer).moveCursor()
}

// TypedRune receives text input events when the Entry widget is focused.
func (e *Entry) TypedRune(r rune) {
	if e.ReadOnly {
		return
	}

	if e.popUp != nil {
		e.popUp.Hide()
	}

	provider := e.textProvider()

	// if we've typed a character and we're selecting then replace the selection with the character
	if e.selecting {
		e.eraseSelection()
		e.selecting = false
	}

	runes := []rune{r}
	provider.insertAt(e.cursorTextPos(), runes)
	e.Lock()
	e.CursorColumn += len(runes)
	e.Unlock()
	e.updateText(provider.String())
	Renderer(e).(*entryRenderer).moveCursor()
}

// KeyDown handler for keypress events - used to store shift modifier state for text selection
func (e *Entry) KeyDown(key *fyne.KeyEvent) {
	// For keyboard cursor controlled selection we now need to store shift key state and selection "start"
	// Note: selection start is where the highlight started (if the user moves the selection up or left then
	// the selectRow/Column will not match SelectionStart)
	if key.Name == desktop.KeyShiftLeft || key.Name == desktop.KeyShiftRight {
		if e.selecting == false {
			e.selectRow = e.CursorRow
			e.selectColumn = e.CursorColumn
		}
		e.selectKeyDown = true
	}
}

// KeyUp handler for key release events - used to reset shift modifier state for text selection
func (e *Entry) KeyUp(key *fyne.KeyEvent) {
	// Handle shift release for keyboard selection
	// Note: if shift is released then the user may repress it without moving to adjust their old selection
	if key.Name == desktop.KeyShiftLeft || key.Name == desktop.KeyShiftRight {
		e.selectKeyDown = false
	}
}

// eraseSelection removes the current selected region and moves the cursor
func (e *Entry) eraseSelection() {
	if e.ReadOnly {
		return
	}

	provider := e.textProvider()
	posA, posB := e.selection()

	if posA == posB {
		return
	}

	e.Lock()

	provider.deleteFromTo(posA, posB)
	e.CursorRow, e.CursorColumn = e.rowColFromTextPos(posA)
	e.selectRow, e.selectColumn = e.CursorRow, e.CursorColumn
	e.Unlock()
	e.updateText(provider.String())
	e.selecting = false
}

// selectingKeyHandler performs keypress action in the scenario that a selection
// is either a) in progress or b) about to start
// returns true if the keypress has been fully handled
func (e *Entry) selectingKeyHandler(key *fyne.KeyEvent) bool {

	if e.selectKeyDown && e.selecting == false {
		switch key.Name {
		case fyne.KeyUp, fyne.KeyDown,
			fyne.KeyLeft, fyne.KeyRight,
			fyne.KeyEnd, fyne.KeyHome,
			fyne.KeyPageUp, fyne.KeyPageDown:
			e.selecting = true
		}
	}

	if e.selecting == false {
		return false
	}

	switch key.Name {
	case fyne.KeyBackspace, fyne.KeyDelete:
		// clears the selection -- return handled
		e.eraseSelection()
		return true
	case fyne.KeyReturn, fyne.KeyEnter:
		// clear the selection -- return unhandled to add the newline
		e.eraseSelection()
		return false
	}

	if e.selectKeyDown == false {
		switch key.Name {
		case fyne.KeyLeft:
			// seek to the start of the selection -- return handled
			selectStart, _ := e.selection()
			e.Lock()
			e.CursorRow, e.CursorColumn = e.rowColFromTextPos(selectStart)
			e.Unlock()
			e.selecting = false
			return true
		case fyne.KeyRight:
			// seek to the end of the selection -- return handled
			_, selectEnd := e.selection()
			e.Lock()
			e.CursorRow, e.CursorColumn = e.rowColFromTextPos(selectEnd)
			e.Unlock()
			e.selecting = false
			return true
		case fyne.KeyUp, fyne.KeyDown, fyne.KeyEnd, fyne.KeyHome:
			// cursor movement without left or right shift -- clear selection and return unhandled
			e.selecting = false
			return false
		}
	}

	return false
}

// TypedKey receives key input events when the Entry widget is focused.
func (e *Entry) TypedKey(key *fyne.KeyEvent) {
	if e.ReadOnly {
		return
	}

	provider := e.textProvider()

	if e.selectKeyDown || e.selecting {
		if e.selectingKeyHandler(key) {
			e.updateText(provider.String())
			Renderer(e).(*entryRenderer).moveCursor()
			return
		}
	}

	switch key.Name {
	case fyne.KeyBackspace:
		e.RLock()
		isEmpty := provider.len() == 0 || (e.CursorColumn == 0 && e.CursorRow == 0)
		e.RUnlock()
		if isEmpty {
			return
		}
		pos := e.cursorTextPos()
		e.Lock()
		deleted := provider.deleteFromTo(pos-1, pos)
		if deleted[0] == '\n' {
			e.CursorRow--
			rowLength := provider.rowLength(e.CursorRow)
			e.CursorColumn = rowLength
		} else {
			e.CursorColumn--
		}
		e.Unlock()
	case fyne.KeyDelete:
		pos := e.cursorTextPos()
		if provider.len() == 0 || pos == provider.len() {
			return
		}
		provider.deleteFromTo(pos, pos+1)
	case fyne.KeyReturn, fyne.KeyEnter:
		if !e.MultiLine {
			return
		}
		provider.insertAt(e.cursorTextPos(), []rune("\n"))
		e.Lock()
		e.CursorColumn = 0
		e.CursorRow++
		e.Unlock()
	case fyne.KeyUp:
		if !e.MultiLine {
			return
		}

		e.Lock()
		if e.CursorRow > 0 {
			e.CursorRow--
		}

		rowLength := provider.rowLength(e.CursorRow)
		if e.CursorColumn > rowLength {
			e.CursorColumn = rowLength
		}
		e.Unlock()
	case fyne.KeyDown:
		if !e.MultiLine {
			return
		}

		e.Lock()
		if e.CursorRow < provider.rows()-1 {
			e.CursorRow++
		}

		rowLength := provider.rowLength(e.CursorRow)
		if e.CursorColumn > rowLength {
			e.CursorColumn = rowLength
		}
		e.Unlock()
	case fyne.KeyLeft:
		e.Lock()
		if e.CursorColumn > 0 {
			e.CursorColumn--
		} else if e.MultiLine && e.CursorRow > 0 {
			e.CursorRow--
			e.CursorColumn = provider.rowLength(e.CursorRow)
		}
		e.Unlock()
	case fyne.KeyRight:
		e.Lock()
		if e.MultiLine {
			rowLength := provider.rowLength(e.CursorRow)
			if e.CursorColumn < rowLength {
				e.CursorColumn++
			} else if e.CursorRow < provider.rows()-1 {
				e.CursorRow++
				e.CursorColumn = 0
			}
		} else if e.CursorColumn < provider.len() {
			e.CursorColumn++
		}
		e.Unlock()
	case fyne.KeyEnd:
		e.Lock()
		if e.MultiLine {
			e.CursorColumn = provider.rowLength(e.CursorRow)
		} else {
			e.CursorColumn = provider.len()
		}
		e.Unlock()
	case fyne.KeyHome:
		e.Lock()
		e.CursorColumn = 0
		e.Unlock()
	default:
		return
	}

	e.updateText(provider.String())
	Renderer(e).(*entryRenderer).moveCursor()
}

// TypedShortcut implements the Shortcutable interface
func (e *Entry) TypedShortcut(shortcut fyne.Shortcut) bool {
	return e.shortcut.TypedShortcut(shortcut)
}

// textProvider returns the text handler for this entry
func (e *Entry) textProvider() *textProvider {
	return Renderer(e).(*entryRenderer).text
}

// placeholderProvider returns the placeholder text handler for this entry
func (e *Entry) placeholderProvider() *textProvider {
	return Renderer(e).(*entryRenderer).placeholder
}

// textAlign tells the rendering textProvider our alignment
func (e *Entry) textAlign() fyne.TextAlign {
	return fyne.TextAlignLeading
}

// textStyle tells the rendering textProvider our style
func (e *Entry) textStyle() fyne.TextStyle {
	return fyne.TextStyle{}
}

// textColor tells the rendering textProvider our color
func (e *Entry) textColor() color.Color {
	return theme.TextColor()
}

// password tells the rendering textProvider if we are a password field
func (e *Entry) password() bool {
	return e.Password
}

// object returns the root object of the widget so it can be referenced
func (e *Entry) object() fyne.Widget {
	return nil
}

type placeholderPresenter struct {
	e *Entry
}

// textAlign tells the rendering textProvider our alignment
func (p *placeholderPresenter) textAlign() fyne.TextAlign {
	return fyne.TextAlignLeading
}

// textStyle tells the rendering textProvider our style
func (p *placeholderPresenter) textStyle() fyne.TextStyle {
	return fyne.TextStyle{}
}

// textColor tells the rendering textProvider our color
func (p *placeholderPresenter) textColor() color.Color {
	return theme.PlaceHolderColor()
}

// password tells the rendering textProvider if we are a password field
// placeholder text is not obfuscated, returning false
func (p *placeholderPresenter) password() bool {
	return false
}

// object returns the root object of the widget so it can be referenced
func (p *placeholderPresenter) object() fyne.Widget {
	return nil
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (e *Entry) CreateRenderer() fyne.WidgetRenderer {
	text := newTextProvider(e.Text, e)
	placeholder := newTextProvider(e.PlaceHolder, &placeholderPresenter{e})

	line := canvas.NewRectangle(theme.ButtonColor())
	cursor := canvas.NewRectangle(theme.FocusColor())

	return &entryRenderer{&text, &placeholder, line, cursor, []fyne.CanvasObject{},
		[]fyne.CanvasObject{line, &placeholder, &text, cursor}, e}
}

func (e *Entry) registerShortcut() {
	e.shortcut.AddShortcut(&fyne.ShortcutCut{}, func(se fyne.Shortcut) {
		if e.Password {
			return
		}

		cut := se.(*fyne.ShortcutCut)
		text := e.selectedText()
		e.eraseSelection()

		cut.Clipboard.SetContent(text)
	})
	e.shortcut.AddShortcut(&fyne.ShortcutCopy{}, func(se fyne.Shortcut) {
		if e.Password {
			return
		}

		copy := se.(*fyne.ShortcutCopy)
		text := e.selectedText()

		copy.Clipboard.SetContent(text)
	})
	e.shortcut.AddShortcut(&fyne.ShortcutPaste{}, func(se fyne.Shortcut) {
		paste := se.(*fyne.ShortcutPaste)
		e.pasteFromClipboard(paste.Clipboard)
	})
	e.shortcut.AddShortcut(&fyne.ShortcutSelectAll{}, func(se fyne.Shortcut) {
		e.selectAll()
	})
}

// NewEntry creates a new single line entry widget.
func NewEntry() *Entry {
	e := &Entry{}
	e.registerShortcut()
	Refresh(e)
	return e
}

// NewMultiLineEntry creates a new entry that allows multiple lines
func NewMultiLineEntry() *Entry {
	e := &Entry{MultiLine: true}
	e.registerShortcut()
	Refresh(e)
	return e
}

// NewPasswordEntry creates a new entry password widget
func NewPasswordEntry() *Entry {
	e := &Entry{Password: true}
	e.registerShortcut()
	Refresh(e)
	return e
}

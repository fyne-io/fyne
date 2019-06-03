package widget

import (
	"image/color"
	"math"
	"strings"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

const (
	multiLineRows = 3
)

type entryRenderer struct {
	text         *textProvider
	placeholder  *textProvider
	line, cursor *canvas.Rectangle
	selection    []*canvas.Rectangle

	entry *Entry
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
	curRow := e.entry.CursorRow
	curCol := e.entry.CursorColumn
	selRow := -1
	selCol := -1
	if e.entry.selecting {
		selRow = e.entry.selectRow
		selCol = e.entry.selectColumn
	}
	e.entry.RUnlock()

	textRenderer := Renderer(e.text).(*textRenderer)

	// Convert column, row into x,y
	getCoordinates := func(column int, row int) (int, int) {
		sz := textRenderer.lineSize(column, row)
		return sz.Width + theme.Padding()*2, sz.Height*row + theme.Padding()*2
	}

	lineHeight := e.text.charMinSize().Height

	// if we have a selection then we should calculate the set of boxes and add them to e.selection
	if selRow != -1 {

		minmax := func(a, b int) (int, int) {
			if a < b {
				return a, b
			}
			return b, a
		}

		ssRow, seRow := minmax(selRow, curRow)
		ssCol, seCol := minmax(selCol, curCol)
		if selRow < curRow {
			ssCol, seCol = selCol, curCol
		}
		if selRow > curRow {
			ssCol, seCol = curCol, selCol
		}
		rows := seRow - ssRow + 1

		if len(e.selection) > rows {
			e.selection = e.selection[:rows]
		}

		for i := 0; i < rows; i++ {
			if len(e.selection) <= i {
				box := canvas.NewRectangle(theme.ButtonColor())
				e.selection = append(e.selection, box)
			}

			row := ssRow + i
			ss, se := ssCol, seCol
			if ssRow < row {
				ss = 0
			}
			if seRow > row {
				se = textRenderer.provider.rowLength(row)
			}
			x1, y1 := getCoordinates(ss, row)
			x2, _ := getCoordinates(se, row)

			e.selection[i].Resize(fyne.NewSize(x2-x1+1, lineHeight))
			e.selection[i].Move(fyne.NewPos(x1-1, y1))
			e.selection[i].Show()
		}
	} else {
		e.selection = e.selection[:0]
	}
}

func (e *entryRenderer) moveCursor() {
	textRenderer := Renderer(e.text).(*textRenderer)
	e.entry.RLock()
	size := textRenderer.lineSize(e.entry.CursorColumn, e.entry.CursorRow)
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

	canvas.Refresh(e.entry)
}

func (e *entryRenderer) Objects() []fyne.CanvasObject {
	// Objects are generated dynamically force selection rectangles to appear underneath the text
	objs := []fyne.CanvasObject{}
	for _, o := range e.selection {
		objs = append(objs, o)
	}
	objs = append(objs, e.line, e.placeholder, e.text, e.cursor)
	return objs
}

func (e *entryRenderer) Destroy() {
}

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

	// Selection start/end are private as they represent where the selection started, which
	// could be confused with the "start of the selection"
	// The API user should use SelectionStart and SelectionEnd
	selectRow, selectColumn int
	selectKeyDown           bool
	selecting               bool
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
}

// Hide this widget, if it was previously visible
func (e *Entry) Hide() {
	if e.focused {
		fyne.CurrentApp().Driver().CanvasForObject(e).Focus(nil)
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

// SelectionStart returns the first text position of the span of text covered by the selection highlight
// eg: "T e s[t i]n g" == 3 (whitespace for clarity)
func (e *Entry) SelectionStart() int {
	e.RLock()
	defer e.RUnlock()

	if e.selecting == false {
		return -1
	}

	row := e.CursorRow
	col := e.CursorColumn
	if row == e.selectRow && col == e.selectColumn {
		return -1
	}

	if row > e.selectRow {
		row = e.selectRow
		col = e.selectColumn
	} else if row == e.selectRow && col > e.selectColumn {
		row = e.selectRow
		col = e.selectColumn
	}
	return e.textPosFromRowCol(row, col)
}

// SelectionEnd returns the last text position of the span of text covered by the selection highlight
// eg: "T e s[t i]n g" == 5 (whitespace for clarity)
func (e *Entry) SelectionEnd() int {
	e.RLock()
	defer e.RUnlock()

	if e.selecting == false {
		return -1
	}

	row := e.CursorRow
	col := e.CursorColumn
	if row == e.selectRow && col == e.selectColumn {
		return -1
	}

	if row < e.selectRow {
		row = e.selectRow
		col = e.selectColumn
	} else if row == e.selectRow && col < e.selectColumn {
		row = e.selectRow
		col = e.selectColumn
	}

	return e.textPosFromRowCol(row, col)
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
		wid := textMinSize(str, theme.TextSize(), e.textStyle()).Width
		if wid > pos.X {
			return i
		}
	}
	return len(text)
}

// Tapped is called when this entry has been tapped so we should update the cursor position.
func (e *Entry) Tapped(ev *fyne.PointEvent) {
	if !e.focused {
		e.FocusGained()
	}
	if e.selecting && e.selectKeyDown == false {
		e.selecting = false
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
	e.Unlock()
	Renderer(e).(*entryRenderer).moveCursor()
}

// TappedSecondary is called when right or alternative tap is invoked - this is currently ignored.
func (e *Entry) TappedSecondary(_ *fyne.PointEvent) {
}

// TypedRune receives text input events when the Entry widget is focused.
func (e *Entry) TypedRune(r rune) {
	if e.ReadOnly {
		return
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
		e.selecting = true
		e.selectKeyDown = true
	}
}

// KeyUp handler for key release events - used to reset shift modifier state for text selection
func (e *Entry) KeyUp(key *fyne.KeyEvent) {
	// Handle shift release for keyboard selection
	// Note: if shift is released then the user may repress it without moving to adjust their old selection
	if key.Name == desktop.KeyShiftLeft || key.Name == desktop.KeyShiftRight {
		e.selectKeyDown = false
	} else {
		if e.selectKeyDown == false {
			e.selecting = false
		}
	}
}

// eraseSelection removes the current selected region and moves the cursor
func (e *Entry) eraseSelection() {
	if e.ReadOnly {
		return
	}

	provider := e.textProvider()
	posA := e.SelectionStart()
	posB := e.SelectionEnd()
	e.Lock()

	provider.deleteFromTo(posA, posB)
	e.CursorRow, e.CursorColumn = e.rowColFromTextPos(posA)
	e.selectRow, e.selectColumn = e.CursorRow, e.CursorColumn
	e.Unlock()

	// e.selecting = false
}

// TypedKey receives key input events when the Entry widget is focused.
func (e *Entry) TypedKey(key *fyne.KeyEvent) {
	if e.ReadOnly {
		return
	}

	provider := e.textProvider()

	// seeks to the start/end of the selection - used by: up, down, left, right
	seekSelection := func(start bool) {
		var pos int
		if start {
			pos = e.SelectionStart()
		} else {
			pos = e.SelectionEnd()
		}
		e.Lock()
		e.CursorRow, e.CursorColumn = e.rowColFromTextPos(pos)
		e.Unlock()
		e.selecting = false
	}

	switch key.Name {
	case fyne.KeyBackspace:
		if e.selecting {
			e.eraseSelection() // clears the current selection (exactly like delete)
		} else {
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
		}
	case fyne.KeyDelete:
		if e.selecting {
			e.eraseSelection() // clears the selection (exactly like backspace)
		} else {
			pos := e.cursorTextPos()
			if provider.len() == 0 || pos == provider.len() {
				return
			}
			provider.deleteFromTo(pos, pos+1)
		}
	case fyne.KeyReturn, fyne.KeyEnter:
		if !e.MultiLine {
			return
		}
		if e.selecting {
			e.eraseSelection() // clear the selection and fallthrough to add the newline
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

		if e.selecting && e.selectKeyDown == false {
			seekSelection(true) // seek to the start of the selection and fallthrough to move up
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
		if e.selecting && e.selectKeyDown == false {
			seekSelection(false) // seek to the end of the selection and fallthrough to move down
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
		if e.selecting && e.selectKeyDown == false {
			seekSelection(true) // seek to the start of the selection (no fallthrough)
		} else {
			e.Lock()
			if e.CursorColumn > 0 {
				e.CursorColumn--
			} else if e.MultiLine && e.CursorRow > 0 {
				e.CursorRow--
				e.CursorColumn = provider.rowLength(e.CursorRow)
			}
			e.Unlock()
		}
	case fyne.KeyRight:
		if e.selecting && e.selectKeyDown == false {
			seekSelection(false) // seek to the end of the selection (no fallthrough)
		} else {
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
		}
	case fyne.KeyEnd:
		e.Lock()
		// if the user pressed end and isn't holding shift then end selection (fallthrough)
		if e.selecting && e.selectKeyDown == false {
			e.selecting = false
		}
		if e.MultiLine {
			e.CursorColumn = provider.rowLength(e.CursorRow)
		} else {
			e.CursorColumn = provider.len()
		}
		e.Unlock()
	case fyne.KeyHome:
		e.Lock()
		// if the user pressed home and isn't holding shift then end selection (fallthrough)
		if e.selecting && e.selectKeyDown == false {
			e.selecting = false
		}
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

	return &entryRenderer{&text, &placeholder, line, cursor, nil, e}
}

func (e *Entry) registerShortcut() {
	scPaste := &fyne.ShortcutPaste{}
	e.shortcut.AddShortcut(scPaste, func(se fyne.Shortcut) {
		scPaste = se.(*fyne.ShortcutPaste)
		text := scPaste.Clipboard.Content()
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

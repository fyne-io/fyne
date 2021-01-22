package widget

import (
	"image/color"
	"math"
	"strings"
	"time"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

const (
	multiLineRows            = 3
	doubleClickWordSeperator = "`~!@#$%^&*()-=+[{]}\\|;:'\",.<>/?"
)

// Declare conformity with interfaces
var _ fyne.Disableable = (*Entry)(nil)
var _ fyne.Draggable = (*Entry)(nil)
var _ fyne.Focusable = (*Entry)(nil)
var _ fyne.Tappable = (*Entry)(nil)
var _ fyne.Widget = (*Entry)(nil)
var _ desktop.Mouseable = (*Entry)(nil)
var _ desktop.Keyable = (*Entry)(nil)
var _ mobile.Keyboardable = (*Entry)(nil)

// Entry widget allows simple text to be input when focused.
type Entry struct {
	DisableableWidget
	shortcut fyne.ShortcutHandler
	Text     string
	// Since: 2.0
	TextStyle   fyne.TextStyle
	PlaceHolder string
	OnChanged   func(string) `json:"-"`
	// Since: 2.0
	OnSubmitted func(string) `json:"-"`
	Password    bool
	MultiLine   bool
	Wrapping    fyne.TextWrap

	// Set a validator that this entry will check against
	// Since: 1.4
	Validator           fyne.StringValidator
	validationStatus    *validationStatus
	onValidationChanged func(error)
	validationError     error

	CursorRow, CursorColumn int
	OnCursorChanged         func() `json:"-"`

	focused     bool
	text        *textProvider
	placeholder *textProvider
	content     *entryContent
	scroll      *widget.Scroll

	// selectRow and selectColumn represent the selection start location
	// The selection will span from selectRow/Column to CursorRow/Column -- note that the cursor
	// position may occur before or after the select start position in the text.
	selectRow, selectColumn int

	// selectKeyDown indicates whether left shift or right shift is currently held down
	selectKeyDown bool

	// selecting indicates whether the cursor has moved since it was at the selection start location
	selecting bool
	popUp     *PopUpMenu
	// TODO: Add OnSelectChanged

	// ActionItem is a small item which is displayed at the outer right of the entry (like a password revealer)
	ActionItem   fyne.CanvasObject
	textSource   binding.String
	textListener binding.DataListener
}

// NewEntry creates a new single line entry widget.
func NewEntry() *Entry {
	e := &Entry{Wrapping: fyne.TextTruncate}
	e.ExtendBaseWidget(e)
	return e
}

// NewEntryWithData returns an Entry widget connected to the specified data source.
//
// Since: 2.0
func NewEntryWithData(data binding.String) *Entry {
	entry := NewEntry()
	entry.Bind(data)

	return entry
}

// NewMultiLineEntry creates a new entry that allows multiple lines
func NewMultiLineEntry() *Entry {
	e := &Entry{MultiLine: true, Wrapping: fyne.TextTruncate}
	e.ExtendBaseWidget(e)
	return e
}

// NewPasswordEntry creates a new entry password widget
func NewPasswordEntry() *Entry {
	e := &Entry{Password: true, Wrapping: fyne.TextTruncate}
	e.ExtendBaseWidget(e)
	e.ActionItem = newPasswordRevealer(e)
	return e
}

// Bind connects the specified data source to this Entry.
// The current value will be displayed and any changes in the data will cause the widget to update.
// User interactions with this Entry will set the value into the data source.
//
// Since: 2.0
func (e *Entry) Bind(data binding.String) {
	e.Unbind()
	e.textSource = data

	var convertErr error
	e.Validator = func(string) error {
		return convertErr
	}
	e.textListener = binding.NewDataListener(func() {
		val, err := data.Get()
		if err != nil {
			convertErr = err
			e.SetValidationError(err)
			return
		}
		e.Text = val
		convertErr = nil
		e.Refresh()
		if cache.IsRendered(e) {
			e.Refresh()
		}
	})
	data.AddListener(e.textListener)

	e.OnChanged = func(s string) {
		convertErr = data.Set(s)
		e.SetValidationError(convertErr)
	}
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
//
// Implements: fyne.Widget
func (e *Entry) CreateRenderer() fyne.WidgetRenderer {
	e.ExtendBaseWidget(e)

	// initialise
	e.textProvider()
	e.placeholderProvider()

	box := canvas.NewRectangle(theme.InputBackgroundColor())
	line := canvas.NewRectangle(theme.ShadowColor())

	e.content = &entryContent{entry: e}
	e.scroll = widget.NewScroll(e.content)
	objects := []fyne.CanvasObject{box, line, e.scroll}
	e.content.scroll = e.scroll

	if e.Password && e.ActionItem == nil {
		// An entry widget has been created via struct setting manually
		// the Password field to true. Going to enable the password revealer.
		e.ActionItem = newPasswordRevealer(e)
	}

	if e.ActionItem != nil {
		objects = append(objects, e.ActionItem)
	}

	return &entryRenderer{box, line, e.scroll, objects, e}
}

// Cursor returns the cursor type of this widget
//
// Implements: desktop.Cursorable
func (e *Entry) Cursor() desktop.Cursor {
	return desktop.TextCursor
}

// Disable this widget so that it cannot be interacted with, updating any style appropriately.
//
// Implements: fyne.Disableable
func (e *Entry) Disable() {
	e.DisableableWidget.Disable()
}

// Disabled returns whether the entry is disabled or read-only.
//
// Implements: fyne.Disableable
func (e *Entry) Disabled() bool {
	return e.DisableableWidget.disabled
}

// DoubleTapped is called when this entry has been double tapped so we should select text below the pointer
//
// Implements: fyne.DoubleTappable
func (e *Entry) DoubleTapped(p *fyne.PointEvent) {
	row := e.textProvider().row(e.CursorRow)
	start, end := getTextWhitespaceRegion(row, e.CursorColumn)
	if start == -1 || end == -1 {
		return
	}

	e.setFieldsAndRefresh(func() {
		if !e.selectKeyDown {
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
	})
}

// DragEnd is called at end of a drag event. It does nothing.
//
// Implements: fyne.Draggable
func (e *Entry) DragEnd() {
}

// Dragged is called when the pointer moves while a button is held down.
// It updates the selection accordingly.
//
// Implements: fyne.Draggable
func (e *Entry) Dragged(d *fyne.DragEvent) {
	if !e.selecting {
		e.selectRow, e.selectColumn = e.getRowCol(&d.PointEvent)

		e.selecting = true
	}
	e.updateMousePointer(&d.PointEvent, false)
}

// Enable this widget, updating any style or features appropriately.
//
// Implements: fyne.Disableable
func (e *Entry) Enable() {
	e.DisableableWidget.Enable()
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
func (e *Entry) ExtendBaseWidget(wid fyne.Widget) {
	impl := e.getImpl()
	if impl != nil {
		return
	}

	e.propertyLock.Lock()
	defer e.propertyLock.Unlock()
	e.BaseWidget.impl = wid
	e.registerShortcut()
}

// FocusGained is called when the Entry has been given focus.
//
// Implements: fyne.Focusable
func (e *Entry) FocusGained() {
	if e.Disabled() {
		return
	}
	e.setFieldsAndRefresh(func() {
		e.focused = true
	})
}

// FocusLost is called when the Entry has had focus removed.
//
// Implements: fyne.Focusable
func (e *Entry) FocusLost() {
	e.setFieldsAndRefresh(func() {
		e.focused = false
	})
}

// Hide hides the entry.
//
// Implements: fyne.Widget
func (e *Entry) Hide() {
	if e.popUp != nil {
		e.popUp.Hide()
		e.popUp = nil
	}
	e.DisableableWidget.Hide()
}

// Keyboard implements the Keyboardable interface
//
// Implements: mobile.Keyboardable
func (e *Entry) Keyboard() mobile.KeyboardType {
	e.propertyLock.RLock()
	defer e.propertyLock.RUnlock()

	if e.MultiLine {
		return mobile.DefaultKeyboard
	}

	return mobile.SingleLineKeyboard
}

// KeyDown handler for keypress events - used to store shift modifier state for text selection
//
// Implements: desktop.Keyable
func (e *Entry) KeyDown(key *fyne.KeyEvent) {
	// For keyboard cursor controlled selection we now need to store shift key state and selection "start"
	// Note: selection start is where the highlight started (if the user moves the selection up or left then
	// the selectRow/Column will not match SelectionStart)
	if key.Name == desktop.KeyShiftLeft || key.Name == desktop.KeyShiftRight {
		if !e.selecting {
			e.selectRow = e.CursorRow
			e.selectColumn = e.CursorColumn
		}
		e.selectKeyDown = true
	}
}

// KeyUp handler for key release events - used to reset shift modifier state for text selection
//
// Implements: desktop.Keyable
func (e *Entry) KeyUp(key *fyne.KeyEvent) {
	// Handle shift release for keyboard selection
	// Note: if shift is released then the user may repress it without moving to adjust their old selection
	if key.Name == desktop.KeyShiftLeft || key.Name == desktop.KeyShiftRight {
		e.selectKeyDown = false
	}
}

// MinSize returns the size that this widget should not shrink below.
//
// Implements: fyne.Widget
func (e *Entry) MinSize() fyne.Size {
	e.ExtendBaseWidget(e)

	min := e.BaseWidget.MinSize()
	if e.ActionItem != nil {
		min = min.Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0))
	}
	if e.Validator != nil {
		min = min.Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0))
	}

	return min
}

// MouseDown called on mouse click, this triggers a mouse click which can move the cursor,
// update the existing selection (if shift is held), or start a selection dragging operation.
//
// Implements: desktop.Mouseable
func (e *Entry) MouseDown(m *desktop.MouseEvent) {
	e.propertyLock.Lock()
	if e.selectKeyDown {
		e.selecting = true
	}
	if e.selecting && !e.selectKeyDown && m.Button == desktop.MouseButtonPrimary {
		e.selecting = false
	}
	e.propertyLock.Unlock()

	e.updateMousePointer(&m.PointEvent, m.Button == desktop.MouseButtonSecondary)
}

// MouseUp called on mouse release
// If a mouse drag event has completed then check to see if it has resulted in an empty selection,
// if so, and if a text select key isn't held, then disable selecting
//
// Implements: desktop.Mouseable
func (e *Entry) MouseUp(m *desktop.MouseEvent) {
	start, _ := e.selection()

	e.propertyLock.Lock()
	defer e.propertyLock.Unlock()
	if start == -1 && e.selecting && !e.selectKeyDown {
		e.selecting = false
	}
}

// SelectedText returns the text currently selected in this Entry.
// If there is no selection it will return the empty string.
func (e *Entry) SelectedText() string {
	e.propertyLock.RLock()
	selecting := e.selecting
	e.propertyLock.RUnlock()
	if !selecting {
		return ""
	}

	start, stop := e.selection()
	e.propertyLock.RLock()
	defer e.propertyLock.RUnlock()
	return string(e.textProvider().buffer[start:stop])
}

// SetPlaceHolder sets the text that will be displayed if the entry is otherwise empty
func (e *Entry) SetPlaceHolder(text string) {
	e.propertyLock.Lock()
	e.PlaceHolder = text
	e.propertyLock.Unlock()

	e.placeholderProvider().setText(text) // refreshes
}

// SetText manually sets the text of the Entry to the given text value.
func (e *Entry) SetText(text string) {
	e.textProvider().setText(text)
	e.updateText(text)

	if text == "" {
		e.setFieldsAndRefresh(func() {
			e.CursorColumn = 0
			e.CursorRow = 0
		})
		return
	}
	e.propertyLock.Lock()
	defer e.propertyLock.Unlock()
	if e.CursorRow >= e.textProvider().rows() {
		e.CursorRow = e.textProvider().rows() - 1
	}
	rowLength := e.textProvider().rowLength(e.CursorRow)
	if e.CursorColumn >= rowLength {
		e.CursorColumn = rowLength
	}
}

// Tapped is called when this entry has been tapped so we should update the cursor position.
//
// Implements: fyne.Tappable
func (e *Entry) Tapped(ev *fyne.PointEvent) {
	if fyne.CurrentDevice().IsMobile() && e.selecting {
		e.selecting = false
	}
	e.updateMousePointer(ev, false)
}

// TappedSecondary is called when right or alternative tap is invoked.
//
// Opens the PopUpMenu with `Paste` item to paste text from the clipboard.
//
// Implements: fyne.SecondaryTappable
func (e *Entry) TappedSecondary(pe *fyne.PointEvent) {
	if e.Disabled() && e.concealed() {
		return // no popup options for a disabled concealed field
	}

	cutItem := fyne.NewMenuItem("Cut", func() {
		clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
		e.cutToClipboard(clipboard)
	})
	copyItem := fyne.NewMenuItem("Copy", func() {
		clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
		e.copyToClipboard(clipboard)
	})
	pasteItem := fyne.NewMenuItem("Paste", func() {
		clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
		e.pasteFromClipboard(clipboard)
	})
	selectAllItem := fyne.NewMenuItem("Select all", e.selectAll)

	super := e.super()
	entryPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(super)
	popUpPos := entryPos.Add(fyne.NewPos(pe.Position.X, pe.Position.Y))
	c := fyne.CurrentApp().Driver().CanvasForObject(super)

	var menu *fyne.Menu
	if e.Disabled() {
		menu = fyne.NewMenu("", copyItem, selectAllItem)
	} else if e.concealed() {
		menu = fyne.NewMenu("", pasteItem, selectAllItem)
	} else {
		menu = fyne.NewMenu("", cutItem, copyItem, pasteItem, selectAllItem)
	}

	e.popUp = NewPopUpMenu(menu, c)
	e.popUp.ShowAtPosition(popUpPos)
}

// TypedKey receives key input events when the Entry widget is focused.
//
// Implements: fyne.Focusable
func (e *Entry) TypedKey(key *fyne.KeyEvent) {
	if e.Disabled() {
		return
	}

	e.propertyLock.RLock()
	provider := e.textProvider()
	onSubmitted := e.OnSubmitted
	multiLine := e.MultiLine
	selectDown := e.selectKeyDown
	text := e.Text
	e.propertyLock.RUnlock()

	if e.selectKeyDown || e.selecting {
		if e.selectingKeyHandler(key) {
			e.Refresh()
			return
		}
	}

	switch key.Name {
	case fyne.KeyBackspace:
		e.propertyLock.RLock()
		isEmpty := provider.len() == 0 || (e.CursorColumn == 0 && e.CursorRow == 0)
		e.propertyLock.RUnlock()
		if isEmpty {
			return
		}

		e.propertyLock.Lock()
		pos := e.cursorTextPos()
		provider.deleteFromTo(pos-1, pos)
		e.CursorRow, e.CursorColumn = e.rowColFromTextPos(pos - 1)
		e.propertyLock.Unlock()
	case fyne.KeyDelete:
		pos := e.cursorTextPos()
		if provider.len() == 0 || pos == provider.len() {
			return
		}

		e.propertyLock.Lock()
		provider.deleteFromTo(pos, pos+1)
		e.propertyLock.Unlock()
	case fyne.KeyReturn, fyne.KeyEnter:
		if !multiLine {
			// Single line doesn't support newline.
			// Call submitted callback, if any.
			if onSubmitted != nil {
				onSubmitted(text)
			}
			return
		} else if selectDown && onSubmitted != nil {
			// Multiline supports newline, unless shift is held and OnSubmitted is set.
			onSubmitted(text)
			return
		}
		e.propertyLock.Lock()
		provider.insertAt(e.cursorTextPos(), []rune("\n"))
		e.CursorColumn = 0
		e.CursorRow++
		e.propertyLock.Unlock()
	case fyne.KeyTab:
		e.TypedRune('\t')
	case fyne.KeyUp:
		if !multiLine {
			return
		}

		e.propertyLock.Lock()
		if e.CursorRow > 0 {
			e.CursorRow--
		}

		rowLength := provider.rowLength(e.CursorRow)
		if e.CursorColumn > rowLength {
			e.CursorColumn = rowLength
		}
		e.propertyLock.Unlock()
	case fyne.KeyDown:
		if !multiLine {
			return
		}

		e.propertyLock.Lock()
		if e.CursorRow < provider.rows()-1 {
			e.CursorRow++
		}

		rowLength := provider.rowLength(e.CursorRow)
		if e.CursorColumn > rowLength {
			e.CursorColumn = rowLength
		}
		e.propertyLock.Unlock()
	case fyne.KeyLeft:
		e.propertyLock.Lock()
		if e.CursorColumn > 0 {
			e.CursorColumn--
		} else if e.MultiLine && e.CursorRow > 0 {
			e.CursorRow--
			e.CursorColumn = provider.rowLength(e.CursorRow)
		}
		e.propertyLock.Unlock()
	case fyne.KeyRight:
		e.propertyLock.Lock()
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
		e.propertyLock.Unlock()
	case fyne.KeyEnd:
		e.propertyLock.Lock()
		if e.MultiLine {
			e.CursorColumn = provider.rowLength(e.CursorRow)
		} else {
			e.CursorColumn = provider.len()
		}
		e.propertyLock.Unlock()
	case fyne.KeyHome:
		e.propertyLock.Lock()
		e.CursorColumn = 0
		e.propertyLock.Unlock()
	case fyne.KeyPageUp:
		e.propertyLock.Lock()
		if e.MultiLine {
			e.CursorRow = 0
		}
		e.CursorColumn = 0
		e.propertyLock.Unlock()
	case fyne.KeyPageDown:
		e.propertyLock.Lock()
		if e.MultiLine {
			e.CursorRow = provider.rows() - 1
			e.CursorColumn = provider.rowLength(e.CursorRow)
		} else {
			e.CursorColumn = provider.len()
		}
		e.propertyLock.Unlock()
	default:
		return
	}

	e.propertyLock.Lock()
	if e.CursorRow == e.selectRow && e.CursorColumn == e.selectColumn {
		e.selecting = false
	}
	e.propertyLock.Unlock()
	e.updateText(provider.String())
}

// TypedRune receives text input events when the Entry widget is focused.
//
// Implements: fyne.Focusable
func (e *Entry) TypedRune(r rune) {
	if e.Disabled() {
		return
	}

	e.propertyLock.Lock()
	if e.popUp != nil {
		e.popUp.Hide()
	}

	selecting := e.selecting
	e.propertyLock.Unlock()

	// if we've typed a character and we're selecting then replace the selection with the character
	if selecting {
		e.eraseSelection()
	}

	e.propertyLock.Lock()
	provider := e.textProvider()
	e.selecting = false

	runes := []rune{r}
	pos := e.cursorTextPos()
	provider.insertAt(pos, runes)
	e.CursorRow, e.CursorColumn = e.rowColFromTextPos(pos + len(runes))

	content := provider.String()
	e.propertyLock.Unlock()
	e.updateText(content)
}

// TypedShortcut implements the Shortcutable interface
//
// Implements: fyne.Shortcutable
func (e *Entry) TypedShortcut(shortcut fyne.Shortcut) {
	e.shortcut.TypedShortcut(shortcut)
}

// Unbind disconnects any configured data source from this Entry.
// The current value will remain at the last value of the data source.
//
// Since: 2.0
func (e *Entry) Unbind() {
	e.OnChanged = nil
	if e.textSource == nil || e.textListener == nil {
		return
	}

	e.Validator = nil
	e.textSource.RemoveListener(e.textListener)
	e.textListener = nil
	e.textSource = nil
}

// concealed tells the rendering textProvider if we are a concealed field
func (e *Entry) concealed() bool {
	return e.Password
}

// copyToClipboard copies the current selection to a given clipboard.
// This does nothing if it is a concealed entry.
func (e *Entry) copyToClipboard(clipboard fyne.Clipboard) {
	if !e.selecting || e.concealed() {
		return
	}

	clipboard.SetContent(e.SelectedText())
}

func (e *Entry) cursorColAt(text []rune, pos fyne.Position) int {
	for i := 0; i < len(text); i++ {
		str := string(text[0 : i+1])
		wid := fyne.MeasureText(str, theme.TextSize(), e.textStyle()).Width + theme.Padding()
		if wid > pos.X+theme.Padding() {
			return i
		}
	}
	return len(text)
}

func (e *Entry) cursorTextPos() (pos int) {
	return e.textPosFromRowCol(e.CursorRow, e.CursorColumn)
}

// copyToClipboard copies the current selection to a given clipboard and then removes the selected text.
// This does nothing if it is a concealed entry.
func (e *Entry) cutToClipboard(clipboard fyne.Clipboard) {
	if !e.selecting || e.concealed() {
		return
	}

	e.copyToClipboard(clipboard)
	e.eraseSelection()
}

// eraseSelection removes the current selected region and moves the cursor
func (e *Entry) eraseSelection() {
	if e.Disabled() {
		return
	}

	provider := e.textProvider()
	posA, posB := e.selection()

	if posA == posB {
		return
	}

	e.propertyLock.Lock()
	provider.deleteFromTo(posA, posB)
	e.CursorRow, e.CursorColumn = e.rowColFromTextPos(posA)
	e.selectRow, e.selectColumn = e.CursorRow, e.CursorColumn
	e.selecting = false
	e.propertyLock.Unlock()
	e.updateText(provider.String())
}

func (e *Entry) getRowCol(ev *fyne.PointEvent) (int, int) {
	e.propertyLock.RLock()
	defer e.propertyLock.RUnlock()

	rowHeight := e.textProvider().charMinSize().Height
	row := int(math.Floor(float64(ev.Position.Y+e.scroll.Offset.Y-theme.Padding()) / float64(rowHeight)))
	col := 0
	if row < 0 {
		row = 0
	} else if row >= e.textProvider().rows() {
		row = e.textProvider().rows() - 1
		col = 0
	} else {
		col = e.cursorColAt(e.textProvider().row(row), ev.Position.Add(e.scroll.Offset))
	}

	return row, col
}

// object returns the root object of the widget so it can be referenced
func (e *Entry) object() fyne.Widget {
	return nil
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
		lastNewlineIndex := 0
		for i, r := range runes {
			if r == '\n' {
				lastNewlineIndex = i
			}
		}
		e.CursorColumn = len(runes) - lastNewlineIndex - 1
	}
	e.updateText(provider.String())
	e.Refresh()
}

// placeholderProvider returns the placeholder text handler for this entry
func (e *Entry) placeholderProvider() *textProvider {
	if e.placeholder != nil {
		return e.placeholder
	}

	text := newTextProvider(e.PlaceHolder, &placeholderPresenter{e})
	text.ExtendBaseWidget(text)
	text.extraPad = fyne.NewSize(theme.Padding(), theme.InputBorderSize())
	e.placeholder = text
	return e.placeholder
}

func (e *Entry) registerShortcut() {
	e.shortcut.AddShortcut(&fyne.ShortcutCut{}, func(se fyne.Shortcut) {
		cut := se.(*fyne.ShortcutCut)
		e.cutToClipboard(cut.Clipboard)
	})
	e.shortcut.AddShortcut(&fyne.ShortcutCopy{}, func(se fyne.Shortcut) {
		cpy := se.(*fyne.ShortcutCopy)
		e.copyToClipboard(cpy.Clipboard)
	})
	e.shortcut.AddShortcut(&fyne.ShortcutPaste{}, func(se fyne.Shortcut) {
		paste := se.(*fyne.ShortcutPaste)
		e.pasteFromClipboard(paste.Clipboard)
	})
	e.shortcut.AddShortcut(&fyne.ShortcutSelectAll{}, func(se fyne.Shortcut) {
		e.selectAll()
	})
}

// Obtains row,col from a given textual position
// expects a read or write lock to be held by the caller
func (e *Entry) rowColFromTextPos(pos int) (row int, col int) {
	provider := e.textProvider()
	canWrap := e.Wrapping == fyne.TextWrapBreak || e.Wrapping == fyne.TextWrapWord
	for i := 0; i < provider.rows(); i++ {
		b := provider.rowBoundary(i)
		if b[0] <= pos {
			if b[1] < pos {
				row++
			}
			col = pos - b[0]
			if canWrap && b[0] == pos && col == 0 && pos != 0 {
				row++
			}
		} else {
			break
		}
	}
	return
}

// selectAll selects all text in entry
func (e *Entry) selectAll() {
	if e.textProvider().len() == 0 {
		return
	}
	e.setFieldsAndRefresh(func() {
		e.selectRow = 0
		e.selectColumn = 0

		lastRow := e.textProvider().rows() - 1
		e.CursorColumn = e.textProvider().rowLength(lastRow)
		e.CursorRow = lastRow
		e.selecting = true
	})
}

// selectingKeyHandler performs keypress action in the scenario that a selection
// is either a) in progress or b) about to start
// returns true if the keypress has been fully handled
func (e *Entry) selectingKeyHandler(key *fyne.KeyEvent) bool {

	if e.selectKeyDown && !e.selecting {
		switch key.Name {
		case fyne.KeyUp, fyne.KeyDown,
			fyne.KeyLeft, fyne.KeyRight,
			fyne.KeyEnd, fyne.KeyHome,
			fyne.KeyPageUp, fyne.KeyPageDown:
			e.selecting = true
		}
	}

	if !e.selecting {
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

	if !e.selectKeyDown {
		switch key.Name {
		case fyne.KeyLeft:
			// seek to the start of the selection -- return handled
			selectStart, _ := e.selection()
			e.propertyLock.Lock()
			e.CursorRow, e.CursorColumn = e.rowColFromTextPos(selectStart)
			e.selecting = false
			e.propertyLock.Unlock()
			return true
		case fyne.KeyRight:
			// seek to the end of the selection -- return handled
			_, selectEnd := e.selection()
			e.propertyLock.Lock()
			e.CursorRow, e.CursorColumn = e.rowColFromTextPos(selectEnd)
			e.selecting = false
			e.propertyLock.Unlock()
			return true
		case fyne.KeyUp, fyne.KeyDown, fyne.KeyEnd, fyne.KeyHome, fyne.KeyPageUp, fyne.KeyPageDown:
			// cursor movement without left or right shift -- clear selection and return unhandled
			e.selecting = false
			return false
		}
	}

	return false
}

// selection returns the start and end text positions for the selected span of text
// Note: this functionality depends on the relationship between the selection start row/col and
// the current cursor row/column.
// eg: (whitespace for clarity, '_' denotes cursor)
//   "T  e  s [t  i]_n  g" == 3, 5
//   "T  e  s_[t  i] n  g" == 3, 5
//   "T  e_[s  t  i] n  g" == 2, 5
func (e *Entry) selection() (int, int) {
	e.propertyLock.RLock()
	noSelection := !e.selecting || (e.CursorRow == e.selectRow && e.CursorColumn == e.selectColumn)
	e.propertyLock.RUnlock()

	if noSelection {
		return -1, -1
	}

	e.propertyLock.Lock()
	defer e.propertyLock.Unlock()
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

// textAlign tells the rendering textProvider our alignment
func (e *Entry) textAlign() fyne.TextAlign {
	return fyne.TextAlignLeading
}

// textColor tells the rendering textProvider our color
func (e *Entry) textColor() color.Color {
	if e.Disabled() {
		return theme.DisabledColor()
	}
	return theme.ForegroundColor()
}

// Obtains textual position from a given row and col
// expects a read or write lock to be held by the caller
func (e *Entry) textPosFromRowCol(row, col int) int {
	return e.textProvider().rowBoundary(row)[0] + col
}

// textProvider returns the text handler for this entry
func (e *Entry) textProvider() *textProvider {
	if e.text != nil {
		return e.text
	}

	text := newTextProvider(e.Text, e)
	text.ExtendBaseWidget(text)
	text.extraPad = fyne.NewSize(theme.Padding(), theme.InputBorderSize())
	e.text = text
	return e.text
}

// textStyle tells the rendering textProvider our style
func (e *Entry) textStyle() fyne.TextStyle {
	return e.TextStyle
}

// textWrap tells the rendering textProvider our wrapping
func (e *Entry) textWrap() fyne.TextWrap {
	if e.Wrapping == fyne.TextTruncate { // this is now the default - but we scroll around this large content
		return fyne.TextWrapOff
	}

	if !e.MultiLine && (e.Wrapping == fyne.TextWrapBreak || e.Wrapping == fyne.TextWrapWord) {
		fyne.LogError("Entry cannot wrap single line", nil)
		e.Wrapping = fyne.TextTruncate
	}
	return e.Wrapping
}

func (e *Entry) updateMousePointer(ev *fyne.PointEvent, rightClick bool) {
	row, col := e.getRowCol(ev)
	e.setFieldsAndRefresh(func() {
		if !rightClick || rightClick && !e.selecting {
			e.CursorRow = row
			e.CursorColumn = col
		}

		if !e.selecting {
			e.selectRow = row
			e.selectColumn = col
		}
	})
}

// updateText updates the internal text to the given value
func (e *Entry) updateText(text string) {
	var callback func(string)
	e.setFieldsAndRefresh(func() {
		changed := e.Text != text
		e.Text = text

		if changed {
			callback = e.OnChanged
		}
	})

	if validate := e.Validator; validate != nil {
		e.SetValidationError(validate(text))
	}

	if callback != nil {
		callback(text)
	}
}

var _ fyne.WidgetRenderer = (*entryRenderer)(nil)

type entryRenderer struct {
	box, line *canvas.Rectangle
	scroll    *widget.Scroll

	objects []fyne.CanvasObject
	entry   *Entry
}

func (r *entryRenderer) Destroy() {
}

func (r *entryRenderer) Layout(size fyne.Size) {
	r.line.Resize(fyne.NewSize(size.Width, theme.InputBorderSize()))
	r.line.Move(fyne.NewPos(0, size.Height-theme.InputBorderSize()))
	r.box.Resize(size.Subtract(fyne.NewSize(0, theme.InputBorderSize()*2)))
	r.box.Move(fyne.NewPos(0, theme.InputBorderSize()))

	actionIconSize := fyne.NewSize(0, 0)
	xInset := float32(0)
	if r.entry.ActionItem != nil {
		actionIconSize = fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
		xInset = theme.IconInlineSize() + 2*theme.Padding()

		r.entry.ActionItem.Resize(actionIconSize)
		r.entry.ActionItem.Move(fyne.NewPos(size.Width-actionIconSize.Width-2*theme.Padding(), theme.Padding()*2))
	}

	validatorIconSize := fyne.NewSize(0, 0)
	if r.entry.Validator != nil {
		validatorIconSize = fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())

		r.ensureValidationSetup()
		r.entry.validationStatus.Resize(validatorIconSize)

		if r.entry.ActionItem == nil {
			r.entry.validationStatus.Move(fyne.NewPos(size.Width-validatorIconSize.Width-2*theme.Padding(), theme.Padding()*2))
			xInset = theme.IconInlineSize() + 2*theme.Padding()
		} else {
			r.entry.validationStatus.Move(fyne.NewPos(size.Width-validatorIconSize.Width-actionIconSize.Width-3*theme.Padding(), theme.Padding()*2))
			xInset += theme.IconInlineSize() + theme.Padding()
		}
	}

	entrySize := size.Subtract(fyne.NewSize(xInset, theme.InputBorderSize()*2))
	entryPos := fyne.NewPos(0, theme.InputBorderSize())
	r.scroll.Resize(entrySize)
	r.scroll.Move(entryPos)
}

// MinSize calculates the minimum size of an entry widget.
// This is based on the contained text with a standard amount of padding added.
// If MultiLine is true then we will reserve space for at leasts 3 lines
func (r *entryRenderer) MinSize() fyne.Size {
	if r.scroll.Direction == widget.ScrollNone {
		return r.scroll.MinSize().Add(fyne.NewSize(0, theme.InputBorderSize()*2))
	}

	minSize := r.entry.placeholderProvider().charMinSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))

	if r.entry.MultiLine {
		// ensure multiline height is at least charMinSize * multilineRows
		rowHeight := r.entry.text.charMinSize().Height * multiLineRows
		minSize.Height = fyne.Max(minSize.Height, rowHeight+(multiLineRows-1)*theme.Padding())
	}

	return minSize.Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
}

func (r *entryRenderer) Objects() []fyne.CanvasObject {
	r.entry.propertyLock.RLock()
	defer r.entry.propertyLock.RUnlock()

	return r.objects
}

func (r *entryRenderer) Refresh() {
	r.entry.propertyLock.RLock()
	provider := r.entry.textProvider()
	content := r.entry.Text
	focused := r.entry.focused
	r.entry.propertyLock.RUnlock()

	if content != string(provider.buffer) {
		r.entry.SetText(content)
		return
	}

	r.box.FillColor = theme.InputBackgroundColor()
	if focused {
		r.line.FillColor = theme.PrimaryColor()
	} else {
		if r.entry.Disabled() {
			r.line.FillColor = theme.DisabledColor()
		} else {
			r.line.FillColor = theme.ShadowColor()
		}
	}

	r.entry.text.propertyLock.Lock()
	r.entry.text.updateRowBounds()
	r.entry.text.propertyLock.Unlock()
	r.entry.placeholder.propertyLock.Lock()
	r.entry.placeholder.updateRowBounds()
	r.entry.placeholder.propertyLock.Unlock()

	r.entry.text.Refresh()
	r.entry.placeholder.Refresh()
	if r.entry.ActionItem != nil {
		r.entry.ActionItem.Refresh()
	}

	if r.entry.Validator != nil {
		if !r.entry.focused && r.entry.Text != "" && r.entry.validationError != nil {
			r.line.FillColor = theme.ErrorColor()
		}
		r.ensureValidationSetup()
		r.entry.validationStatus.Refresh()
	} else if r.entry.validationStatus != nil {
		r.entry.validationStatus.Hide()
	}

	cache.Renderer(r.scroll.Content.(*entryContent)).Refresh()
	canvas.Refresh(r.entry.super())
}

func (r *entryRenderer) ensureValidationSetup() {
	if r.entry.validationStatus == nil {
		r.entry.validationStatus = newValidationStatus(r.entry)
		r.objects = append(r.objects, r.entry.validationStatus)
		r.Layout(r.entry.size)
		r.Refresh()
	}
}

var _ fyne.Widget = (*entryContent)(nil)

type entryContent struct {
	BaseWidget

	entry  *Entry
	scroll *widget.Scroll
}

func (e *entryContent) CreateRenderer() fyne.WidgetRenderer {
	e.ExtendBaseWidget(e)

	cursor := canvas.NewRectangle(color.Transparent)
	cursor.Hide()

	e.entry.propertyLock.Lock()
	defer e.entry.propertyLock.Unlock()
	provider := e.entry.textProvider()
	placeholder := e.entry.placeholderProvider()
	if provider.len() != 0 {
		placeholder.Hide()
	}
	objects := []fyne.CanvasObject{placeholder, provider, cursor}

	r := &entryContentRenderer{cursor, []fyne.CanvasObject{}, nil, objects,
		provider, placeholder, e}
	r.updateScrollDirections()
	r.Layout(e.size)
	return r
}

// DragEnd is called at end of a drag event. It does nothing.
//
// Implements: fyne.Draggable
func (e *entryContent) DragEnd() {
	impl := e.entry.super()
	// we need to propagate the focus, top level widget handles focus APIs
	fyne.CurrentApp().Driver().CanvasForObject(impl).Focus(impl.(interface{}).(fyne.Focusable))

	e.entry.DragEnd()
}

// Dragged is called when the pointer moves while a button is held down.
// It updates the selection accordingly.
//
// Implements: fyne.Draggable
func (e *entryContent) Dragged(d *fyne.DragEvent) {
	e.entry.Dragged(d)
}

var _ fyne.WidgetRenderer = (*entryContentRenderer)(nil)

type entryContentRenderer struct {
	cursor     *canvas.Rectangle
	selection  []fyne.CanvasObject
	cursorAnim *fyne.Animation
	objects    []fyne.CanvasObject

	provider, placeholder *textProvider
	content               *entryContent
}

func (r *entryContentRenderer) Destroy() {
	r.cursorAnim.Stop()
}

func (r *entryContentRenderer) Layout(size fyne.Size) {
	r.provider.Resize(size)
	r.placeholder.Resize(size)
}

func (r *entryContentRenderer) MinSize() fyne.Size {
	minSize := r.content.entry.placeholderProvider().MinSize()

	if r.content.entry.textProvider().len() > 0 {
		minSize = r.content.entry.text.MinSize()
	}

	return minSize
}

func (r *entryContentRenderer) Objects() []fyne.CanvasObject {
	r.content.entry.propertyLock.RLock()
	defer r.content.entry.propertyLock.RUnlock()
	// Objects are generated dynamically force selection rectangles to appear underneath the text
	if r.content.entry.selecting {
		return append(r.selection, r.objects...)
	}
	return r.objects
}

func (r *entryContentRenderer) Refresh() {
	r.content.entry.propertyLock.RLock()
	provider := r.content.entry.textProvider()
	placeholder := r.content.entry.placeholderProvider()
	content := r.content.entry.Text
	focused := r.content.entry.focused
	selections := r.selection
	r.updateScrollDirections()
	r.content.entry.propertyLock.RUnlock()

	if content != string(provider.buffer) {
		return
	}

	if provider.len() == 0 {
		placeholder.Show()
	} else if placeholder.Visible() {
		placeholder.Hide()
	}

	if focused {
		r.cursor.Show()
		if r.cursorAnim == nil {
			r.cursorAnim = makeCursorAnimation(r.cursor)
			r.cursorAnim.Start()
		}
	} else {
		if r.cursorAnim != nil {
			r.cursorAnim.Stop()
			r.cursorAnim = nil
		}
		r.cursor.Hide()
	}
	r.moveCursor()

	for _, selection := range selections {
		selection.(*canvas.Rectangle).Hidden = !r.content.entry.focused && !r.content.entry.disabled
		selection.(*canvas.Rectangle).FillColor = theme.PrimaryColor()
	}

	canvas.Refresh(r.content)
}

// This process builds a slice of rectangles:
// - one entry per row of text
// - ordered by row order as they occur in multiline text
// This process could be optimized in the scenario where the user is selecting upwards:
// If the upwards case instead produces an order-reversed slice then only the newest rectangle would
// require movement and resizing. The existing solution creates a new rectangle and then moves/resizes
// all rectangles to comply with the occurrence order as stated above.
func (r *entryContentRenderer) buildSelection() {
	r.content.entry.propertyLock.RLock()
	cursorRow, cursorCol := r.content.entry.CursorRow, r.content.entry.CursorColumn
	selectRow, selectCol := -1, -1
	if r.content.entry.selecting {
		selectRow = r.content.entry.selectRow
		selectCol = r.content.entry.selectColumn
	}
	r.content.entry.propertyLock.RUnlock()

	if selectRow == -1 {
		r.selection = r.selection[:0]

		return
	}

	provider := r.content.entry.textProvider()
	// Convert column, row into x,y
	getCoordinates := func(column int, row int) (float32, float32) {
		sz := provider.lineSizeToColumn(column, row)
		return sz.Width + theme.Padding(), sz.Height*float32(row) + theme.Padding()
	}

	lineHeight := r.content.entry.text.charMinSize().Height

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
	if len(r.selection) > rowCount {
		r.selection = r.selection[:rowCount]
	}

	r.content.entry.propertyLock.Lock()
	defer r.content.entry.propertyLock.Unlock()
	// build a rectangle for each row and add it to r.selection
	for i := 0; i < rowCount; i++ {
		if len(r.selection) <= i {
			box := canvas.NewRectangle(theme.PrimaryColor())
			r.selection = append(r.selection, box)
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
		r.selection[i].Resize(fyne.NewSize(x2-x1+1, lineHeight))
		r.selection[i].Move(fyne.NewPos(x1-1, y1+theme.InputBorderSize()))
	}
}

func (r *entryContentRenderer) ensureCursorVisible() {
	cx1 := r.cursor.Position().X
	cy1 := r.cursor.Position().Y
	cx2 := cx1 + r.cursor.Size().Width
	cy2 := cy1 + r.cursor.Size().Height
	offset := r.content.scroll.Offset
	size := r.content.scroll.Size()

	if offset.X <= cx1 && cx2 < offset.X+size.Width &&
		offset.Y <= cy1 && cy2 < offset.Y+size.Height {
		return
	}

	move := fyne.NewDelta(0, 0)
	if cx1 < offset.X {
		move.DX -= offset.X - cx1
	} else if cx2 >= offset.X+size.Width {
		move.DX += cx2 - (offset.X + size.Width)
	}
	if cy1 < offset.Y {
		move.DY -= offset.Y - cy1
	} else if cy2 >= offset.X+size.Height {
		move.DY += cy2 - (offset.Y + size.Height)
	}
	r.content.scroll.Offset = r.content.scroll.Offset.Add(move)
	r.content.scroll.Refresh()
}

func (r *entryContentRenderer) moveCursor() {
	// build r.selection[] if the user has made a selection
	r.buildSelection()
	r.content.entry.propertyLock.RLock()
	provider := r.content.entry.textProvider()
	provider.propertyLock.RLock()
	size := provider.lineSizeToColumn(r.content.entry.CursorColumn, r.content.entry.CursorRow)
	provider.propertyLock.RUnlock()
	xPos := size.Width
	yPos := size.Height * float32(r.content.entry.CursorRow)
	r.content.entry.propertyLock.RUnlock()

	r.content.entry.propertyLock.Lock()
	lineHeight := r.content.entry.text.charMinSize().Height
	r.cursor.Resize(fyne.NewSize(2, lineHeight))
	r.cursor.Move(fyne.NewPos(xPos-1+theme.Padding(), yPos+theme.Padding()+theme.InputBorderSize()))

	callback := r.content.entry.OnCursorChanged
	r.content.entry.propertyLock.Unlock()
	r.ensureCursorVisible()

	if callback != nil {
		callback()
	}
}

func (r *entryContentRenderer) updateScrollDirections() {
	switch r.content.entry.Wrapping {
	case fyne.TextWrapOff:
		r.content.scroll.Direction = widget.ScrollNone
	case fyne.TextTruncate: // this is now the default - but we scroll
		r.content.scroll.Direction = widget.ScrollBoth
	default: // fyne.TextWrapBreak, fyne.TextWrapWord
		r.content.scroll.Direction = widget.ScrollVerticalOnly
	}
}

type placeholderPresenter struct {
	e *Entry
}

// concealed tells the rendering textProvider if we are a concealed field
// placeholder text is not obfuscated, returning false
func (p *placeholderPresenter) concealed() bool {
	return false
}

// object returns the root object of the widget so it can be referenced
func (p *placeholderPresenter) object() fyne.Widget {
	return nil
}

// textAlign tells the rendering textProvider our alignment
func (p *placeholderPresenter) textAlign() fyne.TextAlign {
	return fyne.TextAlignLeading
}

// textColor tells the rendering textProvider our color
func (p *placeholderPresenter) textColor() color.Color {
	if p.e.Disabled() {
		return theme.DisabledColor()
	}
	return theme.PlaceHolderColor()
}

// textStyle tells the rendering textProvider our style
func (p *placeholderPresenter) textStyle() fyne.TextStyle {
	return fyne.TextStyle{}
}

// textWrap tells the rendering textProvider our wrapping
func (p *placeholderPresenter) textWrap() fyne.TextWrap {
	return p.e.Wrapping
}

// getTextWhitespaceRegion returns the start/end markers for selection highlight on starting from col
// and expanding to the start and end of the whitespace or text underneath the specified position.
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

func makeCursorAnimation(cursor *canvas.Rectangle) *fyne.Animation {
	cursorOpaque := theme.PrimaryColor()
	r, g, b, _ := theme.PrimaryColor().RGBA()
	cursorDim := color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0x16}
	anim := canvas.NewColorRGBAAnimation(cursorDim, cursorOpaque, time.Second/2, func(c color.Color) {
		cursor.FillColor = c
		cursor.Refresh()
	})
	anim.RepeatCount = fyne.AnimationRepeatForever
	anim.AutoReverse = true

	return anim
}

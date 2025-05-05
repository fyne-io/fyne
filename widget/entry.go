package widget

import (
	"image/color"
	"runtime"
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
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
)

const (
	bindIgnoreDelay = time.Millisecond * 100 // ignore incoming DataItem fire after we have called Set
	multiLineRows   = 3
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
var _ mobile.Touchable = (*Entry)(nil)
var _ fyne.Tabbable = (*Entry)(nil)

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

	// Scroll can be used to turn off the scrolling of our entry when Wrapping is WrapNone.
	//
	// Since: 2.4
	Scroll fyne.ScrollDirection

	// Set a validator that this entry will check against
	// Since: 1.4
	Validator           fyne.StringValidator `json:"-"`
	validationStatus    *validationStatus
	onValidationChanged func(error)
	validationError     error

	CursorRow, CursorColumn int
	OnCursorChanged         func() `json:"-"`

	cursorAnim *entryCursorAnimation

	dirty       bool
	focused     bool
	text        RichText
	placeholder RichText
	content     *entryContent
	scroll      *widget.Scroll

	// useful for Form validation (as the error text should only be shown when
	// the entry is unfocused)
	onFocusChanged func(bool)

	// selectKeyDown indicates whether left shift or right shift is currently held down
	selectKeyDown bool

	sel   *selectable
	popUp *PopUpMenu
	// TODO: Add OnSelectChanged

	// ActionItem is a small item which is displayed at the outer right of the entry (like a password revealer)
	ActionItem      fyne.CanvasObject `json:"-"`
	binder          basicBinder
	conversionError error
	minCache        fyne.Size
	multiLineRows   int // override global default number of visible lines

	// undoStack stores the data necessary for undo/redo functionality
	// See entryUndoStack for implementation details.
	undoStack entryUndoStack
}

// NewEntry creates a new single line entry widget.
func NewEntry() *Entry {
	e := &Entry{Wrapping: fyne.TextWrap(fyne.TextTruncateClip)}
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
	e := &Entry{MultiLine: true, Wrapping: fyne.TextWrap(fyne.TextTruncateClip)}
	e.ExtendBaseWidget(e)
	return e
}

// NewPasswordEntry creates a new entry password widget
func NewPasswordEntry() *Entry {
	e := &Entry{Password: true, Wrapping: fyne.TextWrap(fyne.TextTruncateClip)}
	e.ExtendBaseWidget(e)
	e.ActionItem = newPasswordRevealer(e)
	return e
}

// AcceptsTab returns if Entry accepts the Tab key or not.
//
// Implements: fyne.Tabbable
//
// Since: 2.1
func (e *Entry) AcceptsTab() bool {
	return e.MultiLine
}

// Bind connects the specified data source to this Entry.
// The current value will be displayed and any changes in the data will cause the widget to update.
// User interactions with this Entry will set the value into the data source.
//
// Since: 2.0
func (e *Entry) Bind(data binding.String) {
	e.binder.SetCallback(e.updateFromData)
	e.binder.Bind(data)

	e.Validator = func(string) error {
		return e.conversionError
	}
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
//
// Implements: fyne.Widget
func (e *Entry) CreateRenderer() fyne.WidgetRenderer {
	th := e.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	e.ExtendBaseWidget(e)

	// initialise
	e.textProvider()
	e.placeholderProvider()
	e.syncSelectable()

	box := canvas.NewRectangle(th.Color(theme.ColorNameInputBackground, v))
	box.CornerRadius = th.Size(theme.SizeNameInputRadius)
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = th.Size(theme.SizeNameInputBorder)
	border.StrokeColor = th.Color(theme.ColorNameInputBorder, v)
	border.CornerRadius = th.Size(theme.SizeNameInputRadius)
	cursor := canvas.NewRectangle(color.Transparent)
	cursor.Hide()

	e.cursorAnim = newEntryCursorAnimation(cursor)
	e.content = &entryContent{entry: e}
	e.scroll = widget.NewScroll(nil)
	objects := []fyne.CanvasObject{box, border}
	if e.Wrapping != fyne.TextWrapOff || e.Scroll != widget.ScrollNone {
		e.scroll.Content = e.content
		objects = append(objects, e.scroll)
	} else {
		e.scroll.Hide()
		objects = append(objects, e.content)
	}
	e.content.scroll = e.scroll

	if e.Password && e.ActionItem == nil {
		// An entry widget has been created via struct setting manually
		// the Password field to true. Going to enable the password revealer.
		e.ActionItem = newPasswordRevealer(e)
	}

	if e.ActionItem != nil {
		objects = append(objects, e.ActionItem)
	}

	e.syncSegments()
	return &entryRenderer{box, border, e.scroll, objects, e}
}

// Cursor returns the cursor type of this widget
//
// Implements: desktop.Cursorable
func (e *Entry) Cursor() desktop.Cursor {
	return desktop.TextCursor
}

// DoubleTapped is called when this entry has been double tapped so we should select text below the pointer
//
// Implements: fyne.DoubleTappable
func (e *Entry) DoubleTapped(_ *fyne.PointEvent) {
	e.focused = true
	e.syncSelectable()
	e.sel.doubleTappedAtUnixMillis = time.Now().UnixMilli()
	row := e.textProvider().row(e.CursorRow)
	start, end := getTextWhitespaceRegion(row, e.CursorColumn, false)
	if start == -1 || end == -1 {
		return
	}

	e.setFieldsAndRefresh(func() {
		if !e.selectKeyDown {
			e.sel.selectRow = e.CursorRow
			e.sel.selectColumn = start
		}
		// Always aim to maximise the selected region
		if e.sel.selectRow > e.CursorRow || (e.sel.selectRow == e.CursorRow && e.sel.selectColumn > e.CursorColumn) {
			e.CursorColumn = start
		} else {
			e.CursorColumn = end
		}

		e.syncSelectable()
		e.sel.selecting = true
	})
}

// DragEnd is called at end of a drag event.
//
// Implements: fyne.Draggable
func (e *Entry) DragEnd() {
	e.syncSelectable()

	if e.CursorColumn == e.sel.selectColumn && e.CursorRow == e.sel.selectRow {
		e.sel.selecting = false
	}
}

// Dragged is called when the pointer moves while a button is held down.
// It updates the selection accordingly.
//
// Implements: fyne.Draggable
func (e *Entry) Dragged(d *fyne.DragEvent) {
	d.Position = d.Position.Add(fyne.NewPos(0, e.Theme().Size(theme.SizeNameInputBorder)))
	e.sel.dragged(d, false)
	e.updateMousePointer(d.Position, false)
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
func (e *Entry) ExtendBaseWidget(wid fyne.Widget) {
	e.BaseWidget.ExtendBaseWidget(wid)
	e.registerShortcut()
}

// FocusGained is called when the Entry has been given focus.
//
// Implements: fyne.Focusable
func (e *Entry) FocusGained() {
	e.setFieldsAndRefresh(func() {
		e.dirty = true
		e.focused = true
	})
	if e.onFocusChanged != nil {
		e.onFocusChanged(true)
	}
}

// FocusLost is called when the Entry has had focus removed.
//
// Implements: fyne.Focusable
func (e *Entry) FocusLost() {
	e.setFieldsAndRefresh(func() {
		e.focused = false
		e.selectKeyDown = false
	})
	if e.onFocusChanged != nil {
		e.onFocusChanged(false)
	}
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
	if e.MultiLine {
		return mobile.DefaultKeyboard
	} else if e.Password {
		return mobile.PasswordKeyboard
	}

	return mobile.SingleLineKeyboard
}

// KeyDown handler for keypress events - used to store shift modifier state for text selection
//
// Implements: desktop.Keyable
func (e *Entry) KeyDown(key *fyne.KeyEvent) {
	if e.Disabled() {
		return
	}
	// For keyboard cursor controlled selection we now need to store shift key state and selection "start"
	// Note: selection start is where the highlight started (if the user moves the selection up or left then
	// the selectRow/Column will not match SelectionStart)
	if key.Name == desktop.KeyShiftLeft || key.Name == desktop.KeyShiftRight {
		if !e.sel.selecting {
			e.sel.selectRow = e.CursorRow
			e.sel.selectColumn = e.CursorColumn
		}
		e.selectKeyDown = true
	}
}

// KeyUp handler for key release events - used to reset shift modifier state for text selection
//
// Implements: desktop.Keyable
func (e *Entry) KeyUp(key *fyne.KeyEvent) {
	if e.Disabled() {
		return
	}
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
	cached := e.minCache
	if !cached.IsZero() {
		return cached
	}

	e.ExtendBaseWidget(e)
	min := e.BaseWidget.MinSize()

	e.minCache = min
	return min
}

// MouseDown called on mouse click, this triggers a mouse click which can move the cursor,
// update the existing selection (if shift is held), or start a selection dragging operation.
//
// Implements: desktop.Mouseable
func (e *Entry) MouseDown(m *desktop.MouseEvent) {
	e.requestFocus()
	e.syncSelectable()

	if isTripleTap(e.sel.doubleTappedAtUnixMillis, time.Now().UnixMilli()) {
		e.sel.selectCurrentRow(false)
		e.CursorColumn = e.sel.cursorColumn
		e.Refresh()
		return
	}
	if e.selectKeyDown {
		e.sel.selecting = true
	}
	if e.sel.selecting && !e.selectKeyDown && m.Button == desktop.MouseButtonPrimary {
		e.sel.selecting = false
	}

	e.updateMousePointer(m.Position.Add(e.scroll.Offset), m.Button == desktop.MouseButtonSecondary)

	if !e.Disabled() {
		e.requestFocus()
	}
}

// MouseUp called on mouse release
// If a mouse drag event has completed then check to see if it has resulted in an empty selection,
// if so, and if a text select key isn't held, then disable selecting
//
// Implements: desktop.Mouseable
func (e *Entry) MouseUp(m *desktop.MouseEvent) {
	e.syncSelectable()
	start, _ := e.sel.selection()
	if start == -1 && e.sel.selecting && !e.selectKeyDown {
		e.sel.selecting = false
	}
}

// Redo un-does the last undo action.
//
// Since: 2.5
func (e *Entry) Redo() {
	newText, action := e.undoStack.Redo(e.Text)
	modify, ok := action.(*entryModifyAction)
	if !ok {
		return
	}
	pos := modify.Position
	if !modify.Delete {
		pos += len(modify.Text)
	}
	e.updateText(newText, false)
	e.CursorRow, e.CursorColumn = e.rowColFromTextPos(pos)
	e.syncSelectable()
	if e.OnChanged != nil {
		e.OnChanged(newText)
	}
	e.Refresh()
}

func (e *Entry) Refresh() {
	e.minCache = fyne.Size{}

	if e.sel != nil {
		e.sel.style = e.TextStyle
		e.sel.theme = e.Theme()
		e.sel.focussed = e.focused
		e.sel.Refresh()
	}
	e.BaseWidget.Refresh()
}

// SelectedText returns the text currently selected in this Entry.
// If there is no selection it will return the empty string.
func (e *Entry) SelectedText() string {
	return e.sel.SelectedText()
}

// SetMinRowsVisible forces a multi-line entry to show `count` number of rows without scrolling.
// This is not a validation or requirement, it just impacts the minimum visible size.
// Use this carefully as Fyne apps can run on small screens so you may wish to add a scroll container if
// this number is high. Default is 3.
//
// Since: 2.2
func (e *Entry) SetMinRowsVisible(count int) {
	e.multiLineRows = count
	e.Refresh()
}

// SetPlaceHolder sets the text that will be displayed if the entry is otherwise empty
func (e *Entry) SetPlaceHolder(text string) {
	e.Theme() // setup theme cache before locking

	e.PlaceHolder = text

	e.placeholderProvider().Segments[0].(*TextSegment).Text = text
	e.placeholder.updateRowBounds()
	e.placeholderProvider().Refresh()
}

// SetText manually sets the text of the Entry to the given text value.
// Calling SetText resets all undo history.
func (e *Entry) SetText(text string) {
	e.setText(text, false)
}

func (e *Entry) setText(text string, fromBinding bool) {
	e.Theme() // setup theme cache before locking
	e.updateTextAndRefresh(text, fromBinding)
	e.updateCursorAndSelection()

	e.undoStack.Clear()
}

// Append appends the text to the end of the entry.
//
// Since: 2.4
func (e *Entry) Append(text string) {
	provider := e.textProvider()
	provider.insertAt(provider.len(), []rune(text))
	content := provider.String()
	changed := e.updateText(content, false)
	cb := e.OnChanged
	e.undoStack.Clear()

	if changed {
		e.validate()
		if cb != nil {
			cb(content)
		}
	}
	e.Refresh()
}

// Tapped is called when this entry has been tapped. We update the cursor position in
// device-specific callbacks (MouseDown() and TouchDown()).
//
// Implements: fyne.Tappable
func (e *Entry) Tapped(ev *fyne.PointEvent) {

	if fyne.CurrentDevice().IsMobile() && e.sel.selecting {
		e.sel.selecting = false
	}
}

// TappedSecondary is called when right or alternative tap is invoked.
//
// Opens the PopUpMenu with `Paste` item to paste text from the clipboard.
//
// Implements: fyne.SecondaryTappable
func (e *Entry) TappedSecondary(pe *fyne.PointEvent) {
	if e.Disabled() && e.Password {
		return // no popup options for a disabled concealed field
	}

	e.requestFocus()
	clipboard := fyne.CurrentApp().Clipboard()
	super := e.super()

	undoItem := fyne.NewMenuItem(lang.L("Undo"), e.Undo)
	redoItem := fyne.NewMenuItem(lang.L("Redo"), e.Redo)
	cutItem := fyne.NewMenuItem(lang.L("Cut"), func() {
		super.(fyne.Shortcutable).TypedShortcut(&fyne.ShortcutCut{Clipboard: clipboard})
	})
	copyItem := fyne.NewMenuItem(lang.L("Copy"), func() {
		super.(fyne.Shortcutable).TypedShortcut(&fyne.ShortcutCopy{Clipboard: clipboard})
	})
	pasteItem := fyne.NewMenuItem(lang.L("Paste"), func() {
		super.(fyne.Shortcutable).TypedShortcut(&fyne.ShortcutPaste{Clipboard: clipboard})
	})
	selectAllItem := fyne.NewMenuItem(lang.L("Select all"), e.selectAll)

	entryPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(super)
	popUpPos := entryPos.Add(fyne.NewPos(pe.Position.X, pe.Position.Y))
	c := fyne.CurrentApp().Driver().CanvasForObject(super)

	var menu *fyne.Menu
	if e.Disabled() {
		menu = fyne.NewMenu("", copyItem, selectAllItem)
	} else if e.Password {
		menu = fyne.NewMenu("", pasteItem, selectAllItem)
	} else {
		var menuItems []*fyne.MenuItem
		canUndo, canRedo := e.undoStack.CanUndo(), e.undoStack.CanRedo()
		if canUndo {
			menuItems = append(menuItems, undoItem)
		}
		if canRedo {
			menuItems = append(menuItems, redoItem)
		}
		if canUndo || canRedo {
			menuItems = append(menuItems, fyne.NewMenuItemSeparator())
		}
		menuItems = append(menuItems, cutItem, copyItem, pasteItem, selectAllItem)
		menu = fyne.NewMenu("", menuItems...)
	}

	e.popUp = NewPopUpMenu(menu, c)
	e.popUp.ShowAtPosition(popUpPos)
}

// TouchDown is called when this entry gets a touch down event on mobile device, we ensure we have focus.
//
// Since: 2.1
//
// Implements: mobile.Touchable
func (e *Entry) TouchDown(ev *mobile.TouchEvent) {
	now := time.Now().UnixMilli()
	e.syncSegments()
	if !e.Disabled() {
		e.requestFocus()
	}
	if isTripleTap(e.sel.doubleTappedAtUnixMillis, now) {
		e.sel.selectCurrentRow(false)
		e.CursorColumn = e.sel.cursorColumn
		e.Refresh()
		return
	}

	e.updateMousePointer(ev.Position, false)
}

// TouchUp is called when this entry gets a touch up event on mobile device.
//
// Since: 2.1
//
// Implements: mobile.Touchable
func (e *Entry) TouchUp(*mobile.TouchEvent) {
}

// TouchCancel is called when this entry gets a touch cancel event on mobile device (app was removed from focus).
//
// Since: 2.1
//
// Implements: mobile.Touchable
func (e *Entry) TouchCancel(*mobile.TouchEvent) {
}

// TypedKey receives key input events when the Entry widget is focused.
//
// Implements: fyne.Focusable
func (e *Entry) TypedKey(key *fyne.KeyEvent) {
	if e.Disabled() {
		return
	}
	if e.cursorAnim != nil {
		e.cursorAnim.interrupt()
	}
	provider := e.textProvider()
	multiLine := e.MultiLine

	if e.selectKeyDown || e.sel.selecting {
		if e.selectingKeyHandler(key) {
			e.Refresh()
			return
		}
	}

	switch key.Name {
	case fyne.KeyBackspace:
		isEmpty := provider.len() == 0 || (e.CursorColumn == 0 && e.CursorRow == 0)
		if isEmpty {
			return
		}

		pos := e.cursorTextPos()
		deletedText := provider.deleteFromTo(pos-1, pos)
		e.CursorRow, e.CursorColumn = e.rowColFromTextPos(pos - 1)
		e.syncSelectable()
		e.undoStack.MergeOrAdd(&entryModifyAction{
			Delete:   true,
			Position: pos - 1,
			Text:     deletedText,
		})
	case fyne.KeyDelete:
		pos := e.cursorTextPos()
		if provider.len() == 0 || pos == provider.len() {
			return
		}

		deletedText := provider.deleteFromTo(pos, pos+1)
		e.undoStack.MergeOrAdd(&entryModifyAction{
			Delete:   true,
			Position: pos,
			Text:     deletedText,
		})
	case fyne.KeyReturn, fyne.KeyEnter:
		e.typedKeyReturn(provider, multiLine)
	case fyne.KeyTab:
		e.typedKeyTab()
	case fyne.KeyUp:
		e.typedKeyUp(provider)
	case fyne.KeyDown:
		e.typedKeyDown(provider)
	case fyne.KeyLeft:
		e.typedKeyLeft(provider)
	case fyne.KeyRight:
		e.typedKeyRight(provider)
	case fyne.KeyEnd:
		e.typedKeyEnd(provider)
	case fyne.KeyHome:
		e.typedKeyHome()
	case fyne.KeyPageUp:
		if e.MultiLine {
			e.CursorRow = 0
		}
		e.CursorColumn = 0
		e.syncSelectable()
	case fyne.KeyPageDown:
		if e.MultiLine {
			e.CursorRow = provider.rows() - 1
			e.CursorColumn = provider.rowLength(e.CursorRow)
		} else {
			e.CursorColumn = provider.len()
		}
		e.syncSelectable()
	default:
		return
	}

	content := provider.String()
	changed := e.updateText(content, false)
	if e.CursorRow == e.sel.selectRow && e.CursorColumn == e.sel.selectColumn {
		e.sel.selecting = false
	}
	cb := e.OnChanged
	if changed {
		e.validate()
		if cb != nil {
			cb(content)
		}
	}
	e.Refresh()
}

// Undo un-does the last modifying user-action.
//
// Since: 2.5
func (e *Entry) Undo() {
	newText, action := e.undoStack.Undo(e.Text)
	modify, ok := action.(*entryModifyAction)
	if !ok {
		return
	}
	pos := modify.Position
	if modify.Delete {
		pos += len(modify.Text)
	}
	e.updateText(newText, false)
	e.CursorRow, e.CursorColumn = e.rowColFromTextPos(pos)
	e.syncSelectable()
	if e.OnChanged != nil {
		e.OnChanged(newText)
	}
	e.Refresh()
}

func (e *Entry) typedKeyUp(provider *RichText) {
	if e.CursorRow > 0 {
		e.CursorRow--
	} else {
		e.CursorColumn = 0
	}

	rowLength := provider.rowLength(e.CursorRow)
	if e.CursorColumn > rowLength {
		e.CursorColumn = rowLength
	}
	e.syncSelectable()
}

func (e *Entry) typedKeyDown(provider *RichText) {
	rowLength := provider.rowLength(e.CursorRow)

	if e.CursorRow < provider.rows()-1 {
		e.CursorRow++
		rowLength = provider.rowLength(e.CursorRow)
	} else {
		e.CursorColumn = rowLength
	}

	if e.CursorColumn > rowLength {
		e.CursorColumn = rowLength
	}
	e.syncSelectable()
}

func (e *Entry) typedKeyLeft(provider *RichText) {
	if e.CursorColumn > 0 {
		e.CursorColumn--
	} else if e.MultiLine && e.CursorRow > 0 {
		e.CursorRow--
		e.CursorColumn = provider.rowLength(e.CursorRow)
	}
	e.syncSelectable()
}

func (e *Entry) typedKeyRight(provider *RichText) {
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
	e.syncSelectable()
}

func (e *Entry) typedKeyHome() {
	e.CursorColumn = 0
}

func (e *Entry) typedKeyEnd(provider *RichText) {
	if e.MultiLine {
		e.CursorColumn = provider.rowLength(e.CursorRow)
	} else {
		e.CursorColumn = provider.len()
	}
}

// handler for Ctrl+[backspace/delete] - delete the word
// to the left or right of the cursor
func (e *Entry) deleteWord(right bool) {
	provider := e.textProvider()
	cursorRow, cursorCol := e.CursorRow, e.CursorColumn

	// start, end relative to text row
	start, end := getTextWhitespaceRegion(provider.row(cursorRow), cursorCol, true)
	if right {
		start = cursorCol
	} else {
		end = cursorCol
	}
	if start == -1 || end == -1 {
		return
	}

	// convert start, end to absolute text position
	b := provider.rowBoundary(cursorRow)
	if b != nil {
		start += b.begin
		end += b.begin
	}

	erased := provider.deleteFromTo(start, end)
	e.undoStack.MergeOrAdd(&entryModifyAction{
		Delete:   true,
		Position: start,
		Text:     erased,
	})

	if !right {
		e.CursorColumn = cursorCol - (end - start)
	}
	e.updateTextAndRefresh(provider.String(), false)
}

func (e *Entry) typedKeyTab() {
	if dd, ok := fyne.CurrentApp().Driver().(desktop.Driver); ok {
		if dd.CurrentKeyModifiers()&fyne.KeyModifierShift != 0 {
			return // don't insert a tab when Shift+Tab typed
		}
	}
	e.TypedRune('\t')
}

// TypedRune receives text input events when the Entry widget is focused.
//
// Implements: fyne.Focusable
func (e *Entry) TypedRune(r rune) {
	if e.Disabled() {
		return
	}

	e.syncSelectable()
	if e.popUp != nil {
		e.popUp.Hide()
	}

	// if we've typed a character and we're selecting then replace the selection with the character
	cb := e.OnChanged
	if e.sel.selecting {
		e.eraseSelection()
	}

	runes := []rune{r}
	pos := e.cursorTextPos()

	provider := e.textProvider()
	provider.insertAt(pos, runes)

	content := provider.String()
	e.updateText(content, false)
	e.CursorRow, e.CursorColumn = e.rowColFromTextPos(pos + len(runes))
	e.syncSelectable()

	e.undoStack.MergeOrAdd(&entryModifyAction{
		Position: pos,
		Text:     runes,
	})

	e.validate()
	if cb != nil {
		cb(content)
	}
	e.Refresh()
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
	e.Validator = nil
	e.binder.Unbind()
}

// copyToClipboard copies the current selection to a given clipboard.
// This does nothing if it is a concealed entry.
func (e *Entry) copyToClipboard(clipboard fyne.Clipboard) {
	if !e.sel.selecting || e.Password {
		return
	}

	clipboard.SetContent(e.sel.SelectedText())
}

func (e *Entry) cursorTextPos() (pos int) {
	return textPosFromRowCol(e.CursorRow, e.CursorColumn, e.textProvider())
}

// cutToClipboard copies the current selection to a given clipboard and then removes the selected text.
// This does nothing if it is a concealed entry.
func (e *Entry) cutToClipboard(clipboard fyne.Clipboard) {
	if !e.sel.selecting || e.Password {
		return
	}

	e.copyToClipboard(clipboard)
	e.eraseSelectionAndUpdate()
	content := e.Text
	cb := e.OnChanged

	e.validate()
	if cb != nil {
		cb(content)
	}
	e.Refresh()
}

// eraseSelection deletes the selected text and moves the cursor but does not update the text field.
func (e *Entry) eraseSelection() bool {
	if e.Disabled() {
		return false
	}

	provider := e.textProvider()
	posA, posB := e.sel.selection()

	if posA == posB {
		return false
	}

	erasedText := provider.deleteFromTo(posA, posB)
	e.CursorRow, e.CursorColumn = e.rowColFromTextPos(posA)
	e.syncSelectable()
	e.sel.selectRow, e.sel.selectColumn = e.CursorRow, e.CursorColumn
	e.sel.selecting = false

	e.undoStack.MergeOrAdd(&entryModifyAction{
		Delete:   true,
		Position: posA,
		Text:     erasedText,
	})

	return true
}

// eraseSelectionAndUpdate removes the current selected region and moves the cursor.
// It also updates the text if something has been erased.
func (e *Entry) eraseSelectionAndUpdate() {
	if e.eraseSelection() {
		e.updateText(e.textProvider().String(), false)
	}
}

// pasteFromClipboard inserts text from the clipboard content,
// starting from the cursor position.
func (e *Entry) pasteFromClipboard(clipboard fyne.Clipboard) {
	e.syncSelectable()
	text := clipboard.Content()
	if text == "" {
		changed := e.sel.selecting && e.eraseSelection()

		if changed {
			e.Refresh()
		}

		return // Nothing to paste into the text content.
	}

	if !e.MultiLine {
		// format clipboard content to be compatible with single line entry
		text = strings.Replace(text, "\n", " ", -1)
	}

	if e.sel.selecting {
		e.eraseSelection()
	}

	runes := []rune(text)
	pos := e.cursorTextPos()
	provider := e.textProvider()
	provider.insertAt(pos, runes)

	e.undoStack.Add(&entryModifyAction{
		Position: pos,
		Text:     runes,
	})
	content := provider.String()
	e.updateText(content, false)
	e.CursorRow, e.CursorColumn = e.rowColFromTextPos(pos + len(runes))
	e.syncSelectable()
	cb := e.OnChanged

	e.validate()
	if cb != nil {
		cb(content) // We know that the text has changed.
	}

	e.Refresh() // placing the cursor (and refreshing) happens last
}

// placeholderProvider returns the placeholder text handler for this entry
func (e *Entry) placeholderProvider() *RichText {
	if len(e.placeholder.Segments) > 0 {
		return &e.placeholder
	}

	e.placeholder.Scroll = widget.ScrollNone
	e.placeholder.inset = fyne.NewSize(0, e.Theme().Size(theme.SizeNameInputBorder))

	style := RichTextStyleInline
	style.ColorName = theme.ColorNamePlaceHolder
	style.TextStyle = e.TextStyle

	e.placeholder.Segments = []RichTextSegment{
		&TextSegment{
			Style: style,
			Text:  e.PlaceHolder,
		},
	}

	return &e.placeholder
}

func (e *Entry) registerShortcut() {
	e.shortcut.AddShortcut(&fyne.ShortcutUndo{}, func(se fyne.Shortcut) {
		e.Undo()
	})
	e.shortcut.AddShortcut(&fyne.ShortcutRedo{}, func(se fyne.Shortcut) {
		e.Redo()
	})
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

	moveWord := func(s fyne.Shortcut) {
		row := e.textProvider().row(e.CursorRow)
		start, end := getTextWhitespaceRegion(row, e.CursorColumn, true)
		if start == -1 || end == -1 {
			return
		}

		e.setFieldsAndRefresh(func() {
			if s.(*desktop.CustomShortcut).KeyName == fyne.KeyLeft {
				if e.CursorColumn == 0 {
					if e.CursorRow > 0 {
						e.CursorRow--
						e.CursorColumn = len(e.textProvider().row(e.CursorRow))
					}
				} else {
					e.CursorColumn = start
				}
			} else {
				if e.CursorColumn == len(e.textProvider().row(e.CursorRow)) {
					if e.CursorRow < e.textProvider().rows()-1 {
						e.CursorRow++
						e.CursorColumn = 0
					}
				} else {
					e.CursorColumn = end
				}
			}
			e.syncSelectable()
		})
	}
	selectMoveWord := func(se fyne.Shortcut) {
		if !e.sel.selecting {
			e.sel.selectColumn = e.CursorColumn
			e.sel.selectRow = e.CursorRow
			e.sel.selecting = true
		}
		moveWord(se)
	}
	unselectMoveWord := func(se fyne.Shortcut) {
		e.sel.selecting = false
		moveWord(se)
	}

	moveWordModifier := fyne.KeyModifierShortcutDefault
	if runtime.GOOS == "darwin" {
		moveWordModifier = fyne.KeyModifierAlt

		// Cmd+left, Cmd+right shortcuts behave like Home and End keys on Mac OS
		shortcutHomeEnd := func(s fyne.Shortcut) {
			e.sel.selecting = false
			if s.(*desktop.CustomShortcut).KeyName == fyne.KeyLeft {
				e.typedKeyHome()
			} else {
				e.typedKeyEnd(e.textProvider())
			}
			e.Refresh()
		}
		e.shortcut.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyLeft, Modifier: fyne.KeyModifierSuper}, shortcutHomeEnd)
		e.shortcut.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyRight, Modifier: fyne.KeyModifierSuper}, shortcutHomeEnd)
	}

	e.shortcut.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyLeft, Modifier: moveWordModifier}, unselectMoveWord)
	e.shortcut.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyLeft, Modifier: moveWordModifier | fyne.KeyModifierShift}, selectMoveWord)
	e.shortcut.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyRight, Modifier: moveWordModifier}, unselectMoveWord)
	e.shortcut.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyRight, Modifier: moveWordModifier | fyne.KeyModifierShift}, selectMoveWord)

	e.shortcut.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyBackspace, Modifier: moveWordModifier},
		func(fyne.Shortcut) { e.deleteWord(false) })
	e.shortcut.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyDelete, Modifier: moveWordModifier},
		func(fyne.Shortcut) { e.deleteWord(true) })
}

func (e *Entry) requestFocus() {
	impl := e.super()
	if c := fyne.CurrentApp().Driver().CanvasForObject(impl); c != nil {
		c.Focus(impl.(fyne.Focusable))
	}
}

// Obtains row,col from a given textual position
// expects a read or write lock to be held by the caller
func (e *Entry) rowColFromTextPos(pos int) (row int, col int) {
	provider := e.textProvider()
	canWrap := e.Wrapping == fyne.TextWrapBreak || e.Wrapping == fyne.TextWrapWord
	totalRows := provider.rows()
	for i := 0; i < totalRows; i++ {
		b := provider.rowBoundary(i)
		if b == nil {
			continue
		}
		if b.begin <= pos {
			if b.end < pos {
				row++
			}
			col = pos - b.begin
			// if this gap is at `pos` and is a line wrap, increment (safe to access boundary i-1)
			if canWrap && b.begin == pos && pos != 0 && provider.rowBoundary(i-1).end == b.begin && row < (totalRows-1) {
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
		e.sel.selectRow = 0
		e.sel.selectColumn = 0

		lastRow := e.textProvider().rows() - 1
		e.CursorColumn = e.textProvider().rowLength(lastRow)
		e.CursorRow = lastRow
		e.syncSelectable()
		e.sel.selecting = true
	})
}

// selectingKeyHandler performs keypress action in the scenario that a selection
// is either a) in progress or b) about to start
// returns true if the keypress has been fully handled
func (e *Entry) selectingKeyHandler(key *fyne.KeyEvent) bool {

	if e.selectKeyDown && !e.sel.selecting {
		switch key.Name {
		case fyne.KeyUp, fyne.KeyDown,
			fyne.KeyLeft, fyne.KeyRight,
			fyne.KeyEnd, fyne.KeyHome,
			fyne.KeyPageUp, fyne.KeyPageDown:
			e.sel.selecting = true
		}
	}

	if !e.sel.selecting {
		return false
	}

	switch key.Name {
	case fyne.KeyBackspace, fyne.KeyDelete:
		// clears the selection -- return handled
		e.eraseSelectionAndUpdate()
		content := e.Text
		cb := e.OnChanged

		e.validate()
		if cb != nil {
			cb(content)
		}
		e.Refresh()
		return true
	case fyne.KeyReturn, fyne.KeyEnter:
		if e.MultiLine {
			// clear the selection -- return unhandled to add the newline
			e.setFieldsAndRefresh(e.eraseSelectionAndUpdate)
		}
		return false
	}

	if !e.selectKeyDown {
		switch key.Name {
		case fyne.KeyLeft:
			// seek to the start of the selection -- return handled
			selectStart, _ := e.sel.selection()
			e.CursorRow, e.CursorColumn = e.rowColFromTextPos(selectStart)
			e.syncSelectable()
			e.sel.selecting = false
			return true
		case fyne.KeyRight:
			// seek to the end of the selection -- return handled
			_, selectEnd := e.sel.selection()
			e.CursorRow, e.CursorColumn = e.rowColFromTextPos(selectEnd)
			e.syncSelectable()
			e.sel.selecting = false
			return true
		case fyne.KeyUp, fyne.KeyDown, fyne.KeyEnd, fyne.KeyHome, fyne.KeyPageUp, fyne.KeyPageDown:
			// cursor movement without left or right shift -- clear selection and return unhandled
			e.sel.selecting = false
			return false
		}
	}

	return false
}

func (e *Entry) syncSegments() {
	colName := theme.ColorNameForeground
	wrap := e.textWrap()
	disabled := e.Disabled()
	if disabled {
		colName = theme.ColorNameDisabled
	}

	text := e.textProvider()
	text.Wrapping = wrap

	textSegment := text.Segments[0].(*TextSegment)
	textSegment.Text = e.Text
	textSegment.Style.ColorName = colName
	textSegment.Style.concealed = e.Password
	textSegment.Style.TextStyle = e.TextStyle

	colName = theme.ColorNamePlaceHolder
	if disabled {
		colName = theme.ColorNameDisabled
	}

	placeholder := e.placeholderProvider()
	placeholder.Wrapping = wrap

	textSegment = placeholder.Segments[0].(*TextSegment)
	textSegment.Style.ColorName = colName
	textSegment.Style.TextStyle = e.TextStyle
	textSegment.Text = e.PlaceHolder
}

func (e *Entry) syncSelectable() {
	if e.sel == nil {
		e.sel = &selectable{theme: e.Theme(), provider: e.textProvider(), focus: e, password: e.Password, style: e.TextStyle}
		e.sel.ExtendBaseWidget(e.sel)
	}

	e.sel.cursorRow, e.sel.cursorColumn = e.CursorRow, e.CursorColumn
}

// textProvider returns the text handler for this entry
func (e *Entry) textProvider() *RichText {
	if len(e.text.Segments) > 0 {
		return &e.text
	}

	if e.Text != "" {
		e.dirty = true
	}

	e.text.Scroll = widget.ScrollNone
	e.text.inset = fyne.NewSize(0, e.Theme().Size(theme.SizeNameInputBorder))
	e.text.Segments = []RichTextSegment{&TextSegment{Style: RichTextStyleInline, Text: e.Text}}
	return &e.text
}

// textWrap calculates the wrapping that we should apply.
func (e *Entry) textWrap() fyne.TextWrap {
	if e.Wrapping == fyne.TextWrap(fyne.TextTruncateClip) { // this is now the default - but we scroll around this large content
		return fyne.TextWrapOff
	}

	if !e.MultiLine && (e.Wrapping == fyne.TextWrapBreak || e.Wrapping == fyne.TextWrapWord) {
		fyne.LogError("Entry cannot wrap single line", nil)
		e.Wrapping = fyne.TextWrap(fyne.TextTruncateClip)
		return fyne.TextWrapOff
	}
	return e.Wrapping
}

func (e *Entry) updateCursorAndSelection() {
	e.CursorRow, e.CursorColumn = e.truncatePosition(e.CursorRow, e.CursorColumn)

	e.syncSelectable()
	e.sel.selectRow, e.sel.selectColumn = e.truncatePosition(e.sel.selectRow, e.sel.selectColumn)
}

func (e *Entry) updateFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	textSource, ok := data.(binding.String)
	if !ok {
		return
	}

	val, err := textSource.Get()
	e.conversionError = err
	e.validate()
	if err != nil {
		return
	}
	e.setText(val, true)
}

func (e *Entry) truncatePosition(row, col int) (int, int) {
	if e.Text == "" {
		return 0, 0
	}
	newRow := row
	newCol := col
	if row >= e.textProvider().rows() {
		newRow = e.textProvider().rows() - 1
	}
	rowLength := e.textProvider().rowLength(newRow)
	if (newCol >= rowLength) || (newRow < row) {
		newCol = rowLength
	}
	return newRow, newCol
}

func (e *Entry) updateMousePointer(p fyne.Position, rightClick bool) {
	row, col := e.sel.getRowCol(p)

	if !rightClick || !e.sel.selecting {
		e.CursorRow = row
		e.CursorColumn = col

		e.syncSelectable()
	}

	if !e.sel.selecting {
		e.sel.selectRow = row
		e.sel.selectColumn = col
	}

	r := cache.Renderer(e.content)
	if r != nil {
		r.(*entryContentRenderer).moveCursor()
	}
}

// updateText updates the internal text to the given value.
// It assumes that a lock exists on the widget.
func (e *Entry) updateText(text string, fromBinding bool) bool {
	changed := e.Text != text
	e.Text = text
	e.syncSegments()
	e.text.updateRowBounds()

	if e.Text != "" {
		e.dirty = true
	}

	if changed && !fromBinding {
		if e.binder.dataListenerPair.listener != nil {
			e.binder.SetCallback(nil)
			e.binder.CallWithData(e.writeData)
			e.binder.SetCallback(e.updateFromData)
		}
	}
	return changed
}

// updateTextAndRefresh updates the internal text to the given value then refreshes it.
// This should not be called under a property lock
func (e *Entry) updateTextAndRefresh(text string, fromBinding bool) {
	var callback func(string)

	changed := e.updateText(text, fromBinding)

	if changed {
		callback = e.OnChanged
	}

	e.validate()
	if callback != nil {
		callback(text)
	}
	e.Refresh()
}

func (e *Entry) writeData(data binding.DataItem) {
	if data == nil {
		return
	}
	textTarget, ok := data.(binding.String)
	if !ok {
		return
	}
	curValue, err := textTarget.Get()
	if err == nil && curValue == e.Text {
		e.conversionError = nil
		return
	}
	e.conversionError = textTarget.Set(e.Text)
}

func (e *Entry) typedKeyReturn(provider *RichText, multiLine bool) {
	onSubmitted := e.OnSubmitted
	selectDown := e.selectKeyDown
	text := e.Text

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
	s := []rune("\n")
	pos := e.cursorTextPos()
	provider.insertAt(pos, s)
	e.undoStack.MergeOrAdd(&entryModifyAction{
		Position: pos,
		Text:     s,
	})
	e.CursorColumn = 0
	e.CursorRow++
	e.syncSelectable()
}

func (e *Entry) setFieldsAndRefresh(f func()) {
	f()

	impl := e.super()
	if impl == nil {
		return
	}
	impl.Refresh()
}

var _ fyne.WidgetRenderer = (*entryRenderer)(nil)

type entryRenderer struct {
	box, border *canvas.Rectangle
	scroll      *widget.Scroll

	objects []fyne.CanvasObject
	entry   *Entry
}

func (r *entryRenderer) Destroy() {
}

func (r *entryRenderer) trailingInset() float32 {
	th := r.entry.Theme()
	xInset := float32(0)

	if r.entry.ActionItem != nil {
		xInset = r.entry.ActionItem.MinSize().Width
	}

	if r.entry.Validator != nil {
		iconSpace := th.Size(theme.SizeNameInlineIcon) + th.Size(theme.SizeNameLineSpacing)
		if r.entry.ActionItem == nil {
			xInset = iconSpace + th.Size(theme.SizeNameInnerPadding)
		} else {
			xInset += iconSpace
		}
	}

	return xInset
}

func (r *entryRenderer) Layout(size fyne.Size) {
	th := r.entry.Theme()
	borderSize := th.Size(theme.SizeNameInputBorder)
	iconSize := th.Size(theme.SizeNameInlineIcon)
	innerPad := th.Size(theme.SizeNameInnerPadding)
	inputBorder := th.Size(theme.SizeNameInputBorder)

	// 0.5 is removed so on low DPI it rounds down on the trailing edge
	r.border.Resize(fyne.NewSize(size.Width-borderSize-.5, size.Height-borderSize-.5))
	r.border.StrokeWidth = borderSize
	r.border.Move(fyne.NewSquareOffsetPos(borderSize / 2))
	r.box.Resize(size.Subtract(fyne.NewSquareSize(borderSize * 2)))
	r.box.Move(fyne.NewSquareOffsetPos(borderSize))

	pad := theme.InputBorderSize()
	actionIconSize := fyne.NewSize(0, size.Height-pad*2)
	if r.entry.ActionItem != nil {
		actionIconSize.Width = r.entry.ActionItem.MinSize().Width
		r.entry.ActionItem.Resize(actionIconSize)
		r.entry.ActionItem.Move(fyne.NewPos(size.Width-actionIconSize.Width-pad, pad))
	}

	validatorIconSize := fyne.NewSize(0, 0)
	if r.entry.Validator != nil {
		validatorIconSize = fyne.NewSquareSize(iconSize)

		r.ensureValidationSetup()
		r.entry.validationStatus.Resize(validatorIconSize)

		if r.entry.ActionItem == nil {
			r.entry.validationStatus.Move(fyne.NewPos(size.Width-validatorIconSize.Width-innerPad, innerPad))
		} else {
			r.entry.validationStatus.Move(fyne.NewPos(size.Width-validatorIconSize.Width-actionIconSize.Width, innerPad))
		}
	}

	r.entry.textProvider().inset = fyne.NewSize(0, inputBorder)
	r.entry.placeholderProvider().inset = fyne.NewSize(0, inputBorder)
	entrySize := size.Subtract(fyne.NewSize(r.trailingInset(), inputBorder*2))
	entryPos := fyne.NewPos(0, inputBorder)

	prov := r.entry.textProvider()
	textPos := textPosFromRowCol(r.entry.CursorRow, r.entry.CursorColumn, prov)
	selectPos := textPosFromRowCol(r.entry.sel.selectRow, r.entry.sel.selectColumn, prov)
	if r.entry.Wrapping == fyne.TextWrapOff && r.entry.Scroll == widget.ScrollNone {
		r.entry.content.Resize(entrySize)
		r.entry.content.Move(entryPos)
	} else {
		r.scroll.Resize(entrySize)
		r.scroll.Move(entryPos)
	}

	resizedTextPos := textPosFromRowCol(r.entry.CursorRow, r.entry.CursorColumn, prov)
	if textPos != resizedTextPos {
		r.entry.setFieldsAndRefresh(func() {
			r.entry.CursorRow, r.entry.CursorColumn = r.entry.rowColFromTextPos(textPos)
			r.entry.sel.cursorRow, r.entry.sel.cursorRow = r.entry.CursorRow, r.entry.CursorColumn

			if r.entry.sel.selecting {
				r.entry.sel.selectRow, r.entry.sel.selectColumn = r.entry.rowColFromTextPos(selectPos)
			}
		})
	}
}

// MinSize calculates the minimum size of an entry widget.
// This is based on the contained text with a standard amount of padding added.
// If MultiLine is true then we will reserve space for at leasts 3 lines
func (r *entryRenderer) MinSize() fyne.Size {
	if rend := cache.Renderer(r.entry.content); rend != nil {
		rend.(*entryContentRenderer).updateScrollDirections()
	}

	th := r.entry.Theme()
	minSize := fyne.Size{}

	if r.scroll.Direction == widget.ScrollNone {
		minSize = r.entry.content.MinSize().AddWidthHeight(0, th.Size(theme.SizeNameInputBorder)*2)
	} else {
		innerPadding := th.Size(theme.SizeNameInnerPadding)
		textSize := th.Size(theme.SizeNameText)
		charMin := r.entry.placeholderProvider().charMinSize(r.entry.Password, r.entry.TextStyle, textSize)
		minSize = charMin.Add(fyne.NewSquareSize(innerPadding))

		if r.entry.MultiLine {
			count := r.entry.multiLineRows
			if count <= 0 {
				count = multiLineRows
			}

			minSize.Height = charMin.Height*float32(count) + innerPadding
		}

		minSize = minSize.AddWidthHeight(innerPadding*2, innerPadding)
	}

	iconSpace := th.Size(theme.SizeNameInlineIcon) + th.Size(theme.SizeNameLineSpacing)
	if r.entry.ActionItem != nil {
		minSize.Width += iconSpace
	}
	if r.entry.Validator != nil {
		minSize.Width += iconSpace
	}

	return minSize
}

func (r *entryRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *entryRenderer) Refresh() {
	content := r.entry.content
	focusedAppearance := r.entry.focused && !r.entry.Disabled()
	scroll := r.entry.Scroll
	wrapping := r.entry.Wrapping

	r.entry.syncSegments()
	r.entry.text.updateRowBounds()
	r.entry.placeholder.updateRowBounds()
	r.entry.text.Refresh()
	r.entry.placeholder.Refresh()

	th := r.entry.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	inputBorder := th.Size(theme.SizeNameInputBorder)

	// correct our scroll wrappers if the wrap mode changed
	entrySize := r.entry.Size().Subtract(fyne.NewSize(r.trailingInset(), inputBorder*2))
	if wrapping == fyne.TextWrapOff && scroll == widget.ScrollNone && r.scroll.Content != nil {
		r.scroll.Hide()
		r.scroll.Content = nil
		content.Move(fyne.NewPos(0, inputBorder))
		content.Resize(entrySize)

		for i, o := range r.objects {
			if o == r.scroll {
				r.objects[i] = content
				break
			}
		}
	} else if (wrapping != fyne.TextWrapOff || scroll != widget.ScrollNone) && r.scroll.Content == nil {
		r.scroll.Content = content
		content.Move(fyne.NewPos(0, 0))
		r.scroll.Move(fyne.NewPos(0, inputBorder))
		r.scroll.Resize(entrySize)
		r.scroll.Show()

		for i, o := range r.objects {
			if o == content {
				r.objects[i] = r.scroll
				break
			}
		}
	}
	r.entry.updateCursorAndSelection()

	r.box.FillColor = th.Color(theme.ColorNameInputBackground, v)
	r.box.CornerRadius = th.Size(theme.SizeNameInputRadius)
	r.border.CornerRadius = r.box.CornerRadius
	if focusedAppearance {
		r.border.StrokeColor = th.Color(theme.ColorNamePrimary, v)
	} else {
		if r.entry.Disabled() {
			r.border.StrokeColor = th.Color(theme.ColorNameDisabled, v)
		} else {
			r.border.StrokeColor = th.Color(theme.ColorNameInputBorder, v)
		}
	}
	if r.entry.ActionItem != nil {
		r.entry.ActionItem.Refresh()
	}

	if r.entry.Validator != nil {
		if !r.entry.focused && !r.entry.Disabled() && r.entry.dirty && r.entry.validationError != nil {
			r.border.StrokeColor = th.Color(theme.ColorNameError, v)
		}
		r.ensureValidationSetup()
		r.entry.validationStatus.Refresh()
	} else if r.entry.validationStatus != nil {
		r.entry.validationStatus.Hide()
	}

	r.entry.sel.Hidden = !r.entry.focused

	cache.Renderer(r.entry.content).Refresh()
	canvas.Refresh(r.entry.super())
}

func (r *entryRenderer) ensureValidationSetup() {
	if r.entry.validationStatus == nil {
		r.entry.validationStatus = newValidationStatus(r.entry)
		r.objects = append(r.objects, r.entry.validationStatus)
		r.Layout(r.entry.Size())

		r.entry.validate()
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

	provider := e.entry.textProvider()
	placeholder := e.entry.placeholderProvider()
	if provider.len() != 0 {
		placeholder.Hide()
	}
	objects := []fyne.CanvasObject{placeholder, provider, e.entry.cursorAnim.cursor}

	r := &entryContentRenderer{e.entry.cursorAnim.cursor, objects,
		provider, placeholder, e}
	r.updateScrollDirections()
	r.Layout(e.Size())
	return r
}

// DragEnd is called at end of a drag event.
//
// Implements: fyne.Draggable
func (e *entryContent) DragEnd() {
	// we need to propagate the focus, top level widget handles focus APIs
	e.entry.requestFocus()

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
	cursor  *canvas.Rectangle
	objects []fyne.CanvasObject

	provider, placeholder *RichText
	content               *entryContent
}

func (r *entryContentRenderer) Destroy() {
	r.content.entry.cursorAnim.stop()
}

func (r *entryContentRenderer) Layout(size fyne.Size) {
	r.provider.Resize(size)
	r.placeholder.Resize(size)
}

func (r *entryContentRenderer) MinSize() fyne.Size {
	r.content.Theme() // setup theme cache before locking
	minSize := r.content.entry.placeholderProvider().MinSize()

	if r.content.entry.textProvider().len() > 0 {
		minSize = r.content.entry.text.MinSize()
	}

	return minSize
}

func (r *entryContentRenderer) Objects() []fyne.CanvasObject {
	// Objects are generated dynamically force selection rectangles to appear underneath the text
	if r.content.entry.sel.selecting {
		return append([]fyne.CanvasObject{r.content.entry.sel}, r.objects...)
	}
	return r.objects
}

func (r *entryContentRenderer) Refresh() {
	provider := r.content.entry.textProvider()
	placeholder := r.content.entry.placeholderProvider()
	focused := r.content.entry.focused
	focusedAppearance := focused && !r.content.entry.Disabled()
	r.updateScrollDirections()

	if provider.len() == 0 {
		placeholder.Show()
	} else if placeholder.Visible() {
		placeholder.Hide()
	}

	th := r.content.entry.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	if focusedAppearance {
		if fyne.CurrentApp().Settings().ShowAnimations() {
			r.content.entry.cursorAnim.start()
		} else {
			r.cursor.FillColor = th.Color(theme.ColorNamePrimary, v)
		}
		r.cursor.Show()
	} else {
		r.content.entry.cursorAnim.stop()
		r.cursor.Hide()
	}
	r.moveCursor()

	canvas.Refresh(r.content)
}

func (r *entryContentRenderer) ensureCursorVisible() {
	th := r.content.entry.Theme()
	lineSpace := th.Size(theme.SizeNameLineSpacing)

	letter := fyne.MeasureText("e", th.Size(theme.SizeNameText), r.content.entry.TextStyle)
	padX := letter.Width*2 + lineSpace
	padY := letter.Height - lineSpace
	cx := r.cursor.Position().X
	cy := r.cursor.Position().Y
	cx1 := cx - padX
	cy1 := cy - padY
	cx2 := cx + r.cursor.Size().Width + padX
	cy2 := cy + r.cursor.Size().Height + padY
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
	} else if cy2 >= offset.Y+size.Height {
		move.DY += cy2 - (offset.Y + size.Height)
	}
	if r.content.scroll.Content != nil {
		r.content.scroll.ScrollToOffset(r.content.scroll.Offset.Add(move))
	}
}

func (r *entryContentRenderer) moveCursor() {
	// build r.selection[] if the user has made a selection
	r.content.entry.sel.Refresh()

	th := r.content.entry.Theme()
	textSize := th.Size(theme.SizeNameText)
	provider := r.content.entry.textProvider()
	innerPad := th.Size(theme.SizeNameInnerPadding)
	inputBorder := th.Size(theme.SizeNameInputBorder)
	size := provider.lineSizeToColumn(r.content.entry.CursorColumn, r.content.entry.CursorRow, textSize, innerPad)
	xPos := size.Width
	yPos := size.Height * float32(r.content.entry.CursorRow)

	lineHeight := r.content.entry.text.charMinSize(r.content.entry.Password, r.content.entry.TextStyle, textSize).Height
	r.cursor.Resize(fyne.NewSize(inputBorder, lineHeight))
	r.cursor.Move(fyne.NewPos(xPos-(inputBorder/2), yPos+innerPad-inputBorder))

	callback := r.content.entry.OnCursorChanged
	r.ensureCursorVisible()

	if callback != nil {
		callback()
	}
}

func (r *entryContentRenderer) updateScrollDirections() {
	if r.content.scroll == nil { // not scrolling
		return
	}

	switch r.content.entry.Wrapping {
	case fyne.TextWrapOff:
		r.content.scroll.Direction = r.content.entry.Scroll
	case fyne.TextWrap(fyne.TextTruncateClip): // this is now the default - but we scroll
		r.content.scroll.Direction = widget.ScrollBoth
	default: // fyne.TextWrapBreak, fyne.TextWrapWord
		r.content.scroll.Direction = widget.ScrollVerticalOnly
	}
}

// getTextWhitespaceRegion returns the start/end markers for selection highlight on starting from col
// and expanding to the start and end of the whitespace or text underneath the specified position.
// Pass `true` for `expand` if you want whitespace selection to extend to the neighboring words.
func getTextWhitespaceRegion(row []rune, col int, expand bool) (int, int) {
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
		// If this rune is a typical word separator then classify it as whitespace
		if isWordSeparator(r) {
			return ' '
		}
		return '-'
	}
	toks := strings.Map(space, string(row))
	c := byte(' ')

	startCheck := col
	endCheck := col
	if expand {
		if col > 0 && toks[col-1] == ' ' { // ignore the prior whitespace then count
			startCheck = strings.LastIndexByte(toks[:startCheck], '-')
			if startCheck == -1 {
				startCheck = 0
			}
		}
		if toks[col] == ' ' { // ignore the current whitespace then count
			endCheck = col + strings.IndexByte(toks[endCheck:], '-')
		}
	} else if toks[col] == ' ' {
		c = byte('-')
	}

	// LastIndexByte + 1 ensures that the position of the unwanted character ' ' is excluded
	// +1 also has the added side effect whereby if ' ' isn't found then -1 is snapped to 0
	start := strings.LastIndexByte(toks[:startCheck], c) + 1

	// IndexByte will find the position of the next unwanted character, this is to be the end
	// marker for the selection
	end := -1
	if endCheck != -1 {
		end = strings.IndexByte(toks[endCheck:], c)
	}

	if end == -1 {
		end = len(toks) // snap end to len(toks) if it results in -1
	} else {
		end += endCheck // otherwise include the text slice position
	}
	return start, end
}

func isWordSeparator(r rune) bool {
	return unicode.IsSpace(r) ||
		strings.ContainsRune("`~!@#$%^&*()-=+[{]}\\|;:'\",.<>/?", r)
}

// entryUndoAction represents a single user action that can be undone
type entryUndoAction interface {
	Undo(string) string
	Redo(string) string
}

// entryMergeableUndoAction is like entryUndoAction, but the undoStack
// can try to merge it with the next action (see TryMerge).
// This is useful because it allows grouping together actions like
// entering every single characters in a word. We don't want to have to
// undo every single character addition.
type entryMergeableUndoAction interface {
	entryUndoAction
	// TryMerge attempts to merge the current action
	// with the next action. It returns true if successful.
	// If it fails, the undoStack will simply add the next
	// item without merging.
	TryMerge(next entryMergeableUndoAction) bool
}

// Declare conformity with entryMergeableUndoAction interface
var _ entryMergeableUndoAction = (*entryModifyAction)(nil)

// entryModifyAction implements entryMergeableUndoAction.
// It represents the insertion/deletion of a single string at a
// position (e.g. "Hello" => "Hello, world", or "Hello" => "He").
type entryModifyAction struct {
	// Delete is true if this action deletes Text, and false if it inserts Text
	Delete bool
	// Position represents the start position of Text
	Position int
	// Text is the text that is inserted or deleted at Position
	Text []rune
}

func (i *entryModifyAction) Undo(s string) string {
	if i.Delete {
		return i.add(s)
	} else {
		return i.sub(s)
	}
}

func (i *entryModifyAction) Redo(s string) string {
	if i.Delete {
		return i.sub(s)
	} else {
		return i.add(s)
	}
}

// Inserts Text
func (i *entryModifyAction) add(s string) string {
	runes := []rune(s)
	return string(runes[:i.Position]) + string(i.Text) + string(runes[i.Position:])
}

// Deletes Text
func (i *entryModifyAction) sub(s string) string {
	runes := []rune(s)
	return string(runes[:i.Position]) + string(runes[i.Position+len(i.Text):])
}

func (i *entryModifyAction) TryMerge(other entryMergeableUndoAction) bool {
	if other, ok := other.(*entryModifyAction); ok {
		// Don't merge two different types of modifyAction
		if i.Delete != other.Delete {
			return false
		}

		// Don't merge two separate words
		wordSeparators := func(s []rune) (num int, onlyWordSeparators bool) {
			onlyWordSeparators = true
			for _, r := range s {
				if isWordSeparator(r) {
					num++
					onlyWordSeparators = false
				}
			}
			return
		}
		selfNumWS, _ := wordSeparators(i.Text)
		otherNumWS, otherOnlyWS := wordSeparators(other.Text)
		if !((selfNumWS == 0 && otherNumWS == 0) ||
			(selfNumWS > 0 && otherOnlyWS)) {
			return false
		}

		if i.Delete {
			if i.Position == other.Position+len(other.Text) {
				i.Position = other.Position
				i.Text = append(other.Text, i.Text...)
				return true
			}
		} else {
			if i.Position+len(i.Text) == other.Position {
				i.Text = append(i.Text, other.Text...)
				return true
			}
		}
		return false
	}
	return false
}

// entryUndoStack stores the information necessary for textual undo/redo functionality.
type entryUndoStack struct {
	// items is the stack for storing the history of user actions.
	items []entryUndoAction
	// index is the size of the current effective undo stack.
	// items[index-1] and below are the possible undo actions.
	// items[index] and above are the possible redo actions.
	index int
}

// Applies the undo action to s and returns the result along with the action performed
func (u *entryUndoStack) Undo(s string) (newS string, action entryUndoAction) {
	if !u.CanUndo() {
		return s, nil
	}
	u.index--
	action = u.items[u.index]
	return action.Undo(s), action
}

// Applies the redo action to s and returns the result along with the action performed
func (u *entryUndoStack) Redo(s string) (newS string, action entryUndoAction) {
	if !u.CanRedo() {
		return s, nil
	}
	action = u.items[u.index]
	res := action.Redo(s)
	u.index++
	return res, action
}

// Returns true if an undo action is available
func (u *entryUndoStack) CanUndo() bool {
	return u.index != 0
}

// Returns true if an redo action is available
func (u *entryUndoStack) CanRedo() bool {
	return u.index != len(u.items)
}

// Adds the action to the stack, which can later be undone by calling Undo()
func (u *entryUndoStack) Add(a entryUndoAction) {
	u.items = u.items[:u.index]
	u.items = append(u.items, a)
	u.index++
}

// Tries to merge the action with the last item on the undo stack.
// If it can't be merged, it calls Add().
func (u *entryUndoStack) MergeOrAdd(a entryUndoAction) {
	u.items = u.items[:u.index]
	if u.index == 0 {
		u.Add(a)
		return
	}
	ma, ok := a.(entryMergeableUndoAction)
	if !ok {
		u.Add(a)
		return
	}
	mprev, ok := u.items[u.index-1].(entryMergeableUndoAction)
	if !ok {
		u.Add(a)
		return
	}
	if !mprev.TryMerge(ma) {
		u.Add(a)
		return
	}
}

// Removes all items from the undo stack
func (u *entryUndoStack) Clear() {
	u.items = nil
	u.index = 0
}

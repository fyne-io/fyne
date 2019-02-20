package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

const (
	multiLineRows = 3
)

type entryRenderer struct {
	text        *textProvider
	placeholder *textProvider
	box, cursor *canvas.Rectangle

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

func (e *entryRenderer) moveCursor() {
	textRenderer := Renderer(e.text).(*textRenderer)
	size := textRenderer.lineSize(e.entry.CursorColumn, e.entry.CursorRow)
	xPos := size.Width
	yPos := size.Height * e.entry.CursorRow

	lineHeight := e.text.charMinSize().Height
	e.cursor.Resize(fyne.NewSize(2, lineHeight))
	e.cursor.Move(fyne.NewPos(xPos+theme.Padding()*2, yPos+theme.Padding()*2))

	canvas.Refresh(e.cursor)
}

// Layout the components of the entry widget.
func (e *entryRenderer) Layout(size fyne.Size) {
	e.box.Resize(size.Subtract(fyne.NewSize(theme.Padding(), theme.Padding())))
	e.box.Move(fyne.NewPos(theme.Padding()/2, theme.Padding()/2))

	e.text.Resize(size.Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
	e.text.Move(fyne.NewPos(theme.Padding(), theme.Padding()))

	e.placeholder.Resize(size.Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
	e.placeholder.Move(fyne.NewPos(theme.Padding(), theme.Padding()))

	e.moveCursor()
}

// ApplyTheme is called when the Entry may need to update it's look.
func (e *entryRenderer) ApplyTheme() {
	Renderer(e.text).ApplyTheme()
	e.box.FillColor = theme.BackgroundColor()
	e.Refresh()
}

func (e *entryRenderer) BackgroundColor() color.Color {
	if e.entry.focused {
		return theme.FocusColor()
	}

	return theme.ButtonColor()
}

func (e *entryRenderer) Refresh() {
	e.placeholder.refreshTextRenderer()
	e.placeholder.Hide()
	if e.text.len() == 0 && e.entry.Visible() {
		e.placeholder.Show()
	}

	e.text.refreshTextRenderer()
	if e.entry.focused {
		e.cursor.FillColor = theme.FocusColor()
	} else {
		e.cursor.FillColor = color.RGBA{0, 0, 0, 0}
	}

	canvas.Refresh(e.entry)
}

func (e *entryRenderer) Objects() []fyne.CanvasObject {
	return e.objects
}

// Entry widget allows simple text to be input when focused.
type Entry struct {
	baseWidget

	Text        string
	PlaceHolder string
	OnChanged   func(string) `json:"-"`
	Password    bool
	ReadOnly    bool
	MultiLine   bool

	CursorRow, CursorColumn int

	focused bool
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (e *Entry) Resize(size fyne.Size) {
	e.resize(size, e)
}

// Move the widget to a new position, relative to it's parent.
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
	e.hide(e)
}

// SetText manually sets the text of the Entry to the given text value.
func (e *Entry) SetText(text string) {
	e.textProvider().SetText(text)
	e.updateText(text)
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
	e.Text = text
	if e.OnChanged != nil {
		e.OnChanged(text)
	}

	Refresh(e)
}

func (e *Entry) cursorTextPos() int {
	pos := 0
	provider := e.textProvider()
	for i := 0; i < e.CursorRow; i++ {
		rowLength := provider.rowLength(i)
		pos += rowLength + 1
	}
	pos += e.CursorColumn
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

// TypedRune receives text input events when the Entry widget is focused.
func (e *Entry) TypedRune(r rune) {
	if e.ReadOnly {
		return
	}
	provider := e.textProvider()

	runes := []rune{r}
	provider.insertAt(e.cursorTextPos(), runes)
	e.CursorColumn += len(runes)

	e.updateText(provider.String())
	Renderer(e).(*entryRenderer).moveCursor()
}

// TypedKey receives key input events when the Entry widget is focused.
func (e *Entry) TypedKey(key *fyne.KeyEvent) {
	if e.ReadOnly {
		return
	}
	provider := e.textProvider()
	switch key.Name {
	case fyne.KeyBackspace:
		if provider.len() == 0 || (e.CursorColumn == 0 && e.CursorRow == 0) {
			return
		}
		pos := e.cursorTextPos()
		deleted := provider.deleteFromTo(pos-1, pos)
		if deleted[0] == '\n' {
			e.CursorRow--
			rowLength := provider.rowLength(e.CursorRow)
			e.CursorColumn = rowLength
			break
		}
		e.CursorColumn--
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
		e.CursorColumn = 0
		e.CursorRow++
	case fyne.KeyUp:
		if !e.MultiLine {
			return
		}

		if e.CursorRow > 0 {
			e.CursorRow--
		}

		rowLength := provider.rowLength(e.CursorRow)
		if e.CursorColumn > rowLength {
			e.CursorColumn = rowLength
		}
	case fyne.KeyDown:
		if !e.MultiLine {
			return
		}

		if e.CursorRow < provider.rows()-1 {
			e.CursorRow++
		}

		rowLength := provider.rowLength(e.CursorRow)
		if e.CursorColumn > rowLength {
			e.CursorColumn = rowLength
		}
	case fyne.KeyLeft:
		if e.CursorColumn > 0 {
			e.CursorColumn--
			break
		}

		if e.MultiLine && e.CursorRow > 0 {
			e.CursorRow--
			rowLength := provider.rowLength(e.CursorRow)
			e.CursorColumn = rowLength
		}
	case fyne.KeyRight:
		if e.MultiLine {
			rowLength := provider.rowLength(e.CursorRow)
			if e.CursorColumn < rowLength {
				e.CursorColumn++
				break
			}
			if e.CursorRow < provider.rows()-1 {
				e.CursorRow++
				e.CursorColumn = 0
			}
			break
		}
		if e.CursorColumn < provider.len() {
			e.CursorColumn++
		}
	case fyne.KeyEnd:
		e.CursorColumn = provider.len()
	case fyne.KeyHome:
		e.CursorColumn = 0
	default:
		return
	}

	e.updateText(provider.String())
	Renderer(e).(*entryRenderer).moveCursor()
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
func (p *placeholderPresenter) password() bool {
	return p.e.Password
}

// object returns the root object of the widget so it can be referenced
func (p *placeholderPresenter) object() fyne.Widget {
	return nil
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (e *Entry) CreateRenderer() fyne.WidgetRenderer {
	text := &textProvider{buffer: []rune(e.Text), presenter: e}
	placeholder := &textProvider{presenter: &placeholderPresenter{e}, buffer: []rune(e.PlaceHolder)}

	box := canvas.NewRectangle(theme.BackgroundColor())
	cursor := canvas.NewRectangle(theme.BackgroundColor())

	return &entryRenderer{text, placeholder, box, cursor,
		[]fyne.CanvasObject{box, placeholder, text, cursor}, e}
}

// NewEntry creates a new single line entry widget.
func NewEntry() *Entry {
	e := &Entry{}

	Refresh(e)
	return e
}

// NewMultiLineEntry creates a new entry that allows multiple lines
func NewMultiLineEntry() *Entry {
	e := &Entry{MultiLine: true}

	Refresh(e)
	return e
}

// NewPasswordEntry creates a new entry password widget
func NewPasswordEntry() *Entry {
	e := &Entry{Password: true}

	Refresh(e)
	return e
}

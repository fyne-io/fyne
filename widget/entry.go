package widget

import (
	"image/color"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/fyne-io/fyne/theme"
)

const (
	passwordChar = "*"
)

type entryRenderer struct {
	label       *Label
	box, cursor *canvas.Rectangle

	objects []fyne.CanvasObject
	entry   *Entry
}

func emptyTextMinSize(style fyne.TextStyle) fyne.Size {
	return textMinSize("M", theme.TextSize(), style)
}

func textMinSize(text string, size int, style fyne.TextStyle) fyne.Size {
	label := canvas.NewText(text, color.Black)
	label.TextSize = size
	label.TextStyle = style
	return label.MinSize()
}

// MinSize calculates the minimum size of an entry widget.
// This is based on the contained text with a standard amount of padding added.
func (e *entryRenderer) MinSize() fyne.Size {
	var textSize fyne.Size
	if e.label.Text == "" {
		textSize = emptyTextMinSize(e.label.TextStyle)
	} else {
		textSize = e.label.MinSize()
	}

	return textSize.Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
}

func (e *entryRenderer) cursorPosition() (int, int) {
	renderlabel := Renderer(e.label).(*labelRenderer).texts[0]
	lineHeight := emptyTextMinSize(e.label.TextStyle).Height

	runes := []rune(Renderer(e.label).(*labelRenderer).texts[e.entry.CursorRow].Text)
	// sanity check, as the underlying entry text can actually change
	if e.entry.CursorColumn > len(runes) {
		e.entry.CursorColumn = len(runes)
	}
	substr := string(runes[0:e.entry.CursorColumn])
	subSize := textMinSize(substr, renderlabel.TextSize, e.label.TextStyle)

	return subSize.Width, e.entry.CursorRow * lineHeight
}

func (e *entryRenderer) moveCursor() {
	xPos, yPos := e.cursorPosition()
	lineHeight := emptyTextMinSize(e.label.TextStyle).Height

	e.cursor.Resize(fyne.NewSize(2, lineHeight))
	e.cursor.Move(fyne.NewPos(xPos+theme.Padding()*2, yPos+theme.Padding()*2))

	canvas.Refresh(e.cursor)
}

// Layout the components of the entry widget.
func (e *entryRenderer) Layout(size fyne.Size) {
	e.box.Resize(size.Subtract(fyne.NewSize(theme.Padding(), theme.Padding())))
	e.box.Move(fyne.NewPos(theme.Padding()/2, theme.Padding()/2))

	e.label.Resize(size.Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
	e.label.Move(fyne.NewPos(theme.Padding(), theme.Padding()))

	e.moveCursor()
}

// ApplyTheme is called when the Entry may need to update it's look.
func (e *entryRenderer) ApplyTheme() {
	Renderer(e.label).ApplyTheme()
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
	var text string
	if e.entry.password {
		text = strings.Repeat(passwordChar, utf8.RuneCountInString(e.entry.Text))
	} else {
		text = e.entry.Text
	}

	e.label.SetText(text)

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

	Text      string
	OnChanged func(string) `json:"-"`

	CursorRow, CursorColumn int

	focused  bool
	password bool
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
}

// Hide this widget, if it was previously visible
func (e *Entry) Hide() {
	e.hide(e)
}

// SetText manually sets the text of the Entry to the given text value.
func (e *Entry) SetText(text string) {
	e.updateText(text)
}

// updateText updates the internal text to the given value
func (e *Entry) updateText(text string) {
	e.Text = text
	if e.OnChanged != nil {
		e.OnChanged(text)
	}

	Renderer(e).Refresh()
}

func (e *Entry) cursorTextPos() int {
	pos := 0
	texts := Renderer(e.label()).(*labelRenderer).texts
	for i := 0; i < e.CursorRow; i++ {
		line := texts[i].Text
		pos += len(line) + 1
	}
	pos += e.CursorColumn

	// Some sanity checks here
	if pos > len(e.Text) {
		pos = 0
		e.CursorColumn = 0
		e.CursorRow = 0
	}

	return pos
}

func (e *Entry) insertAtCursor(text string) {
	pos := e.cursorTextPos()
	runes := []rune(e.Text)
	runes = append(runes[:pos], append([]rune(text), runes[pos:]...)...)
	e.updateText(string(runes))
}

func (e *Entry) deleteFromTo(from int, to int) {
	runes := []rune(e.Text)
	runes = append(runes[:from], runes[to:]...)
	e.updateText(string(runes))
}

// OnFocusGained is called when the Entry has been given focus.
func (e *Entry) OnFocusGained() {
	e.focused = true

	Renderer(e).Refresh()
}

// OnFocusLost is called when the Entry has had focus removed.
func (e *Entry) OnFocusLost() {
	e.focused = false

	Renderer(e).Refresh()
}

// Focused returns whether or not this Entry has focus.
func (e *Entry) Focused() bool {
	return e.focused
}

// OnKeyDown receives key input events when the Entry widget is focused.
func (e *Entry) OnKeyDown(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyBackspace {
		if len(e.Text) == 0 || (e.CursorColumn == 0 && e.CursorRow == 0) {
			return
		}

		pos := e.cursorTextPos()
		if e.Text[pos-1] == '\n' {
			e.CursorRow--
			e.CursorColumn = e.label().RowLength(e.CursorRow)
		} else {
			if e.CursorColumn > 0 {
				e.CursorColumn--
			}
		}
		e.deleteFromTo(pos-1, pos)
	} else if key.Name == fyne.KeyDelete {
		texts := Renderer(e.label()).(*labelRenderer).texts
		if len(e.Text) == 0 || (e.CursorRow == len(texts)-1 && e.CursorColumn == len(texts[e.CursorRow].Text)) {
			return
		}

		pos := e.cursorTextPos()
		e.deleteFromTo(pos, pos+1)
	} else if key.Name == fyne.KeyReturn && !e.password {
		e.insertAtCursor("\n")

		e.CursorColumn = 0
		e.CursorRow++
	} else if key.Name == fyne.KeyUp {
		if e.CursorRow > 0 {
			e.CursorRow--
		}

		if e.CursorColumn > e.label().RowLength(e.CursorRow) {
			e.CursorColumn = e.label().RowLength(e.CursorRow)
		}
	} else if key.Name == fyne.KeyDown {
		if e.CursorRow < e.label().Rows()-1 {
			e.CursorRow++
		}

		if e.CursorColumn > e.label().RowLength(e.CursorRow) {
			e.CursorColumn = e.label().RowLength(e.CursorRow)
		}
	} else if key.Name == fyne.KeyLeft {
		if e.CursorColumn > 0 {
			e.CursorColumn--
		} else if e.CursorRow > 0 {
			e.CursorRow--
			e.CursorColumn = e.label().RowLength(e.CursorRow)
		}
	} else if key.Name == fyne.KeyRight {
		if e.CursorColumn < e.label().RowLength(e.CursorRow) {
			e.CursorColumn++
		} else if e.CursorRow < e.label().Rows()-1 {
			e.CursorRow++
			e.CursorColumn = 0
		}

	} else if key.String != "" {
		e.insertAtCursor(key.String)

		e.CursorColumn += len(key.String)
	} else {
		log.Println("Unhandled key press", key.String)
	}

	Renderer(e).(*entryRenderer).moveCursor()
}

func (e *Entry) label() *Label {
	return Renderer(e).(*entryRenderer).label
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (e *Entry) CreateRenderer() fyne.WidgetRenderer {
	text := NewLabel(e.Text)
	box := canvas.NewRectangle(theme.BackgroundColor())
	cursor := canvas.NewRectangle(theme.BackgroundColor())

	return &entryRenderer{text, box, cursor,
		[]fyne.CanvasObject{box, text, cursor}, e}
}

// NewEntry creates a new entry widget.
func NewEntry() *Entry {
	e := &Entry{}

	Renderer(e).Layout(e.MinSize())
	return e
}

// NewPasswordEntry creates a new entry password widget
func NewPasswordEntry() *Entry {
	e := &Entry{password: true}

	Renderer(e).Layout(e.MinSize())
	return e
}

package widget

import (
	"fmt"
	"image/color"
)
import "log"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/theme"

type entryRenderer struct {
	label           *Label
	bg, box, cursor *canvas.Rectangle

	objects []fyne.CanvasObject
	entry   *Entry
}

func emptyTextMinSize(style fyne.TextStyle) fyne.Size {
	return fyne.GetDriver().RenderedTextSize("M", theme.TextSize(), style)
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
	renderlabel := e.label.renderer.(*labelRenderer).texts[0]
	lineHeight := emptyTextMinSize(e.label.TextStyle).Height

	str := e.label.renderer.(*labelRenderer).texts[e.entry.CursorRow].Text
	// sanity check, as the underlying entry text can actually change
	if e.entry.CursorColumn > len(str) {
		e.entry.CursorColumn = len(str)
	}
	substr := str[0:e.entry.CursorColumn]
	subSize := fyne.GetDriver().RenderedTextSize(substr, renderlabel.TextSize, e.label.TextStyle)

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
	e.bg.Resize(size)

	e.box.Resize(size.Subtract(fyne.NewSize(theme.Padding(), theme.Padding())))
	e.box.Move(fyne.NewPos(theme.Padding()/2, theme.Padding()/2))

	e.label.Resize(size.Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
	e.label.Move(fyne.NewPos(theme.Padding(), theme.Padding()))

	e.moveCursor()
}

// ApplyTheme is called when the Entry may need to update it's look.
func (e *entryRenderer) ApplyTheme() {
	e.label.ApplyTheme()
	e.box.FillColor = theme.BackgroundColor()
	e.Refresh()
}

func (e *entryRenderer) Refresh() {
	e.label.SetText(e.entry.Text)

	if e.entry.focused {
		e.bg.FillColor = theme.FocusColor()
		e.cursor.FillColor = theme.FocusColor()
	} else {
		e.bg.FillColor = theme.ButtonColor()
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

	focused bool
}

// SetText manually sets the text of the Entry to the given text value.
func (e *Entry) SetText(text string) {
	e.Text = text

	if e.OnChanged != nil {
		e.OnChanged(text)
	}

	e.Renderer().Refresh()
}

func (e *Entry) cursorTextPos() int {
	pos := 0
	texts := e.label().renderer.(*labelRenderer).texts
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
	newText := fmt.Sprintf("%s%s%s", e.Text[:pos], text, e.Text[pos:])

	e.SetText(newText)
}

// OnFocusGained is called when the Entry has been given focus.
func (e *Entry) OnFocusGained() {
	e.focused = true

	e.Renderer().Refresh()
}

// OnFocusLost is called when the Entry has had focus removed.
func (e *Entry) OnFocusLost() {
	e.focused = false

	e.Renderer().Refresh()
}

// Focused returns whether or not this Entry has focus.
func (e *Entry) Focused() bool {
	return e.focused
}

// OnKeyDown receives key input events when the Entry widget is focused.
func (e *Entry) OnKeyDown(key *fyne.KeyEvent) {
	if key.Name == "BackSpace" {
		if len(e.Text) == 0 || (e.CursorColumn == 0 && e.CursorRow == 0) {
			return
		}

		runes := []rune(e.Text)
		pos := e.cursorTextPos()
		substr := fmt.Sprintf("%s%s", string(runes[:pos-1]), string(runes[pos:]))

		if runes[pos-1] == '\n' {
			e.CursorRow--
			e.CursorColumn = e.label().RowLength(e.CursorRow)
		} else {
			if e.CursorColumn > 0 {
				e.CursorColumn--
			}
		}

		e.SetText(substr)
	} else if key.Name == "Delete" {
		texts := e.label().renderer.(*labelRenderer).texts
		if len(e.Text) == 0 || (e.CursorRow == len(texts)-1 && e.CursorColumn == len(texts[e.CursorRow].Text)) {
			return
		}

		runes := []rune(e.Text)
		pos := e.cursorTextPos()
		substr := fmt.Sprintf("%s%s", string(runes[:pos]), string(runes[pos+1:]))

		e.SetText(substr)
	} else if key.Name == "Return" {
		e.insertAtCursor("\n")

		e.CursorColumn = 0
		e.CursorRow++
	} else if key.Name == "Up" {
		if e.CursorRow > 0 {
			e.CursorRow--
		}

		if e.CursorColumn > e.label().RowLength(e.CursorRow) {
			e.CursorColumn = e.label().RowLength(e.CursorRow)
		}
	} else if key.Name == "Down" {
		if e.CursorRow < e.label().Rows()-1 {
			e.CursorRow++
		}

		if e.CursorColumn > e.label().RowLength(e.CursorRow) {
			e.CursorColumn = e.label().RowLength(e.CursorRow)
		}
	} else if key.Name == "Left" {
		if e.CursorColumn > 0 {
			e.CursorColumn--
		} else if e.CursorRow > 0 {
			e.CursorRow--
			e.CursorColumn = e.label().RowLength(e.CursorRow)
		}
	} else if key.Name == "Right" {
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

	e.Renderer().(*entryRenderer).moveCursor()
}

func (e *Entry) label() *Label {
	return e.Renderer().(*entryRenderer).label
}

func (e *Entry) createRenderer() fyne.WidgetRenderer {
	text := NewLabel(e.Text)
	bg := canvas.NewRectangle(theme.ButtonColor())
	box := canvas.NewRectangle(theme.BackgroundColor())
	cursor := canvas.NewRectangle(theme.BackgroundColor())

	return &entryRenderer{text, bg, box, cursor,
		[]fyne.CanvasObject{bg, box, text, cursor}, e}
}

// Renderer is a private method to Fyne which links this widget to it's renderer
func (e *Entry) Renderer() fyne.WidgetRenderer {
	if e.renderer == nil {
		e.renderer = e.createRenderer()
	}

	return e.renderer
}

// NewEntry creates a new entry widget.
func NewEntry() *Entry {
	e := &Entry{}

	e.Renderer().Layout(e.MinSize())
	return e
}

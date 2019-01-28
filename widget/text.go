package widget

import (
	"image/color"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

const (
	passwordChar = "*"
)

type textParent interface {
	textAlign() fyne.TextAlign
	textStyle() fyne.TextStyle
	textColor() color.Color

	password() bool

	object() fyne.Widget
}

// textWidget represents the base element for text based widget.
type textWidget struct {
	baseWidget
	parent textParent

	buffer    []rune
	rowBounds [][2]int
}

// NewText returns a new Text with the given text and default settings.
func newTextWidget(text string, parent textParent) textWidget {
	t := textWidget{
		buffer: []rune(text),
		parent: parent,
	}
	t.updateRowBounds()
	return t
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (t *textWidget) Resize(size fyne.Size) {
	t.resize(size, t)
}

// Move the widget to a new position, relative to it's parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (t *textWidget) Move(pos fyne.Position) {
	t.move(pos, t)
}

// MinSize returns the smallest size this widget can shrink to
func (t *textWidget) MinSize() fyne.Size {
	return t.minSize(t)
}

// Show this widget, if it was previously hidden
func (t *textWidget) Show() {
	t.show(t)
}

// Hide this widget, if it was previously visible
func (t *textWidget) Hide() {
	t.hide(t)
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (t *textWidget) CreateRenderer() fyne.WidgetRenderer {
	if t.parent == nil {
		panic("Cannot create a textWidget without a parent")
	}
	r := &textRenderer{textWidget: t}

	t.updateRowBounds() // set up the initial text layout etc
	r.Refresh()
	return r
}

// updateRowBounds updates the row bounds used to render properly the text widget.
// updateRowBounds should be invoked every time t.buffer changes.
func (t *textWidget) updateRowBounds() {
	var lowBound, highBound int
	t.rowBounds = [][2]int{}

	if len(t.buffer) == 0 {
		t.rowBounds = append(t.rowBounds, [2]int{lowBound, highBound})
		return
	}

	for i, r := range t.buffer {
		highBound = i
		if r != '\n' {
			continue
		}
		t.rowBounds = append(t.rowBounds, [2]int{lowBound, highBound})
		lowBound = i + 1
	}
	//first or last line, increase the highBound index to include the last char
	highBound++
	t.rowBounds = append(t.rowBounds, [2]int{lowBound, highBound})
}

// refreshTextRenderer refresh the textRenderer canvas objects
// this method should be invoked every time the t.buffer changes
// example:
// t.buffer = []rune("new text")
// t.updateRowBounds()
// t.refreshTextRenderer()
func (t *textWidget) refreshTextRenderer() {
	var obj fyne.Widget
	if t.parent != nil && t.parent.object() != nil {
		obj = t.parent.object()
	} else {
		obj = t
	}

	Renderer(obj).Refresh()
}

// SetText sets the text of the widget
func (t *textWidget) SetText(text string) {
	t.buffer = []rune(text)
	t.updateRowBounds()
	t.refreshTextRenderer()
}

// String returns the text widget buffer as string
func (t *textWidget) String() string {
	return string(t.buffer)
}

// Len returns the text widget buffer length
func (t *textWidget) len() int {
	return len(t.buffer)
}

// insertAt inserts the text at the specified position
func (t *textWidget) insertAt(pos int, runes []rune) {
	t.buffer = append(t.buffer[:pos], append(runes, t.buffer[pos:]...)...)
	t.updateRowBounds()
	t.refreshTextRenderer()
}

// deleteFromTo removes the text between the specified positions
func (t *textWidget) deleteFromTo(lowBound int, highBound int) []rune {
	deleted := make([]rune, highBound-lowBound)
	copy(deleted, t.buffer[lowBound:highBound])
	t.buffer = append(t.buffer[:lowBound], t.buffer[highBound:]...)
	t.updateRowBounds()
	t.refreshTextRenderer()
	return deleted
}

// rows returns the number of text rows in this text entry.
// The entry may be longer than required to show this amount of content.
func (t *textWidget) rows() int {
	return len(t.rowBounds)
}

// Row returns the characters in the row specified.
// The row parameter should be between 0 and t.Rows()-1.
func (t *textWidget) row(row int) []rune {
	bounds := t.rowBounds[row]
	return t.buffer[bounds[0]:bounds[1]]
}

// RowLength returns the number of visible characters in the row specified.
// The row parameter should be between 0 and t.Rows()-1.
func (t *textWidget) rowLength(row int) int {
	return len(t.row(row))
}

// CharMinSize returns the average char size to use for internal computation
func (t *textWidget) charMinSize() fyne.Size {
	defaultChar := "M"
	if t.parent.password() {
		defaultChar = passwordChar
	}
	return textMinSize(defaultChar, theme.TextSize(), t.parent.textStyle())
}

// Renderer
type textRenderer struct {
	objects []fyne.CanvasObject

	texts []*canvas.Text

	*textWidget
}

// MinSize calculates the minimum size of a label.
// This is based on the contained text with a standard amount of padding added.
func (r *textRenderer) MinSize() fyne.Size {
	height := 0
	width := 0
	for i := 0; i < len(r.texts); i++ {
		min := r.texts[i].MinSize()
		if r.texts[i].Text == "" {
			min = r.textWidget.charMinSize()
		}
		height += min.Height
		width = fyne.Max(width, min.Width)
	}

	return fyne.NewSize(width, height).Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
}

func (r *textRenderer) Layout(size fyne.Size) {
	yPos := theme.Padding()
	lineHeight := r.textWidget.charMinSize().Height
	lineSize := fyne.NewSize(size.Width-theme.Padding()*2, lineHeight)
	for i := 0; i < len(r.texts); i++ {
		text := r.texts[i]
		text.Resize(lineSize)
		text.Move(fyne.NewPos(theme.Padding(), yPos))
		yPos += lineHeight
	}
}

func (r *textRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// ApplyTheme is called when the Label may need to update it's look
func (r *textRenderer) ApplyTheme() {
	c := theme.TextColor()
	if r.textWidget.parent.textColor() != nil {
		c = r.textWidget.parent.textColor()
	}
	for _, text := range r.texts {
		text.Color = c
	}
}

func (r *textRenderer) Refresh() {
	r.texts = []*canvas.Text{}
	r.objects = []fyne.CanvasObject{}
	for index := 0; index < r.textWidget.rows(); index++ {
		var line string
		row := r.textWidget.row(index)
		if r.textWidget.parent.password() {
			line = strings.Repeat(passwordChar, len(row))
		} else {
			line = string(row)
		}
		textCanvas := canvas.NewText(line, theme.TextColor())
		textCanvas.Alignment = r.textWidget.parent.textAlign()
		textCanvas.TextStyle = r.textWidget.parent.textStyle()
		r.texts = append(r.texts, textCanvas)
		r.objects = append(r.objects, textCanvas)
	}

	r.ApplyTheme()
	r.Layout(r.textWidget.Size())
	if r.textWidget.parent.object() == nil {
		canvas.Refresh(r.textWidget)
	} else {
		canvas.Refresh(r.textWidget.parent.object())
	}
}

func (r *textRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

// lineSize returns the rendered size for the line specified by col and row
func (r *textRenderer) lineSize(col, row int) (size fyne.Size) {
	text := r.textWidget

	line := text.row(row)

	if col >= len(line) {
		col = len(line)
	}
	lineCopy := *r.texts[row]
	if r.textWidget.parent.password() {
		lineCopy.Text = strings.Repeat(passwordChar, col)
	} else {
		lineCopy.Text = string(line[0:col])
	}

	return lineCopy.MinSize()
}

func textMinSize(text string, size int, style fyne.TextStyle) fyne.Size {
	t := canvas.NewText(text, color.Black)
	t.TextSize = size
	t.TextStyle = style
	return t.MinSize()
}

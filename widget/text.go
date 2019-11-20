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

// textPresenter provides the widget specific information to a generic text provider
type textPresenter interface {
	textAlign() fyne.TextAlign
	textStyle() fyne.TextStyle
	textColor() color.Color

	password() bool

	object() fyne.Widget
}

// textProvider represents the base element for text based widget.
type textProvider struct {
	BaseWidget
	presenter textPresenter

	buffer    []rune
	rowBounds [][2]int
}

// newTextProvider returns a new textProvider with the given text and settings from the passed textPresenter.
func newTextProvider(text string, pres textPresenter) textProvider {
	if pres == nil {
		panic("textProvider requires a presenter")
	}
	t := textProvider{
		buffer:    []rune(text),
		presenter: pres,
	}
	t.updateRowBounds()
	return t
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (t *textProvider) CreateRenderer() fyne.WidgetRenderer {
	if t.presenter == nil {
		panic("Cannot render a textProvider without a presenter")
	}

	if t.presenter.object() == nil {
		t.ExtendBaseWidget(t)
	} else {
		t.ExtendBaseWidget(t.presenter.object())
	}
	r := &textRenderer{provider: t}

	t.updateRowBounds() // set up the initial text layout etc
	r.Refresh()
	return r
}

// updateRowBounds updates the row bounds used to render properly the text widget.
// updateRowBounds should be invoked every time t.buffer changes.
func (t *textProvider) updateRowBounds() {
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
func (t *textProvider) refreshTextRenderer() {
	if t.presenter == nil {
		return // not yet shown
	}
	obj := t.presenter.object()
	if obj == nil {
		obj = t
	}

	obj.Refresh()
}

// SetText sets the text of the widget
func (t *textProvider) SetText(text string) {
	t.buffer = []rune(text)
	t.updateRowBounds()

	t.refreshTextRenderer()
}

// String returns the text widget buffer as string
func (t *textProvider) String() string {
	return string(t.buffer)
}

// Len returns the text widget buffer length
func (t *textProvider) len() int {
	return len(t.buffer)
}

// insertAt inserts the text at the specified position
func (t *textProvider) insertAt(pos int, runes []rune) {
	// edge case checking
	if len(t.buffer) < pos {
		// append to the end if our position was out of sync
		t.buffer = append(t.buffer, runes...)
	} else {
		t.buffer = append(t.buffer[:pos], append(runes, t.buffer[pos:]...)...)
	}
	t.updateRowBounds()
	t.refreshTextRenderer()
}

// deleteFromTo removes the text between the specified positions
func (t *textProvider) deleteFromTo(lowBound int, highBound int) []rune {
	deleted := make([]rune, highBound-lowBound)
	copy(deleted, t.buffer[lowBound:highBound])
	t.buffer = append(t.buffer[:lowBound], t.buffer[highBound:]...)
	t.updateRowBounds()
	t.refreshTextRenderer()
	return deleted
}

// rows returns the number of text rows in this text entry.
// The entry may be longer than required to show this amount of content.
func (t *textProvider) rows() int {
	return len(t.rowBounds)
}

// Row returns the characters in the row specified.
// The row parameter should be between 0 and t.Rows()-1.
func (t *textProvider) row(row int) []rune {
	if row < 0 || row >= t.rows() {
		return nil
	}
	bounds := t.rowBounds[row]
	from := bounds[0]
	to := bounds[1]
	if from < 0 || to > len(t.buffer) {
		return nil
	}
	if to < from {
		return nil
	}
	return t.buffer[from:to]
}

// RowLength returns the number of visible characters in the row specified.
// The row parameter should be between 0 and t.Rows()-1.
func (t *textProvider) rowLength(row int) int {
	return len(t.row(row))
}

// CharMinSize returns the average char size to use for internal computation
func (t *textProvider) charMinSize() fyne.Size {
	defaultChar := "M"
	if t.presenter.password() {
		defaultChar = passwordChar
	}
	return textMinSize(defaultChar, theme.TextSize(), t.presenter.textStyle())
}

// lineSizeToColumn returns the rendered size for the line specified by row up to the col position
func (t *textProvider) lineSizeToColumn(col, row int) (size fyne.Size) {
	line := t.row(row)
	if line == nil {
		return fyne.NewSize(0, 0)
	}

	if col >= len(line) {
		col = len(line)
	}

	measureText := string(line[0:col])
	if t.presenter.password() {
		measureText = strings.Repeat(passwordChar, col)
	}

	label := canvas.NewText(measureText, theme.TextColor())
	label.TextStyle = t.presenter.textStyle()
	return label.MinSize()
}

// Renderer
type textRenderer struct {
	objects []fyne.CanvasObject

	texts []*canvas.Text

	provider *textProvider
}

// MinSize calculates the minimum size of a label.
// This is based on the contained text with a standard amount of padding added.
func (r *textRenderer) MinSize() fyne.Size {
	height := 0
	width := 0
	for i := 0; i < fyne.Min(len(r.texts), r.provider.rows()); i++ {
		min := r.texts[i].MinSize()
		if r.texts[i].Text == "" {
			min = r.provider.charMinSize()
		}
		height += min.Height
		width = fyne.Max(width, min.Width)
	}

	return fyne.NewSize(width, height).Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
}

func (r *textRenderer) Layout(size fyne.Size) {
	yPos := theme.Padding()
	lineHeight := r.provider.charMinSize().Height
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

// applyTheme updates the label to match the current theme.
func (r *textRenderer) applyTheme() {
	c := theme.TextColor()
	if r.provider.presenter.textColor() != nil {
		c = r.provider.presenter.textColor()
	}
	for _, text := range r.texts {
		text.Color = c
		text.TextSize = theme.TextSize()
	}
}

func (r *textRenderer) Refresh() {
	index := 0
	for ; index < r.provider.rows(); index++ {
		var line string
		row := r.provider.row(index)
		if r.provider.presenter.password() {
			line = strings.Repeat(passwordChar, len(row))
		} else {
			line = string(row)
		}

		var textCanvas *canvas.Text
		add := false
		if index >= len(r.texts) {
			add = true
			textCanvas = canvas.NewText(line, theme.TextColor())
		} else {
			textCanvas = r.texts[index]
			textCanvas.Text = line
		}

		textCanvas.Alignment = r.provider.presenter.textAlign()
		textCanvas.TextStyle = r.provider.presenter.textStyle()
		textCanvas.Hidden = r.provider.Hidden

		if add {
			r.texts = append(r.texts, textCanvas)
			r.objects = append(r.objects, textCanvas)
		}
	}

	for ; index < len(r.texts); index++ {
		r.texts[index].Text = ""
	}

	r.applyTheme()
	r.Layout(r.provider.Size())
	if r.provider.presenter.object() == nil {
		canvas.Refresh(r.provider)
	} else {
		canvas.Refresh(r.provider.presenter.object())
	}
}

func (r *textRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *textRenderer) Destroy() {
}

func textMinSize(text string, size int, style fyne.TextStyle) fyne.Size {
	t := canvas.NewText(text, color.Black)
	t.TextSize = size
	t.TextStyle = style
	return t.MinSize()
}

package widget

import (
	"image/color"
	"strings"
	"unicode"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/theme"
)

const (
	passwordChar = "•"
)

// textPresenter provides the widget specific information to a generic text provider
type textPresenter interface {
	textAlign() fyne.TextAlign
	textWrap() fyne.TextWrap
	textStyle() fyne.TextStyle
	textColor() color.Color

	concealed() bool

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

func (t *textProvider) Resize(size fyne.Size) {
	if t.size == size {
		return
	}
	t.size = size
	t.updateRowBounds()
	if t.presenter != nil {
		t.refreshTextRenderer()
		cache.Renderer(t).Layout(size)
	}
}

// updateRowBounds updates the row bounds used to render properly the text widget.
// updateRowBounds should be invoked every time t.buffer or viewport changes.
func (t *textProvider) updateRowBounds() {
	if t.presenter == nil {
		t.rowBounds = [][2]int{}
		return // not yet shown
	}
	textWrap := t.presenter.textWrap()
	textStyle := t.presenter.textStyle()
	textSize := theme.TextSize()
	maxWidth := t.Size().Width - 2*theme.Padding()

	t.rowBounds = lineBounds(t.buffer, textWrap, maxWidth, func(text []rune) int {
		return fyne.MeasureText(string(text), textSize, textStyle).Width
	})
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

// RowBoundary returns the boundary of the row specified.
// The row parameter should be between 0 and t.Rows()-1.
func (t *textProvider) rowBoundary(row int) [2]int {
	if row < 0 || row >= t.rows() {
		return [2]int{0, 0}
	}
	return t.rowBounds[row]
}

// RowLength returns the number of visible characters in the row specified.
// The row parameter should be between 0 and t.Rows()-1.
func (t *textProvider) rowLength(row int) int {
	return len(t.row(row))
}

// CharMinSize returns the average char size to use for internal computation
func (t *textProvider) charMinSize() fyne.Size {
	defaultChar := "M"
	if t.presenter.concealed() {
		defaultChar = passwordChar
	}
	return fyne.MeasureText(defaultChar, theme.TextSize(), t.presenter.textStyle())
}

// lineSizeToColumn returns the rendered size for the line specified by row up to the col position
func (t *textProvider) lineSizeToColumn(col, row int) fyne.Size {
	line := t.row(row)
	if line == nil {
		return fyne.NewSize(0, 0)
	}

	if col >= len(line) {
		col = len(line)
	}

	measureText := string(line[0:col])
	if t.presenter.concealed() {
		measureText = strings.Repeat(passwordChar, col)
	}

	label := canvas.NewText(measureText, theme.TextColor())
	label.TextStyle = t.presenter.textStyle()
	return label.MinSize()
}

// Renderer
type textRenderer struct {
	baseRenderer
	texts    []*canvas.Text
	provider *textProvider
}

// MinSize calculates the minimum size of a label.
// This is based on the contained text with a standard amount of padding added.
func (r *textRenderer) MinSize() fyne.Size {
	wrap := r.provider.presenter.textWrap()
	charMinSize := r.provider.charMinSize()
	height := 0
	width := 0
	i := 0
	for ; i < fyne.Min(len(r.texts), r.provider.rows()); i++ {
		min := r.texts[i].MinSize()
		if r.texts[i].Text == "" {
			min = charMinSize
		}
		if wrap == fyne.TextWrapOff {
			width = fyne.Max(width, min.Width)
		}
		height += min.Height
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
		if r.provider.presenter.concealed() {
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
			r.setObjects(append(r.Objects(), textCanvas))
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

// splitLines accepts a slice of runes and returns a slice containing the
// start and end indicies of each line delimited by the newline character.
func splitLines(text []rune) [][2]int {
	var low, high int
	var lines [][2]int
	length := len(text)
	for i := 0; i < length; i++ {
		if text[i] == '\n' {
			high = i
			lines = append(lines, [2]int{low, high})
			low = i + 1
		}
	}
	return append(lines, [2]int{low, length})
}

// lineBounds accepts a slice of runes, a wrapping mode, a maximum line width and a function to measure line width.
// lineBounds returns a slice containing the start and end indicies of each line with the given wrapping applied.
func lineBounds(text []rune, wrap fyne.TextWrap, maxWidth int, measurer func([]rune) int) [][2]int {

	binarySearch := func(low int, maxHigh int) int {
		if low >= maxHigh {
			return low
		}
		if measurer(text[low:maxHigh]) <= maxWidth {
			return maxHigh
		}
		high := low
		delta := maxHigh - low
		for delta > 0 {
			delta /= 2
			if measurer(text[low:high+delta]) < maxWidth {
				high += delta
			}
		}
		for (high < maxHigh) && measurer(text[low:high+1]) <= maxWidth {
			high++
		}
		return high
	}

	lines := splitLines(text)
	if maxWidth <= 0 || wrap == fyne.TextWrapOff {
		return lines
	}
	var bounds [][2]int
	for _, l := range lines {
		low := l[0]
		high := l[1]
		if low == high {
			bounds = append(bounds, l)
			continue
		}
		switch wrap {
		case fyne.TextTruncate:
			high = binarySearch(low, high)
			bounds = append(bounds, [2]int{low, high})
		case fyne.TextWrapBreak:
			for low < high {
				if measurer(text[low:high]) <= maxWidth {
					bounds = append(bounds, [2]int{low, high})
					low = high
					high = l[1]
				} else {
					high = binarySearch(low, high)
				}
			}
		case fyne.TextWrapWord:
			for low < high {
				sub := text[low:high]
				if measurer(sub) <= maxWidth {
					bounds = append(bounds, [2]int{low, high})
					low = high
					high = l[1]
					if low < high && unicode.IsSpace(text[low]) {
						low++
					}
				} else {
					last := len(sub) - 1
					findSpaceIndex := func(last int) int {
						curIndex := last
						for ; curIndex >= 0; curIndex-- {
							if unicode.IsSpace(sub[curIndex]) {
								break
							}
						}
						if curIndex < 0 {
							return last
						}
						return curIndex
					}
					high = low + findSpaceIndex(binarySearch(low, low+last)-low)
				}
			}
		}
	}
	return bounds
}

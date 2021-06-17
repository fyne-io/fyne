package widget

import (
	"image/color"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

const (
	passwordChar = "•"
)

var (
	// RichTextStyleEmphasis represents standard text that can be surrounded by other elements.
	//
	// Since: 2.1
	RichTextStyleEmphasis = RichTextStyle{
		ColorName: theme.ColorNameForeground,
		Inline:    true,
		SizeName:  theme.SizeNameText,
		TextStyle: fyne.TextStyle{Italic: true},
	}
	// RichTextStyleHeading represents a heading text that stands on its own line.
	//
	// Since: 2.1
	RichTextStyleHeading = RichTextStyle{
		ColorName: theme.ColorNameForeground,
		Inline:    false,
		SizeName:  theme.SizeNameHeadingText,
		TextStyle: fyne.TextStyle{Bold: true},
	}
	// RichTextStyleInline represents standard text that can be surrounded by other elements.
	//
	// Since: 2.1
	RichTextStyleInline = RichTextStyle{
		ColorName: theme.ColorNameForeground,
		Inline:    true,
		SizeName:  theme.SizeNameText,
	}
	// RichTextStyleParagraph represents standard text that should appear separate from other text.
	//
	// Since: 2.1
	RichTextStyleParagraph = RichTextStyle{
		ColorName: theme.ColorNameForeground,
		Inline:    false,
		SizeName:  theme.SizeNameText,
	}
	// RichTextStylePassword represents standard sized text where the characters are obscured.
	//
	// Since: 2.1
	RichTextStylePassword = RichTextStyle{
		ColorName: theme.ColorNameForeground,
		Inline:    true,
		SizeName:  theme.SizeNameText,
		concealed: true,
	}
	// RichTextStyleStrong represents standard text that can be surrounded by other elements.
	//
	// Since: 2.1
	RichTextStyleStrong = RichTextStyle{
		ColorName: theme.ColorNameForeground,
		Inline:    true,
		SizeName:  theme.SizeNameText,
		TextStyle: fyne.TextStyle{Bold: true},
	}
	// RichTextStyleSubHeading represents a sub-heading text that stands on its own line.
	//
	// Since: 2.1
	RichTextStyleSubHeading = RichTextStyle{
		ColorName: theme.ColorNameForeground,
		Inline:    false,
		SizeName:  theme.SizeNameSubHeadingText,
		TextStyle: fyne.TextStyle{Bold: true},
	}
)

// RichTextStyle describes the details of a text object inside a RichText widget.
//
// Since: 2.1
type RichTextStyle struct {
	Alignment fyne.TextAlign
	ColorName fyne.ThemeColorName
	Inline    bool
	SizeName  fyne.ThemeSizeName
	TextStyle fyne.TextStyle

	// an internal detail where we obscure password fields
	concealed bool
}

// RichTextSegment describes any element that can be rendered in a RichText widget.
//
// Since: 2.1
type RichTextSegment interface {
	Inline() bool
	Textual() string
	Visual() fyne.CanvasObject

	Select(pos1, pos2 fyne.Position)
	SelectedText() string
	Unselect()
}

// TextSegment represents the styling for a segment of rich text.
//
// Since: 2.1
type TextSegment struct {
	Style RichTextStyle
	Text  string
}

// Inline should return true if this text can be included within other elements, or false if it creates a new block.
func (t *TextSegment) Inline() bool {
	return t.Style.Inline
}

// Textual returns the content of this segment rendered to plain text.
func (t *TextSegment) Textual() string {
	return t.Text
}

// Visual returns the graphical elements required to render this segment.
func (t *TextSegment) Visual() fyne.CanvasObject {
	obj := canvas.NewText(t.Text, t.color())

	obj.Alignment = t.Style.Alignment
	obj.TextStyle = t.Style.TextStyle
	obj.TextSize = t.size()
	return obj
}

// Select tells the segment that the user is selecting the content between the two positions.
func (t *TextSegment) Select(begin, end fyne.Position) {
	// no-op: this will be added when we progress to editor
}

// SelectedText should return the text representation of any content currently selected through the Select call.
func (t *TextSegment) SelectedText() string {
	// no-op: this will be added when we progress to editor
	return ""
}

// Unselect tells the segment that the user is has cancelled the previous selection.
func (t *TextSegment) Unselect() {
	// no-op: this will be added when we progress to editor
}

func (t TextSegment) color() color.Color {
	if t.Style.ColorName != "" {
		return fyne.CurrentApp().Settings().Theme().Color(t.Style.ColorName, fyne.CurrentApp().Settings().ThemeVariant())
	}

	return theme.ForegroundColor()
}

func (t TextSegment) size() float32 {
	if t.Style.SizeName != "" {
		return fyne.CurrentApp().Settings().Theme().Size(t.Style.SizeName)
	}

	return theme.TextSize()
}

// RichText represents the base element for a rich text-based widget.
//
// Since: 2.1
type RichText struct {
	BaseWidget
	Segments []RichTextSegment
	Wrapping fyne.TextWrap

	inset     fyne.Size     // this varies due to how the widget works (entry with scroller vs others with padding)
	rowBounds []rowBoundary // cache for boundaries
}

// NewRichText returns a new RichText widget that renders the given text and segments.
// If no segments are specified it will be converted to a single segment using the default text settings.
//
// Since: 2.1
func NewRichText(segments ...RichTextSegment) *RichText {
	t := &RichText{Segments: segments}
	t.updateRowBounds()
	return t
}

// NewRichTextWithText returns a new RichText widget that renders the given text.
// The string will be converted to a single text segment using the default text settings.
//
// Since: 2.1
func NewRichTextWithText(text string) *RichText {
	return NewRichText(&TextSegment{
		Style: RichTextStyleInline,
		Text:  text,
	})
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (t *RichText) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	r := &textRenderer{obj: t}

	t.updateRowBounds() // set up the initial text layout etc
	r.Refresh()
	return r
}

// Refresh triggers a redraw of the rich text.
//
// Implements: fyne.Widget
func (t *RichText) Refresh() {
	t.updateRowBounds()

	t.BaseWidget.Refresh()
}

// Resize sets a new size for the rich text.
// This should only be called if it is not in a container with a layout manager.
//
// Implements: fyne.Widget
func (t *RichText) Resize(size fyne.Size) {
	t.propertyLock.RLock()
	baseSize := t.size
	t.propertyLock.RUnlock()
	if baseSize == size {
		return
	}

	t.propertyLock.Lock()
	t.size = size
	t.propertyLock.Unlock()
	t.updateRowBounds()

	t.Refresh()
	cache.Renderer(t).Layout(size)
}

// String returns the text widget buffer as string
func (t *RichText) String() string {
	ret := strings.Builder{}
	for _, seg := range t.Segments {
		ret.WriteString(seg.Textual())
	}
	return ret.String()
}

// CharMinSize returns the average char size to use for internal computation
func (t *RichText) charMinSize(concealed bool, style fyne.TextStyle) fyne.Size {
	defaultChar := "M"
	if concealed {
		defaultChar = passwordChar
	}

	return fyne.MeasureText(defaultChar, theme.TextSize(), style)
}

// deleteFromTo removes the text between the specified positions
func (t *RichText) deleteFromTo(lowBound int, highBound int) string {
	// TODO handle start portion, whole elements and end portion!
	index := 0
	start := 0
	var from *TextSegment
	for i, seg := range t.Segments {
		if _, ok := seg.(*TextSegment); !ok {
			continue
		}
		end := start + len([]rune(seg.(*TextSegment).Text))
		from = seg.(*TextSegment)
		index = i
		if end > lowBound {
			break
		}

		start = end
	}

	if from == nil {
		return ""
	}
	deleted := make([]rune, highBound-lowBound)
	r := ([]rune)(from.Text)
	copy(deleted, r[lowBound:highBound])
	if highBound > len(r) {
		highBound = len(r) // TODO remove this workaround and delete all segments)
	}
	r2 := append(r[:lowBound], r[highBound:]...)
	from.Text = string(r2)
	t.Segments[index] = from
	t.Refresh()
	return string(deleted)
}

// insertAt inserts the text at the specified position
func (t *RichText) insertAt(pos int, runes string) {
	index := 0
	start := 0
	var into *TextSegment
	for i, seg := range t.Segments {
		if _, ok := seg.(*TextSegment); !ok {
			continue
		}
		end := start + len([]rune(seg.(*TextSegment).Text))
		into = seg.(*TextSegment)
		index = i
		if end > pos {
			break
		}

		start = end
	}

	if into == nil {
		return
	}
	r := ([]rune)(into.Text)
	r2 := append(r[:pos], append([]rune(runes), r[pos:]...)...)
	into.Text = string(r2)
	t.Segments[index] = into

	t.Refresh()
}

// Len returns the text widget buffer length
func (t *RichText) len() int {
	ret := 0
	for _, seg := range t.Segments {
		ret += len([]rune(seg.Textual()))
	}
	return ret
}

// lineSizeToColumn returns the rendered size for the line specified by row up to the col position
func (t *RichText) lineSizeToColumn(col, row int) fyne.Size {
	bound := t.rowBoundary(row)
	total := fyne.NewSize(0, 0)
	counted := 0
	last := false
	for i, seg := range bound.segments {
		var size fyne.Size
		if text, ok := seg.(*TextSegment); ok {
			start := 0
			if i == 0 {
				start = bound.begin
			}
			measureText := []rune(text.Text)[start:]
			if col < counted+len(measureText) {
				measureText = measureText[0 : col-counted]
				last = true
			}
			if concealed(seg) {
				measureText = []rune(strings.Repeat(passwordChar, len(measureText)))
			}
			counted += len(measureText)

			label := canvas.NewText(string(measureText), color.Black)
			label.TextStyle = text.Style.TextStyle
			label.TextSize = text.size()

			size = label.MinSize()
		} else {
			size = seg.Visual().MinSize()
		}

		total.Width += size.Width
		total.Height = fyne.Max(total.Height, size.Height)
		if last {
			break
		}
	}
	return total.Add(fyne.NewSize(theme.Padding()*2-t.inset.Width, 0))
}

// Row returns the characters in the row specified.
// The row parameter should be between 0 and t.Rows()-1.
func (t *RichText) row(row int) []rune {
	if row < 0 || row >= t.rows() {
		return nil
	}
	bound := t.rowBounds[row]
	var ret []rune
	for i, seg := range bound.segments {
		if text, ok := seg.(*TextSegment); ok {
			if i == 0 {
				if len(bound.segments) == 1 {
					ret = append(ret, []rune(text.Text)[bound.begin:bound.end]...)
				} else {
					ret = append(ret, []rune(text.Text)[bound.begin:]...)
				}
			} else if i == len(bound.segments)-1 && len(bound.segments) > 1 && bound.end != 0 {
				ret = append(ret, []rune(text.Text)[:bound.end]...)
			}
		}
	}
	return ret
}

// RowBoundary returns the boundary of the row specified.
// The row parameter should be between 0 and t.Rows()-1.
func (t *RichText) rowBoundary(row int) *rowBoundary {
	t.propertyLock.RLock()
	defer t.propertyLock.RUnlock()
	if row < 0 || row >= t.rows() {
		return nil
	}
	return &t.rowBounds[row]
}

// RowLength returns the number of visible characters in the row specified.
// The row parameter should be between 0 and t.Rows()-1.
func (t *RichText) rowLength(row int) int {
	return len(t.row(row))
}

// rows returns the number of text rows in this text entry.
// The entry may be longer than required to show this amount of content.
func (t *RichText) rows() int {
	return len(t.rowBounds)
}

// updateRowBounds updates the row bounds used to render properly the text widget.
// updateRowBounds should be invoked every time a segment Text, widget Wrapping or size changes.
func (t *RichText) updateRowBounds() {
	t.propertyLock.RLock()
	var bounds []rowBoundary
	maxWidth := t.size.Width - 4*theme.Padding() - 2*t.inset.Width
	wrapWidth := maxWidth

	var currentBound *rowBoundary
	for _, seg := range t.Segments {
		if _, ok := seg.(*TextSegment); !ok {
			if currentBound == nil {
				bound := rowBoundary{segments: []RichTextSegment{seg}}
				bounds = append(bounds, bound)
				currentBound = &bound
			} else {
				bounds[len(bounds)-1].segments = append(bounds[len(bounds)-1].segments, seg)
			}
			if seg.Inline() {
				wrapWidth -= seg.Visual().MinSize().Width
			} else {
				currentBound = nil
			}
			continue
		}
		textSeg := seg.(*TextSegment)
		textStyle := textSeg.Style.TextStyle
		textSize := textSeg.size()

		retBounds := lineBounds(textSeg, t.Wrapping, wrapWidth, maxWidth, func(text []rune) float32 {
			return fyne.MeasureText(string(text), textSize, textStyle).Width
		})
		if currentBound != nil {
			if len(retBounds) > 0 {
				bounds[len(bounds)-1].end = retBounds[0].end // invalidate row ending as we have more content
				bounds[len(bounds)-1].segments = append(bounds[len(bounds)-1].segments, seg)
				bounds = append(bounds, retBounds[1:]...)
			}
		} else {
			bounds = append(bounds, retBounds...)
		}
		currentBound = &bounds[len(bounds)-1]
		if seg.Inline() {
			last := bounds[len(bounds)-1]
			begin := 0
			if len(last.segments) == 1 {
				begin = last.begin
			}
			text := string([]rune(textSeg.Text)[begin:last.end])
			lastWidth := fyne.MeasureText(text, textSeg.size(), textSeg.Style.TextStyle).Width
			if len(retBounds) == 1 {
				wrapWidth -= lastWidth
			} else {
				wrapWidth = maxWidth - lastWidth
			}
		} else {
			currentBound = nil
			wrapWidth = maxWidth
		}
	}
	t.propertyLock.RUnlock()

	t.propertyLock.Lock()
	t.rowBounds = bounds
	t.propertyLock.Unlock()
}

// Renderer
type textRenderer struct {
	widget.BaseRenderer
	obj *RichText
}

func (r *textRenderer) Layout(size fyne.Size) {
	r.obj.propertyLock.RLock()
	bounds := r.obj.rowBounds
	objs := r.Objects()
	r.obj.propertyLock.RUnlock()

	left := theme.Padding()*2 - r.obj.inset.Width
	yPos := theme.Padding()*2 - r.obj.inset.Height
	lineWidth := size.Width - yPos*2
	var rowItems []fyne.CanvasObject
	rowAlign := fyne.TextAlignLeading
	i := 0
	for _, bound := range bounds {
		for segI := range bound.segments {
			inline := segI < len(bound.segments)-1
			obj := objs[i]
			i++
			_, isText := obj.(*canvas.Text)
			if !isText && !inline {
				height := obj.MinSize().Height

				obj.Move(fyne.NewPos(left, yPos))
				obj.Resize(fyne.NewSize(lineWidth, height))
				yPos += height
				continue
			}
			rowItems = append(rowItems, obj)
			if inline {
				continue
			}

			if len(rowItems) == 1 { // TODO align link
				rowAlign = bound.segments[len(bound.segments)-1].(*TextSegment).Style.Alignment
			}
			yPos += r.layoutRow(rowItems, rowAlign, left, yPos, lineWidth)
			if !inline && bound.end == len(bound.segments[len(bound.segments)-1].(*TextSegment).Text) && i < len(objs)-1 {
				yPos += theme.Padding()
			}
			rowItems = nil
		}
	}
}

// MinSize calculates the minimum size of a rich text widget.
// This is based on the contained text with a standard amount of padding added.
func (r *textRenderer) MinSize() fyne.Size {
	r.obj.propertyLock.RLock()
	bounds := r.obj.rowBounds
	wrap := r.obj.Wrapping
	objs := r.Objects()
	r.obj.propertyLock.RUnlock()

	height := float32(0)
	width := float32(0)
	rowHeight := float32(0)
	rowWidth := float32(0)

	i := 0
	for _, bound := range bounds {
		for range bound.segments {
			obj := objs[i]
			i++

			min := obj.MinSize()
			rowHeight = fyne.Max(rowHeight, min.Height)
			rowWidth += min.Width
		}

		if wrap == fyne.TextWrapOff {
			width = fyne.Max(width, rowWidth)
		}
		height += rowHeight
		rowHeight = 0
		rowWidth = 0
	}

	if height == 0 {
		charMinSize := r.obj.charMinSize(false, fyne.TextStyle{})
		height = charMinSize.Height
	}
	return fyne.NewSize(width, height).
		Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*4).Subtract(r.obj.inset).Subtract(r.obj.inset))
}

func (r *textRenderer) Refresh() {
	r.obj.propertyLock.RLock()
	bounds := r.obj.rowBounds
	r.obj.propertyLock.RUnlock()

	var objs []fyne.CanvasObject
	for _, bound := range bounds {
		for i, seg := range bound.segments {
			if _, ok := seg.(*TextSegment); !ok {
				objs = append(objs, seg.Visual())
				continue
			}

			txt := seg.Visual().(*canvas.Text)
			textSeg := seg.(*TextSegment)
			runes := []rune(textSeg.Text)

			if i == 0 {
				if len(bound.segments) == 1 {
					txt.Text = string(runes[bound.begin:bound.end])
				} else {
					txt.Text = string(runes[bound.begin:])
				}
			} else if i == len(bound.segments)-1 && len(bound.segments) > 1 {
				txt.Text = string(runes[:bound.end])
			}
			if concealed(seg) {
				txt.Text = strings.Repeat(passwordChar, len(runes))
			}

			objs = append(objs, txt)
		}
	}

	r.obj.propertyLock.Lock()
	r.SetObjects(objs)
	r.obj.propertyLock.Unlock()

	r.Layout(r.obj.Size())
	canvas.Refresh(r.obj)
}

func (r *textRenderer) layoutRow(texts []fyne.CanvasObject, align fyne.TextAlign, xPos, yPos, lineWidth float32) float32 {
	if len(texts) == 1 {
		texts[0].Resize(fyne.NewSize(lineWidth, texts[0].MinSize().Height))
		texts[0].Move(fyne.NewPos(xPos, yPos))
		return texts[0].MinSize().Height
	}
	height := float32(0)
	for _, text := range texts {
		size := text.MinSize()

		text.Resize(size)
		text.Move(fyne.NewPos(xPos, yPos)) // TODO also baseline align for height (need new measure info)
		xPos += size.Width
		height = fyne.Max(height, size.Height)
	}
	spare := lineWidth - xPos
	switch align {
	case fyne.TextAlignTrailing:
		first := texts[0]
		first.Resize(fyne.NewSize(first.Size().Width+spare, height))
		setAlign(first, fyne.TextAlignTrailing)

		for _, text := range texts[1:] {
			text.Move(text.Position().Add(fyne.NewPos(spare, 0)))
		}
	case fyne.TextAlignCenter:
		pad := spare / 2
		first := texts[0]
		first.Resize(fyne.NewSize(first.Size().Width+pad, height))
		setAlign(first, fyne.TextAlignTrailing)
		last := texts[len(texts)-1]
		last.Resize(fyne.NewSize(last.Size().Width+pad, height))
		setAlign(last, fyne.TextAlignLeading)

		for _, text := range texts[1:] {
			text.Move(text.Position().Add(fyne.NewPos(pad, 0)))
		}
	default:
		last := texts[len(texts)-1]
		last.Resize(fyne.NewSize(last.Size().Width+spare, height))
		setAlign(last, fyne.TextAlignLeading)
	}

	return height
}

// binarySearch accepts a function that checks if the text width less the maximum width and the start and end rune index
// binarySearch returns the index of rune located as close to the maximum line width as possible
func binarySearch(lessMaxWidth func(int, int) bool, low int, maxHigh int) int {
	if low >= maxHigh {
		return low
	}
	if lessMaxWidth(low, maxHigh) {
		return maxHigh
	}
	high := low
	delta := maxHigh - low
	for delta > 0 {
		delta /= 2
		if lessMaxWidth(low, high+delta) {
			high += delta
		}
	}
	for (high < maxHigh) && lessMaxWidth(low, high+1) {
		high++
	}
	return high
}

// concealed returns true if the segment represents a password, meaning the text should be obscured.
func concealed(seg RichTextSegment) bool {
	if text, ok := seg.(*TextSegment); ok {
		return text.Style.concealed
	}

	return false
}

// findSpaceIndex accepts a slice of runes and a fallback index
// findSpaceIndex returns the index of the last space in the text, or fallback if there are no spaces
func findSpaceIndex(text []rune, fallback int) int {
	curIndex := fallback
	for ; curIndex >= 0; curIndex-- {
		if unicode.IsSpace(text[curIndex]) {
			break
		}
	}
	if curIndex < 0 {
		return fallback
	}
	return curIndex
}

// lineBounds accepts a slice of Segments, a wrapping mode, a maximum line width and a function to measure line width.
// lineBounds returns a slice containing the boundary metadata of each line with the given wrapping applied.
func lineBounds(seg *TextSegment, wrap fyne.TextWrap, firstWidth, maxWidth float32, measurer func([]rune) float32) []rowBoundary {
	lines := splitLines(seg)
	if maxWidth <= 0 || wrap == fyne.TextWrapOff {
		return lines
	}

	measureWidth := firstWidth
	text := []rune(seg.Text)
	checker := func(low int, high int) bool {
		return measurer(text[low:high]) <= measureWidth
	}

	var bounds []rowBoundary
	for _, l := range lines {
		low := l.begin
		high := l.end
		if low == high {
			bounds = append(bounds, l)
			continue
		}
		switch wrap {
		case fyne.TextTruncate:
			high = binarySearch(checker, low, high)
			bounds = append(bounds, rowBoundary{[]RichTextSegment{seg}, low, high})
		case fyne.TextWrapBreak:
			for low < high {
				if measurer(text[low:high]) <= measureWidth {
					bounds = append(bounds, rowBoundary{[]RichTextSegment{seg}, low, high})
					low = high
					high = l.end
					measureWidth = maxWidth
				} else {
					high = binarySearch(checker, low, high)
				}
			}
		case fyne.TextWrapWord:
			for low < high {
				sub := text[low:high]
				if measurer(sub) <= measureWidth {
					bounds = append(bounds, rowBoundary{[]RichTextSegment{seg}, low, high})
					low = high
					high = l.end
					if low < high && unicode.IsSpace(text[low]) {
						low++
					}
					measureWidth = maxWidth
				} else {
					oldHigh := high
					last := low + len(sub) - 1
					fallback := binarySearch(checker, low, last) - low
					high = low + findSpaceIndex(sub, fallback)
					if high == fallback && measurer(sub) <= maxWidth { // add a newline as there is more space on next
						bounds = append(bounds, rowBoundary{[]RichTextSegment{seg}, low, low})
						high = oldHigh
						measureWidth = maxWidth
						continue
					}
				}
			}
		}
	}
	return bounds
}

func setAlign(obj fyne.CanvasObject, align fyne.TextAlign) {
	if text, ok := obj.(*canvas.Text); ok {
		text.Alignment = align
		return
	}
	if c, ok := obj.(*fyne.Container); ok {
		wid := c.Objects[0]
		if link, ok := wid.(*Hyperlink); ok {
			link.Alignment = align
		}
	}
}

// splitLines accepts a text segment and returns a slice of boundary metadata denoting the
// start and end indices of each line delimited by the newline character.
func splitLines(seg *TextSegment) []rowBoundary {
	var low, high int
	var lines []rowBoundary
	text := []rune(seg.Text)
	length := len(text)
	for i := 0; i < length; i++ {
		if text[i] == '\n' {
			high = i
			lines = append(lines, rowBoundary{[]RichTextSegment{seg}, low, high})
			low = i + 1
		}
	}
	return append(lines, rowBoundary{[]RichTextSegment{seg}, low, length})
}

type rowBoundary struct {
	segments   []RichTextSegment
	begin, end int
}

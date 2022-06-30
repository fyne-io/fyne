package widget

import (
	"image/color"
	"math"
	"strings"
	"sync"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

const (
	passwordChar = "•"
)

// RichText represents the base element for a rich text-based widget.
//
// Since: 2.1
type RichText struct {
	BaseWidget
	Segments []RichTextSegment
	Wrapping fyne.TextWrap
	Scroll   widget.ScrollDirection

	inset     fyne.Size     // this varies due to how the widget works (entry with scroller vs others with padding)
	rowBounds []rowBoundary // cache for boundaries
	scr       *widget.Scroll
	prop      *canvas.Rectangle // used to apply text minsize to the scroller `scr`, if present - TODO improve #2464

	visualCache map[RichTextSegment][]fyne.CanvasObject
	cacheLock   sync.Mutex
}

// NewRichText returns a new RichText widget that renders the given text and segments.
// If no segments are specified it will be converted to a single segment using the default text settings.
//
// Since: 2.1
func NewRichText(segments ...RichTextSegment) *RichText {
	t := &RichText{Segments: segments}
	t.Scroll = widget.ScrollNone
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
	if t.scr == nil && t.Scroll != widget.ScrollNone {
		t.prop = canvas.NewRectangle(color.Transparent)
		t.scr = widget.NewScroll(&fyne.Container{Layout: layout.NewMaxLayout(), Objects: []fyne.CanvasObject{
			t.prop, &fyne.Container{}}})
	}

	t.ExtendBaseWidget(t)
	r := &textRenderer{obj: t}

	t.updateRowBounds() // set up the initial text layout etc
	r.Refresh()
	return r
}

// MinSize calculates the minimum size of a rich text widget.
// This is based on the contained text with a standard amount of padding added.
func (t *RichText) MinSize() fyne.Size {
	t.ExtendBaseWidget(t)

	return t.BaseWidget.MinSize()
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
	start := 0
	var ret []rune
	deleting := false
	var segs []RichTextSegment
	for i, seg := range t.Segments {
		if _, ok := seg.(*TextSegment); !ok {
			if !deleting {
				segs = append(segs, seg)
			}
			continue
		}
		end := start + len([]rune(seg.(*TextSegment).Text))
		if end < lowBound {
			segs = append(segs, seg)
			start = end
			continue
		}

		startOff := int(math.Max(float64(lowBound-start), 0))
		endOff := int(math.Min(float64(end), float64(highBound))) - start
		deleted := make([]rune, endOff-startOff)
		r := ([]rune)(seg.(*TextSegment).Text)
		copy(deleted, r[startOff:endOff])
		ret = append(ret, deleted...)
		r2 := append(r[:startOff], r[endOff:]...)
		seg.(*TextSegment).Text = string(r2)
		segs = append(segs, seg)

		// prepare next iteration
		start = end
		if start >= highBound {
			segs = append(segs, t.Segments[i+1:]...)
			break
		} else if start >= lowBound {
			deleting = true
		}
	}
	t.Segments = segs
	t.Refresh()
	return string(ret)
}

// cachedSegmentVisual returns a cached segment visual representation.
// The offset value is > 0 if the segment had been split and so we need multiple objects.
func (t *RichText) cachedSegmentVisual(seg RichTextSegment, offset int) fyne.CanvasObject {
	t.cacheLock.Lock()
	defer t.cacheLock.Unlock()
	if t.visualCache == nil {
		t.visualCache = make(map[RichTextSegment][]fyne.CanvasObject)
	}

	if vis, ok := t.visualCache[seg]; ok && offset < len(vis) {
		return vis[offset]
	}

	vis := seg.Visual()
	if offset < len(t.visualCache[seg]) {
		t.visualCache[seg][offset] = vis
	} else {
		t.visualCache[seg] = append(t.visualCache[seg], vis)
	}
	return vis
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
			size = t.cachedSegmentVisual(seg, bound.firstSegmentReuse).MinSize()
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
	maxWidth := t.size.Width - 4*theme.Padding() + 2*t.inset.Width
	wrapWidth := maxWidth

	var iterateSegments func(segList []RichTextSegment)
	iterateSegments = func(segList []RichTextSegment) {
		var currentBound *rowBoundary
		for _, seg := range segList {
			if parent, ok := seg.(RichTextBlock); ok {
				iterateSegments(parent.Segments())
				if !seg.Inline() {
					wrapWidth = maxWidth
				}
				continue
			}
			if _, ok := seg.(*TextSegment); !ok {
				if currentBound == nil {
					bound := rowBoundary{segments: []RichTextSegment{seg}}
					bounds = append(bounds, bound)
					currentBound = &bound
				} else {
					bounds[len(bounds)-1].segments = append(bounds[len(bounds)-1].segments, seg)
				}
				if seg.Inline() {
					wrapWidth -= t.cachedSegmentVisual(seg, 0).MinSize().Width
				} else {
					wrapWidth = maxWidth
					currentBound = nil
				}
				continue
			}
			textSeg := seg.(*TextSegment)
			textStyle := textSeg.Style.TextStyle
			textSize := textSeg.size()

			leftPad := float32(0)
			if textSeg.Style == RichTextStyleBlockquote {
				leftPad = theme.Padding() * 4
			}
			retBounds := lineBounds(textSeg, t.Wrapping, wrapWidth-leftPad, maxWidth, func(text []rune) float32 {
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
				runes := []rune(textSeg.Text)
				// check ranges - as we resize it can be wrong?
				if begin > len(runes) {
					begin = len(runes)
				}
				end := last.end
				if end > len(runes) {
					end = len(runes)
				}
				text := string(runes[begin:end])
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
	}

	iterateSegments(t.Segments)
	t.propertyLock.RUnlock()

	t.propertyLock.Lock()
	t.rowBounds = bounds
	t.propertyLock.Unlock()
}

// RichTextBlock is an extension of a text segment that contains other segments
//
// Since: 2.1
type RichTextBlock interface {
	Segments() []RichTextSegment
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
	if r.obj.scr != nil {
		r.obj.scr.Resize(size)
		objs = r.obj.scr.Content.(*fyne.Container).Objects[1].(*fyne.Container).Objects
	}
	r.obj.propertyLock.RUnlock()

	left := theme.Padding()*2 - r.obj.inset.Width
	yPos := theme.Padding()*2 - r.obj.inset.Height
	lineWidth := size.Width - left*2
	var rowItems []fyne.CanvasObject
	rowAlign := fyne.TextAlignLeading
	i := 0
	for row, bound := range bounds {
		for segI := range bound.segments {
			if i == len(objs) {
				break // Refresh may not have created all objects for all rows yet...
			}
			inline := segI < len(bound.segments)-1
			obj := objs[i]
			i++
			_, isText := obj.(*canvas.Text)
			if !isText && !inline {
				if len(rowItems) != 0 {
					width, _ := r.layoutRow(rowItems, rowAlign, left, yPos, lineWidth)
					left += width
				}
				height := obj.MinSize().Height

				obj.Move(fyne.NewPos(left, yPos))
				obj.Resize(fyne.NewSize(lineWidth, height))
				yPos += height + theme.Padding()
				continue
			}
			rowItems = append(rowItems, obj)
			if inline {
				continue
			}

			leftPad := float32(0)
			if text, ok := bound.segments[0].(*TextSegment); ok {
				rowAlign = text.Style.Alignment
				if text.Style == RichTextStyleBlockquote {
					leftPad = theme.Padding() * 4
				}
			} else if link, ok := bound.segments[0].(*HyperlinkSegment); ok {
				rowAlign = link.Alignment
			}
			_, y := r.layoutRow(rowItems, rowAlign, left+leftPad, yPos, lineWidth-leftPad)
			yPos += y
			rowItems = nil
		}

		lastSeg := bound.segments[len(bound.segments)-1]
		if !lastSeg.Inline() && row < len(bounds)-1 && bounds[row+1].segments[0] != lastSeg { // ignore wrapped lines etc
			yPos += theme.Padding()
		}
	}
}

// MinSize calculates the minimum size of a rich text widget.
// This is based on the contained text with a standard amount of padding added.
func (r *textRenderer) MinSize() fyne.Size {
	r.obj.propertyLock.RLock()
	bounds := r.obj.rowBounds
	wrap := r.obj.Wrapping
	scroll := r.obj.Scroll
	objs := r.Objects()
	if r.obj.scr != nil {
		objs = r.obj.scr.Content.(*fyne.Container).Objects[1].(*fyne.Container).Objects
	}
	r.obj.propertyLock.RUnlock()

	height := float32(0)
	width := float32(0)
	rowHeight := float32(0)
	rowWidth := float32(0)

	i := 0
	for row, bound := range bounds {
		for range bound.segments {
			if i == len(objs) {
				break // Refresh may not have created all objects for all rows yet...
			}
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

		lastSeg := bound.segments[len(bound.segments)-1]
		if !lastSeg.Inline() && row < len(bounds)-1 && bounds[row+1].segments[0] != lastSeg { // ignore wrapped lines etc
			height += theme.Padding()
		}
	}

	if height == 0 {
		charMinSize := r.obj.charMinSize(false, fyne.TextStyle{})
		height = charMinSize.Height
	}
	min := fyne.NewSize(width, height).
		Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*4).Subtract(r.obj.inset).Subtract(r.obj.inset))

	if r.obj.scr != nil {
		r.obj.prop.SetMinSize(min)
	}

	switch scroll {
	case widget.ScrollBoth:
		return fyne.NewSize(32, 32)
	case widget.ScrollHorizontalOnly:
		return fyne.NewSize(32, min.Height)
	case widget.ScrollVerticalOnly:
		return fyne.NewSize(min.Width, 32)
	default:
		return min
	}
}

func (r *textRenderer) Refresh() {
	r.obj.propertyLock.RLock()
	bounds := r.obj.rowBounds
	scroll := r.obj.Scroll
	r.obj.propertyLock.RUnlock()

	var objs []fyne.CanvasObject
	for _, bound := range bounds {
		for i, seg := range bound.segments {
			if _, ok := seg.(*TextSegment); !ok {
				obj := r.obj.cachedSegmentVisual(seg, 0)
				seg.Update(obj)
				objs = append(objs, obj)
				continue
			}

			obj := r.obj.cachedSegmentVisual(seg, bound.firstSegmentReuse)
			seg.Update(obj)
			txt := obj.(*canvas.Text)
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
	if r.obj.scr != nil {
		r.obj.scr.Content = &fyne.Container{Layout: layout.NewMaxLayout(), Objects: []fyne.CanvasObject{
			r.obj.prop, &fyne.Container{Objects: objs}}}
		r.obj.scr.Direction = scroll
		r.SetObjects([]fyne.CanvasObject{r.obj.scr})
		r.obj.scr.Refresh()
	} else {
		r.SetObjects(objs)
	}
	r.obj.propertyLock.Unlock()

	r.Layout(r.obj.Size())
	canvas.Refresh(r.obj)
}

func (r *textRenderer) layoutRow(texts []fyne.CanvasObject, align fyne.TextAlign, xPos, yPos, lineWidth float32) (float32, float32) {
	initialX := xPos
	if len(texts) == 1 {
		texts[0].Resize(fyne.NewSize(lineWidth, texts[0].MinSize().Height))
		texts[0].Move(fyne.NewPos(xPos, yPos))
		return texts[0].MinSize().Width, texts[0].MinSize().Height
	}
	height := float32(0)
	tallestBaseline := float32(0)
	realign := false
	baselines := make([]float32, len(texts))
	for i, text := range texts {
		var size fyne.Size
		if txt, ok := text.(*canvas.Text); ok {
			s, base := fyne.CurrentApp().Driver().RenderedTextSize(txt.Text, txt.TextSize, txt.TextStyle)
			if base > tallestBaseline {
				if tallestBaseline > 0 {
					realign = true
				}
				tallestBaseline = base
			}
			size = s
			baselines[i] = base
		} else if c, ok := text.(*fyne.Container); ok {
			wid := c.Objects[0]
			if link, ok := wid.(*Hyperlink); ok {
				s, base := fyne.CurrentApp().Driver().RenderedTextSize(link.Text, theme.TextSize(), link.TextStyle)
				if base > tallestBaseline {
					if tallestBaseline > 0 {
						realign = true
					}
					tallestBaseline = base
				}
				size = s
				baselines[i] = base
			}
		}
		if size.IsZero() {
			size = text.MinSize()
		}
		text.Resize(size)
		text.Move(fyne.NewPos(xPos, yPos))

		xPos += size.Width
		if height == 0 {
			height = size.Height
		} else if height != size.Height {
			height = fyne.Max(height, size.Height)
			realign = true
		}
	}

	if realign {
		for i, text := range texts {
			delta := tallestBaseline - baselines[i]
			text.Move(fyne.NewPos(text.Position().X, yPos+delta))
		}
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

	return xPos - initialX, height
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
	if wrap == fyne.TextWrapOff {
		return lines
	}

	measureWidth := firstWidth
	text := []rune(seg.Text)
	checker := func(low int, high int) bool {
		return measurer(text[low:high]) <= measureWidth
	}

	reuse := 0
	var bounds []rowBoundary
	for _, l := range lines {
		low := l.begin
		high := l.end
		if low == high {
			l.firstSegmentReuse = reuse
			reuse++
			bounds = append(bounds, l)
			continue
		}
		switch wrap {
		case fyne.TextTruncate:
			high = binarySearch(checker, low, high)
			bounds = append(bounds, rowBoundary{[]RichTextSegment{seg}, reuse, low, high})
			reuse++
		case fyne.TextWrapBreak:
			for low < high {
				if measurer(text[low:high]) <= measureWidth {
					bounds = append(bounds, rowBoundary{[]RichTextSegment{seg}, reuse, low, high})
					reuse++
					low = high
					high = l.end
					measureWidth = maxWidth
				} else {
					newHigh := binarySearch(checker, low, high)
					if newHigh <= low {
						bounds = append(bounds, rowBoundary{[]RichTextSegment{seg}, reuse, low, low + 1})
						reuse++
						low++
					} else {
						high = newHigh
					}
				}
			}
		case fyne.TextWrapWord:
			for low < high {
				sub := text[low:high]
				subWidth := measurer(sub)
				if subWidth <= measureWidth {
					bounds = append(bounds, rowBoundary{[]RichTextSegment{seg}, reuse, low, high})
					reuse++
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

					if fallback < 1 { // even a character won't fit
						bounds = append(bounds, rowBoundary{[]RichTextSegment{seg}, reuse, low, low + 1})
						low++
						high = low + 1
						reuse++

						if high > l.end {
							return bounds
						}
					} else {
						high = low + findSpaceIndex(sub, fallback)
					}
					if high == fallback && subWidth <= maxWidth { // add a newline as there is more space on next
						bounds = append(bounds, rowBoundary{[]RichTextSegment{seg}, reuse, low, low})
						reuse++
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
			link.Refresh()
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
			lines = append(lines, rowBoundary{[]RichTextSegment{seg}, len(lines), low, high})
			low = i + 1
		}
	}
	return append(lines, rowBoundary{[]RichTextSegment{seg}, len(lines), low, length})
}

type rowBoundary struct {
	segments          []RichTextSegment
	firstSegmentReuse int
	begin, end        int
}

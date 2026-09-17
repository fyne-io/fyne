package widget

import (
	"fmt"
	"image/color"
	"math"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/go-text/typesetting/di"
	"github.com/go-text/typesetting/shaping"
	"golang.org/x/image/math/fixed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	paint "fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

const (
	averageChar  = "M"
	passwordChar = "•"
)

var _ fyne.Widget = (*RichText)(nil)

// RichText represents the base element for a rich text-based widget.
//
// Since: 2.1
type RichText struct {
	BaseWidget
	Segments []RichTextSegment
	Wrapping fyne.TextWrap
	Scroll   fyne.ScrollDirection

	// The truncation mode of the text
	//
	// Since: 2.4
	Truncation fyne.TextTruncation

	inset     fyne.Size     // this varies due to how the widget works (entry with scroller vs others with padding)
	rowBounds []rowBoundary // cache for boundaries
	scr       *widget.Scroll
	prop      *canvas.Rectangle // used to apply text minsize to the scroller `scr`, if present - TODO improve #2464

	visualCache    map[RichTextSegment]visualCacheEntry
	visualCacheGen int64
	minCache       fyne.Size
	geometryValid  bool // whether rowBounds carries up to date row positions
	decor          []rowDecoration

	// highlight is the selection whose rectangles are drawn over this content,
	// above the panels that blocks sit on but below the text
	highlight  *selectable
	highlights []fyne.CanvasObject
}

// rowDecoration is a graphical element drawn behind a run of rows, such as the
// panel that sits under the lines of a code block.
type rowDecoration struct {
	obj      fyne.CanvasObject
	from, to int // the first and last row that this decoration covers
}

// panelSegment is a block segment whose content is drawn on a panel, so that the
// rows of the block appear inside it.
type panelSegment interface {
	RichTextSegment
	RichTextBlock

	panel() fyne.CanvasObject
}

type visualCacheEntry struct {
	gen int64
	obj []fyne.CanvasObject
}

// NewRichText returns a new RichText widget that renders the given text and segments.
// If no segments are specified it will be converted to a single segment using the default text settings.
//
// Since: 2.1
func NewRichText(segments ...RichTextSegment) *RichText {
	t := &RichText{Segments: segments}
	t.Scroll = widget.ScrollNone
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
	t.prop = canvas.NewRectangle(color.Transparent)
	if t.scr == nil && t.Scroll != widget.ScrollNone {
		t.scr = widget.NewScroll(&fyne.Container{Layout: layout.NewStackLayout(), Objects: []fyne.CanvasObject{
			t.prop, &fyne.Container{},
		}})
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
	// We return the minCache here which might be outdated if internal segments were changed.
	// Users must call Refresh() to force an update after any changes to t.
	t.ExtendBaseWidget(t)

	if t.minCache.IsZero() {
		minSize := t.BaseWidget.MinSize()
		t.minCache = minSize
	}
	return t.minCache
}

// Refresh triggers a redraw of the rich text.
func (t *RichText) Refresh() {
	t.minCache = fyne.Size{}
	t.updateRowBounds()

	for _, s := range t.Segments {
		switch seg := s.(type) {
		case *TextSegment:
			seg.parent = t
		case *listMarkerSegment:
			seg.parent = t
		}
	}

	t.BaseWidget.Refresh()
}

// Resize sets a new size for the rich text.
// This should only be called if it is not in a container with a layout manager.
func (t *RichText) Resize(size fyne.Size) {
	if size == t.Size() {
		return
	}

	t.size = size

	skipResize := !t.minCache.IsZero() && size.Width >= t.minCache.Width && size.Height >= t.minCache.Height && t.Wrapping == fyne.TextWrapOff && t.Truncation == fyne.TextTruncateOff

	if skipResize {
		if len(t.Segments) < 2 { // we can simplify :)
			cache.Renderer(t).Layout(size)
			return
		}
	}

	t.Refresh()
}

// String returns the text widget buffer as string
func (t *RichText) String() string {
	ret := strings.Builder{}
	for _, seg := range t.contentSegments() {
		ret.WriteString(seg.Textual())
	}
	return ret.String()
}

// contentSegments returns the segments that carry the content of this rich text,
// in the order that they are laid out. Blocks that keep their content in child
// segments, such as lists, are replaced by the segments that they hold.
func (t *RichText) contentSegments() []RichTextSegment {
	for _, seg := range t.Segments {
		if holdsContent(seg) {
			return appendContentSegments(t.Segments, make([]RichTextSegment, 0, len(t.Segments)+2))
		}
	}

	return t.Segments // the common case of content that is not nested
}

func appendContentSegments(in, out []RichTextSegment) []RichTextSegment {
	for _, seg := range in {
		if holdsContent(seg) {
			out = appendContentSegments(seg.(RichTextBlock).Segments(), out)
			continue
		}

		out = append(out, seg)
	}
	return out
}

// holdsContent reports whether a segment keeps its content in the segments inside
// it. A block such as code, that holds the text itself, does not.
func holdsContent(seg RichTextSegment) bool {
	if _, ok := seg.(textHolder); ok {
		return false
	}

	_, ok := seg.(RichTextBlock)
	return ok
}

// contentIs reports whether the segments spell out exactly the given text.
// It avoids building a string to compare against, as this runs on every edit.
func (t *RichText) contentIs(text string) bool {
	for _, seg := range t.contentSegments() {
		content := seg.Textual()
		if !strings.HasPrefix(text, content) {
			return false
		}
		text = text[len(content):]
	}
	return text == ""
}

// charMinSize returns the average char size to use for internal computation
func (*RichText) charMinSize(concealed bool, style fyne.TextStyle, textSize float32) fyne.Size {
	defaultChar := averageChar
	if concealed {
		defaultChar = passwordChar
	}

	return fyne.MeasureText(defaultChar, textSize, style)
}

// textHolder is content that an editor can add text to and remove text from.
// A text segment is the usual case, a code block holds the lines of code that it
// draws on its panel.
type textHolder interface {
	content() string
	setContent(string)
}

// deleteFromTo removes the text between the specified positions.
// Positions are rune offsets into the whole content, spanning all segments.
func (t *RichText) deleteFromTo(lowBound int, highBound int) []rune {
	if lowBound >= highBound {
		return []rune{}
	}

	start := 0
	ret := make([]rune, 0, highBound-lowBound)
	var dropped []RichTextSegment
	for _, seg := range t.contentSegments() {
		end := start + utf8.RuneCountInString(seg.Textual())

		if end <= lowBound || start >= highBound { // wholly outside the deleted range
			start = end
			continue
		}

		if holder, ok := seg.(textHolder); ok {
			r := []rune(holder.content())
			from := max(lowBound-start, 0)
			to := min(highBound-start, len(r))
			ret = append(ret, r[from:to]...)
			holder.setContent(string(r[:from]) + string(r[to:]))
		} else { // an object cannot be split, so it goes entirely
			ret = append(ret, []rune(seg.Textual())...)
			dropped = append(dropped, seg)
		}
		start = end
	}

	if len(dropped) > 0 {
		t.Segments = removeSegments(t.Segments, dropped)
	}

	t.Refresh()
	return ret
}

// removeSegments takes the given segments out of the content, looking inside the
// blocks that hold them.
func removeSegments(in []RichTextSegment, drop []RichTextSegment) []RichTextSegment {
	out := make([]RichTextSegment, 0, len(in))
	for _, seg := range in {
		if segmentIn(seg, drop) {
			continue
		}

		switch block := seg.(type) {
		case *ListSegment:
			block.Items = removeListItems(block, drop)
			if len(block.Items) == 0 {
				continue
			}
		case *ParagraphSegment:
			block.Texts = removeSegments(block.Texts, drop)
		}

		out = append(out, seg)
	}
	return out
}

// removeListItems drops the items of a list whose bullet was deleted, keeping the
// content of each one by moving it into the item before.
func removeListItems(l *ListSegment, drop []RichTextSegment) []RichTextSegment {
	items := make([]RichTextSegment, 0, len(l.Items))
	for i, item := range l.Items {
		if i > 0 && i < len(l.markers) && segmentIn(l.markers[i], drop) && len(items) > 0 {
			last := items[len(items)-1]
			items[len(items)-1] = &ParagraphSegment{Texts: append(blockContent(last), blockContent(item)...)}
			continue
		}

		items = append(items, removeSegments([]RichTextSegment{item}, drop)...)
	}

	l.markers = nil // the bullets are made again to match the items that are left
	return items
}

// blockContent returns the segments inside a block, or the segment itself if it
// is not one that holds others.
func blockContent(seg RichTextSegment) []RichTextSegment {
	if block, ok := seg.(*ParagraphSegment); ok {
		return block.Texts
	}

	return []RichTextSegment{seg}
}

func segmentIn(seg RichTextSegment, list []RichTextSegment) bool {
	for _, in := range list {
		if in == seg {
			return true
		}
	}
	return false
}

// cachedSegmentVisual returns a cached segment visual representation.
// The offset value is > 0 if the segment had been split and so we need multiple objects.
func (t *RichText) cachedSegmentVisual(seg RichTextSegment, offset int) fyne.CanvasObject {
	if t.visualCache == nil {
		t.visualCache = make(map[RichTextSegment]visualCacheEntry)
	}

	if vis, ok := t.visualCache[seg]; ok && offset < len(vis.obj) {
		return vis.obj[offset]
	}

	vis := seg.Visual()
	if offset < len(t.visualCache[seg].obj) {
		t.visualCache[seg].obj[offset] = vis
	} else {
		entry := t.visualCache[seg]
		entry.obj = append(entry.obj, vis)
		t.visualCache[seg] = entry
	}
	return vis
}

func (t *RichText) cleanVisualCache() {
	if len(t.visualCache) <= len(t.Segments) {
		return
	}

	// mark cache entries that are still valid
	t.visualCacheGen++
	mark := func(seg RichTextSegment) {
		if c, ok := t.visualCache[seg]; ok {
			c.gen = t.visualCacheGen
			t.visualCache[seg] = c
		}
	}
	for _, seg := range t.Segments {
		mark(seg)
	}
	for i := range t.rowBounds { // content inside blocks is only found in the rows
		for _, seg := range t.rowBounds[i].segments {
			mark(seg)
		}
	}

	// delete entries that are not marked as valid
	var deletingSegs []RichTextSegment
	for seg1, c := range t.visualCache {
		if c.gen != t.visualCacheGen {
			deletingSegs = append(deletingSegs, seg1)
		}
	}
	for _, seg := range deletingSegs {
		delete(t.visualCache, seg)
	}
}

// insertAt inserts the text at the specified position.
// The position is a rune offset into the whole content, spanning all segments.
func (t *RichText) insertAt(pos int, runes []rune) {
	// Find best segment if multiple match.
	var contains, empty, before, after textHolder
	beforeEndsLine := false
	offset, start := 0, 0
	for _, seg := range t.contentSegments() {
		if start > pos {
			break
		}
		end := start + utf8.RuneCountInString(seg.Textual())

		if holder, ok := seg.(textHolder); ok {
			switch {
			case start < pos && pos < end:
				if contains == nil {
					contains, offset = holder, pos-start
				}
			case start == pos && end == pos:
				if empty == nil {
					empty = holder
				}
			case end == pos:
				before = holder
				// content that closes a line hands this position to what follows it
				beforeEndsLine = !seg.Inline() || strings.HasSuffix(holder.content(), "\n")
			case start == pos:
				if after == nil {
					after = holder
				}
			}
		} else if start == pos && end == pos {
			// a decoration such as a list bullet sits here, so text typed at this
			// position belongs to the content that follows it
			empty, before = nil, nil
		}
		start = end
	}

	into := contains
	switch {
	case into != nil:
	case empty != nil:
		into, offset = empty, 0
	case before != nil && !beforeEndsLine:
		into, offset = before, utf8.RuneCountInString(before.content())
	case after != nil:
		into, offset = after, 0
	case before != nil: // there is nothing after it, so the line grows instead
		into, offset = before, utf8.RuneCountInString(before.content())
	}

	if into == nil { // no text segment covers this position, so start a new one
		t.Segments = append(t.Segments, &TextSegment{Style: RichTextStyleInline, Text: string(runes), parent: t})
		return
	}

	r := []rune(into.content())
	offset = min(offset, len(r))
	r2 := make([]rune, 0, len(r)+len(runes))
	r2 = append(r2, r[:offset]...)
	r2 = append(r2, runes...)
	r2 = append(r2, r[offset:]...)
	into.setContent(string(r2))
}

// Len returns the text widget buffer length
func (t *RichText) len() int {
	ret := 0
	for _, seg := range t.contentSegments() {
		ret += utf8.RuneCountInString(seg.Textual())
	}
	return ret
}

// lineSizeToColumn returns the rendered size for the line specified by row up to the col position
func (t *RichText) lineSizeToColumn(col, row int, textSize, innerPad float32) fyne.Size {
	if row < 0 {
		row = 0
	}
	if col < 0 {
		col = 0
	}
	bound := t.rowBoundary(row)
	total := fyne.NewSize(0, 0)
	counted := 0
	last := false
	if bound == nil {
		return t.charMinSize(false, fyne.TextStyle{}, textSize)
	}

	leftPad, _ := t.rowPaddingAndAlign(*bound, t.Theme().Size(theme.SizeNameLineSpacing), fyne.TextAlignLeading)
	for i, seg := range bound.segments {
		var size fyne.Size
		measureText := rowSegmentRunes(bound, i)
		partial := false
		if col < counted+len(measureText) {
			measureText = measureText[0 : col-counted]
			partial = true
			last = true
		}
		counted += len(measureText)

		if text, ok := seg.(*TextSegment); ok {
			if concealed(seg) {
				measureText = []rune(strings.Repeat(passwordChar, len(measureText)))
			}

			size, _ = fyne.CurrentApp().Driver().RenderedTextSize(string(measureText), text.size(), text.Style.TextStyle, nil)
		} else if link, ok := seg.(*HyperlinkSegment); ok {
			sizeName := link.SizeName
			if sizeName == "" {
				sizeName = theme.SizeNameText
			}
			size, _ = fyne.CurrentApp().Driver().RenderedTextSize(string(measureText), t.Theme().Size(sizeName), link.TextStyle, nil)
		} else if partial {
			size = fyne.Size{} // the cursor is before this object, so it adds no width
		} else {
			size = t.cachedSegmentVisual(seg, 0).MinSize()
		}

		total.Width += size.Width
		total.Height = fyne.Max(total.Height, size.Height)
		if last {
			break
		}
	}
	return total.Add(fyne.NewSize(innerPad-t.inset.Width+leftPad, 0))
}

// Row returns the characters in the row specified.
// The row parameter should be between 0 and t.Rows()-1.
func (t *RichText) row(row int) []rune {
	if row < 0 || row >= t.rows() {
		return nil
	}
	bound := &t.rowBounds[row]
	var ret []rune
	for i := range bound.segments {
		ret = append(ret, rowSegmentRunes(bound, i)...)
	}
	return ret
}

// RowBoundary returns the boundary of the row specified.
// The row parameter should be between 0 and t.Rows()-1.
func (t *RichText) rowBoundary(row int) *rowBoundary {
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
	if t.rowBounds == nil { // if the widget API is used before it is shown
		t.updateRowBounds()
	}
	return len(t.rowBounds)
}

// updateRowBounds updates the row bounds used to render properly the text widget.
// updateRowBounds should be invoked every time a segment Text, widget Wrapping or size changes.
func (t *RichText) updateRowBounds() {
	b := newRowBoundsBuilder(t)
	b.walk(t.Segments, 0)

	t.rowBounds = b.bounds
	t.geometryValid = false // calculated on demand, as this runs on every edit
}

// segmentsEndRow reports whether these segments closed the row that they were laid
// out on, so that whatever follows them starts a new one.
func segmentsEndRow(segs []RichTextSegment) bool {
	if len(segs) == 0 {
		return false
	}

	last := segs[len(segs)-1]
	if block, ok := last.(RichTextBlock); ok {
		inner := block.Segments()
		if segmentsEndRow(inner) {
			return true
		}

		// a block stands on its own line unless its content ended the row already
		return !endsWithNewline(inner)
	}
	return !last.Inline()
}

// markPanelRows records that the rows a block just added are drawn on its panel.
// The block may have started on a row that was already open, so the row before
// the ones it added is included when it holds content of the block.
func markPanelRows(bounds []rowBoundary, first int, segs []RichTextSegment, block panelSegment) {
	if first > 0 && rowHoldsAny(&bounds[first-1], segs) {
		first--
	}

	for i := first; i < len(bounds); i++ {
		bounds[i].panel = block
	}
}

// rowHoldsAny reports whether any of the given segments puts content on a row.
func rowHoldsAny(bound *rowBoundary, segs []RichTextSegment) bool {
	for _, in := range bound.segments {
		for _, seg := range segs {
			if in == seg {
				return true
			}
		}
	}
	return false
}

// ensureRowGeometry calculates the row positions if they are not already known.
func (t *RichText) ensureRowGeometry() {
	if t.geometryValid {
		return
	}

	t.geometryValid = true
	t.updateRowGeometry()
}

// updateRowGeometry records the vertical offset and height of each row.
func (t *RichText) updateRowGeometry() {
	if t.uniformRowGeometry() {
		return
	}

	th := t.Theme()
	lineSpacing := th.Size(theme.SizeNameLineSpacing)
	textSize := th.Size(theme.SizeNameText)

	yPos := float32(0)
	for i := range t.rowBounds {
		bound := &t.rowBounds[i]
		height := float32(0)
		for j, seg := range bound.segments {
			var segHeight float32
			switch s := seg.(type) {
			case *TextSegment:
				segHeight = fyne.MeasureText(string(rowSegmentRunes(bound, j)), s.size(), s.Style.TextStyle).Height
			case *HyperlinkSegment:
				sizeName := s.SizeName
				if sizeName == "" {
					sizeName = theme.SizeNameText
				}
				segHeight = fyne.MeasureText(string(rowSegmentRunes(bound, j)), th.Size(sizeName), s.TextStyle).Height
			default:
				segHeight = t.cachedSegmentVisual(seg, 0).MinSize().Height
			}
			height = fyne.Max(height, segHeight)
		}
		if height == 0 {
			height = fyne.MeasureText(averageChar, textSize, fyne.TextStyle{}).Height
		}

		bound.yPos = yPos
		bound.height = height
		yPos += height

		lastSeg := bound.segments[len(bound.segments)-1]
		if !lastSeg.Inline() && i < len(t.rowBounds)-1 && t.rowBounds[i+1].segments[0] != lastSeg {
			yPos += lineSpacing
		}
	}
}

// uniformRowGeometry handles the common case of content that is one run of text,
// returning true if the geometry was calculated.
func (t *RichText) uniformRowGeometry() bool {
	if len(t.Segments) != 1 {
		return false
	}
	text, ok := t.Segments[0].(*TextSegment)
	if !ok {
		return false
	}

	height := fyne.MeasureText(averageChar, text.size(), text.Style.TextStyle).Height
	yPos := float32(0)
	for i := range t.rowBounds {
		t.rowBounds[i].yPos = yPos
		t.rowBounds[i].height = height
		yPos += height
	}
	return true
}

// rowFirstVisibleSegment returns the first segment that puts content on this row.
func rowFirstVisibleSegment(bound *rowBoundary) RichTextSegment {
	for i, seg := range bound.segments {
		if i == 0 && len(bound.segments) > 1 && bound.segBegin >= utf8.RuneCountInString(seg.Textual()) {
			continue
		}
		return seg
	}
	return nil
}

// rowAt returns the index of the row rendered at the specified vertical offset.
func (t *RichText) rowAt(y float32) int {
	rows := t.rows()
	t.ensureRowGeometry()

	return sort.Search(rows, func(i int) bool {
		bound := &t.rowBounds[i]
		return y < bound.yPos+bound.height
	})
}

// rowGeometry returns the vertical offset and height of the specified row.
func (t *RichText) rowGeometry(row int) (y, height float32) {
	if row < 0 || row >= t.rows() {
		th := t.Theme()
		return 0, fyne.MeasureText(averageChar, th.Size(theme.SizeNameText), fyne.TextStyle{}).Height
	}
	t.ensureRowGeometry()

	bound := &t.rowBounds[row]
	if bound.height == 0 {
		th := t.Theme()
		return bound.yPos, fyne.MeasureText(averageChar, th.Size(theme.SizeNameText), fyne.TextStyle{}).Height
	}
	return bound.yPos, bound.height
}

// rowSegmentRunes returns the runes of the segment at index i within the given
// row that are visible on that row.
func rowSegmentRunes(bound *rowBoundary, i int) []rune {
	runes := []rune(bound.segments[i].Textual())
	last := len(bound.segments) - 1

	if i == 0 {
		begin := min(bound.segBegin, len(runes))
		if last == 0 {
			return runes[begin:max(min(bound.segEnd, len(runes)), begin)]
		}
		return runes[begin:]
	}
	if i == last {
		return runes[:min(bound.segEnd, len(runes))]
	}
	return runes
}

// RichTextBlock is an extension of a text segment that contains other segments
//
// Since: 2.1
type RichTextBlock interface {
	Segments() []RichTextSegment
}

// setHighlights records the rectangles that a selection wants drawn over this
// content.
func (t *RichText) setHighlights(sel *selectable, objs []fyne.CanvasObject) {
	t.highlight, t.highlights = sel, objs

	canvas.Refresh(t.super())
}

// highlightObjects returns the selection rectangles to draw over this content.
func (t *RichText) highlightObjects() []fyne.CanvasObject {
	if t.highlight == nil || !t.highlight.selecting {
		return nil
	}

	return t.highlights
}

// updateDecorations prepares the graphical elements drawn behind rows. It returns the
// objects to draw, which are added before the text so that they appear behind it.
func (t *RichText) updateDecorations() []fyne.CanvasObject {
	var decor []rowDecoration
	var current panelSegment
	for row := range t.rowBounds {
		block := t.rowBounds[row].panel
		if block == nil {
			current = nil
			continue
		}

		if block == current { // another row of the block we are already drawing
			decor[len(decor)-1].to = row
			continue
		}

		current = block
		decor = append(decor, rowDecoration{obj: block.panel(), from: row, to: row})
	}

	objs := make([]fyne.CanvasObject, len(decor))
	for i, d := range decor {
		objs[i] = d.obj
		d.obj.Refresh()
	}
	t.decor = decor
	return objs
}

// Renderer
type textRenderer struct {
	widget.BaseRenderer
	obj *RichText
}

// codeInlineText returns the text inside an inline-code container, identified by
// its codeInlineLayout, or a bare *canvas.Text as-is. It returns false for any
// other object.
func codeInlineText(obj fyne.CanvasObject) (*canvas.Text, bool) {
	switch o := obj.(type) {
	case *canvas.Text:
		return o, true
	case *fyne.Container:
		if _, ok := o.Layout.(*codeInlineLayout); ok {
			t, _ := o.Objects[1].(*canvas.Text)
			return t, true
		}
	}
	return nil, false
}

// textObjects returns the visuals of the rendered segments, leaving out the
// decorations that are drawn behind them.
func (r *textRenderer) textObjects() []fyne.CanvasObject {
	objs := r.BaseRenderer.Objects()
	if r.obj.scr != nil {
		objs = r.obj.scr.Content.(*fyne.Container).Objects[1].(*fyne.Container).Objects
	}

	if len(r.obj.decor) > len(objs) { // a refresh has not caught up with the decorations yet
		return nil
	}
	return objs[len(r.obj.decor):]
}

// Objects returns the visuals of this rich text, with any selection highlights
// slotted in above the panels that blocks sit on and below the text itself.
func (r *textRenderer) Objects() []fyne.CanvasObject {
	objs := r.BaseRenderer.Objects()
	if len(r.obj.decor) == 0 {
		return objs
	}

	highlights := r.obj.highlightObjects()
	if len(highlights) == 0 {
		return objs
	}

	at := min(len(r.obj.decor), len(objs))
	out := make([]fyne.CanvasObject, 0, len(objs)+len(highlights))
	out = append(out, objs[:at]...)
	out = append(out, highlights...)
	return append(out, objs[at:]...)
}

func (r *textRenderer) Layout(size fyne.Size) {
	th := r.obj.Theme()
	bounds := r.obj.rowBounds
	if r.obj.scr != nil {
		r.obj.scr.Resize(size)
	}
	objs := r.textObjects()

	// Accessing theme here is slow, so we cache the value
	innerPadding := th.Size(theme.SizeNameInnerPadding)
	lineSpacing := th.Size(theme.SizeNameLineSpacing)

	xInset := innerPadding - r.obj.inset.Width
	left := xInset
	yPos := innerPadding - r.obj.inset.Height
	lineWidth := size.Width - left*2
	var rowItems []fyne.CanvasObject
	rowAlign := fyne.TextAlignLeading
	i := 0
	for row, bound := range bounds {
		leftPad, align := r.obj.rowPaddingAndAlign(bound, lineSpacing, rowAlign)
		rowAlign = align
		rowY := yPos

		for segI := range bound.segments {
			if i == len(objs) {
				break // Refresh may not have created all objects for all rows yet...
			}
			inline := segI < len(bound.segments)-1
			obj := objs[i]
			i++
			_, isText := codeInlineText(obj) // code-inline containers are text-like, not blocks
			if !isText && !inline {
				if len(rowItems) != 0 {
					width, _ := r.layoutRow(rowItems, rowAlign, left+leftPad, yPos, lineWidth-leftPad)
					left += width
					rowItems = nil
				}
				height := obj.MinSize().Height

				obj.Move(fyne.NewPos(left+leftPad, yPos))
				obj.Resize(fyne.NewSize(lineWidth-leftPad, height))
				yPos += height
				left = xInset
				continue
			}
			rowItems = append(rowItems, obj)
			if inline {
				continue
			}

			_, y := r.layoutRow(rowItems, rowAlign, left+leftPad, yPos, lineWidth-leftPad)
			yPos += y
			rowItems = nil
		}

		// record where this row landed so a cursor or selection can be placed quickly.
		bounds[row].yPos = rowY - (innerPadding - r.obj.inset.Height)
		bounds[row].height = yPos - rowY

		lastSeg := bound.segments[len(bound.segments)-1]
		if !lastSeg.Inline() && row < len(bounds)-1 && bounds[row+1].segments[0] != lastSeg { // ignore wrapped lines etc
			yPos += lineSpacing
		}
	}
	r.obj.geometryValid = true

	r.layoutDecorations(bounds, xInset, lineWidth, innerPadding-r.obj.inset.Height, lineSpacing)
}

// layoutDecorations places the elements drawn behind rows, now that the rows
// they cover have been positioned.
func (r *textRenderer) layoutDecorations(bounds []rowBoundary, xInset, lineWidth, yOffset, lineSpacing float32) {
	if len(r.obj.decor) == 0 {
		return
	}

	inset := r.obj.Theme().Size(theme.SizeNameInnerPadding)
	for _, d := range r.obj.decor {
		if d.to >= len(bounds) {
			continue
		}

		top, bottom := bounds[d.from], bounds[d.to]
		leftPad, _ := r.obj.rowPaddingAndAlign(top, lineSpacing, fyne.TextAlignLeading)
		leftPad -= inset // the panel surrounds the text rather than starting at it

		d.obj.Move(fyne.NewPos(xInset+leftPad, top.yPos+yOffset-lineSpacing/2))
		d.obj.Resize(fyne.NewSize(lineWidth-leftPad, bottom.yPos+bottom.height-top.yPos+lineSpacing))
	}
}

// MinSize calculates the minimum size of a rich text widget.
// This is based on the contained text with a standard amount of padding added.
func (r *textRenderer) MinSize() fyne.Size {
	th := r.obj.Theme()
	textSize := th.Size(theme.SizeNameText)
	innerPad := th.Size(theme.SizeNameInnerPadding)

	bounds := r.obj.rowBounds
	wrap := r.obj.Wrapping
	trunc := r.obj.Truncation
	scroll := r.obj.Scroll
	objs := r.textObjects()

	charMinSize := r.obj.charMinSize(false, fyne.TextStyle{}, textSize)
	minSize := r.calculateMin(bounds, wrap, objs, charMinSize, th)
	if r.obj.scr != nil {
		r.obj.prop.SetMinSize(minSize)
	}

	if trunc != fyne.TextTruncateOff && r.obj.Scroll == widget.ScrollNone {
		minBounds := charMinSize
		if wrap == fyne.TextWrapOff {
			minBounds.Height = minSize.Height
		} else {
			minBounds = minBounds.Add(fyne.NewSquareSize(innerPad * 2).Subtract(r.obj.inset).Subtract(r.obj.inset))
		}
		if trunc == fyne.TextTruncateClip {
			return minBounds
		}
		if trunc == fyne.TextTruncateEllipsis {
			ellipsisSize := fyne.MeasureText("…", th.Size(theme.SizeNameText), fyne.TextStyle{})
			return minBounds.AddWidthHeight(ellipsisSize.Width, 0)
		}
	}

	const minScrolledSize = 32
	switch scroll {
	case widget.ScrollBoth:
		return fyne.NewSize(minScrolledSize, minScrolledSize)
	case widget.ScrollHorizontalOnly:
		return fyne.NewSize(minScrolledSize, minSize.Height)
	case widget.ScrollVerticalOnly:
		return fyne.NewSize(minSize.Width, minScrolledSize)
	default:
		return minSize
	}
}

func (r *textRenderer) calculateMin(bounds []rowBoundary, wrap fyne.TextWrap, objs []fyne.CanvasObject,
	charMinSize fyne.Size, th fyne.Theme,
) fyne.Size {
	height := float32(0)
	width := float32(0)
	rowHeight := float32(0)
	rowWidth := float32(0)
	trunc := r.obj.Truncation
	innerPad := th.Size(theme.SizeNameInnerPadding)

	// Accessing the theme here is slow, so we cache the value
	lineSpacing := th.Size(theme.SizeNameLineSpacing)

	i := 0
	for row, bound := range bounds {
		for range bound.segments {
			if i == len(objs) {
				break // Refresh may not have created all objects for all rows yet...
			}
			obj := objs[i]
			i++

			minSize := obj.MinSize()
			if img, ok := obj.(*richImage); ok {
				if newMin := img.MinSize(); newMin != img.oldMin {
					img.oldMin = newMin

					minSize := r.calculateMin(bounds, wrap, objs, charMinSize, th)
					if r.obj.scr != nil {
						r.obj.prop.SetMinSize(minSize)
					}
					r.Refresh() // TODO resolve this in a similar way to #2991
				}
			}
			rowHeight = fyne.Max(rowHeight, minSize.Height)
			rowWidth += minSize.Width
		}

		if wrap == fyne.TextWrapOff && trunc == fyne.TextTruncateOff {
			width = fyne.Max(width, rowWidth)
		}
		height += rowHeight
		rowHeight = 0
		rowWidth = 0

		lastSeg := bound.segments[len(bound.segments)-1]
		if !lastSeg.Inline() && row < len(bounds)-1 && bounds[row+1].segments[0] != lastSeg { // ignore wrapped lines etc
			height += lineSpacing
		}
	}

	if height == 0 {
		height = charMinSize.Height
	}
	return fyne.NewSize(width, height).
		Add(fyne.NewSquareSize(innerPad * 2).Subtract(r.obj.inset).Subtract(r.obj.inset))
}

func (r *textRenderer) Refresh() {
	bounds := r.obj.rowBounds
	scroll := r.obj.Scroll

	objs := r.obj.updateDecorations()
	for _, bound := range bounds {
		for i, seg := range bound.segments {
			_, isText := seg.(*TextSegment)
			hlSeg, isHyperlink := seg.(*HyperlinkSegment)
			if !isText && !isHyperlink {
				obj := r.obj.cachedSegmentVisual(seg, 0)
				seg.Update(obj)
				objs = append(objs, obj)
				continue
			}

			reuse := 0
			if i == 0 {
				reuse = bound.firstSegmentReuse
			}
			obj := r.obj.cachedSegmentVisual(seg, reuse)
			seg.Update(obj)
			var txt string
			runes := []rune(seg.Textual())

			if i == 0 {
				if len(bound.segments) == 1 {
					txt = string(runes[bound.segBegin:bound.segEnd])
				} else {
					txt = string(runes[bound.segBegin:])
				}
			} else if i == len(bound.segments)-1 && len(bound.segments) > 1 {
				txt = string(runes[:bound.segEnd])
			} else {
				txt = string(runes)
			}
			if bound.ellipsis && i == len(bound.segments)-1 {
				txt = txt + "…"
			}

			if concealed(seg) {
				txt = strings.Repeat(passwordChar, len(runes))
			}

			if isText {
				to, _ := codeInlineText(obj)
				to.Text = txt
			} else if isHyperlink {
				hl, _ := obj.(*fyne.Container).Objects[0].(*Hyperlink)
				hl.Text = txt
				r.associateSiblings(hl, hlSeg, reuse)
				hl.Refresh()
			}
			objs = append(objs, obj)
		}
	}

	if r.obj.scr != nil {
		if inner := scrollInnerContainer(r.obj.scr); inner != nil {
			if inner.Objects == nil {
				r.obj.scr.Content = &fyne.Container{Layout: layout.NewStackLayout(), Objects: []fyne.CanvasObject{
					r.obj.prop, &fyne.Container{Objects: objs},
				}}
				r.obj.scr.Direction = scroll
				r.SetObjects([]fyne.CanvasObject{r.obj.scr})
			} else {
				inner.Objects = objs
			}
		}
		r.obj.scr.Refresh()
	} else {
		r.SetObjects(objs)
	}

	r.Layout(r.obj.Size())
	canvas.Refresh(r.obj.super())

	r.obj.cleanVisualCache()
}

func (r *textRenderer) associateSiblings(hl *Hyperlink, hlSeg *HyperlinkSegment, reuse int) {
	hl.siblings = hl.siblings[:0]
	for prev := 0; prev < reuse; prev++ {
		prevHL, _ := r.obj.cachedSegmentVisual(hlSeg, prev).(*fyne.Container).Objects[0].(*Hyperlink)
		prevHL.siblings = append(prevHL.siblings, hl)
		hl.siblings = append(hl.siblings, prevHL)
	}
}

func (r *textRenderer) layoutRow(texts []fyne.CanvasObject, align fyne.TextAlign, xPos, yPos, lineWidth float32) (x, height float32) {
	initialX := xPos
	if len(texts) == 1 {
		minSize := texts[0].MinSize()
		if text, ok := codeInlineText(texts[0]); ok {
			texts[0].Resize(minSize)
			xPad := float32(0)
			switch text.Alignment {
			case fyne.TextAlignLeading:
			case fyne.TextAlignTrailing:
				xPad = lineWidth - minSize.Width
			case fyne.TextAlignCenter:
				xPad = (lineWidth - minSize.Width) / 2
			}
			texts[0].Move(fyne.NewPos(xPos+xPad, yPos))
		} else {
			texts[0].Resize(fyne.NewSize(lineWidth, minSize.Height))
			texts[0].Move(fyne.NewPos(xPos, yPos))
		}
		return minSize.Width, minSize.Height
	}
	height = float32(0)
	tallestBaseline := float32(0)
	realign := false
	baselines := make([]float32, len(texts))

	driver := fyne.CurrentApp().Driver()
	for i, text := range texts {
		var size fyne.Size
		if txt, ok := codeInlineText(text); ok { // bare text or an inline-code container
			s, base := driver.RenderedTextSize(txt.Text, txt.TextSize, txt.TextStyle, txt.FontSource)
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
				sizeName := link.SizeName
				if sizeName == "" {
					sizeName = theme.SizeNameText
				}
				textSize := theme.SizeForWidget(sizeName, r.obj)
				s, base := driver.RenderedTextSize(link.Text, textSize, link.TextStyle, nil)
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

	innerPadding := r.obj.Theme().Size(theme.SizeNameInnerPadding)
	spare := lineWidth - xPos
	switch align {
	case fyne.TextAlignTrailing:
		spare += innerPadding
		first := texts[0]
		first.Resize(fyne.NewSize(first.Size().Width+spare, height))
		setAlign(first, fyne.TextAlignTrailing)

		for _, text := range texts[1:] {
			text.Move(text.Position().Add(fyne.NewPos(spare, 0)))
		}
	case fyne.TextAlignCenter:
		spare += innerPadding
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

// scrollInnerContainer returns the container holding the visual objects of a RichText
// scroll content, or nil if the scroll structure is not the expected one.
func scrollInnerContainer(o *widget.Scroll) *fyne.Container {
	if c, ok := o.Content.(*fyne.Container); ok {
		if len(c.Objects) == 2 {
			if inner, ok := c.Objects[1].(*fyne.Container); ok {
				return inner
			}
		}
	}
	return nil
}

// howManyRunesFit accepts a rune slice, an available width, an average
// character width, and a function that calculates the (pixel) size of a given
// rune slice.
// howManyRunesFit returns how many runes fit into the available width.
func howManyRunesFit(runes []rune, availableWidth float32, charWidth float32, measurer func([]rune) fyne.Size) int {
	if availableWidth <= 0 {
		return 0
	}
	length := len(runes)
	fits := 0
	tooLong := length + 1
	estimation := int(availableWidth / charWidth)
	if estimation > length {
		estimation = length
	}
	for tooLong-fits > 1 {
		subWidth := measurer(runes[:estimation]).Width
		if subWidth <= availableWidth {
			fits = estimation
		} else {
			tooLong = estimation
		}
		estimation = int(float32(estimation) * availableWidth / subWidth)
		if estimation >= tooLong {
			estimation = tooLong - 1
		}
		if estimation <= fits {
			estimation = fits + 1
		}
	}
	return fits
}

// concealed returns true if the segment represents a password, meaning the text should be obscured.
func concealed(seg RichTextSegment) bool {
	if text, ok := seg.(*TextSegment); ok {
		return text.Style.concealed
	}

	return false
}

func ellipsisPriorBound(bounds []rowBoundary, trunc fyne.TextTruncation, width float32, charWidth float32, measurer func([]rune) fyne.Size) []rowBoundary {
	if trunc != fyne.TextTruncateEllipsis || len(bounds) == 0 {
		return bounds
	}

	prior := bounds[len(bounds)-1]
	seg, ok := prior.segments[0].(*TextSegment)
	if !ok {
		fyne.LogError(fmt.Sprintf("unexpected rich text segment: %#v", prior.segments[0]), nil)
		return bounds
	}

	ellipsisSize := fyne.MeasureText("…", seg.size(), seg.Style.TextStyle) //revive:disable-line:add-constant

	fitCount := howManyRunesFit([]rune(seg.Text)[prior.segBegin:prior.segEnd], width-ellipsisSize.Width, charWidth, measurer)
	prior.segEnd = prior.segBegin + fitCount

	prior.ellipsis = true
	bounds[len(bounds)-1] = prior
	return bounds
}

// findSpaceIndex accepts a slice of runes and a start position index
// findSpaceIndex returns the index of the last space in the text, or -1 if there are no spaces
func findSpaceIndex(text []rune, curIndex int) int {
	for ; curIndex >= 0; curIndex-- {
		if unicode.IsSpace(text[curIndex]) {
			break
		}
	}
	return curIndex
}

func float32ToFixed266(f float32) fixed.Int26_6 {
	return fixed.Int26_6(float64(f) * (1 << 6))
}

// lineBounds accepts a slice of Segments, a wrapping mode, a maximum size available to display and a function to
// measure text size.
// It will return a slice containing the boundary metadata of each line with the given wrapping applied and the
// total height required to render the boundaries at the given width/height constraints
func lineBounds(t *RichText, seg RichTextSegment, firstWidth float32, maxSize fyne.Size, measurer func([]rune) fyne.Size) ([]rowBoundary, float32) {
	wrap := t.Wrapping
	trunc := t.Truncation
	lines := splitLines(seg)

	if wrap == fyne.TextWrap(fyne.TextTruncateClip) {
		if trunc == fyne.TextTruncateOff {
			trunc = fyne.TextTruncateClip
		}
		wrap = fyne.TextWrapOff
	}

	if maxSize.Width <= 0 || wrap == fyne.TextWrapOff && trunc == fyne.TextTruncateOff {
		return lines, 0 // don't bother returning a calculated height, our MinSize is going to cover it
	}

	measureWidth := float32(math.Min(float64(firstWidth), float64(maxSize.Width)))

	switch wrap {
	case fyne.TextWrapBreak:
		return wrapBreakLines(seg, trunc, measureWidth, maxSize, measurer, lines)
	case fyne.TextWrapWord:
		return wrapWordLines(seg, trunc, measureWidth, maxSize, measurer, lines)
	default:
		return truncateLines(t, seg, trunc, measureWidth, measurer, lines)
	}
}

func wrapBreakLines(seg RichTextSegment, trunc fyne.TextTruncation, measureWidth float32, maxSize fyne.Size, measurer func([]rune) fyne.Size, lines []rowBoundary) ([]rowBoundary, float32) {
	text := []rune(seg.Textual())
	charSize := measurer([]rune("z"))
	charWidth := charSize.Width
	lineHeight := charSize.Height
	reuse := 0
	yPos := float32(0)
	var bounds []rowBoundary
	for _, l := range lines {
		low := l.segBegin
		high := l.segEnd
		if low == high {
			l.firstSegmentReuse = reuse
			reuse++
			bounds = append(bounds, l)
			continue
		}
		for low < high {
			if yPos+lineHeight > maxSize.Height && trunc != fyne.TextTruncateOff {
				return ellipsisPriorBound(bounds, trunc, measureWidth, charWidth, measurer), yPos
			}

			fitCount := howManyRunesFit(text[low:high], measureWidth, charWidth, measurer)
			switch fitCount {
			case high - low: // all characters fit on this line
				bounds = append(bounds, rowBoundary{segments: []RichTextSegment{seg}, firstSegmentReuse: reuse, segBegin: low, segEnd: high, ellipsis: false})
				reuse++
				low = high
				high = l.segEnd
				measureWidth = maxSize.Width
				yPos += lineHeight
			case 0: // even a character won't fit
				bounds = append(bounds, rowBoundary{segments: []RichTextSegment{seg}, firstSegmentReuse: reuse, segBegin: low, segEnd: low + 1, ellipsis: false})
				reuse++
				low++
				yPos += lineHeight
			default:
				high = low + fitCount
			}
		}
	}
	return bounds, yPos
}

func wrapWordLines(seg RichTextSegment, trunc fyne.TextTruncation, measureWidth float32, maxSize fyne.Size, measurer func([]rune) fyne.Size, lines []rowBoundary) ([]rowBoundary, float32) {
	text := []rune(seg.Textual())
	charSize := measurer([]rune("z"))
	charWidth := charSize.Width
	lineHeight := charSize.Height
	reuse := 0
	yPos := float32(0)
	var bounds []rowBoundary
	for _, l := range lines {
		low := l.segBegin
		high := l.segEnd
		if low == high {
			l.firstSegmentReuse = reuse
			reuse++
			bounds = append(bounds, l)
			continue
		}
		for low < high {
			if yPos+lineHeight > maxSize.Height && trunc != fyne.TextTruncateOff {
				return ellipsisPriorBound(bounds, trunc, measureWidth, charWidth, measurer), yPos
			}

			sub := text[low:high]
			fitCount := howManyRunesFit(sub, measureWidth, charWidth, measurer)
			if fitCount == high-low { // all characters fit on this line
				bounds = append(bounds, rowBoundary{segments: []RichTextSegment{seg}, firstSegmentReuse: reuse, segBegin: low, segEnd: high, ellipsis: false})
				reuse++
				low = high
				high = l.segEnd
				if low < high && unicode.IsSpace(text[low]) {
					low++
				}
				measureWidth = maxSize.Width

				yPos += lineHeight
				continue
			}
			if fitCount == 0 { // even a character won't fit
				if measureWidth < maxSize.Width {
					bounds = append(bounds, rowBoundary{segments: []RichTextSegment{seg}, firstSegmentReuse: reuse, segBegin: low, segEnd: low, ellipsis: false})
					reuse++
					measureWidth = maxSize.Width
					yPos += lineHeight
					continue
				}
				include := 1
				ellipsis := false
				if trunc == fyne.TextTruncateEllipsis {
					include = 0
					ellipsis = true
				}
				bounds = append(bounds, rowBoundary{segments: []RichTextSegment{seg}, firstSegmentReuse: reuse, segBegin: low, segEnd: low + include, ellipsis: ellipsis})
				low++
				high = low + 1
				reuse++

				yPos += lineHeight
				if high > l.segEnd {
					return bounds, yPos
				}
				continue
			}
			spaceIndex := findSpaceIndex(sub, fitCount)
			if spaceIndex >= 0 {
				if spaceIndex == 0 {
					spaceIndex = 1
				}
				high = low + spaceIndex
				continue
			}
			oldHigh := high
			high = low + fitCount
			if low == 0 && measureWidth < maxSize.Width { // add a newline as there is more space on next
				bounds = append(bounds, rowBoundary{segments: []RichTextSegment{seg}, firstSegmentReuse: reuse, segBegin: low, segEnd: low, ellipsis: false})
				reuse++
				high = oldHigh
				measureWidth = maxSize.Width

				yPos += lineHeight
			}
		}
	}
	return bounds, yPos
}

func truncateLines(t *RichText, seg RichTextSegment, trunc fyne.TextTruncation, measureWidth float32, measurer func([]rune) fyne.Size, lines []rowBoundary) ([]rowBoundary, float32) {
	text := []rune(seg.Textual())
	yPos := float32(0)
	var bounds []rowBoundary
	charSize := measurer([]rune("z")) //revive:disable-line:add-constant -- TODO: clarify whether we want to define a common letter constant for approximate character sizes
	charWidth := charSize.Width
	reuse := 0
	for _, l := range lines {
		low := l.segBegin
		high := l.segEnd
		if low == high {
			l.firstSegmentReuse = reuse
			reuse++
			bounds = append(bounds, l)
			continue
		}
		switch trunc {
		case fyne.TextTruncateEllipsis:
			txt := []rune(seg.Textual())[low:high]
			var textObj *canvas.Text
			switch s := seg.(type) {
			case *TextSegment:
				textObj, _ = codeInlineText(seg.Visual())
			case *HyperlinkSegment:
				textObj = canvas.NewText(string(txt), color.Black)
				textObj.TextStyle = s.TextStyle
				sizeName := s.SizeName
				if sizeName == "" {
					sizeName = theme.SizeNameText
				}
				textObj.TextSize = theme.SizeForWidget(sizeName, t)
			}
			end, full := truncateLimit(string(txt), textObj, int(measureWidth), []rune{'…'})
			high = low + end
			bounds = append(bounds, rowBoundary{segments: []RichTextSegment{seg}, firstSegmentReuse: reuse, segBegin: low, segEnd: high, ellipsis: !full})
			reuse++
		case fyne.TextTruncateClip:
			fitCount := howManyRunesFit(text[low:high], measureWidth, charWidth, measurer)
			high = low + fitCount
			bounds = append(bounds, rowBoundary{segments: []RichTextSegment{seg}, firstSegmentReuse: reuse, segBegin: low, segEnd: high, ellipsis: false})
			reuse++
		case fyne.TextTruncateOff:
			// don’t do anything
		}
	}
	return bounds, yPos
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

// rowPaddingAndAlign returns the left padding and text alignment for a row.
func (t *RichText) rowPaddingAndAlign(bound rowBoundary, lineSpacing float32, currentAlign fyne.TextAlign) (float32, fyne.TextAlign) {
	leftPad := bound.indent
	align := currentAlign
	quoting := 0

	switch first := rowFirstVisibleSegment(&bound).(type) {
	case *TextSegment:
		align = first.Style.Alignment
		quoting = first.Style.QuotingDepth
	case *HyperlinkSegment:
		align = first.Alignment
		quoting = first.quotingLevel
	case *CodeBlockSegment:
		align = fyne.TextAlignLeading
		quoting = first.quotingLevel
	case *listMarkerSegment:
		align = fyne.TextAlignLeading
		quoting = first.quoting
	}

	if quoting > 0 {
		leftPad = lineSpacing * 4 * float32(quoting)
	}
	if bound.panel != nil { // the rows sit inside the panel, not against its edge
		leftPad += theme.SizeForWidget(theme.SizeNameInnerPadding, t)
	}
	return leftPad, align
}

// splitLines accepts a text segment and returns a slice of boundary metadata denoting the
// start and end indices of each line delimited by the newline character.
func splitLines(seg RichTextSegment) []rowBoundary {
	var low, high int
	var lines []rowBoundary
	text := []rune(seg.Textual())
	length := len(text)
	for i := 0; i < length; i++ {
		if text[i] == '\n' {
			high = i
			lines = append(lines, rowBoundary{segments: []RichTextSegment{seg}, firstSegmentReuse: len(lines), segBegin: low, segEnd: high, ellipsis: false})
			low = i + 1
		}
	}
	return append(lines, rowBoundary{segments: []RichTextSegment{seg}, firstSegmentReuse: len(lines), segBegin: low, segEnd: length, ellipsis: false})
}

func truncateLimit(s string, text *canvas.Text, limit int, ellipsis []rune) (int, bool) {
	face := paint.CachedFontFace(text.TextStyle, text.FontSource, text)

	runes := []rune(s)
	in := shaping.Input{
		Text:      ellipsis,
		RunStart:  0,
		RunEnd:    len(ellipsis),
		Direction: di.DirectionLTR,
		Face:      face.Fonts.ResolveFace(ellipsis[0]),
		Size:      float32ToFixed266(text.TextSize),
	}
	shaper := &shaping.HarfbuzzShaper{}
	segmenter := &shaping.Segmenter{}

	conf := shaping.WrapConfig{}
	conf = conf.WithTruncator(shaper, in)
	conf.BreakPolicy = shaping.WhenNecessary
	conf.TruncateAfterLines = 1
	l := shaping.LineWrapper{}

	in.Text = runes
	in.RunEnd = len(runes)
	ins := segmenter.Split(in, face.Fonts)
	outs := make([]shaping.Output, len(ins))
	for i, in := range ins {
		outs[i] = shaper.Shape(in)
	}

	l.Prepare(conf, runes, shaping.NewSliceIterator(outs))
	wrapped, done := l.WrapNextLine(limit)

	count := len(runes)
	if wrapped.Truncated != 0 {
		count -= wrapped.Truncated
		count += len(ellipsis)
	}

	full := done && count == len(runes)
	if !full && len(ellipsis) > 0 {
		count--
	}
	return count, full
}

type rowBoundary struct {
	segments          []RichTextSegment
	firstSegmentReuse int

	// segBegin indexes into the first segment of this row and segEnd into the
	// last, as a row may start or finish part way through a segment that it
	// shares with its neighbours.
	segBegin, segEnd int

	// docBegin and docEnd are the rune offsets of this row within the whole
	// text, which is what a cursor or selection position is measured in.
	docBegin, docEnd int

	ellipsis bool
	indent   float32

	// panel is set when this row is part of a block that draws its content on a
	// panel, such as a code block.
	panel panelSegment

	// yPos and height record where this row was placed by the renderer, so that
	// widgets can position a cursor or selection against rows of differing size.
	yPos, height float32
}

package widget

import (
	"strings"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// rowBoundsBuilder walks a tree of segments to work out the rows that they lay out on.
type rowBoundsBuilder struct {
	text         *RichText
	theme        fyne.Theme
	innerPadding float32

	bounds  []rowBoundary
	fitSize fyne.Size // the space that is left for content to fit within

	maxWidth  float32 // the width available to a row that starts at the left edge
	wrapWidth float32 // the width that is left on the row being built

	rowOpen         bool    // whether the last of bounds can take more inline content
	rowDepth        int     // the nesting depth of the segment that started the open row
	rowIndent       float32 // how far rows that continue this one are indented, -1 until known
	rowMarkerIndent float32 // the width of a list bullet that introduces this row

	docOffset int // rune offset of the current segment within the whole text
}

func newRowBoundsBuilder(t *RichText) *rowBoundsBuilder {
	th := t.Theme()
	innerPadding := th.Size(theme.SizeNameInnerPadding)

	fitSize := t.Size()
	if t.scr != nil {
		fitSize = t.scr.Content.MinSize()
	}
	fitSize.Height -= (innerPadding + t.inset.Height) * 2
	maxWidth := t.Size().Width - 2*innerPadding + 2*t.inset.Width

	return &rowBoundsBuilder{
		text: t, theme: th, innerPadding: innerPadding, fitSize: fitSize,
		maxWidth: maxWidth, wrapWidth: maxWidth, rowIndent: -1,
	}
}

// walk lays out each of the given segments in turn, recursing into any blocks.
func (b *rowBoundsBuilder) walk(segList []RichTextSegment, depth int) {
	for _, seg := range segList {
		switch s := seg.(type) {
		case RichTextBlock:
			b.appendBlock(seg, s, depth)
		case *TextSegment, *HyperlinkSegment:
			b.appendText(seg, depth)
		default:
			b.appendObject(seg, depth)
		}
	}
}

// startRow records that the row being built began at the given depth, with no
// indent or bullet carried over from whatever came before it.
func (b *rowBoundsBuilder) startRow(depth int) {
	b.rowDepth = depth
	b.rowIndent = -1
	b.rowMarkerIndent = 0
}

// closeRow ends the row being built so that the next segment starts a new one.
func (b *rowBoundsBuilder) closeRow(depth int) {
	b.startRow(depth)
	b.rowOpen = false
	b.wrapWidth = b.maxWidth
}

// appendBlock lays out the content of a block, marking the rows it added as
// belonging to it if the block draws a panel behind them.
func (b *rowBoundsBuilder) appendBlock(seg RichTextSegment, block RichTextBlock, depth int) {
	segs := block.Segments()
	first := len(b.bounds)
	b.walk(segs, depth+1)

	if panel, ok := seg.(panelSegment); ok {
		markPanelRows(b.bounds, first, segs, panel)
	}
	if segmentsEndRow(segs) {
		b.closeRow(depth)
	}
	if len(segs) == 0 { // otherwise the content was counted as it was walked
		b.docOffset += utf8.RuneCountInString(seg.Textual())
	}
}

// appendObject lays out a segment that draws an object instead of text, such as
// an image or a list bullet.
func (b *rowBoundsBuilder) appendObject(seg RichTextSegment, depth int) {
	segLen := utf8.RuneCountInString(seg.Textual())
	if b.rowOpen {
		row := &b.bounds[len(b.bounds)-1]
		row.segments = append(row.segments, seg)
		row.docEnd = b.docOffset + segLen
	} else {
		b.bounds = append(b.bounds, rowBoundary{
			segments: []RichTextSegment{seg},
			docBegin: b.docOffset, docEnd: b.docOffset + segLen,
		})
		b.rowOpen = true
		b.rowDepth = depth
	}

	itemMin := b.text.cachedSegmentVisual(seg, 0).MinSize()
	if seg.Inline() {
		b.wrapWidth -= itemMin.Width
		if _, isMarker := seg.(*listMarkerSegment); isMarker {
			// so that the item text wraps in line with itself, not the bullet
			b.rowMarkerIndent = itemMin.Width
		}
	} else {
		b.closeRow(depth)
		b.fitSize.Height -= itemMin.Height + b.theme.Size(theme.SizeNameLineSpacing)
	}
	b.docOffset += segLen
}

// appendText lays out a text or hyperlink segment, wrapping it over as many
// rows as it needs.
func (b *rowBoundsBuilder) appendText(seg RichTextSegment, depth int) {
	style, size, leftPad := b.textAttributes(seg)
	rows, height := lineBounds(b.text, seg, b.wrapWidth-leftPad, fyne.NewSize(b.maxWidth, b.fitSize.Height),
		func(text []rune) fyne.Size {
			return fyne.MeasureText(string(text), size, style)
		})
	for i := range rows {
		rows[i].docBegin = b.docOffset + rows[i].segBegin
		rows[i].docEnd = b.docOffset + rows[i].segEnd
	}

	if b.rowOpen {
		b.continueRow(seg, rows, depth, height)
	} else {
		b.bounds = append(b.bounds, rows...)
		b.fitSize.Height -= height
		b.rowOpen = true
		b.startRow(depth)
	}

	if seg.Inline() {
		b.advanceRow(seg, rows, style, size)
	} else {
		b.closeRow(depth)
	}
	b.docOffset += utf8.RuneCountInString(seg.Textual())
}

// continueRow adds the first of the new rows to the row that is already open,
// as the segment runs on from the content that is there, then appends the rest.
func (b *rowBoundsBuilder) continueRow(seg RichTextSegment, rows []rowBoundary, depth int, height float32) {
	if len(rows) == 0 {
		return
	}

	// this row now runs on into another segment, so segEnd moves to
	// index the new last segment rather than the previous one
	row := &b.bounds[len(b.bounds)-1]
	row.segEnd = rows[0].segEnd
	row.docEnd = b.docOffset + rows[0].segEnd
	row.segments = append(row.segments, seg)

	if depth > b.rowDepth || b.rowMarkerIndent > 0 {
		b.indentContinuation(seg, rows[1:])
	}
	b.bounds = append(b.bounds, rows[1:]...)
	b.fitSize.Height -= height
}

// indentContinuation lines the rows that a wrapped segment ran on to up with the
// text of its first row, rather than with the bullet or quote that introduced it.
func (b *rowBoundsBuilder) indentContinuation(seg RichTextSegment, rows []rowBoundary) {
	if b.rowMarkerIndent > 0 {
		b.rowIndent = b.rowMarkerIndent
	} else if b.rowIndent == -1 {
		b.rowIndent = b.maxWidth - b.wrapWidth
	}
	if b.rowIndent <= 0 {
		return
	}

	runes := []rune(seg.Textual())
	for i := range rows {
		row := &rows[i]
		if row.segBegin > 0 && row.segBegin <= len(runes) && runes[row.segBegin-1] == '\n' {
			continue // a new line starts back at the left edge
		}
		row.indent = b.rowIndent
	}
}

// advanceRow takes the width of the text that an inline segment just placed on
// the open row off the space that is left for whatever follows it.
func (b *rowBoundsBuilder) advanceRow(seg RichTextSegment, rows []rowBoundary, style fyne.TextStyle, size float32) {
	last := b.bounds[len(b.bounds)-1]
	begin := 0
	if len(last.segments) == 1 {
		begin = last.segBegin
	}

	// check ranges - as we resize it can be wrong?
	runes := []rune(seg.Textual())
	begin = min(begin, len(runes))
	end := min(last.segEnd, len(runes))

	lastWidth := fyne.MeasureText(string(runes[begin:end]), size, style).Width
	if len(rows) == 1 {
		b.wrapWidth -= lastWidth
	} else {
		b.wrapWidth = b.maxWidth - lastWidth
	}
	if strings.ContainsRune(seg.Textual(), '\n') {
		b.rowMarkerIndent = 0 // we are past the line that a bullet introduced
	}
}

// textAttributes returns the style and size that a text or hyperlink segment is
// measured with, along with the padding that quoting it adds to its left.
func (b *rowBoundsBuilder) textAttributes(seg RichTextSegment) (style fyne.TextStyle, size, leftPad float32) {
	switch s := seg.(type) {
	case *TextSegment:
		return s.Style.TextStyle, s.size(), b.innerPadding * 2 * float32(s.Style.QuotingDepth)
	case *HyperlinkSegment:
		return s.TextStyle, theme.SizeForWidget(theme.SizeNameText, b.text), b.innerPadding * 2 * float32(s.quotingLevel)
	}
	return style, size, leftPad
}

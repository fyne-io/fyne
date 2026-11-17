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
	rowCut          bool    // whether truncation hid the end of the open row, so more inline content is hidden too
	rowPaid         bool    // whether the height of the open row was taken from fitSize already
	rowDepth        int     // the nesting depth of the segment that started the open row
	rowIndent       float32 // how far rows that continue this one are indented, -1 until known
	rowMarkerIndent float32 // the width of a list bullet that introduces this row

	docOffset int // rune offset of the current segment within the whole text

	truncating bool // whether rows are truncated to the width, rather than wrapped
	full       bool // whether truncation ran out of height, so no more content is shown
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
	truncating := t.Wrapping == fyne.TextWrap(fyne.TextTruncateClip) ||
		t.Wrapping == fyne.TextWrapOff && t.Truncation != fyne.TextTruncateOff

	return &rowBoundsBuilder{
		text: t, theme: th, innerPadding: innerPadding, fitSize: fitSize,
		maxWidth: maxWidth, wrapWidth: maxWidth, rowIndent: -1, truncating: truncating,
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
	b.rowCut = false
	b.rowPaid = false
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

	if len(segs) == 0 { // otherwise the content was counted as it was walked
		b.docOffset += utf8.RuneCountInString(seg.Textual())
	}

	if panel, ok := seg.(panelSegment); ok {
		rows := b.bounds
		if b.endsAfterBreak(first, segs) {
			rows = rows[:len(rows)-1]
			markPanelRows(rows, first, segs, panel)

			b.startRow(depth)
			b.rowOpen = true
			b.wrapWidth = b.maxWidth
			return
		}
		markPanelRows(rows, first, segs, panel)
	}
	if segmentsEndRow(segs) {
		b.closeRow(depth)
	}
}

// endsAfterBreak reports whether the content of a block finished with a line
// break, leaving an empty row after it as the last of those added since first.
func (b *rowBoundsBuilder) endsAfterBreak(first int, segs []RichTextSegment) bool {
	if len(b.bounds) <= first || !endsWithNewline(segs) {
		return false
	}

	last := &b.bounds[len(b.bounds)-1]
	return last.docBegin == b.docOffset && last.docEnd == b.docOffset
}

// appendObject lays out a segment that draws an object instead of text, such as
// an image or a list bullet.
func (b *rowBoundsBuilder) appendObject(seg RichTextSegment, depth int) {
	segLen := utf8.RuneCountInString(seg.Textual())
	if b.full || b.rowOpen && b.rowCut && seg.Inline() { // truncated away
		if !b.full {
			b.bounds[len(b.bounds)-1].hiddenEnd = seg
		}
		b.docOffset += segLen
		return
	}

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
	runes := []rune(seg.Textual())
	if b.full {
		b.docOffset += len(runes)
		return
	}

	style, size, leftPad := b.textAttributes(seg)
	measurer := func(text []rune) fyne.Size {
		return fyne.MeasureText(string(text), size, style)
	}

	// the first line of this text runs on along the open row, which may have paid for its height already
	shared := float32(0)
	if b.rowOpen && b.rowPaid && len(runes) > 0 && runes[0] != '\n' {
		shared = measurer([]rune(averageChar)).Height
	}
	rows, height := lineBounds(b.text, seg, b.wrapWidth-leftPad, fyne.NewSize(b.maxWidth, b.fitSize.Height+shared), measurer)
	for i := range rows {
		rows[i].docBegin = b.docOffset + rows[i].segBegin
		rows[i].docEnd = b.docOffset + rows[i].segEnd
	}
	b.fitSize.Height -= max(height-shared, 0)
	b.docOffset += len(runes)

	if len(rows) == 0 || !b.truncating && rows[len(rows)-1].truncated {
		b.full = true // there was not enough height left to show all of this text
	}
	if len(rows) == 0 {
		if b.rowOpen && b.text.Truncation == fyne.TextTruncateEllipsis {
			b.endRowWithEllipsis()
		}
		return
	}

	if b.rowOpen {
		b.continueRow(seg, rows, depth)
	} else {
		b.bounds = append(b.bounds, rows...)
		b.rowOpen = true
		b.startRow(depth)
	}
	b.updateRowState(runes, rows, height)

	if !seg.Inline() {
		b.closeRow(depth)
	} else if !b.rowCut {
		b.advanceRow(seg, rows, style, size)
	}
}

// continueRow adds the first of the new rows to the row that is already open,
// as the segment runs on from the content that is there, then appends the rest.
func (b *rowBoundsBuilder) continueRow(seg RichTextSegment, rows []rowBoundary, depth int) {
	first := rows[0]
	switch {
	case b.rowCut:
		// the end of the open row is hidden, so this text only shows from its next line
		b.bounds[len(b.bounds)-1].hiddenEnd = seg
	case first.truncated && first.segEnd == first.segBegin:
		// none of this text fits, so the open row ends with what it has already
		if first.ellipsis {
			b.endRowWithEllipsis()
		}
		b.rowCut = true
		b.bounds[len(b.bounds)-1].hiddenEnd = seg
	default:
		// this row now runs on into another segment, so segEnd moves to
		// index the new last segment rather than the previous one
		row := &b.bounds[len(b.bounds)-1]
		row.segEnd = first.segEnd
		row.docEnd = first.docEnd
		row.ellipsis = first.ellipsis
		row.truncated = first.truncated
		row.segments = append(row.segments, seg)
	}

	if depth > b.rowDepth || b.rowMarkerIndent > 0 {
		b.indentContinuation(seg, rows[1:])
	}
	b.bounds = append(b.bounds, rows[1:]...)
}

// updateRowState records whether the open row was truncated and whether its height
// is paid for, now that the given rows of this text were added.
func (b *rowBoundsBuilder) updateRowState(runes []rune, rows []rowBoundary, height float32) {
	last := rows[len(rows)-1]
	if len(rows) > 1 { // the open row is the last of these rows, not the one that they continued
		b.rowCut = false
		b.rowPaid = false
	}

	if b.truncating {
		b.rowCut = b.rowCut || last.truncated
	}
	if height > 0 && (len(runes) == 0 || runes[len(runes)-1] != '\n') {
		b.rowPaid = true
	}
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

// endRowWithEllipsis ends the open row with an ellipsis after the text on it,
// shortening that text if there is not enough room left for the ellipsis.
func (b *rowBoundsBuilder) endRowWithEllipsis() {
	row := &b.bounds[len(b.bounds)-1]
	space := b.wrapWidth
	b.rowCut = true
	for {
		last := len(row.segments) - 1
		if !isTextSegment(row.segments[last]) {
			return // the ellipsis is drawn at the end of text, so there is nowhere to put it
		}

		style, size, _ := b.textAttributes(row.segments[last])
		measurer := func(text []rune) fyne.Size {
			return fyne.MeasureText(string(text), size, style)
		}
		runes := rowSegmentRunes(row, last)
		width := measurer(runes).Width
		fit := howManyRunesFit(runes, width+space-measurer([]rune(ellipsisChar)).Width, measurer([]rune(averageChar)).Width, measurer)
		if fit > 0 || last == 0 || !isTextSegment(row.segments[last-1]) {
			row.segEnd -= len(runes) - fit
			row.docEnd -= len(runes) - fit
			row.ellipsis = true
			return
		}

		// none of this text fits, so the ellipsis goes on the text before it instead
		space += width
		row.docEnd -= len(runes)
		row.segments = row.segments[:last]
		row.segEnd = utf8.RuneCountInString(row.segments[last-1].Textual()) // which was shown in full
	}
}

func isTextSegment(seg RichTextSegment) bool {
	switch seg.(type) {
	case *TextSegment, *HyperlinkSegment:
		return true
	}
	return false
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

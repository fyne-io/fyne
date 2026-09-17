package widget

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Declare conformity with interfaces
var (
	_ fyne.Widget      = (*RichTextEntry)(nil)
	_ fyne.Focusable   = (*RichTextEntry)(nil)
	_ fyne.Disableable = (*RichTextEntry)(nil)
)

const (
	markdownBold      = "**"
	markdownCodeFence = "```"
	markdownQuote     = "> "
	markdownStrike    = "~~"
)

// RichTextEntry widget allows styled text to be edited when focused.
// It behaves like an [Entry] - but the content is held as a list of [RichTextSegment],
// so different runs of text can use their own style and objects can appear inline.
//
// Bulleted lists and fenced code blocks keep the bullets and panel that they are
// drawn with, whilst the text inside them is edited like the text around it.
//
// The content can be set from markdown using [RichTextEntry.ParseMarkdown] or [NewRichTextEntryFromMarkdown].
// Setting [RichTextEntry.TypeMarkdown] additionally converts markdown into the matching style as it is typed.
//
// Since: 2.9
type RichTextEntry struct {
	Entry

	// TypeMarkdown converts markdown into the style that it describes as soon as it has been typed.
	TypeMarkdown bool
}

// NewRichTextEntry creates a new rich text entry widget.
// The returned entry accepts multiple lines and wraps on word boundaries.
//
// Since: 2.9
func NewRichTextEntry() *RichTextEntry {
	e := &RichTextEntry{}
	e.MultiLine = true
	e.Wrapping = fyne.TextWrapWord
	e.ExtendBaseWidget(e)
	return e
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
func (e *RichTextEntry) ExtendBaseWidget(wid fyne.Widget) {
	e.richProvider()
	e.Entry.ExtendBaseWidget(wid)
}

// NewRichTextEntryFromMarkdown creates a new rich text entry widget with the
// content described by the specified markdown.
//
// Since: 2.9
func NewRichTextEntryFromMarkdown(content string) *RichTextEntry {
	e := NewRichTextEntry()
	e.ParseMarkdown(content)
	return e
}

// Segments returns the styled segments that make up the content of this entry.
// The result should not be modified directly, use [RichTextEntry.SetSegments].
//
// Since: 2.9
func (e *RichTextEntry) Segments() []RichTextSegment {
	return e.richProvider().Segments
}

// SetSegments replaces the content of this entry with the specified segments.
// Bulleted lists and fenced code blocks are kept, so that they still draw their
// bullets and panel, with the text inside them editable like any other.
// Segments that cannot be edited in place are converted to styled text.
//
// Since: 2.9
func (e *RichTextEntry) SetSegments(segments []RichTextSegment) {
	provider := e.richProvider()
	provider.Segments = editableSegments(segments)

	e.CursorRow, e.CursorColumn = 0, 0
	e.ClearSelection()
	e.undoStack.Clear()
	e.updateTextAndRefresh(provider.String(), false)
	e.updateCursorAndSelection()
}

// SetText replaces the content of this entry with unstyled text.
//
// Since: 2.9
func (e *RichTextEntry) SetText(text string) {
	e.SetSegments([]RichTextSegment{&TextSegment{Style: RichTextStyleInline, Text: text}})
}

// ParseMarkdown replaces the content of this entry with the result of parsing
// the specified markdown.
//
// Since: 2.9
func (e *RichTextEntry) ParseMarkdown(content string) {
	e.SetSegments(parseMarkdown(content))
}

// AppendMarkdown parses the specified markdown and adds the result to the end of
// the current content.
//
// Since: 2.9
func (e *RichTextEntry) AppendMarkdown(content string) {
	provider := e.richProvider()
	segments := provider.Segments
	if current := provider.String(); current != "" && !strings.HasSuffix(current, newLineChar) {
		segments = append(segments, &TextSegment{Style: RichTextStyleInline, Text: newLineChar})
	}
	segments = append(segments, editableSegments(parseMarkdown(content))...)

	pos := e.CursorTextOffset()
	provider.Segments = mergeSegments(segments)
	e.updateTextAndRefresh(provider.String(), false)
	e.setCursorOffset(pos)
}

// Markdown returns the content of this entry serialised back to markdown.
//
// Since: 2.9
func (e *RichTextEntry) Markdown() string {
	return segmentsToMarkdown(e.richProvider().Segments)
}

// SetStyleForRange applies the specified style to the text between the two rune
// offsets, splitting segments where required. Positions outside the content are
// clamped and a start at or after the end does nothing.
//
// Since: 2.9
func (e *RichTextEntry) SetStyleForRange(start, end int, style RichTextStyle) {
	style.Inline = true
	e.styleRange(start, end, func(s *RichTextStyle) {
		*s = style
	})

	pos := e.CursorTextOffset()
	provider := e.richProvider()
	provider.Segments = mergeSegments(provider.Segments)
	e.Refresh()
	e.setCursorOffset(pos)
}

// SetStyleForSelection applies the specified style to the currently selected
// text. It does nothing if there is no selection.
//
// Since: 2.9
func (e *RichTextEntry) SetStyleForSelection(style RichTextStyle) {
	e.syncSelectable()
	start, end := e.sel.selection()
	if start < 0 || start == end {
		return
	}

	e.SetStyleForRange(start, end, style)
}

// StyleAtCursor returns the style that newly typed text will use.
// Any current selection is ignored.
//
// Since: 2.9
func (e *RichTextEntry) StyleAtCursor() RichTextStyle {
	pos := e.CursorTextOffset()
	style := RichTextStyleInline

	off := 0
	for _, seg := range e.richProvider().contentSegments() {
		length := utf8.RuneCountInString(seg.Textual())
		text, isText := seg.(*TextSegment)
		if isText && ((off < pos && pos <= off+length) || (off == pos && length == 0)) {
			style = text.Style
		}
		if off > pos {
			break
		}
		off += length
	}
	return style
}

// TypedRune receives text input events when this widget is focused.
func (e *RichTextEntry) TypedRune(r rune) {
	e.Entry.TypedRune(r)

	if e.TypeMarkdown && e.styleMarkdownAtCursor(r) {
		e.Refresh()
	}
}

// TypedKey receives key input events when this widget is focused.
func (e *RichTextEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyBackspace:
		if e.hasSelection() {
			break
		}

		// a bullet is removed before the text that it introduces
		pos := e.CursorTextOffset()
		if e.removeBullet(pos) || e.joinIntoBlockAbove(pos) || e.clearEmptyStyle() {
			e.Refresh()
			return
		}
	case fyne.KeyDelete:
		if e.hasSelection() {
			break
		}

		// deleting the break at the end of an item joins the one below it
		pos := e.CursorTextOffset()
		if e.removeBullet(pos + 1) {
			e.setCursorOffset(pos)
			e.Refresh()
			return
		}
		if e.clearEmptyStyle() {
			e.Refresh()
			return
		}
	case fyne.KeyReturn, fyne.KeyEnter:
		if e.MultiLine && !e.hasSelection() && e.startsBlockLine() {
			e.Refresh()
			return
		}
	}

	e.Entry.TypedKey(key)

	switch key.Name {
	case fyne.KeyBackspace, fyne.KeyDelete:
		if e.pruneEmptySegments() {
			e.Refresh()
		}
	case fyne.KeyReturn, fyne.KeyEnter:
		if e.MultiLine && isBlockStyle(e.StyleAtCursor()) {
			// headings and quotes cover a single line, so the next one starts
			// again from the base style
			e.insertEmptySegmentAt(e.CursorTextOffset(), RichTextStyleInline)
			e.Refresh()
		}
	}
}

// startsBlockLine handles a return that does more than break the line, either
// carrying a list on to a new bullet or opening and closing a code block.
// It reports whether the content was changed.
func (e *RichTextEntry) startsBlockLine() bool {
	if e.TypeMarkdown && e.toggleCodeFence() {
		return true
	}

	return e.splitListItem(e.CursorTextOffset())
}

// hasSelection reports whether some of the content is currently selected.
func (e *RichTextEntry) hasSelection() bool {
	e.syncSelectable()

	start, end := e.sel.selection()
	return start >= 0 && start != end
}

// isBlockStyle reports whether a style applies to a whole line, rather than to a
// run of text within one.
func isBlockStyle(style RichTextStyle) bool {
	if style.QuotingDepth > 0 {
		return true
	}
	return style.SizeName != "" && style.SizeName != theme.SizeNameText
}

// Undo un-does the last modifying user-action.
//
// Since: 2.9
func (e *RichTextEntry) Undo() {
	_, action := e.undoStack.Undo(e.Text)
	e.applyUndoAction(action, true)
}

// Redo re-applies the last undone user-action.
//
// Since: 2.9
func (e *RichTextEntry) Redo() {
	_, action := e.undoStack.Redo(e.Text)
	e.applyUndoAction(action, false)
}

// applyUndoAction runs an undo stack entry against the segments to retain style.
func (e *RichTextEntry) applyUndoAction(action entryUndoAction, undo bool) {
	modify, ok := action.(*entryModifyAction)
	if !ok {
		return
	}

	provider := e.richProvider()
	pos := modify.Position
	if modify.Delete == undo { // put the text back
		provider.insertAt(pos, modify.Text)
		pos += len(modify.Text)
	} else {
		provider.deleteFromTo(pos, pos+len(modify.Text))
	}

	e.pruneEmptySegments()
	content := provider.String()
	e.updateText(content, false)
	e.setCursorOffset(pos)

	if e.OnChanged != nil {
		e.OnChanged(content)
	}
	e.validate()
	e.Refresh()
}

// setCursorOffset moves the cursor to the specified rune offset in the content.
func (e *RichTextEntry) setCursorOffset(pos int) {
	e.CursorRow, e.CursorColumn = e.rowColFromTextPos(pos)
	e.syncSelectable()
}

// richProvider returns the rich text that holds the content of this entry,
// setting up the first segment if it has not been prepared yet.
func (e *RichTextEntry) richProvider() *RichText {
	if !e.rich {
		e.rich = true
		e.initTextProvider()
		e.text.Segments = []RichTextSegment{&TextSegment{Style: RichTextStyleInline, Text: e.Text}}
	}
	return &e.text
}

// blockContainer returns the segments held inside a block, so that content can
// be added to or removed from it. It returns nil for anything else, including
// blocks such as code that hold their content as text of their own.
func blockContainer(seg RichTextSegment) *[]RichTextSegment {
	switch block := seg.(type) {
	case *ParagraphSegment:
		return &block.Texts
	case *ListSegment:
		return &block.Items
	}

	return nil
}

// splitAt divides segments as required so that the specified rune offset falls on
// a segment boundary. It returns the list of segments holding that boundary, which
// may be inside a block, and the index of the segment that starts there.
func (e *RichTextEntry) splitAt(pos int) (*[]RichTextSegment, int) {
	provider := e.richProvider()
	if list, i, _ := splitSegmentsAt(&provider.Segments, pos, 0); list != nil {
		return list, i
	}

	return &provider.Segments, len(provider.Segments)
}

// splitSegmentsAt looks for the rune offset in the given segments, dividing a text
// segment where the offset falls inside one. It returns the list that holds the
// boundary with the index of the segment starting there, or a nil list and the
// offset reached when the position is beyond these segments.
func splitSegmentsAt(list *[]RichTextSegment, pos, off int) (out *[]RichTextSegment, index int, offset int) {
	for i := 0; i < len(*list); i++ {
		seg := (*list)[i]
		if inner := blockContainer(seg); inner != nil {
			found, index, next := splitSegmentsAt(inner, pos, off)
			if found != nil {
				return found, index, next
			}

			off = next
			continue
		}

		length := utf8.RuneCountInString(seg.Textual())
		if pos == off {
			if _, isText := seg.(*TextSegment); !isText && length == 0 {
				continue // a bullet introduces the text that follows it
			}
			return list, i, off
		}
		if pos < off+length {
			text, ok := seg.(*TextSegment)
			if !ok {
				return list, i + 1, off + length // objects cannot be divided
			}

			runes := []rune(text.Text)
			cut := pos - off
			tail := &TextSegment{Style: text.Style, Text: string(runes[cut:])}
			text.Text = string(runes[:cut])

			segments := make([]RichTextSegment, 0, len(*list)+1)
			segments = append(segments, (*list)[:i+1]...)
			segments = append(segments, tail)
			segments = append(segments, (*list)[i+1:]...)
			*list = segments
			return list, i + 1, pos
		}
		off += length
	}

	return nil, 0, off
}

// styleRange updates the style of every text segment between the two rune
// offsets, dividing segments at the boundaries where required.
func (e *RichTextEntry) styleRange(start, end int, apply func(*RichTextStyle)) {
	provider := e.richProvider()
	start = max(start, 0)
	end = min(end, provider.len())
	if start >= end {
		return
	}

	e.splitAt(start)
	e.splitAt(end)

	off := 0
	for _, seg := range provider.contentSegments() {
		length := utf8.RuneCountInString(seg.Textual())
		if length > 0 && off >= start && off+length <= end {
			if text, ok := seg.(*TextSegment); ok {
				apply(&text.Style)
			}
		}
		off += length
	}
}

// insertEmptySegmentAt places an empty segment with the given style at the rune offset.
func (e *RichTextEntry) insertEmptySegmentAt(pos int, style RichTextStyle) {
	style.Inline = true
	e.dropEmptySegmentsAt(pos) // only one segment may claim the text typed here
	list, i := e.splitAt(pos)

	segments := make([]RichTextSegment, 0, len(*list)+1)
	segments = append(segments, (*list)[:i]...)
	segments = append(segments, &TextSegment{Style: style})
	segments = append(segments, (*list)[i:]...)
	*list = segments
}

// dropEmptySegmentsAt removes any empty text segments sitting at the given rune
// offset, so that a replacement can take ownership of text typed there.
func (e *RichTextEntry) dropEmptySegmentsAt(pos int) {
	provider := e.richProvider()
	dropEmptySegments(&provider.Segments, 0, false, func(off int) bool { return off == pos })
}

// pruneEmptySegments drops empty text segments, apart from any at the cursor -
// where one may have been placed to hold a style that is about to be typed.
// Returns true if any segment was removed.
func (e *RichTextEntry) pruneEmptySegments() bool {
	provider := e.richProvider()
	cursor := e.CursorTextOffset()
	_, removed := dropEmptySegments(&provider.Segments, 0, false, func(off int) bool { return off != cursor })

	if len(provider.Segments) == 0 {
		provider.Segments = []RichTextSegment{&TextSegment{Style: RichTextStyleInline}}
	}
	return removed
}

// dropEmptySegments removes the empty text segments that drop reports for, looking
// inside the blocks that hold content. It returns the offset reached and whether any
// segment was removed.
func dropEmptySegments(list *[]RichTextSegment, off int, inBlock bool, drop func(off int) bool) (int, bool) {
	removed := false
	segments := make([]RichTextSegment, 0, len(*list))
	for _, seg := range *list {
		if inner := blockContainer(seg); inner != nil {
			next, innerRemoved := dropEmptySegments(inner, off, true, drop)
			off, removed = next, removed || innerRemoved

			segments = append(segments, seg)
			continue
		}

		lastInBlock := inBlock && len(*list) == 1
		if text, ok := seg.(*TextSegment); ok && text.Text == "" && !lastInBlock && drop(off) {
			removed = true
			continue
		}

		segments = append(segments, seg)
		off += utf8.RuneCountInString(seg.Textual())
	}

	if removed {
		*list = segments
	}
	return off, removed
}

// mergeSegments joins neighbouring text segments that share a style, keeping the
// segment list as short as the content allows.
func mergeSegments(in []RichTextSegment) []RichTextSegment {
	out := make([]RichTextSegment, 0, len(in))
	for _, seg := range in {
		text, isText := seg.(*TextSegment)
		if !isText {
			out = append(out, seg)
			continue
		}

		if len(out) > 0 {
			if prev, ok := out[len(out)-1].(*TextSegment); ok && prev.Style == text.Style {
				prev.Text += text.Text
				continue
			}
		}
		if text.Text == "" && len(in) > 1 {
			continue
		}
		out = append(out, &TextSegment{Style: text.Style, Text: text.Text})
	}

	if len(out) == 0 {
		out = append(out, &TextSegment{Style: RichTextStyleInline})
	}
	return out
}

// editableSegments converts parsed rich text into a form where every character
// can be edited: a flat run of inline segments, with block boundaries expressed
// as newline characters instead of block styling.
// Lists and code blocks are kept, as they draw the bullets and panel that their
// content sits on, with the content inside them flattened in the same way.
func editableSegments(in []RichTextSegment) []RichTextSegment {
	return mergeSegments(trimTrailingNewline(flattenSegments(in, nil)))
}

// trimTrailingNewline removes the line break that closes the final block, as
// there is no following content for it to separate.
func trimTrailingNewline(segments []RichTextSegment) []RichTextSegment {
	content := appendContentSegments(segments, nil)
	for i := len(content) - 1; i >= 0; i-- {
		holder, ok := content[i].(textHolder)
		if !ok {
			return segments
		}
		if holder.content() == "" {
			continue
		}

		holder.setContent(strings.TrimSuffix(holder.content(), newLineChar))
		return segments
	}
	return segments
}

// endsWithNewline reports whether the flattened content already ends with a
// line break, so that nested blocks do not each add one of their own.
func endsWithNewline(segments []RichTextSegment) bool {
	content := appendContentSegments(segments, nil)
	for i := len(content) - 1; i >= 0; i-- {
		text := content[i].Textual()
		if text == "" {
			continue
		}
		return strings.HasSuffix(text, newLineChar)
	}
	return false
}

func flattenSegments(in []RichTextSegment, out []RichTextSegment) []RichTextSegment {
	appendText := func(style RichTextStyle, text string) {
		style.Inline = true
		out = append(out, &TextSegment{Style: style, Text: text})
	}
	newline := func() {
		if len(out) == 0 || endsWithNewline(out) {
			return // nothing to break away from, or a break is already there
		}
		appendText(RichTextStyleInline, newLineChar)
	}

	for _, seg := range in {
		switch t := seg.(type) {
		case *TextSegment:
			if t.Text != "" {
				appendText(t.Style, t.Text)
			}
			if !t.Style.Inline {
				newline()
			}
		case *HyperlinkSegment:
			out = append(out, t)
		case *ImageSegment:
			out = append(out, t)
			newline()
		case *SeparatorSegment:
			newline()
			appendText(RichTextStyleInline, "---")
			newline()
		case *CodeBlockSegment:
			newline()
			// the last line of the block ends where the panel does
			t.Text = strings.TrimSuffix(t.Text, newLineChar) + newLineChar
			out = append(out, t)
		case *CheckBoxSegment:
			marker := "- [ ] "
			if t.Checked {
				marker = "- [x] "
			}
			appendText(RichTextStyleInline, marker+t.Text)
			newline()
		case *ListSegment:
			newline()
			out = append(out, editableList(t))
		case *TableSegment:
			out = flattenTable(t, out)
		case *ParagraphSegment:
			out = flattenSegments(t.Texts, out)
			newline()
		default:
			if block, ok := seg.(RichTextBlock); ok {
				out = flattenSegments(block.Segments(), out)
				if !seg.Inline() {
					newline()
				}
				continue
			}
			out = append(out, seg)
		}
	}
	return out
}

// editableList rewrites the items of a list so that each one holds a flat run of
// inline segments, ending in the line break that closes the item.
func editableList(l *ListSegment) *ListSegment {
	items := make([]RichTextSegment, 0, len(l.Items))
	for _, item := range l.Items {
		if sub, ok := item.(*ListSegment); ok { // a nested list follows its parent item
			items = append(items, editableList(sub))
			continue
		}

		texts := mergeSegments(flattenSegments([]RichTextSegment{item}, nil))
		if !endsWithNewline(texts) {
			texts = append(texts, &TextSegment{Style: RichTextStyleInline, Text: newLineChar})
		}
		items = append(items, &ParagraphSegment{Texts: texts})
	}

	l.Items = items
	l.markers = nil
	return l
}

func flattenTable(t *TableSegment, out []RichTextSegment) []RichTextSegment {
	row := func(cells [][]RichTextSegment, style RichTextStyle) {
		line := &strings.Builder{}
		line.WriteString("| ")
		for i, cell := range cells {
			if i > 0 {
				line.WriteString(" | ")
			}
			for _, seg := range cell {
				line.WriteString(seg.Textual())
			}
		}
		line.WriteString(" |\n")

		style.Inline = true
		out = append(out, &TextSegment{Style: style, Text: line.String()})
	}

	if len(t.Headers) > 0 {
		row(t.Headers, RichTextStyleStrong)
	}
	for _, r := range t.Rows {
		row(r, RichTextStyleInline)
	}
	return out
}

// segmentsToMarkdown writes styled segments back out as markdown. Each line of
// content becomes its own block, so that parsing the result returns content
// with the same line structure.
func segmentsToMarkdown(segments []RichTextSegment) string {
	out := &strings.Builder{}
	pending := 0 // line breaks to write before the next content
	prefix := "" // written at the start of the next line, such as a list bullet
	atLineStart := true

	// startLine writes what separates this line from the last, along with the
	// marker that introduces it.
	startLine := func(style RichTextStyle) {
		if pending > 0 {
			if out.Len() > 0 {
				out.WriteString(strings.Repeat(newLineChar, pending))
			}
			pending, atLineStart = 0, true
		}
		if !atLineStart {
			return
		}

		out.WriteString(prefix)
		out.WriteString(markdownBlockPrefix(style))
		prefix, atLineStart = "", false
	}

	write := func(text string, style RichTextStyle) {
		if text == "" {
			return
		}
		startLine(style)

		marks := markdownInlineMarks(style)
		out.WriteString(marks)
		out.WriteString(text)
		out.WriteString(reverseMarks(marks))
	}

	writeCode := func(code *CodeBlockSegment) {
		quote := strings.Repeat(markdownQuote, code.quotingLevel)
		startLine(RichTextStyleInline)

		out.WriteString(quote + markdownCodeFence + newLineChar)
		for _, line := range strings.Split(strings.TrimSuffix(code.Text, newLineChar), newLineChar) {
			out.WriteString(quote + line + newLineChar)
		}
		out.WriteString(quote + markdownCodeFence)
		atLineStart = false
	}

	var writeSegments func(segs []RichTextSegment)
	var writeList func(list *ListSegment)

	writeList = func(list *ListSegment) {
		number := list.StartNumber()
		for _, item := range list.Items {
			if sub, ok := item.(*ListSegment); ok { // a nested list follows its item
				writeList(sub)
				continue
			}

			bullet := "- "
			if list.Ordered {
				bullet = strconv.Itoa(number) + ". "
				number++
			}
			prefix = strings.Repeat(markdownQuote, list.quotingLevel) +
				strings.Repeat(" ", list.indentationLevel*listIndentSpaces) + bullet

			startLine(RichTextStyleInline) // an item with no text still has a bullet
			writeSegments(blockContent(item))
			pending, atLineStart = 1, false // the next item is the next line, not a new block
		}
	}

	writeSegments = func(segs []RichTextSegment) {
		for _, seg := range segs {
			switch t := seg.(type) {
			case *TextSegment:
				for i, line := range strings.Split(t.Text, newLineChar) {
					if i > 0 {
						// each break starts a new block, which markdown separates by a blank line
						pending += 2
					}
					write(line, t.Style)
				}
			case *HyperlinkSegment:
				link := ""
				if t.URL != nil {
					link = t.URL.String()
				}
				write("["+t.Text+"]("+link+")", RichTextStyleInline)
			case *ImageSegment:
				source := ""
				if t.Source != nil {
					source = t.Source.String()
				}
				write("!["+t.Title+"]("+source+")", RichTextStyleInline)
			case *ListSegment:
				writeList(t)
				pending = 2
			case *CodeBlockSegment:
				writeCode(t)
				pending = 2
			case *ParagraphSegment:
				writeSegments(t.Texts)
			default:
				write(seg.Textual(), RichTextStyleInline)
			}
		}
	}

	writeSegments(segments)
	return out.String()
}

func markdownBlockPrefix(style RichTextStyle) string {
	prefix := ""
	if style.QuotingDepth > 0 {
		prefix = strings.Repeat(markdownQuote, style.QuotingDepth)
	}

	switch style.SizeName {
	case theme.SizeNameHeadingText:
		return prefix + "# "
	case theme.SizeNameSubHeadingText:
		return prefix + "## "
	}
	return prefix
}

func markdownInlineMarks(style RichTextStyle) string {
	if style.SizeName == theme.SizeNameHeadingText || style.SizeName == theme.SizeNameSubHeadingText {
		return "" // the heading prefix already covers the emphasis
	}

	marks := ""
	if style.TextStyle.Strikethrough {
		marks += markdownStrike
	}
	if style.TextStyle.Bold {
		marks += markdownBold
	}
	if style.TextStyle.Italic && style.QuotingDepth == 0 {
		marks += "*"
	}
	if style.codeInline {
		marks += "`"
	}
	return marks
}

// reverseMarks flips opening markdown markers into the matching closers.
func reverseMarks(marks string) string {
	var out []string
	for i := 0; i < len(marks); {
		switch {
		case strings.HasPrefix(marks[i:], markdownStrike):
			out = append(out, markdownStrike)
			i += 2
		case strings.HasPrefix(marks[i:], markdownBold):
			out = append(out, markdownBold)
			i += 2
		default:
			out = append(out, marks[i:i+1])
			i++
		}
	}

	reversed := &strings.Builder{}
	for i := len(out) - 1; i >= 0; i-- {
		reversed.WriteString(out[i])
	}
	return reversed.String()
}

var markdownInlineStyles = []struct {
	marker string
	apply  func(*RichTextStyle)
}{
	{markdownStrike, func(s *RichTextStyle) { s.TextStyle.Strikethrough = true }},
	{markdownBold, func(s *RichTextStyle) { s.TextStyle.Bold = true }},
	{"__", func(s *RichTextStyle) { s.TextStyle.Bold = true }},
	{"`", func(s *RichTextStyle) { s.TextStyle.Monospace = true; s.codeInline = true }},
	{"*", func(s *RichTextStyle) { s.TextStyle.Italic = true }},
	{"_", func(s *RichTextStyle) { s.TextStyle.Italic = true }},
}

// styleMarkdownAtCursor looks for markdown that was completed by the rune just
// typed and, if it finds any, replaces it with the style that it describes.
// It reports whether the content was changed.
func (e *RichTextEntry) styleMarkdownAtCursor(typed rune) bool {
	pos := e.CursorTextOffset()
	runes := []rune(e.Text)
	if pos > len(runes) {
		return false
	}

	lineStart := 0
	for i := pos - 1; i >= 0; i-- {
		if runes[i] == '\n' {
			lineStart = i + 1
			break
		}
	}

	if typed == ' ' {
		return e.styleMarkdownBlock(runes, lineStart, pos)
	}
	return e.styleMarkdownInline(runes, lineStart, pos)
}

func (e *RichTextEntry) styleMarkdownInline(runes []rune, lineStart, pos int) bool {
	line := runes[lineStart:pos]

	for _, m := range markdownInlineStyles {
		marker := []rune(m.marker)
		if len(line) <= len(marker) || string(line[len(line)-len(marker):]) != m.marker {
			continue
		}

		body := line[:len(line)-len(marker)]
		open := lastIndexRunes(body, marker)
		if open < 0 {
			continue
		}

		content := body[open+len(marker):]
		if len(content) == 0 || strings.TrimSpace(string(content)) == "" ||
			content[0] == ' ' || content[len(content)-1] == ' ' {
			continue
		}
		if len(marker) == 1 && m.marker != "`" {
			// a single marker must not be half of a "**" or "__" pair
			if open > 0 && body[open-1] == marker[0] {
				continue
			}
			if content[len(content)-1] == marker[0] {
				continue
			}
		}

		openPos := lineStart + open
		base := e.StyleAtCursor()
		provider := e.richProvider()

		provider.deleteFromTo(pos-len(marker), pos)
		provider.deleteFromTo(openPos, openPos+len(marker))

		end := openPos + len(content)
		e.styleRange(openPos, end, m.apply)
		provider.Segments = mergeSegments(provider.Segments)

		// text typed after the styled run returns to the style around it
		e.insertEmptySegmentAt(end, base)

		e.finishStyling(end)
		return true
	}
	return false
}

func (e *RichTextEntry) styleMarkdownBlock(runes []rune, lineStart, pos int) bool {
	prefix := string(runes[lineStart : pos-1])
	if ordered, number, indent, ok := markdownListPrefix(prefix); ok {
		return e.startListAt(lineStart, pos, ordered, number, indent)
	}

	style, ok := markdownBlockStyle(prefix)
	if !ok {
		return false
	}

	lineEnd := lineStart
	for lineEnd < len(runes) && runes[lineEnd] != '\n' {
		lineEnd++
	}
	lineEnd -= pos - lineStart // the marker and its trailing space are removed

	provider := e.richProvider()
	provider.deleteFromTo(lineStart, pos)

	if lineEnd > lineStart {
		e.styleRange(lineStart, lineEnd, func(s *RichTextStyle) { *s = style })
	}
	provider.Segments = mergeSegments(provider.Segments)

	// so that the rest of the line is typed in the new style
	e.insertEmptySegmentAt(lineStart, style)

	e.finishStyling(lineStart)
	return true
}

func markdownBlockStyle(prefix string) (RichTextStyle, bool) {
	quoting := 0
	for strings.HasPrefix(prefix, ">") {
		quoting++
		prefix = strings.TrimPrefix(strings.TrimPrefix(prefix, ">"), " ")
	}

	var style RichTextStyle
	switch prefix {
	case "":
		if quoting == 0 {
			return style, false
		}
		style = RichTextStyleBlockquote
	case "#":
		style = RichTextStyleHeading
	case "##":
		style = RichTextStyleSubHeading
	case "###":
		style = RichTextStyleStrong
	default:
		return style, false
	}

	style.Inline = true
	style.QuotingDepth = quoting
	if quoting > 0 {
		style.TextStyle.Italic = true
	}
	return style, true
}

// finishStyling updates the entry state after segments have been restyled.
func (e *RichTextEntry) finishStyling(cursor int) {
	provider := e.richProvider()

	content := provider.String()
	e.updateText(content, false)
	e.setCursorOffset(cursor)

	// the undo stack holds plain text, which can no longer describe this change
	e.undoStack.Clear()

	if e.OnChanged != nil {
		e.OnChanged(content)
	}
	e.validate()
}

func lastIndexRunes(haystack, needle []rune) int {
	for i := len(haystack) - len(needle); i >= 0; i-- {
		if string(haystack[i:i+len(needle)]) == string(needle) {
			return i
		}
	}
	return -1
}

// richListItem describes where a rune offset falls inside a list.
type richListItem struct {
	list  *ListSegment
	owner *[]RichTextSegment // the segments that hold the list
	index int                // which item of the list holds the offset
	start int                // the rune offset that the item content begins at
}

// listItemAt returns the list item holding the given rune offset. An offset in a
// nested list belongs to the innermost list that holds it.
func (e *RichTextEntry) listItemAt(pos int) (richListItem, bool) {
	provider := e.richProvider()
	item, _, ok := findListItem(&provider.Segments, pos, 0)
	return item, ok
}

// findListItem walks segments in document order looking for the list item that
// holds the offset. It returns the item, the offset reached and whether the item
// was found.
func findListItem(owner *[]RichTextSegment, pos, off int) (richListItem, int, bool) {
	for _, seg := range *owner {
		if list, ok := seg.(*ListSegment); ok {
			item, next, found := findInList(list, owner, pos, off)
			if found {
				return item, next, true
			}

			off = next
			continue
		}

		if inner := blockContainer(seg); inner != nil {
			item, next, found := findListItem(inner, pos, off)
			if found {
				return item, next, true
			}

			off = next
			continue
		}

		off += utf8.RuneCountInString(seg.Textual())
	}

	return richListItem{}, off, false
}

func findInList(list *ListSegment, owner *[]RichTextSegment, pos, off int) (richListItem, int, bool) {
	for i, item := range list.Items {
		if sub, ok := item.(*ListSegment); ok {
			found, next, ok := findInList(sub, &list.Items, pos, off)
			if ok {
				return found, next, true
			}

			off = next
			continue
		}

		length := contentLength(item)
		end := off + length
		if !endsWithNewline([]RichTextSegment{item}) {
			end++ // the last item of the content has no line break to close it
		}
		if pos >= off && pos < end {
			return richListItem{list: list, owner: owner, index: i, start: off}, off, true
		}
		off += length
	}

	return richListItem{}, off, false
}

// contentLength returns how many runes of content a segment holds, including any
// that are inside it.
func contentLength(seg RichTextSegment) int {
	total := 0
	for _, leaf := range appendContentSegments([]RichTextSegment{seg}, nil) {
		total += utf8.RuneCountInString(leaf.Textual())
	}
	return total
}

// splitListItem divides the list item at the given offset in two, so that a new
// bullet appears for the text that follows the cursor. Pressing return on an item
// with no content instead closes the list.
// It reports whether the content was changed.
func (e *RichTextEntry) splitListItem(pos int) bool {
	item, ok := e.listItemAt(pos)
	if !ok {
		return false
	}

	content := item.list.Items[item.index]
	if contentLength(content) <= 1 && item.index == len(item.list.Items)-1 {
		return e.leaveList(item, pos) // an empty item at the end closes the list
	}

	texts := blockContent(content)
	closed := endsWithNewline(texts) // the last item of the content is left open

	head, tail := splitContent(texts, pos-item.start)
	head = closeItem(head)
	if closed {
		tail = closeItem(tail)
	} else {
		tail = mergeSegments(tail)
	}

	items := make([]RichTextSegment, 0, len(item.list.Items)+1)
	items = append(items, item.list.Items[:item.index]...)
	items = append(items, &ParagraphSegment{Texts: head}, &ParagraphSegment{Texts: tail})
	items = append(items, item.list.Items[item.index+1:]...)
	item.list.Items = items
	item.list.markers = nil // the bullets are numbered again from the items

	e.finishStyling(pos + 1)
	return true
}

// leaveList takes the empty item out of its list, so that typing carries on below
// the list rather than against a bullet.
func (e *RichTextEntry) leaveList(item richListItem, pos int) bool {
	list := item.list
	list.Items = append(list.Items[:item.index], list.Items[item.index+1:]...)
	list.markers = nil

	segments := *item.owner
	at := indexOfSegment(segments, list) + 1
	if len(list.Items) == 0 { // the list has nothing left in it
		at--
		segments = append(segments[:at], segments[at+1:]...)
	}

	out := make([]RichTextSegment, 0, len(segments)+1)
	out = append(out, segments[:at]...)
	out = append(out, &TextSegment{Style: RichTextStyleInline})
	out = append(out, segments[at:]...)
	*item.owner = out

	e.finishStyling(pos)
	return true
}

// removeBullet takes the item at the cursor out of its list when the cursor is at
// the start of it, so that a backspace removes the bullet before the text.
// It reports whether the content was changed.
func (e *RichTextEntry) removeBullet(pos int) bool {
	item, ok := e.listItemAt(pos)
	if !ok || pos != item.start {
		return false
	}

	list := item.list
	texts := blockContent(list.Items[item.index])
	list.Items = append(list.Items[:item.index], list.Items[item.index+1:]...)
	list.markers = nil

	if item.index > 0 { // the text joins the end of the item above it
		prev := trimTrailingNewline(blockContent(list.Items[item.index-1]))
		list.Items[item.index-1] = &ParagraphSegment{Texts: append(prev, texts...)}

		e.finishStyling(pos - 1)
		return true
	}

	// the first item has no bullet above it to join, so it leaves the list
	segments := *item.owner
	at := indexOfSegment(segments, list)
	if len(list.Items) == 0 {
		segments = append(segments[:at], segments[at+1:]...)
	}

	out := make([]RichTextSegment, 0, len(segments)+len(texts))
	out = append(out, segments[:at]...)
	out = append(out, texts...)
	out = append(out, segments[at:]...)
	*item.owner = out

	e.finishStyling(pos)
	return true
}

// splitContent divides a run of segments at the given rune offset, dividing a text
// segment where the offset falls inside one.
func splitContent(in []RichTextSegment, at int) (head, tail []RichTextSegment) {
	off := 0
	for _, seg := range in {
		length := utf8.RuneCountInString(seg.Textual())
		switch {
		case off+length <= at:
			head = append(head, seg)
		case off >= at:
			tail = append(tail, seg)
		default:
			text, isText := seg.(*TextSegment)
			if !isText { // an object cannot be divided, it goes with the tail
				tail = append(tail, seg)
				break
			}

			runes := []rune(text.Text)
			cut := at - off
			tail = append(tail, &TextSegment{Style: text.Style, Text: string(runes[cut:])})
			text.Text = string(runes[:cut])
			head = append(head, text)
		}
		off += length
	}

	return head, tail
}

// closeItem makes sure that the content of a list item ends in the line break
// that closes it.
func closeItem(texts []RichTextSegment) []RichTextSegment {
	if !endsWithNewline(texts) {
		texts = append(texts, &TextSegment{Style: RichTextStyleInline, Text: newLineChar})
	}

	return mergeSegments(texts)
}

func indexOfSegment(segments []RichTextSegment, seg RichTextSegment) int {
	for i, in := range segments {
		if in == seg {
			return i
		}
	}

	return len(segments)
}

// codeBlockAt returns the code block that holds the given rune offset, with the
// segments that hold the block and its index in them.
func codeBlockAt(owner *[]RichTextSegment, pos, off int) (out *[]RichTextSegment,
	index int, offset int, found bool,
) {
	for i, seg := range *owner {
		if inner := blockContainer(seg); inner != nil {
			list, index, next, ok := codeBlockAt(inner, pos, off)
			if ok {
				return list, index, next, true
			}

			off = next
			continue
		}

		length := utf8.RuneCountInString(seg.Textual())
		if _, ok := seg.(*CodeBlockSegment); ok && pos > off && pos <= off+length {
			return owner, i, off, true
		}
		off += length
	}

	return nil, 0, off, false
}

// toggleCodeFence acts on a "```" that has just been completed by a return,
// opening a code block for the lines that follow or, when the fence was typed
// inside a block, closing it again.
// It reports whether the content was changed.
func (e *RichTextEntry) toggleCodeFence() bool {
	pos := e.CursorTextOffset()
	runes := []rune(e.Text)
	if pos > len(runes) {
		return false
	}

	lineStart := 0
	for i := pos - 1; i >= 0; i-- {
		if runes[i] == '\n' {
			lineStart = i + 1
			break
		}
	}
	if !strings.HasPrefix(strings.TrimSpace(string(runes[lineStart:pos])), markdownCodeFence) {
		return false
	}

	provider := e.richProvider()
	if owner, index, blockStart, ok := codeBlockAt(&provider.Segments, pos, 0); ok {
		// the fence closes the block, so it and the line it is on are removed
		from := lineStart
		if from > blockStart {
			from--
		}
		block, ok := (*owner)[index].(*CodeBlockSegment)
		if !ok {
			return false
		}
		provider.deleteFromTo(from, pos)

		end := blockStart + utf8.RuneCountInString(block.Text)
		e.dropEmptySegmentsAt(end) // only one segment may claim the text typed here
		insertSegmentAt(owner, indexOfSegment(*owner, block)+1, &TextSegment{Style: RichTextStyleInline})

		e.finishStyling(end)
		return true
	}

	provider.deleteFromTo(lineStart, pos)
	e.dropEmptySegmentsAt(lineStart)
	list, at := e.splitAt(lineStart)
	insertSegmentAt(list, at, &CodeBlockSegment{Text: newLineChar})

	e.finishStyling(lineStart)
	return true
}

func insertSegmentAt(list *[]RichTextSegment, index int, seg RichTextSegment) {
	index = min(index, len(*list))
	segments := make([]RichTextSegment, 0, len(*list)+1)
	segments = append(segments, (*list)[:index]...)
	segments = append(segments, seg)
	segments = append(segments, (*list)[index:]...)
	*list = segments
}

// markdownListPrefix reports whether the text typed at the start of a line asks
// for a list, returning what kind of list it describes.
func markdownListPrefix(prefix string) (ordered bool, number, indent int, ok bool) {
	for _, r := range prefix {
		if r == '\t' {
			indent++
			continue
		}
		if r != ' ' {
			break
		}
		number++ // counting spaces for now, the indent is worked out below
	}

	indent += number / listIndentSpaces
	prefix = strings.TrimLeft(prefix, " \t")

	switch prefix {
	case "-", "*", "+":
		return false, 0, indent, true
	}

	if digits, found := strings.CutSuffix(prefix, "."); found {
		if start, err := strconv.Atoi(digits); err == nil {
			return true, start, indent, true
		}
	}
	return false, 0, 0, false
}

// startListAt turns the line at the given offset into the item of a list, joining
// a list that the line follows where there is one.
// It reports whether the content was changed.
func (e *RichTextEntry) startListAt(lineStart, pos int, ordered bool, number, indent int) bool {
	provider := e.richProvider()
	provider.deleteFromTo(lineStart, pos) // the bullet is drawn by the list from here on
	e.dropEmptySegmentsAt(lineStart)

	lineEnd := lineStart
	for runes := []rune(provider.String()); lineEnd < len(runes); lineEnd++ {
		if runes[lineEnd] == '\n' {
			lineEnd++ // the line break closes the item
			break
		}
	}

	list, from := e.splitAt(lineStart)
	other, to := e.splitAt(lineEnd)
	if other != list { // the line ends outside of these segments
		to = len(*list)
	}

	texts := closeItem(append([]RichTextSegment{}, (*list)[from:to]...))
	item := &ParagraphSegment{Texts: texts}

	segments := make([]RichTextSegment, 0, len(*list))
	segments = append(segments, (*list)[:from]...)
	if from > 0 {
		if prev, ok := segments[from-1].(*ListSegment); ok && prev.Ordered == ordered && prev.indentationLevel == indent {
			prev.Items = append(prev.Items, item) // this line carries on the list above it
			prev.markers = nil

			*list = append(segments, (*list)[to:]...)
			e.finishStyling(lineStart)
			return true
		}
	}

	added := &ListSegment{Items: []RichTextSegment{item}, Ordered: ordered, indentationLevel: indent}
	if ordered {
		added.SetStartNumber(number)
	}

	segments = append(segments, added)
	*list = append(segments, (*list)[to:]...)

	e.finishStyling(lineStart)
	return true
}

// clearEmptyStyle takes the styling off a line that has nothing typed on it, so
// that a delete on an empty heading, quote or code block removes the style that
// it was going to be typed in rather than the line break before it.
// It reports whether the content was changed.
func (e *RichTextEntry) clearEmptyStyle() bool {
	provider := e.richProvider()
	bound := provider.rowBoundary(e.CursorRow)
	if bound == nil || len(provider.row(e.CursorRow)) > 0 {
		return false // there is content on this line to remove first
	}

	if block, ok := bound.panel.(*CodeBlockSegment); ok && strings.TrimSuffix(block.Text, newLineChar) == "" {
		return e.removeEmptyBlock(block)
	}

	pos := e.CursorTextOffset()
	style, ok := e.emptyStyleAt(pos)
	if !ok || isPlainStyle(style) {
		return false
	}

	e.insertEmptySegmentAt(pos, RichTextStyleInline)
	e.finishStyling(pos)
	return true
}

// removeEmptyBlock takes a block with nothing in it out of the content, leaving a
// plain line where the block was.
func (e *RichTextEntry) removeEmptyBlock(block RichTextSegment) bool {
	provider := e.richProvider()
	owner, index, start, ok := ownerOf(&provider.Segments, block, 0)
	if !ok {
		return false
	}

	segments := make([]RichTextSegment, 0, len(*owner))
	segments = append(segments, (*owner)[:index]...)
	segments = append(segments, (*owner)[index+1:]...)
	*owner = segments

	e.insertEmptySegmentAt(start, RichTextStyleInline)
	e.finishStyling(start)
	return true
}

// emptyStyleAt returns the style of an empty segment at the given rune offset,
// which is the style that text typed there would take.
func (e *RichTextEntry) emptyStyleAt(pos int) (RichTextStyle, bool) {
	off := 0
	for _, seg := range e.richProvider().contentSegments() {
		length := utf8.RuneCountInString(seg.Textual())
		if text, isText := seg.(*TextSegment); isText && length == 0 && off == pos {
			return text.Style, true
		}

		off += length
		if off > pos {
			break
		}
	}

	return RichTextStyle{}, false
}

// isPlainStyle reports whether a style adds nothing to ordinary body text.
func isPlainStyle(style RichTextStyle) bool {
	return style.TextStyle == (fyne.TextStyle{}) && !isBlockStyle(style)
}

// ownerOf returns the segments that hold the given segment, its index in them and
// the rune offset that it starts at.
func ownerOf(list *[]RichTextSegment, seg RichTextSegment, off int) (out *[]RichTextSegment,
	index int, offset int, found bool,
) {
	for i, in := range *list {
		if in == seg {
			return list, i, off, true
		}

		if inner := blockContainer(in); inner != nil {
			owner, index, next, ok := ownerOf(inner, seg, off)
			if ok {
				return owner, index, next, true
			}

			off = next
			continue
		}

		off += contentLength(in)
	}

	return nil, 0, off, false
}

// joinIntoBlockAbove moves the line at the cursor into the list item, or code
// block, that ends where the line begins.
// It reports whether the content was changed.
func (e *RichTextEntry) joinIntoBlockAbove(pos int) bool {
	runes := []rune(e.Text)
	if pos <= 0 || pos > len(runes) || runes[pos-1] != '\n' {
		return false // this is not the start of a line
	}

	if item, ok := e.listItemAt(pos - 1); ok && item.start+contentLength(item.list.Items[item.index]) == pos {
		texts, taken := e.takeLine(pos)
		if !taken {
			return false
		}

		content := trimTrailingNewline(blockContent(item.list.Items[item.index]))
		item.list.Items[item.index] = &ParagraphSegment{Texts: mergeSegments(append(content, texts...))}

		e.finishStyling(pos - 1)
		return true
	}

	provider := e.richProvider()
	owner, index, _, ok := codeBlockAt(&provider.Segments, pos, 0)
	if !ok {
		return false
	}

	block, ok := (*owner)[index].(*CodeBlockSegment)
	if !ok {
		return false
	}
	if bound := provider.rowBoundary(e.CursorRow); bound == nil || bound.panel == block {
		return false // the cursor is inside the block, not on the line below it
	}

	texts, taken := e.takeLine(pos)
	if !taken {
		return false
	}

	line := &strings.Builder{}
	for _, seg := range appendContentSegments(texts, nil) {
		line.WriteString(seg.Textual())
	}
	block.Text = strings.TrimSuffix(block.Text, newLineChar) + line.String()

	e.finishStyling(pos - 1)
	return true
}

// takeLine removes the segments that make up the line starting at the given rune
// offset from the content, returning them so that they can be added elsewhere.
func (e *RichTextEntry) takeLine(pos int) ([]RichTextSegment, bool) {
	runes := []rune(e.Text)
	lineEnd := pos
	for ; lineEnd < len(runes); lineEnd++ {
		if runes[lineEnd] == '\n' {
			lineEnd++ // the break that closes the line goes with it
			break
		}
	}

	list, from := e.splitAt(pos)
	other, to := e.splitAt(lineEnd)
	if other != list || lineEnd >= len(runes) {
		to = len(*list) // the line runs past the end of these segments
	}
	if from > to {
		return nil, false
	}

	texts := append([]RichTextSegment{}, (*list)[from:to]...)
	segments := make([]RichTextSegment, 0, len(*list))
	segments = append(segments, (*list)[:from]...)
	segments = append(segments, (*list)[to:]...)
	*list = segments
	return texts, true
}

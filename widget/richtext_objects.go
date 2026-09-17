package widget

import (
	"image/color"
	"net/url"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/scale"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

// listIndentSpaces is how many spaces of indentation each level of list nesting adds.
const listIndentSpaces = 4

var (
	// RichTextStyleBlockquote represents a quote presented in an indented block.
	//
	// Since: 2.1
	RichTextStyleBlockquote = RichTextStyle{
		ColorName: theme.ColorNameForeground,
		Inline:    true,
		SizeName:  theme.SizeNameText,
		TextStyle: fyne.TextStyle{Italic: true},
	}
	// RichTextStyleCodeBlock represents a code blog segment.
	//
	// Since: 2.1
	RichTextStyleCodeBlock = RichTextStyle{
		ColorName: theme.ColorNameForeground,
		Inline:    false,
		SizeName:  theme.SizeNameText,
		TextStyle: fyne.TextStyle{Monospace: true},
	}
	// RichTextStyleCodeInline represents an inline code segment.
	//
	// Since: 2.1
	RichTextStyleCodeInline = RichTextStyle{
		ColorName:  theme.ColorNameForeground,
		Inline:     true,
		SizeName:   theme.SizeNameText,
		TextStyle:  fyne.TextStyle{Monospace: true},
		codeInline: true,
	}
	// RichTextStyleEmphasis represents regular text with emphasis.
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
		Inline:    true,
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
	// RichTextStyleStrong represents regular text with a strong emphasis.
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
		Inline:    true,
		SizeName:  theme.SizeNameSubHeadingText,
		TextStyle: fyne.TextStyle{Bold: true},
	}
)

// HyperlinkSegment represents a hyperlink within a rich text widget.
//
// Since: 2.1
type HyperlinkSegment struct {
	Alignment fyne.TextAlign
	Text      string
	URL       *url.URL

	// OnTapped overrides the default `fyne.OpenURL` call when the link is tapped
	//
	// Since: 2.4
	OnTapped func() `json:"-"`

	// Since 2.8
	TextStyle fyne.TextStyle
	// Since 2.8
	SizeName     fyne.ThemeSizeName // The theme name of the text size to use, if blank will be the standard text size
	quotingLevel int
}

// Inline returns true as hyperlinks are inside other elements.
func (*HyperlinkSegment) Inline() bool {
	return true
}

// Textual returns the content of this segment rendered to plain text.
func (h *HyperlinkSegment) Textual() string {
	return h.Text
}

// Visual returns a new instance of a hyperlink widget required to render this segment.
func (h *HyperlinkSegment) Visual() fyne.CanvasObject {
	link := NewHyperlink(h.Text, h.URL)
	link.Alignment = h.Alignment
	link.OnTapped = h.OnTapped
	return &fyne.Container{Layout: &unpadTextWidgetLayout{parent: link}, Objects: []fyne.CanvasObject{link}}
}

// Update applies the current state of this hyperlink segment to an existing visual.
func (h *HyperlinkSegment) Update(o fyne.CanvasObject) {
	link, _ := o.(*fyne.Container).Objects[0].(*Hyperlink)
	link.URL = h.URL
	link.Alignment = h.Alignment
	link.SizeName = h.SizeName
	link.TextStyle = h.TextStyle
	link.OnTapped = h.OnTapped
	link.Refresh()
}

// Select tells the segment that the user is selecting the content between the two positions.
func (*HyperlinkSegment) Select(_, _ fyne.Position) {
	// no-op: this will be added when we progress to editor
}

// SelectedText should return the text representation of any content currently selected through the Select call.
func (*HyperlinkSegment) SelectedText() string {
	// no-op: this will be added when we progress to editor
	return ""
}

// Unselect tells the segment that the user is has cancelled the previous selection.
func (*HyperlinkSegment) Unselect() {
	// no-op: this will be added when we progress to editor
}

// ImageSegment represents an image within a rich text widget.
//
// Since: 2.3
type ImageSegment struct {
	Source fyne.URI
	Title  string

	// Alignment specifies the horizontal alignment of this image segment
	// Since: 2.4
	Alignment fyne.TextAlign
}

// Inline returns false as images in rich text are blocks.
func (*ImageSegment) Inline() bool {
	return false
}

// Textual returns the content of this segment rendered to plain text.
func (i *ImageSegment) Textual() string {
	return "Image " + i.Title
}

// Visual returns a new instance of an image widget required to render this segment.
func (i *ImageSegment) Visual() fyne.CanvasObject {
	return newRichImage(i.Source, i.Alignment)
}

// Update applies the current state of this image segment to an existing visual.
func (i *ImageSegment) Update(o fyne.CanvasObject) {
	newer := canvas.NewImageFromURI(i.Source)
	img, _ := o.(*richImage)

	// one of the following will be used
	img.img.File = newer.File
	img.img.Resource = newer.Resource
	img.setAlign(i.Alignment)

	img.Refresh()
}

// Select tells the segment that the user is selecting the content between the two positions.
func (*ImageSegment) Select(_, _ fyne.Position) {
	// no-op: this will be added when we progress to editor
}

// SelectedText should return the text representation of any content currently selected through the Select call.
func (*ImageSegment) SelectedText() string {
	// no-op: images have no text rendering
	return ""
}

// Unselect tells the segment that the user is has cancelled the previous selection.
func (*ImageSegment) Unselect() {
	// no-op: this will be added when we progress to editor
}

// ListSegment includes an itemised list with the content set using the Items field.
//
// Since: 2.1
type ListSegment struct {
	Items   []RichTextSegment
	Ordered bool

	// startIndex is the starting number - 1 (If it is ordered). Unordered lists
	// ignore startIndex.
	//
	// startIndex is set to start - 1 to allow the empty value of ListSegment to have a starting
	// number of 1, while also allowing the caller to override the starting
	// number to any int, including 0.
	startIndex       int
	indentationLevel int
	quotingLevel     int

	// markers are re-used between calls to Segments
	markers []*listMarkerSegment
}

// SetStartNumber sets the starting number for an ordered list.
// Unordered lists are not affected.
//
// Since: 2.7
func (l *ListSegment) SetStartNumber(s int) {
	l.startIndex = s - 1
}

// StartNumber return the starting number for an ordered list.
//
// Since: 2.7
func (l *ListSegment) StartNumber() int {
	return l.startIndex + 1
}

// Inline returns false as a list should be in a block.
func (*ListSegment) Inline() bool {
	return false
}

// Segments returns the segments required to draw bullets before each item
func (l *ListSegment) Segments() []RichTextSegment {
	out := make([]RichTextSegment, len(l.Items))
	j := l.StartNumber()
	for i, in := range l.Items {
		var texts []RichTextSegment
		if _, ok := in.(*ListSegment); !ok {
			texts = append(texts, l.marker(i, j))
			j++
			if _, ok := in.(*ParagraphSegment); !ok {
				in = &ParagraphSegment{Texts: []RichTextSegment{in}}
			}
		}
		texts = append(texts, in)
		out[i] = &ParagraphSegment{Texts: texts}
	}
	return out
}

// marker returns the bullet, or number, that introduces the item at the given
// index. Markers are kept between calls so that they hold on to their visuals.
func (l *ListSegment) marker(i, number int) *listMarkerSegment {
	for len(l.markers) <= i {
		l.markers = append(l.markers, &listMarkerSegment{})
	}

	marker := l.markers[i]
	marker.ordered = l.Ordered
	marker.number = number
	marker.indent = l.indentationLevel
	marker.quoting = l.quotingLevel
	return marker
}

// Textual returns no content for a list as the content is in sub-segments.
func (*ListSegment) Textual() string {
	return ""
}

// Visual returns no additional elements for this segment.
func (*ListSegment) Visual() fyne.CanvasObject {
	return nil
}

// Update doesn't need to change a list visual.
func (*ListSegment) Update(fyne.CanvasObject) {
}

// Select does nothing for a list container.
func (*ListSegment) Select(_, _ fyne.Position) {
}

// SelectedText returns the empty string for this list.
func (*ListSegment) SelectedText() string {
	return ""
}

// Unselect does nothing for a list container.
func (*ListSegment) Unselect() {
}

// listMarkerSegment draws the bullet, or number, that introduces a list item.
// It adds no characters to the content, so an editor can treat the text of the
// item as ordinary text while the marker is drawn alongside it.
type listMarkerSegment struct {
	ordered bool
	number  int
	indent  int
	quoting int

	colorName fyne.ThemeColorName
	parent    *RichText
}

// Inline returns true as a marker is followed by the text of its item.
func (*listMarkerSegment) Inline() bool {
	return true
}

// Textual returns no content, the marker is a decoration rather than text.
func (*listMarkerSegment) Textual() string {
	return ""
}

// marker returns the text drawn to introduce this item.
func (l *listMarkerSegment) marker() string {
	bullet := "\u2022 "
	if l.ordered {
		bullet = strconv.Itoa(l.number) + "."
	}

	return strings.Repeat(" ", l.indent*listIndentSpaces) + bullet + " "
}

// Visual returns a new text object drawing this marker.
func (l *listMarkerSegment) Visual() fyne.CanvasObject {
	text := canvas.NewText("", color.Transparent)
	l.Update(text)
	return text
}

// Update applies the current state of this marker to an existing visual.
func (l *listMarkerSegment) Update(o fyne.CanvasObject) {
	text, ok := o.(*canvas.Text)
	if !ok {
		return
	}
	text.Text = l.marker()
	col := l.colorName
	if col == "" {
		col = theme.ColorNameForeground
	}
	text.Color = theme.ColorForWidget(col, l.parent)
	text.TextSize = theme.SizeForWidget(theme.SizeNameText, l.parent)
	text.TextStyle = fyne.TextStyle{Bold: true}
	text.Refresh()
}

// Select does nothing for a list marker.
func (*listMarkerSegment) Select(_, _ fyne.Position) {
}

// SelectedText returns the empty string as a marker holds no content.
func (*listMarkerSegment) SelectedText() string {
	return ""
}

// Unselect does nothing for a list marker.
func (*listMarkerSegment) Unselect() {
}

// ParagraphSegment wraps a number of text elements in a paragraph.
// It is similar to using a list of text elements when the final style is RichTextStyleParagraph.
//
// Since: 2.1
type ParagraphSegment struct {
	Texts []RichTextSegment
}

// Inline returns false as a paragraph should be in a block.
func (*ParagraphSegment) Inline() bool {
	return false
}

// Segments returns the list of text elements in this paragraph.
func (p *ParagraphSegment) Segments() []RichTextSegment {
	return p.Texts
}

// Textual returns no content for a paragraph container.
func (*ParagraphSegment) Textual() string {
	return ""
}

// Visual returns the no extra elements.
func (*ParagraphSegment) Visual() fyne.CanvasObject {
	return nil
}

// Update doesn't need to change a paragraph container.
func (*ParagraphSegment) Update(fyne.CanvasObject) {
}

// Select does nothing for a paragraph container.
func (*ParagraphSegment) Select(_, _ fyne.Position) {
}

// SelectedText returns the empty string for this paragraph container.
func (*ParagraphSegment) SelectedText() string {
	return ""
}

// Unselect does nothing for a paragraph container.
func (*ParagraphSegment) Unselect() {
}

// SeparatorSegment includes a horizontal separator in a rich text widget.
//
// Since: 2.1
type SeparatorSegment struct {
	_ bool // Without this a pointer to SeparatorSegment will always be the same.
}

// Inline returns false as a separator should be full width.
func (*SeparatorSegment) Inline() bool {
	return false
}

// Textual returns no content for a separator element.
func (*SeparatorSegment) Textual() string {
	return ""
}

// Visual returns a new instance of a separator widget for this segment.
func (*SeparatorSegment) Visual() fyne.CanvasObject {
	return NewSeparator()
}

// Update doesn't need to change a separator visual.
func (*SeparatorSegment) Update(fyne.CanvasObject) {
}

// Select does nothing for a separator.
func (*SeparatorSegment) Select(_, _ fyne.Position) {
}

// SelectedText returns the empty string for this separator.
func (*SeparatorSegment) SelectedText() string {
	return "" // TODO maybe return "---\n"?
}

// Unselect does nothing for a separator.
func (*SeparatorSegment) Unselect() {
}

// CodeBlockSegment represents a fenced or indented code block. It renders its
// content as monospace text on a panel, so the block stands apart from the
// surrounding prose.
//
// Since: 2.8
type CodeBlockSegment struct {
	Text         string
	quotingLevel int

	body *TextSegment
	bg   *richCodeBlock
}

// Inline returns false as a code block is a full-width block element.
func (*CodeBlockSegment) Inline() bool {
	return false
}

// Textual returns the raw content of this code block.
func (c *CodeBlockSegment) Textual() string {
	return c.Text
}

// content returns the code that this block holds.
func (c *CodeBlockSegment) content() string {
	return c.Text
}

// setContent replaces the code that this block holds.
func (c *CodeBlockSegment) setContent(text string) {
	c.Text = text
}

// Segments returns the content of this block as a run of monospace text, so that
// the lines of code lay out, select and edit like the text around them.
func (c *CodeBlockSegment) Segments() []RichTextSegment {
	if c.body == nil {
		c.body = &TextSegment{Style: RichTextStyleCodeBlock}
	}

	c.body.Text = c.Text
	c.body.Style.QuotingDepth = c.quotingLevel
	return []RichTextSegment{c.body}
}

// panel returns the background that the lines of this block are drawn on.
func (c *CodeBlockSegment) panel() fyne.CanvasObject {
	if c.bg == nil {
		c.bg = newRichCodeBlock()
	}
	return c.bg
}

// Visual returns a new panel widget for the background of this code block.
func (*CodeBlockSegment) Visual() fyne.CanvasObject {
	return newRichCodeBlock()
}

// Update has nothing to change, the content is drawn as text on the panel.
func (*CodeBlockSegment) Update(fyne.CanvasObject) {
}

// Select does nothing for a code block.
func (*CodeBlockSegment) Select(_, _ fyne.Position) {
}

// SelectedText returns the code block content.
func (c *CodeBlockSegment) SelectedText() string {
	return c.Text
}

// Unselect does nothing for a code block.
func (*CodeBlockSegment) Unselect() {
}

// richCodeBlock is the internal widget that draws the panel a code block sits on,
// a rounded and bordered fill behind the rows of code.
type richCodeBlock struct {
	BaseWidget
	bg *canvas.Rectangle
}

func newRichCodeBlock() *richCodeBlock {
	c := &richCodeBlock{}
	c.ExtendBaseWidget(c)
	return c
}

func (c *richCodeBlock) CreateRenderer() fyne.WidgetRenderer {
	c.bg = canvas.NewRectangle(color.Transparent)
	c.applyTheme()
	return NewSimpleRenderer(c.bg)
}

func (c *richCodeBlock) Refresh() {
	if c.bg != nil {
		c.applyTheme()
	}

	c.BaseWidget.Refresh()
}

func (c *richCodeBlock) applyTheme() {
	c.bg.FillColor = theme.ColorForWidget(theme.ColorNameInputBackground, c)
	c.bg.StrokeColor = theme.ColorForWidget(theme.ColorNameInputBorder, c)
	c.bg.StrokeWidth = theme.SizeForWidget(theme.SizeNameInputBorder, c)
	c.bg.CornerRadius = theme.SizeForWidget(theme.SizeNameInputRadius, c)
	c.bg.Refresh()
}

// CheckBoxSegment represents checkbox (with text) in a rich text widget.
//
// Since: 2.8
type CheckBoxSegment struct {
	Checked bool
	Text    string
}

// Inline returns true as a CheckBoxSegment is usually part of a list item.
func (*CheckBoxSegment) Inline() bool {
	return true
}

// Textual returns the content of this segment rendered to plain text.
func (c *CheckBoxSegment) Textual() string {
	if c.Checked {
		return "[x] "
	}
	return "[ ] "
}

// Visual returns a new instance of a check widget for this segment.
func (c *CheckBoxSegment) Visual() fyne.CanvasObject {
	check := NewCheck(c.Text, nil)
	if c.Checked {
		check.SetChecked(true)
	}
	return &fyne.Container{Layout: &unpadTextWidgetLayout{parent: check}, Objects: []fyne.CanvasObject{check}}
}

// Update doesn't need to change a checkbox
func (*CheckBoxSegment) Update(fyne.CanvasObject) {
}

// Select does nothing for a checkbox.
func (*CheckBoxSegment) Select(_, _ fyne.Position) {
}

// SelectedText returns the empty string for a checkbox.
func (*CheckBoxSegment) SelectedText() string {
	return ""
}

// Unselect does nothing for a checkbox.
func (*CheckBoxSegment) Unselect() {
}

// TableSegment represents a table within a rich text widget.
//
// Since: 2.8
type TableSegment struct {
	// Headers holds the cells of the header row, or nil for a header-less table.
	Headers [][]RichTextSegment
	// Rows holds the body rows; each row is a slice of cells, each cell a slice of segments.
	Rows       [][][]RichTextSegment
	Alignments []fyne.TextAlign
}

// Inline returns false as a table is a full-width block element.
func (*TableSegment) Inline() bool {
	return false
}

// Textual returns the table content as tab-separated, newline-delimited text.
func (t *TableSegment) Textual() string {
	var b strings.Builder
	writeRow := func(cells [][]RichTextSegment) {
		for i, cell := range cells {
			if i > 0 {
				b.WriteByte('\t')
			}
			for _, s := range cell {
				b.WriteString(s.Textual())
			}
		}
		b.WriteByte('\n')
	}
	if t.Headers != nil {
		writeRow(t.Headers)
	}
	for _, r := range t.Rows {
		writeRow(r)
	}
	return b.String()
}

func (t *TableSegment) columns() int {
	cols := len(t.Alignments)
	if len(t.Headers) > cols {
		cols = len(t.Headers)
	}
	for _, r := range t.Rows {
		if len(r) > cols {
			cols = len(r)
		}
	}
	return cols
}

func (t *TableSegment) alignFor(col int) fyne.TextAlign {
	if col < len(t.Alignments) {
		return t.Alignments[col]
	}
	return fyne.TextAlignLeading
}

// Visual returns a new grid laying out the table cells.
func (t *TableSegment) Visual() fyne.CanvasObject {
	cols := t.columns()
	if cols == 0 {
		return NewRichText()
	}

	objects := make([]fyne.CanvasObject, 0, cols*(len(t.Rows)+1))
	appendRow := func(cells [][]RichTextSegment, header bool) {
		for c := 0; c < cols; c++ {
			var segs []RichTextSegment
			if c < len(cells) {
				segs = cells[c]
			}
			objects = append(objects, newTableCell(segs, t.alignFor(c), header))
		}
	}
	if t.Headers != nil {
		appendRow(t.Headers, true)
	}
	for _, r := range t.Rows {
		appendRow(r, false)
	}

	grid := &fyne.Container{Layout: &tableSegmentLayout{cols: cols}, Objects: objects}
	border := canvas.NewRectangle(theme.Color(theme.ColorNameInputBorder))
	return widget.NewHScroll(&fyne.Container{Layout: layout.NewStackLayout(), Objects: []fyne.CanvasObject{border, grid}})
}

// Update does nothing; a table visual is rebuilt rather than updated.
func (*TableSegment) Update(fyne.CanvasObject) {
}

// Select does nothing for a table.
func (*TableSegment) Select(_, _ fyne.Position) {
}

// SelectedText returns the table content as text.
func (t *TableSegment) SelectedText() string {
	return t.Textual()
}

// Unselect does nothing for a table.
func (*TableSegment) Unselect() {
}

// newTableCell builds a single table cell: padded rich-text content over a fill,
// so the grid-line colour drawn behind the grid shows through the gaps left by
// tableSegmentLayout.
func newTableCell(segs []RichTextSegment, align fyne.TextAlign, header bool) fyne.CanvasObject {
	fill := theme.Color(theme.ColorNameBackground)
	if header {
		fill = theme.Color(theme.ColorNameHeaderBackground)
	}
	bg := canvas.NewRectangle(fill)

	cell := make([]RichTextSegment, 0, len(segs))
	for _, s := range segs {
		switch seg := s.(type) {
		case *TextSegment:
			seg.Style.Alignment = align
			if header {
				seg.Style.TextStyle.Bold = true
			}
		case *HyperlinkSegment:
			seg.Alignment = align
		}
		cell = append(cell, s)
	}
	if len(cell) == 0 {
		cell = append(cell, &TextSegment{Style: RichTextStyleInline, Text: " "})
	}

	text := NewRichText(cell...)
	text.Wrapping = fyne.TextWrapOff
	padded := &fyne.Container{Layout: layout.NewPaddedLayout(), Objects: []fyne.CanvasObject{text}}
	return &fyne.Container{Layout: layout.NewStackLayout(), Objects: []fyne.CanvasObject{bg, padded}}
}

// tableSegmentLayout arranges cells row-major. Columns are sized to their widest
// cell, any slack width is shared evenly so the table fills the available width,
// and a one-pixel gap is left around each cell so a background drawn behind the
// grid shows through as grid lines.
type tableSegmentLayout struct {
	cols int
}

func (l *tableSegmentLayout) measure(objects []fyne.CanvasObject) (colWidths, rowHeights []float32) {
	rows := (len(objects) + l.cols - 1) / l.cols
	colWidths = make([]float32, l.cols)
	rowHeights = make([]float32, rows)
	for i, o := range objects {
		r, c := i/l.cols, i%l.cols
		m := o.MinSize()
		if m.Width > colWidths[c] {
			colWidths[c] = m.Width
		}
		if m.Height > rowHeights[r] {
			rowHeights[r] = m.Height
		}
	}
	return colWidths, rowHeights
}

func (l *tableSegmentLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	colWidths, rowHeights := l.measure(objects)
	gap := theme.Size(theme.SizeNameSeparatorThickness)
	w := gap
	for _, cw := range colWidths {
		w += cw + gap
	}
	h := gap
	for _, rh := range rowHeights {
		h += rh + gap
	}
	return fyne.NewSize(w, h)
}

func (l *tableSegmentLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	colWidths, rowHeights := l.measure(objects)
	gap := theme.Size(theme.SizeNameSeparatorThickness)

	minWidth := gap
	for _, cw := range colWidths {
		minWidth += cw + gap
	}
	if extra := size.Width - minWidth; extra > 0 && l.cols > 0 {
		share := extra / float32(l.cols)
		for c := range colWidths {
			colWidths[c] += share
		}
	}

	y := gap
	for r, rh := range rowHeights {
		x := gap
		for c := 0; c < l.cols; c++ {
			idx := r*l.cols + c
			if idx >= len(objects) {
				break
			}
			objects[idx].Move(fyne.NewPos(x, y))
			objects[idx].Resize(fyne.NewSize(colWidths[c], rh))
			x += colWidths[c] + gap
		}
		y += rh + gap
	}
}

// RichTextStyle describes the details of a text object inside a RichText widget.
//
// Since: 2.1
type RichTextStyle struct {
	Alignment fyne.TextAlign
	ColorName fyne.ThemeColorName
	Inline    bool
	SizeName  fyne.ThemeSizeName // The theme name of the text size to use, if blank will be the standard text size
	TextStyle fyne.TextStyle
	// Since: 2.8
	QuotingDepth int

	// an internal detail where we obscure password fields
	concealed bool

	// an internal detail marking inline code, which renders on a background fill
	codeInline bool
}

// RichTextSegment describes any element that can be rendered in a RichText widget.
//
// Since: 2.1
type RichTextSegment interface {
	Inline() bool
	Textual() string
	Update(fyne.CanvasObject)
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

	parent *RichText
}

// content returns the text that this segment holds.
func (t *TextSegment) content() string {
	return t.Text
}

// setContent replaces the text that this segment holds.
func (t *TextSegment) setContent(text string) {
	t.Text = text
}

// Inline should return true if this text can be included within other elements, or false if it creates a new block.
func (t *TextSegment) Inline() bool {
	return t.Style.Inline
}

// Textual returns the content of this segment rendered to plain text.
func (t *TextSegment) Textual() string {
	return t.Text
}

// Visual returns a new instance of a graphical element required to render this segment.
func (t *TextSegment) Visual() fyne.CanvasObject {
	text := canvas.NewText(t.Text, t.color())
	if t.Style.codeInline {
		bg := canvas.NewRectangle(theme.ColorForWidget(theme.ColorNameInputBackground, t.parent))
		c := &fyne.Container{Layout: &codeInlineLayout{}, Objects: []fyne.CanvasObject{bg, text}}
		t.Update(c)
		return c
	}

	t.Update(text)
	return text
}

// Update applies the current state of this text segment to an existing visual.
func (t *TextSegment) Update(o fyne.CanvasObject) {
	obj, ok := o.(*canvas.Text)
	if !ok { // inline code container: [background, text]
		c, _ := o.(*fyne.Container)
		bg, _ := c.Objects[0].(*canvas.Rectangle)
		bg.FillColor = theme.ColorForWidget(theme.ColorNameInputBackground, t.parent)
		bg.Refresh()
		obj, _ = c.Objects[1].(*canvas.Text)
	}
	obj.Text = t.Text
	obj.Color = t.color()
	obj.Alignment = t.Style.Alignment
	obj.TextStyle = t.Style.TextStyle
	obj.TextSize = t.size()
	obj.Refresh()
}

// codeInlineLayout keeps the inline-code background tight to the text, so when
// the row layout stretches the container to fill trailing space the fill does
// not stretch with it.
type codeInlineLayout struct{}

func (codeInlineLayout) MinSize(o []fyne.CanvasObject) fyne.Size {
	return o[1].MinSize()
}

func (codeInlineLayout) Layout(o []fyne.CanvasObject, _ fyne.Size) {
	size := o[1].MinSize()
	for _, obj := range o {
		obj.Resize(size)
		obj.Move(fyne.NewPos(0, 0))
	}
}

// Select tells the segment that the user is selecting the content between the two positions.
func (*TextSegment) Select(_, _ fyne.Position) {
	// no-op: this will be added when we progress to editor
}

// SelectedText should return the text representation of any content currently selected through the Select call.
func (*TextSegment) SelectedText() string {
	// no-op: this will be added when we progress to editor
	return ""
}

// Unselect tells the segment that the user is has cancelled the previous selection.
func (*TextSegment) Unselect() {
	// no-op: this will be added when we progress to editor
}

func (t *TextSegment) color() color.Color {
	if t.Style.ColorName != "" {
		return theme.ColorForWidget(t.Style.ColorName, t.parent)
	}

	return theme.ColorForWidget(theme.ColorNameForeground, t.parent)
}

func (t *TextSegment) size() float32 {
	if t.Style.SizeName != "" {
		i := theme.SizeForWidget(t.Style.SizeName, t.parent)
		return i
	}

	i := theme.SizeForWidget(theme.SizeNameText, t.parent)
	return i
}

type richImage struct {
	BaseWidget
	align  fyne.TextAlign
	img    *canvas.Image
	oldMin fyne.Size
	layout *fyne.Container
	min    fyne.Size
}

func newRichImage(u fyne.URI, align fyne.TextAlign) *richImage {
	img := canvas.NewImageFromURI(u)
	img.FillMode = canvas.ImageFillOriginal
	i := &richImage{img: img, align: align}
	i.ExtendBaseWidget(i)
	return i
}

func (r *richImage) CreateRenderer() fyne.WidgetRenderer {
	r.layout = &fyne.Container{Layout: &richImageLayout{r}, Objects: []fyne.CanvasObject{r.img}}
	return NewSimpleRenderer(r.layout)
}

func (r *richImage) MinSize() fyne.Size {
	orig := r.img.MinSize()
	c := fyne.CurrentApp().Driver().CanvasForObject(r)
	if c == nil {
		return r.oldMin // not yet rendered
	}

	// unscale the image so it is not varying based on canvas
	w := scale.ToScreenCoordinate(c, orig.Width)
	h := scale.ToScreenCoordinate(c, orig.Height)
	// we return size / 2 as this assumes a HiDPI / 2x image scaling
	r.min = fyne.NewSize(float32(w)/2, float32(h)/2)
	return r.min
}

func (r *richImage) setAlign(a fyne.TextAlign) {
	if r.layout != nil {
		r.layout.Refresh()
	}
	r.align = a
}

type richImageLayout struct {
	r *richImage
}

func (r *richImageLayout) Layout(_ []fyne.CanvasObject, s fyne.Size) {
	r.r.img.Resize(r.r.min)
	gap := float32(0)

	switch r.r.align {
	case fyne.TextAlignCenter:
		gap = (s.Width - r.r.min.Width) / 2
	case fyne.TextAlignTrailing:
		gap = s.Width - r.r.min.Width
	}

	r.r.img.Move(fyne.NewPos(gap, 0))
}

func (r *richImageLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return r.r.min
}

type unpadTextWidgetLayout struct {
	parent fyne.Widget
}

func (u *unpadTextWidgetLayout) Layout(o []fyne.CanvasObject, s fyne.Size) {
	innerPad := theme.SizeForWidget(theme.SizeNameInnerPadding, u.parent)
	pad := innerPad * -1
	pad2 := pad * -2

	o[0].Move(fyne.NewPos(pad, pad))
	o[0].Resize(s.Add(fyne.NewSize(pad2, pad2)))
}

func (u *unpadTextWidgetLayout) MinSize(o []fyne.CanvasObject) fyne.Size {
	innerPad := theme.SizeForWidget(theme.SizeNameInnerPadding, u.parent)
	pad := innerPad * 2
	return o[0].MinSize().Subtract(fyne.NewSize(pad, pad))
}

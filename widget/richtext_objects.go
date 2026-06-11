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
func (h *HyperlinkSegment) Inline() bool {
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
	link := o.(*fyne.Container).Objects[0].(*Hyperlink)
	link.URL = h.URL
	link.Alignment = h.Alignment
	link.SizeName = h.SizeName
	link.TextStyle = h.TextStyle
	link.OnTapped = h.OnTapped
	link.Refresh()
}

// Select tells the segment that the user is selecting the content between the two positions.
func (h *HyperlinkSegment) Select(begin, end fyne.Position) {
	// no-op: this will be added when we progress to editor
}

// SelectedText should return the text representation of any content currently selected through the Select call.
func (h *HyperlinkSegment) SelectedText() string {
	// no-op: this will be added when we progress to editor
	return ""
}

// Unselect tells the segment that the user is has cancelled the previous selection.
func (h *HyperlinkSegment) Unselect() {
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
func (i *ImageSegment) Inline() bool {
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
	img := o.(*richImage)

	// one of the following will be used
	img.img.File = newer.File
	img.img.Resource = newer.Resource
	img.setAlign(i.Alignment)

	img.Refresh()
}

// Select tells the segment that the user is selecting the content between the two positions.
func (i *ImageSegment) Select(begin, end fyne.Position) {
	// no-op: this will be added when we progress to editor
}

// SelectedText should return the text representation of any content currently selected through the Select call.
func (i *ImageSegment) SelectedText() string {
	// no-op: images have no text rendering
	return ""
}

// Unselect tells the segment that the user is has cancelled the previous selection.
func (i *ImageSegment) Unselect() {
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
func (l *ListSegment) Inline() bool {
	return false
}

// Segments returns the segments required to draw bullets before each item
func (l *ListSegment) Segments() []RichTextSegment {
	out := make([]RichTextSegment, len(l.Items))
	j := l.StartNumber()
	for i, in := range l.Items {
		var texts []RichTextSegment
		if _, ok := in.(*ListSegment); !ok {
			txt := "• "
			if l.Ordered {
				txt = strconv.Itoa(j) + "."
				j++
			}
			indentation := strings.Repeat(" ", l.indentationLevel*4)
			style := RichTextStyleStrong
			style.QuotingDepth = l.quotingLevel
			bullet := &TextSegment{Text: indentation + txt + " ", Style: style}
			texts = append(texts, bullet)
			if _, ok := in.(*ParagraphSegment); !ok {
				in = &ParagraphSegment{Texts: []RichTextSegment{in}}
			}
		}
		texts = append(texts, in)
		out[i] = &ParagraphSegment{Texts: texts}
	}
	return out
}

// Textual returns no content for a list as the content is in sub-segments.
func (l *ListSegment) Textual() string {
	return ""
}

// Visual returns no additional elements for this segment.
func (l *ListSegment) Visual() fyne.CanvasObject {
	return nil
}

// Update doesn't need to change a list visual.
func (l *ListSegment) Update(fyne.CanvasObject) {
}

// Select does nothing for a list container.
func (l *ListSegment) Select(_, _ fyne.Position) {
}

// SelectedText returns the empty string for this list.
func (l *ListSegment) SelectedText() string {
	return ""
}

// Unselect does nothing for a list container.
func (l *ListSegment) Unselect() {
}

// ParagraphSegment wraps a number of text elements in a paragraph.
// It is similar to using a list of text elements when the final style is RichTextStyleParagraph.
//
// Since: 2.1
type ParagraphSegment struct {
	Texts []RichTextSegment
}

// Inline returns false as a paragraph should be in a block.
func (p *ParagraphSegment) Inline() bool {
	return false
}

// Segments returns the list of text elements in this paragraph.
func (p *ParagraphSegment) Segments() []RichTextSegment {
	return p.Texts
}

// Textual returns no content for a paragraph container.
func (p *ParagraphSegment) Textual() string {
	return ""
}

// Visual returns the no extra elements.
func (p *ParagraphSegment) Visual() fyne.CanvasObject {
	return nil
}

// Update doesn't need to change a paragraph container.
func (p *ParagraphSegment) Update(fyne.CanvasObject) {
}

// Select does nothing for a paragraph container.
func (p *ParagraphSegment) Select(_, _ fyne.Position) {
}

// SelectedText returns the empty string for this paragraph container.
func (p *ParagraphSegment) SelectedText() string {
	return ""
}

// Unselect does nothing for a paragraph container.
func (p *ParagraphSegment) Unselect() {
}

// SeparatorSegment includes a horizontal separator in a rich text widget.
//
// Since: 2.1
type SeparatorSegment struct {
	_ bool // Without this a pointer to SeparatorSegment will always be the same.
}

// Inline returns false as a separator should be full width.
func (s *SeparatorSegment) Inline() bool {
	return false
}

// Textual returns no content for a separator element.
func (s *SeparatorSegment) Textual() string {
	return ""
}

// Visual returns a new instance of a separator widget for this segment.
func (s *SeparatorSegment) Visual() fyne.CanvasObject {
	return NewSeparator()
}

// Update doesn't need to change a separator visual.
func (s *SeparatorSegment) Update(fyne.CanvasObject) {
}

// Select does nothing for a separator.
func (s *SeparatorSegment) Select(_, _ fyne.Position) {
}

// SelectedText returns the empty string for this separator.
func (s *SeparatorSegment) SelectedText() string {
	return "" // TODO maybe return "---\n"?
}

// Unselect does nothing for a separator.
func (s *SeparatorSegment) Unselect() {
}

// CodeBlockSegment represents a fenced or indented code block. It renders its
// content as monospace text on a panel, so the block stands apart from the
// surrounding prose.
//
// Since: 2.8
type CodeBlockSegment struct {
	Text         string
	quotingLevel int
}

// Inline returns false as a code block is a full-width block element.
func (c *CodeBlockSegment) Inline() bool {
	return false
}

// Textual returns the raw content of this code block.
func (c *CodeBlockSegment) Textual() string {
	return c.Text
}

// Visual returns a new panel widget rendering this code block.
func (c *CodeBlockSegment) Visual() fyne.CanvasObject {
	return newRichCodeBlock(c.Text)
}

// Update applies the current content of this segment to an existing visual.
func (c *CodeBlockSegment) Update(o fyne.CanvasObject) {
	o.(*richCodeBlock).setText(c.Text)
}

// Select does nothing for a code block.
func (c *CodeBlockSegment) Select(_, _ fyne.Position) {
}

// SelectedText returns the code block content.
func (c *CodeBlockSegment) SelectedText() string {
	return c.Text
}

// Unselect does nothing for a code block.
func (c *CodeBlockSegment) Unselect() {
}

// richCodeBlock is the internal widget that draws a code block: monospace text
// on a rounded, bordered panel.
type richCodeBlock struct {
	BaseWidget
	text  string
	bg    *canvas.Rectangle
	label *Label
}

func newRichCodeBlock(text string) *richCodeBlock {
	c := &richCodeBlock{text: text}
	c.ExtendBaseWidget(c)
	return c
}

func (c *richCodeBlock) setText(text string) {
	c.text = text
	if c.label != nil {
		c.label.SetText(text)
	}
}

func (c *richCodeBlock) CreateRenderer() fyne.WidgetRenderer {
	c.bg = canvas.NewRectangle(theme.Color(theme.ColorNameInputBackground))
	c.bg.StrokeColor = theme.Color(theme.ColorNameInputBorder)
	c.bg.StrokeWidth = 1
	c.bg.CornerRadius = theme.Size(theme.SizeNameInputRadius)
	c.label = NewLabelWithStyle(c.text, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	scroll := widget.NewHScroll(c.label)
	cont := &fyne.Container{Layout: &richCodeBlockLayout{}, Objects: []fyne.CanvasObject{c.bg, scroll}}
	return NewSimpleRenderer(cont)
}

type richCodeBlockLayout struct{}

func (l *richCodeBlockLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return objects[1].MinSize()
}

func (l *richCodeBlockLayout) Layout(objects []fyne.CanvasObject, s fyne.Size) {
	for _, o := range objects {
		o.Move(fyne.NewPos(0, 0))
		o.Resize(s)
	}
}

// CheckBoxSegment represents checkbox (with text) in a rich text widget.
//
// Since: 2.8
type CheckBoxSegment struct {
	Checked bool
	Text    string
}

// Inline returns true as a CheckBoxSegment is usually part of a list item.
func (c *CheckBoxSegment) Inline() bool {
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
func (c *CheckBoxSegment) Update(fyne.CanvasObject) {
}

// Select does nothing for a checkbox.
func (c *CheckBoxSegment) Select(_, _ fyne.Position) {
}

// SelectedText returns the empty string for a checkbox.
func (c *CheckBoxSegment) SelectedText() string {
	return ""
}

// Unselect does nothing for a checkbox.
func (c *CheckBoxSegment) Unselect() {
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
func (t *TableSegment) Inline() bool {
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
	return &fyne.Container{Layout: layout.NewStackLayout(), Objects: []fyne.CanvasObject{border, grid}}
}

// Update does nothing; a table visual is rebuilt rather than updated.
func (t *TableSegment) Update(fyne.CanvasObject) {
}

// Select does nothing for a table.
func (t *TableSegment) Select(_, _ fyne.Position) {
}

// SelectedText returns the table content as text.
func (t *TableSegment) SelectedText() string {
	return t.Textual()
}

// Unselect does nothing for a table.
func (t *TableSegment) Unselect() {
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
		c := o.(*fyne.Container)
		bg := c.Objects[0].(*canvas.Rectangle)
		bg.FillColor = theme.ColorForWidget(theme.ColorNameInputBackground, t.parent)
		bg.Refresh()
		obj = c.Objects[1].(*canvas.Text)
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

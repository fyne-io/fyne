package widget

import (
	"image/color"
	"net/url"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/scale"
	"fyne.io/fyne/v2/theme"
)

var (
	// RichTextStyleBlockquote represents a quote presented in an indented block.
	//
	// Since: 2.1
	RichTextStyleBlockquote = RichTextStyle{
		ColorName: theme.ColorNameForeground,
		Inline:    false,
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
		ColorName: theme.ColorNameForeground,
		Inline:    true,
		SizeName:  theme.SizeNameText,
		TextStyle: fyne.TextStyle{Monospace: true},
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
		Inline:    false,
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
	OnTapped func()
}

// Inline returns true as hyperlinks are inside other elements.
func (h *HyperlinkSegment) Inline() bool {
	return true
}

// Textual returns the content of this segment rendered to plain text.
func (h *HyperlinkSegment) Textual() string {
	return h.Text
}

// Visual returns the hyperlink widget required to render this segment.
func (h *HyperlinkSegment) Visual() fyne.CanvasObject {
	link := NewHyperlink(h.Text, h.URL)
	link.Alignment = h.Alignment
	link.OnTapped = h.OnTapped
	return &fyne.Container{Layout: &unpadTextWidgetLayout{}, Objects: []fyne.CanvasObject{link}}
}

// Update applies the current state of this hyperlink segment to an existing visual.
func (h *HyperlinkSegment) Update(o fyne.CanvasObject) {
	link := o.(*fyne.Container).Objects[0].(*Hyperlink)
	link.Text = h.Text
	link.URL = h.URL
	link.Alignment = h.Alignment
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

// Visual returns the image widget required to render this segment.
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
}

// Inline returns false as a list should be in a block.
func (l *ListSegment) Inline() bool {
	return false
}

// Segments returns the segments required to draw bullets before each item
func (l *ListSegment) Segments() []RichTextSegment {
	out := make([]RichTextSegment, len(l.Items))
	for i, in := range l.Items {
		txt := "• "
		if l.Ordered {
			txt = strconv.Itoa(i+1) + "."
		}
		bullet := &TextSegment{Text: txt + " ", Style: RichTextStyleStrong}
		out[i] = &ParagraphSegment{Texts: []RichTextSegment{
			bullet,
			in,
		}}
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

// Update doesnt need to change a list visual.
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

// Update doesnt need to change a paragraph container.
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

// Visual returns the separator element for this segment.
func (s *SeparatorSegment) Visual() fyne.CanvasObject {
	return NewSeparator()
}

// Update doesnt need to change a separator visual.
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

	t.Update(obj)
	return obj
}

// Update applies the current state of this text segment to an existing visual.
func (t *TextSegment) Update(o fyne.CanvasObject) {
	obj := o.(*canvas.Text)
	obj.Text = t.Text
	obj.Color = t.color()
	obj.Alignment = t.Style.Alignment
	obj.TextStyle = t.Style.TextStyle
	obj.TextSize = t.size()
	obj.Refresh()
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
		return fyne.CurrentApp().Settings().Theme().Color(t.Style.ColorName, fyne.CurrentApp().Settings().ThemeVariant())
	}

	return theme.ForegroundColor()
}

func (t *TextSegment) size() float32 {
	if t.Style.SizeName != "" {
		return fyne.CurrentApp().Settings().Theme().Size(t.Style.SizeName)
	}

	return theme.TextSize()
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
}

func (u *unpadTextWidgetLayout) Layout(o []fyne.CanvasObject, s fyne.Size) {
	pad := theme.InnerPadding() * -1
	pad2 := pad * -2

	o[0].Move(fyne.NewPos(pad, pad))
	o[0].Resize(s.Add(fyne.NewSize(pad2, pad2)))
}

func (u *unpadTextWidgetLayout) MinSize(o []fyne.CanvasObject) fyne.Size {
	pad := theme.InnerPadding() * 2
	return o[0].MinSize().Subtract(fyne.NewSize(pad, pad))
}

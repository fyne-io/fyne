package widget

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// HyperlinkSegment represents a hyperlink within a rich text widget.
//
// Since: 2.1
type HyperlinkSegment struct {
	Text string
	URL  *url.URL
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
	return &fyne.Container{Layout: &unpadTextWidgetLayout{}, Objects: []fyne.CanvasObject{NewHyperlink(h.Text, h.URL)}}
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

// SeparatorSegment includes a horizontal separator in a rich text widget.
//
// Since: 2.1
type SeparatorSegment struct {
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

// Select does nothing for a separator.
func (s *SeparatorSegment) Select(begin, end fyne.Position) {
}

// SelectedText returns the empty string for this separator.
func (s *SeparatorSegment) SelectedText() string {
	return "" // TODO maybe return "---\n"?
}

// Unselect does nothing for a separator.
func (s *SeparatorSegment) Unselect() {
}

type unpadTextWidgetLayout struct {
}

func (u *unpadTextWidgetLayout) Layout(o []fyne.CanvasObject, s fyne.Size) {
	pad2 := theme.Padding() * -2
	pad4 := pad2 * -2

	o[0].Move(fyne.NewPos(pad2, pad2))
	o[0].Resize(s.Add(fyne.NewSize(pad4, pad4)))
}

func (u *unpadTextWidgetLayout) MinSize(o []fyne.CanvasObject) fyne.Size {
	pad4 := theme.Padding() * 4
	return o[0].MinSize().Subtract(fyne.NewSize(pad4, pad4))
}

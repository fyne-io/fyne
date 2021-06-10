package widget

import (
	"fyne.io/fyne/v2"
)

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

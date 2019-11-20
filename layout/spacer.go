package layout

import "fyne.io/fyne"

// SpacerObject is any object that can be used to space out child objects
type SpacerObject interface {
	ExpandVertical() bool
	ExpandHorizontal() bool
}

// Spacer is any simple object that can be used in a box layout to space
// out child objects
type Spacer struct {
	FixHorizontal bool
	FixVertical   bool

	size   fyne.Size
	pos    fyne.Position
	hidden bool
}

// ExpandVertical returns whether or not this spacer expands on the vertical axis
func (s *Spacer) ExpandVertical() bool {
	return !s.FixVertical
}

// ExpandHorizontal returns whether or not this spacer expands on the horizontal axis
func (s *Spacer) ExpandHorizontal() bool {
	return !s.FixHorizontal
}

// Size returns the current size of this Spacer
func (s *Spacer) Size() fyne.Size {
	return s.size
}

// Resize sets a new size for the Spacer - this will be called by the layout
func (s *Spacer) Resize(size fyne.Size) {
	s.size = size
}

// Position returns the current position of this Spacer
func (s *Spacer) Position() fyne.Position {
	return s.pos
}

// Move sets a new position for the Spacer - this will be called by the layout
func (s *Spacer) Move(pos fyne.Position) {
	s.pos = pos
}

// MinSize returns a 0 size as a Spacer can shrink to no actual size
func (s *Spacer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

// Visible returns true if this spacer should affect the layout
func (s *Spacer) Visible() bool {
	return !s.hidden
}

// Show sets the Spacer to be part of the layout calculations
func (s *Spacer) Show() {
	s.hidden = false
}

// Hide removes this Spacer from layout calculations
func (s *Spacer) Hide() {
	s.hidden = true
}

// Refresh does nothing for a spacer but is part of the CanvasObject definition
func (s *Spacer) Refresh() {
}

// NewSpacer returns a spacer object which can fill vertical and horizontal
// space. This is primarily used with a box layout.
func NewSpacer() fyne.CanvasObject {
	return &Spacer{}
}

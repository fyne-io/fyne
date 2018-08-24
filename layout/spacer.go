package layout

import "github.com/fyne-io/fyne"

// SpacerObject is a simple object that can be used in a list layout to space
// out child objects
type SpacerObject interface {
	ExpandVertical() bool
	ExpandHorizontal() bool
}

type spacerObject struct {
	size fyne.Size
	pos  fyne.Position
}

func (s *spacerObject) ExpandVertical() bool {
	return true
}

func (s *spacerObject) ExpandHorizontal() bool {
	return true
}

func (s *spacerObject) CurrentSize() fyne.Size {
	return s.size
}

func (s *spacerObject) Resize(size fyne.Size) {
	s.size = size
}

func (s *spacerObject) CurrentPosition() fyne.Position {
	return s.pos
}

func (s *spacerObject) Move(pos fyne.Position) {
	s.pos = pos
}

func (s *spacerObject) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

// NewSpacer returns a spacer object which can fill vertical and horizontal
// space. This is primarily used with a list layout.
func NewSpacer() fyne.CanvasObject {
	return &spacerObject{}
}

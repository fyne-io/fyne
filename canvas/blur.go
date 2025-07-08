package canvas

import "fyne.io/fyne/v2"

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Blur)(nil)

// Blur creates a rectangular blur region on the output.
// All objects drawn under this will be blurred, any above will not be affected.
//
// Since: 2.7
type Blur struct {
	baseObject

	// Radius refers to how far from a pixel should be used to calculate the blur.
	// It must be greater than 0 but no more than 50.
	Radius float32
}

// Hide will set this blur to not be visible
func (b *Blur) Hide() {
	b.baseObject.Hide()

	repaint(b)
}

// Move the blur to a new position, relative to its parent / canvas
func (b *Blur) Move(pos fyne.Position) {
	b.baseObject.Move(pos)

	repaint(b)
}

// Refresh causes this blur to be redrawn with its configured state.
func (b *Blur) Refresh() {
	Refresh(b)
}

// Resize on a blur updates the new size of this object.
func (b *Blur) Resize(s fyne.Size) {
	if s == b.Size() {
		return
	}

	b.baseObject.Resize(s)
}

// NewBlur returns a new Blur instance
func NewBlur(radius float32) *Blur {
	return &Blur{
		Radius: radius,
	}
}

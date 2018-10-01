// Package canvas contains all of the primitive CanvasObjects that make up a Fyne GUI
package canvas

import "github.com/fyne-io/fyne"

type baseObject struct {
	Size     fyne.Size     // The current size of the Rectangle
	Position fyne.Position // The current position of the Rectangle
	Options  Options       // Options to pass to the renderer

	min fyne.Size // The minimum size this object can be
}

// CurrentSize returns the current size of this rectangle object
func (r *baseObject) CurrentSize() fyne.Size {
	return r.Size
}

// Resize sets a new size for the rectangle object
func (r *baseObject) Resize(size fyne.Size) {
	r.Size = size

	fyne.RefreshObject(r)
}

// CurrentPosition gets the current position of this rectangle object, relative to it's parent / canvas
func (r *baseObject) CurrentPosition() fyne.Position {
	return r.Position
}

// Move the rectangle object to a new position, relative to it's parent / canvas
func (r *baseObject) Move(pos fyne.Position) {
	r.Position = pos

	fyne.RefreshObject(r)
}

// MinSize returns the specified minimum size, if set, or {0, 0} otherwise
func (r *baseObject) MinSize() fyne.Size {
	return r.min
}

// SetMinSize specifies the smallest size this object should be
func (r *baseObject) SetMinSize(size fyne.Size) {
	r.min = size

	fyne.RefreshObject(r)
}

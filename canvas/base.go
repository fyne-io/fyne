// Package canvas contains all of the primitive CanvasObjects that make up a Fyne GUI
package canvas // import "fyne.io/fyne/canvas"

import "fyne.io/fyne"

type baseObject struct {
	size     fyne.Size     // The current size of the Rectangle
	position fyne.Position // The current position of the Rectangle
	Hidden   bool          // Is this object currently hidden

	min fyne.Size // The minimum size this object can be
}

// CurrentSize returns the current size of this rectangle object
func (r *baseObject) Size() fyne.Size {
	return r.size
}

// Resize sets a new size for the rectangle object
func (r *baseObject) Resize(size fyne.Size) {
	r.size = size
}

// CurrentPosition gets the current position of this rectangle object, relative to its parent / canvas
func (r *baseObject) Position() fyne.Position {
	return r.position
}

// Move the rectangle object to a new position, relative to its parent / canvas
func (r *baseObject) Move(pos fyne.Position) {
	r.position = pos
}

// MinSize returns the specified minimum size, if set, or {1, 1} otherwise
func (r *baseObject) MinSize() fyne.Size {
	if r.min.Width == 0 && r.min.Height == 0 {
		return fyne.NewSize(1, 1)
	}

	return r.min
}

// SetMinSize specifies the smallest size this object should be
func (r *baseObject) SetMinSize(size fyne.Size) {
	r.min = size
}

// IsVisible returns true if this object is visible, false otherwise
func (r *baseObject) Visible() bool {
	return !r.Hidden
}

// Show will set this object to be visible
func (r *baseObject) Show() {
	r.Hidden = false
}

// Hide will set this object to not be visible
func (r *baseObject) Hide() {
	r.Hidden = true
}

// Refresh instructs the containing canvas to refresh the specified obj.
func Refresh(obj fyne.CanvasObject) {
	if fyne.CurrentApp() == nil || fyne.CurrentApp().Driver() == nil {
		return
	}

	c := fyne.CurrentApp().Driver().CanvasForObject(obj)
	if c != nil {
		c.Refresh(obj)
	}
}

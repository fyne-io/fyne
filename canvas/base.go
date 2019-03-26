// Package canvas contains all of the primitive CanvasObjects that make up a Fyne GUI
package canvas // import "fyne.io/fyne/canvas"

import (
	"sync"

	"fyne.io/fyne"
)

type baseObject struct {
	size     fyne.Size     // The current size of the Rectangle
	position fyne.Position // The current position of the Rectangle
	Hidden   bool          // Is this object currently hidden

	min     fyne.Size    // The minimum size this object can be
	objLock sync.RWMutex // for atomicity
}

// CurrentSize returns the current size of this rectangle object
func (r *baseObject) Size() fyne.Size {
	r.objLock.RLock()
	defer r.objLock.RUnlock()
	return r.size
}

// Resize sets a new size for the rectangle object
func (r *baseObject) Resize(size fyne.Size) {
	r.objLock.Lock()
	r.size = size
	r.objLock.Unlock()
}

// CurrentPosition gets the current position of this rectangle object, relative to it's parent / canvas
func (r *baseObject) Position() fyne.Position {
	r.objLock.RLock()
	defer r.objLock.RUnlock()
	return r.position
}

// Move the rectangle object to a new position, relative to it's parent / canvas
func (r *baseObject) Move(pos fyne.Position) {
	r.objLock.Lock()
	r.position = pos
	r.objLock.Unlock()
}

// MinSize returns the specified minimum size, if set, or {1, 1} otherwise
func (r *baseObject) MinSize() fyne.Size {
	r.objLock.RLock()
	defer r.objLock.RUnlock()

	if r.min.Width == 0 && r.min.Height == 0 {
		return fyne.NewSize(1, 1)
	}

	return r.min
}

// SetMinSize specifies the smallest size this object should be
func (r *baseObject) SetMinSize(size fyne.Size) {
	r.objLock.Lock()
	r.min = size
	r.objLock.Unlock()
}

// IsVisible returns true if this object is visible, false otherwise
func (r *baseObject) Visible() bool {
	r.objLock.RLock()
	defer r.objLock.RUnlock()
	return !r.Hidden
}

// Show will set this object to be visible
func (r *baseObject) Show() {
	r.objLock.Lock()
	r.Hidden = false
	r.objLock.Unlock()
}

// Hide will set this object to not be visible
func (r *baseObject) Hide() {
	r.objLock.Lock()
	r.Hidden = true
	r.objLock.Unlock()
}

// Refresh instructs the containing canvas to refresh the specified obj.
func Refresh(obj fyne.CanvasObject) {
	c := fyne.CurrentApp().Driver().CanvasForObject(obj)
	if c != nil {
		c.Refresh(obj)
	}
}

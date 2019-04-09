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

	min fyne.Size // The minimum size this object can be

	lck sync.RWMutex
}

// CurrentSize returns the current size of this rectangle object
func (r *baseObject) Size() fyne.Size {
	r.lck.RLock()
	defer r.lck.RUnlock()
	return r.size
}

// Resize sets a new size for the rectangle object
func (r *baseObject) Resize(size fyne.Size) {
	r.lck.Lock()
	defer r.lck.Unlock()
	r.size = size
}

// CurrentPosition gets the current position of this rectangle object, relative to it's parent / canvas
func (r *baseObject) Position() fyne.Position {
	r.lck.RLock()
	defer r.lck.RUnlock()
	return r.position
}

// Move the rectangle object to a new position, relative to it's parent / canvas
func (r *baseObject) Move(pos fyne.Position) {
	r.lck.Lock()
	defer r.lck.Unlock()
	r.position = pos
}

// MinSize returns the specified minimum size, if set, or {1, 1} otherwise
func (r *baseObject) MinSize() fyne.Size {
	r.lck.RLock()
	defer r.lck.RUnlock()
	if r.min.Width == 0 && r.min.Height == 0 {
		return fyne.NewSize(1, 1)
	}

	return r.min
}

// SetMinSize specifies the smallest size this object should be
func (r *baseObject) SetMinSize(size fyne.Size) {
	r.lck.Lock()
	defer r.lck.Unlock()
	r.min = size
}

// IsVisible returns true if this object is visible, false otherwise
func (r *baseObject) Visible() bool {
	r.lck.Lock()
	defer r.lck.Unlock()
	return !r.Hidden
}

// Show will set this object to be visible
func (r *baseObject) Show() {
	r.lck.Lock()
	defer r.lck.Unlock()
	r.Hidden = false
}

// Hide will set this object to not be visible
func (r *baseObject) Hide() {
	r.lck.Lock()
	defer r.lck.Unlock()
	r.Hidden = true
}

// Refresh instructs the containing canvas to refresh the specified obj.
func Refresh(obj fyne.CanvasObject) {
	c := fyne.CurrentApp().Driver().CanvasForObject(obj)
	if c != nil {
		c.Refresh(obj)
	}
}

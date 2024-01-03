// Package canvas contains all of the primitive CanvasObjects that make up a Fyne GUI.
//
// The types implemented in this package are used as building blocks in order
// to build higher order functionality. These types are designed to be
// non-interactive, by design. If additional functionality is required,
// it's usually a sign that this type should be used as part of a custom
// widget.
package canvas // import "fyne.io/fyne/v2/canvas"

import (
	"math"
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2"
)

type baseObject struct {
	size     atomic.Uint64 // The current size of the canvas object
	position atomic.Uint64 // The current position of the object
	Hidden   bool          // Is this object currently hidden

	min atomic.Uint64 // The minimum size this object can be

	propertyLock sync.RWMutex
}

// Hide will set this object to not be visible.
func (o *baseObject) Hide() {
	o.propertyLock.Lock()
	defer o.propertyLock.Unlock()

	o.Hidden = true
}

// MinSize returns the specified minimum size, if set, or {1, 1} otherwise.
func (o *baseObject) MinSize() fyne.Size {
	min := o.min.Load()
	if min == 0 {
		return fyne.Size{Width: 1, Height: 1}
	}

	return fyne.NewSize(twoFloat32FromUint64(min))
}

// Move the object to a new position, relative to its parent.
func (o *baseObject) Move(pos fyne.Position) {
	o.position.Store(uint64fromTwoFloat32(pos.X, pos.Y))
}

// Position gets the current position of this canvas object, relative to its parent.
func (o *baseObject) Position() fyne.Position {
	pos := o.position.Load()
	if pos == 0 {
		return fyne.Position{}
	}

	return fyne.NewPos(twoFloat32FromUint64(pos))
}

// Resize sets a new size for the canvas object.
func (o *baseObject) Resize(size fyne.Size) {
	o.size.Store(uint64fromTwoFloat32(size.Width, size.Height))
}

// SetMinSize specifies the smallest size this object should be.
func (o *baseObject) SetMinSize(size fyne.Size) {
	o.min.Store(uint64fromTwoFloat32(size.Width, size.Height))
}

// Show will set this object to be visible.
func (o *baseObject) Show() {
	o.propertyLock.Lock()
	defer o.propertyLock.Unlock()

	o.Hidden = false
}

// Size returns the current size of this canvas object.
func (o *baseObject) Size() fyne.Size {
	return fyne.NewSize(twoFloat32FromUint64(o.size.Load()))
}

// Visible returns true if this object is visible, false otherwise.
func (o *baseObject) Visible() bool {
	o.propertyLock.RLock()
	defer o.propertyLock.RUnlock()

	return !o.Hidden
}

func uint64fromTwoFloat32(a, b float32) uint64 {
	x := uint64(math.Float32bits(a))
	y := uint64(math.Float32bits(b))
	return (y << 32) | x
}

func twoFloat32FromUint64(combined uint64) (float32, float32) {
	x := uint32(combined & 0x00000000FFFFFFFF)
	y := uint32(combined & 0xFFFFFFFF00000000 >> 32)
	return math.Float32frombits(x), math.Float32frombits(y)
}

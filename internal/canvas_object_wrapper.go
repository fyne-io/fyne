package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

var _ fyne.CanvasObject = (*SimpleCanvasObjectWrapper)(nil)
var _ CanvasObjectWrapper = (*SimpleCanvasObjectWrapper)(nil)

// CanvasObjectWrapper defines an interface for objects wrapping a fyne.CanvasObject.
type CanvasObjectWrapper interface {
	WrappedObject() fyne.CanvasObject
}

// SimpleCanvasObjectWrapper is a simple CanvasObjectWrapper which can be easily extended.
type SimpleCanvasObjectWrapper struct {
	O fyne.CanvasObject

	pos fyne.Position
}

// Hide hides the wrapped object.
//
// Implements: fyne.CanvasObject
func (w *SimpleCanvasObjectWrapper) Hide() {
	w.O.Hide()
}

// MinSize returns the min size of the wrapped object.
//
// Implements: fyne.CanvasObject
func (w *SimpleCanvasObjectWrapper) MinSize() fyne.Size {
	return w.O.MinSize()
}

// Move moves the wrapped object.
//
// Implements: fyne.CanvasObject
func (w *SimpleCanvasObjectWrapper) Move(position fyne.Position) {
	w.pos = position
	w.Refresh()
}

// Position returns the position of the wrapped object.
//
// Implements: fyne.CanvasObject
func (w *SimpleCanvasObjectWrapper) Position() fyne.Position {
	return w.pos
}

// Refresh refreshes the wrapped object.
//
// Implements: fyne.CanvasObject
func (w *SimpleCanvasObjectWrapper) Refresh() {
	w.O.Refresh()
	canvas.Refresh(w)
}

// Resize resizes the wrapped object.
//
// Implements: fyne.CanvasObject
func (w *SimpleCanvasObjectWrapper) Resize(size fyne.Size) {
	w.O.Resize(size)
}

// Show shows the wrapped object.
//
// Implements: fyne.CanvasObject
func (w *SimpleCanvasObjectWrapper) Show() {
	w.O.Show()
}

// Size returns the size of the wrapped object.
//
// Implements: fyne.CanvasObject
func (w *SimpleCanvasObjectWrapper) Size() fyne.Size {
	return w.O.Size()
}

// Visible returns whether the wrapped object is visible or not.
//
// Implements: fyne.CanvasObject
func (w *SimpleCanvasObjectWrapper) Visible() bool {
	return w.O.Visible()
}

// WrappedObject returns the wrapped object.
//
// Implements: CanvasObjectWrapper
func (w *SimpleCanvasObjectWrapper) WrappedObject() fyne.CanvasObject {
	return w.O
}

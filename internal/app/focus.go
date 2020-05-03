package app

import (
	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver"
)

// FocusManager represents a standard manager of input focus for a canvas
type FocusManager struct {
	canvas fyne.Canvas
}

func (f *FocusManager) nextInChain(current fyne.Focusable) fyne.Focusable {
	var first, next fyne.Focusable
	found := current == nil // if we have no starting point then pretend we matched already
	driver.WalkVisibleObjectTree(f.canvas.Content(), func(obj fyne.CanvasObject, _ fyne.Position, _ fyne.Position, _ fyne.Size) bool {
		if w, ok := obj.(fyne.Disableable); ok && w.Disabled() {
			// disabled widget cannot receive focus
			return false
		}

		focus, ok := obj.(fyne.Focusable)
		if !ok {
			return false
		}

		if found {
			next = focus
			return true
		}

		if !found && obj == current.(fyne.CanvasObject) {
			found = true
		}
		if first == nil {
			first = focus
		}

		return false
	}, nil)

	if next != nil {
		return next
	}
	return first
}

func (f *FocusManager) previousInChain(current fyne.Focusable) fyne.Focusable {
	var last, previous fyne.Focusable
	found := false
	driver.WalkVisibleObjectTree(f.canvas.Content(), func(obj fyne.CanvasObject, _ fyne.Position, _ fyne.Position, _ fyne.Size) bool {
		if w, ok := obj.(fyne.Disableable); ok && w.Disabled() {
			// disabled widget cannot receive focus
			return false
		}

		focus, ok := obj.(fyne.Focusable)
		if !ok {
			return false
		}

		if current != nil && obj == current.(fyne.CanvasObject) {
			found = true
		}
		last = focus

		if !found {
			previous = focus
		}

		return false // we cannot exit early - until we make a reverse tree walk...
	}, nil)

	if previous != nil {
		return previous
	}
	return last
}

// FocusNext will find the item after the current that can be focused and focus it.
// If current is nil then the first focusable item in the canvas will be focused.
func (f *FocusManager) FocusNext(current fyne.Focusable) {
	f.canvas.Focus(f.nextInChain(current))
}

// FocusPrevious will find the item before the current that can be focused and focus it.
// If current is nil then the last focusable item in the canvas will be focused.
func (f *FocusManager) FocusPrevious(current fyne.Focusable) {
	f.canvas.Focus(f.previousInChain(current))
}

// NewFocusManager returns a new instance of the standard focus manager for a canvas.
func NewFocusManager(c fyne.Canvas) *FocusManager {
	return &FocusManager{canvas: c}
}

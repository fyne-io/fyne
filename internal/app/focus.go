package app

import (
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver"
)

// FocusManager represents a standard manager of input focus for a canvas
type FocusManager struct {
	sync.RWMutex

	content fyne.CanvasObject
	focused fyne.Focusable
}

// NewFocusManager returns a new instance of the standard focus manager for a canvas.
func NewFocusManager(c fyne.CanvasObject) *FocusManager {
	return &FocusManager{content: c}
}

// Focus focuses the given obj.
func (f *FocusManager) Focus(obj fyne.Focusable) {
	f.Lock()
	defer f.Unlock()
	f.focus(obj)
}

// Focused returns the currently focused object or nil if none.
func (f *FocusManager) Focused() fyne.Focusable {
	f.RLock()
	defer f.RUnlock()
	return f.focused
}

// FocusGained signals to the manager that its content got focus (due to window/overlay switch for instance).
func (f *FocusManager) FocusGained() {
	if focused := f.Focused(); focused != nil {
		focused.FocusGained()
	}
}

// FocusLost signals to the manager that its content lost focus (due to window/overlay switch for instance).
func (f *FocusManager) FocusLost() {
	if focused := f.Focused(); focused != nil {
		focused.FocusLost()
	}
}

// FocusNext will find the item after the current that can be focused and focus it.
// If current is nil then the first focusable item in the canvas will be focused.
func (f *FocusManager) FocusNext() {
	f.Lock()
	defer f.Unlock()
	f.focus(f.nextInChain(f.focused))
}

// FocusPrevious will find the item before the current that can be focused and focus it.
// If current is nil then the last focusable item in the canvas will be focused.
func (f *FocusManager) FocusPrevious() {
	f.Lock()
	defer f.Unlock()
	f.focus(f.previousInChain(f.focused))
}

func (f *FocusManager) focus(obj fyne.Focusable) {
	if f.focused == obj {
		return
	}

	if dis, ok := obj.(fyne.Disableable); ok && dis.Disabled() {
		obj = nil
	}

	if f.focused != nil {
		f.focused.FocusLost()
	}
	f.focused = obj
	if obj != nil {
		obj.FocusGained()
	}
}

func (f *FocusManager) nextInChain(current fyne.Focusable) fyne.Focusable {
	var first, next fyne.Focusable
	found := current == nil // if we have no starting point then pretend we matched already
	driver.WalkVisibleObjectTree(f.content, func(obj fyne.CanvasObject, _ fyne.Position, _ fyne.Position, _ fyne.Size) bool {
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

		if obj == current.(fyne.CanvasObject) {
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
	driver.WalkVisibleObjectTree(f.content, func(obj fyne.CanvasObject, _ fyne.Position, _ fyne.Position, _ fyne.Size) bool {
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

package app

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver"
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
func (f *FocusManager) Focus(obj fyne.Focusable) bool {
	f.Lock()
	defer f.Unlock()
	if obj != nil {
		var hiddenAncestor fyne.CanvasObject
		hidden := false
		found := driver.WalkCompleteObjectTree(
			f.content,
			func(object fyne.CanvasObject, _, _ fyne.Position, _ fyne.Size) bool {
				if hiddenAncestor == nil && !object.Visible() {
					hiddenAncestor = object
				}
				if object == obj.(fyne.CanvasObject) {
					hidden = hiddenAncestor != nil
					return true
				}
				return false
			},
			func(object, _ fyne.CanvasObject) {
				if hiddenAncestor == object {
					hiddenAncestor = nil
				}
			},
		)
		if !found {
			return false
		}
		if hidden {
			return true
		}
		if dis, ok := obj.(fyne.Disableable); ok && dis.Disabled() {
			type selectableText interface {
				SelectedText() string
			}
			if _, isSelectableText := obj.(selectableText); !isSelectableText || fyne.CurrentDevice().IsMobile() {
				return true
			}
		}
	}
	f.focus(obj)
	return true
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

	if f.focused != nil {
		f.focused.FocusLost()
	}
	f.focused = obj
	if obj != nil {
		obj.FocusGained()
	}
}

func (f *FocusManager) nextInChain(current fyne.Focusable) fyne.Focusable {
	return f.nextWithWalker(current, driver.WalkVisibleObjectTree)
}

func (f *FocusManager) nextWithWalker(current fyne.Focusable, walker walkerFunc) fyne.Focusable {
	var next fyne.Focusable
	found := current == nil // if we have no starting point then pretend we matched already
	walker(f.content, func(obj fyne.CanvasObject, _ fyne.Position, _ fyne.Position, _ fyne.Size) bool {
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
		if next == nil {
			next = focus
		}

		if obj == current.(fyne.CanvasObject) {
			found = true
		}

		return false
	}, nil)

	return next
}

func (f *FocusManager) previousInChain(current fyne.Focusable) fyne.Focusable {
	return f.nextWithWalker(current, driver.ReverseWalkVisibleObjectTree)
}

type walkerFunc func(
	fyne.CanvasObject,
	func(fyne.CanvasObject, fyne.Position, fyne.Position, fyne.Size) bool,
	func(fyne.CanvasObject, fyne.CanvasObject),
) bool

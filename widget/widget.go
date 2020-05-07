// Package widget defines the UI widgets within the Fyne toolkit
package widget // import "fyne.io/fyne/widget"

import (
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/cache"
)

// BaseWidget provides a helper that handles basic widget behaviours.
type BaseWidget struct {
	size     fyne.Size
	position fyne.Position
	Hidden   bool

	impl         fyne.Widget
	propertyLock sync.RWMutex
}

// ReadFields provides a guaranteed thread safe way to access widget fields.
func (w *BaseWidget) ReadFields(f func()) {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()

	f()
}

// SetFieldsAndRefresh helps to make changes to a widget that should be followed by a refresh.
// This method is a guaranteed thread-safe way of directly manipulating widget fields.
func (w *BaseWidget) SetFieldsAndRefresh(f func()) {
	w.propertyLock.Lock()
	f()
	w.propertyLock.Unlock()

	w.Refresh()
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
func (w *BaseWidget) ExtendBaseWidget(wid fyne.Widget) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()

	if w.impl != nil {
		return
	}

	w.impl = wid
}

// Size gets the current size of this widget.
func (w *BaseWidget) Size() fyne.Size {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()

	return w.size
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *BaseWidget) Resize(size fyne.Size) {
	var impl fyne.Widget
	w.ReadFields(func() {
		if w.size == size {
			return
		}
		w.size = size

		impl = w.impl
	})
	if impl == nil {
		return
	}
	cache.Renderer(impl).Layout(size)
}

// Position gets the current position of this widget, relative to its parent.
func (w *BaseWidget) Position() fyne.Position {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()

	return w.position
}

// Move the widget to a new position, relative to its parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *BaseWidget) Move(pos fyne.Position) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()

	w.position = pos
}

// MinSize for the widget - it should never be resized below this value.
func (w *BaseWidget) MinSize() fyne.Size {
	w.propertyLock.RLock()
	impl := w.impl
	w.propertyLock.RUnlock()

	r := cache.Renderer(impl)
	if r == nil {
		return fyne.NewSize(0, 0)
	}

	return r.MinSize()
}

// Visible returns whether or not this widget should be visible.
// Note that this may not mean it is currently visible if a parent has been hidden.
func (w *BaseWidget) Visible() bool {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()

	return !w.Hidden
}

// Show this widget so it becomes visible
func (w *BaseWidget) Show() {
	w.SetFieldsAndRefresh(func() {
		if !w.Hidden {
			return
		}

		w.Hidden = false
	})
}

// Hide this widget so it is no longer visible
func (w *BaseWidget) Hide() {
	var impl fyne.Widget
	w.ReadFields(func() {
		if w.Hidden {
			return
		}

		w.Hidden = true
		impl = w.impl
	})
	if impl == nil {
		return
	}
	canvas.Refresh(impl)
}

// Refresh causes this widget to be redrawn in it's current state
func (w *BaseWidget) Refresh() {
	w.propertyLock.RLock()
	impl := w.impl
	w.propertyLock.RUnlock()

	if impl == nil {
		return
	}

	render := cache.Renderer(impl)
	render.Refresh()
}

// super will return the actual object that this represents.
// If extended then this is the extending widget, otherwise it is self.
func (w *BaseWidget) super() fyne.Widget {
	if w.impl == nil {
		var x interface{}
		x = w
		return x.(fyne.Widget)
	}

	return w.impl
}

// DisableableWidget describes an extension to BaseWidget which can be disabled.
// Disabled widgets should have a visually distinct style when disabled, normally using theme.DisabledButtonColor.
type DisableableWidget struct {
	BaseWidget

	disabled bool
}

// Enable this widget, updating any style or features appropriately.
func (w *DisableableWidget) Enable() {
	w.SetFieldsAndRefresh(func() {
		if !w.disabled {
			return
		}

		w.disabled = false
	})
}

// Disable this widget so that it cannot be interacted with, updating any style appropriately.
func (w *DisableableWidget) Disable() {
	w.SetFieldsAndRefresh(func() {
		if w.disabled {
			return
		}

		w.disabled = true
	})
}

// Disabled returns true if this widget is currently disabled or false if it can currently be interacted with.
func (w *DisableableWidget) Disabled() bool {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()

	return w.disabled
}

// Renderer looks up the render implementation for a widget
// Deprecated: Access to widget renderers is being removed, render details should be private to a WidgetRenderer.
func Renderer(wid fyne.Widget) fyne.WidgetRenderer {
	return cache.Renderer(wid)
}

// DestroyRenderer frees a render implementation for a widget.
// This is typically for internal use only.
// Deprecated: Access to widget renderers is being removed, render details should be private to a WidgetRenderer.
func DestroyRenderer(wid fyne.Widget) {
	cache.DestroyRenderer(wid)
}

// Refresh instructs the containing canvas to refresh the specified widget.
// Deprecated: Call Widget.Refresh() instead.
func Refresh(wid fyne.Widget) {
	wid.Refresh()
}

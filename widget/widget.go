// Package widget defines the UI widgets within the Fyne toolkit
package widget // import "fyne.io/fyne/widget"

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/cache"
)

// BaseWidget provides a helper that handles basic widget behaviours.
type BaseWidget struct {
	size     fyne.Size
	position fyne.Position
	Hidden   bool

	impl fyne.Widget
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
func (w *BaseWidget) ExtendBaseWidget(wid fyne.Widget) {
	if w.impl != nil {
		return
	}

	w.impl = wid
}

// Size gets the current size of this widget.
func (w *BaseWidget) Size() fyne.Size {
	return w.size
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *BaseWidget) Resize(size fyne.Size) {
	if w.size == size {
		return
	}
	w.size = size

	if w.impl == nil {
		return
	}
	cache.Renderer(w.impl).Layout(size)
}

// Position gets the current position of this widget, relative to its parent.
func (w *BaseWidget) Position() fyne.Position {
	return w.position
}

// Move the widget to a new position, relative to its parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *BaseWidget) Move(pos fyne.Position) {
	w.position = pos
}

// MinSize for the widget - it should never be resized below this value.
func (w *BaseWidget) MinSize() fyne.Size {
	if w.impl == nil || cache.Renderer(w.impl) == nil {
		return fyne.NewSize(0, 0)
	}
	return cache.Renderer(w.impl).MinSize()
}

// CreateRenderer of BaseWidget does nothing, it must be overridden
func (w *BaseWidget) CreateRenderer() fyne.WidgetRenderer {
	return nil
}

// Visible returns whether or not this widget should be visible.
// Note that this may not mean it is currently visible if a parent has been hidden.
func (w *BaseWidget) Visible() bool {
	return !w.Hidden
}

// Show this widget so it becomes visible
func (w *BaseWidget) Show() {
	if !w.Hidden {
		return
	}

	w.Hidden = false
	if w.impl == nil {
		return
	}
	w.impl.Refresh()
}

// Hide this widget so it is no lonver visible
func (w *BaseWidget) Hide() {
	if w.Hidden {
		return
	}

	w.Hidden = true
	if w.impl == nil {
		return
	}
	canvas.Refresh(w.impl)
}

// Refresh causes this widget to be redrawn in it's current state
func (w *BaseWidget) Refresh() {
	if w.impl == nil {
		return
	}

	w.refresh(w.impl)
}

func (w *BaseWidget) refresh(wid fyne.Widget) {
	render := cache.Renderer(wid)
	render.Refresh()
}

// super will return the actual object that this represents.
// If extended then this is the extending widget, otherwise it is self.
func (w *BaseWidget) super() fyne.Widget {
	if w.impl == nil {
		return w
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
	if !w.disabled {
		return
	}

	w.disabled = false
	w.Refresh()
}

// Disable this widget so that it cannot be interacted with, updating any style appropriately.
func (w *DisableableWidget) Disable() {
	if w.disabled {
		return
	}

	w.disabled = true
	w.Refresh()
}

// Disabled returns true if this widget is currently disabled or false if it can currently be interacted with.
func (w *DisableableWidget) Disabled() bool {
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

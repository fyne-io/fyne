// Package widget defines the UI widgets within the Fyne toolkit
package widget // import "fyne.io/fyne/widget"

import (
	"fyne.io/fyne"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/internal/widget"
)

// BaseWidget provides a helper that handles basic widget behaviours.
type BaseWidget struct {
	widget.BaseWidget
}

// super will return the actual object that this represents.
// If extended then this is the extending widget, otherwise it is self.
func (w *BaseWidget) super() fyne.Widget {
	return widget.Implementor(&w.BaseWidget)
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

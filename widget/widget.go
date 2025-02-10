// Package widget defines the UI widgets within the Fyne toolkit.
package widget // import "fyne.io/fyne/v2/widget"

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	internalWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

// BaseWidget provides a helper that handles basic widget behaviours.
type BaseWidget struct {
	noCopy noCopy // so `go vet` can complain if a widget is passed by value (copied)

	size     fyne.Size
	position fyne.Position
	Hidden   bool

	impl       fyne.Widget
	themeCache fyne.Theme
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
func (w *BaseWidget) ExtendBaseWidget(wid fyne.Widget) {
	if w.super() != nil {
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
	if size == w.Size() {
		return
	}

	w.size = size

	impl := w.super()
	if impl == nil {
		return
	}
	cache.Renderer(impl).Layout(size)
}

// Position gets the current position of this widget, relative to its parent.
func (w *BaseWidget) Position() fyne.Position {
	return w.position
}

// Move the widget to a new position, relative to its parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *BaseWidget) Move(pos fyne.Position) {
	w.position = pos
	internalWidget.Repaint(w.super())
}

// MinSize for the widget - it should never be resized below this value.
func (w *BaseWidget) MinSize() fyne.Size {
	impl := w.super()

	r := cache.Renderer(impl)
	if r == nil {
		return fyne.Size{}
	}

	return r.MinSize()
}

// Visible returns whether or not this widget should be visible.
// Note that this may not mean it is currently visible if a parent has been hidden.
func (w *BaseWidget) Visible() bool {
	return !w.Hidden
}

// Show this widget so it becomes visible
func (w *BaseWidget) Show() {
	if w.Visible() {
		return
	}

	w.Hidden = false

	impl := w.super()
	if impl == nil {
		return
	}
	impl.Refresh()
}

// Hide this widget so it is no longer visible
func (w *BaseWidget) Hide() {
	if !w.Visible() {
		return
	}

	w.Hidden = true

	impl := w.super()
	if impl == nil {
		return
	}
	canvas.Refresh(impl)
}

// Refresh causes this widget to be redrawn in its current state
func (w *BaseWidget) Refresh() {
	impl := w.super()
	if impl == nil {
		return
	}

	w.themeCache = nil

	cache.Renderer(impl).Refresh()
}

// Theme returns a cached Theme instance for this widget (or its extending widget).
// This will be the app theme in most cases, or a widget specific theme if it is inside a ThemeOverride container.
//
// Since: 2.5
func (w *BaseWidget) Theme() fyne.Theme {
	cached := w.themeCache
	if cached != nil {
		return cached
	}

	cached = cache.WidgetTheme(w.super())
	// don't cache the default as it may change
	if cached == nil {
		return theme.Current()
	}

	w.themeCache = cached
	return cached
}

// super will return the actual object that this represents.
// If extended then this is the extending widget, otherwise it is nil.
func (w *BaseWidget) super() fyne.Widget {
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
	if !w.Disabled() {
		return // Enabled already
	}

	w.disabled = false

	impl := w.super()
	if impl == nil {
		return
	}
	impl.Refresh()
}

// Disable this widget so that it cannot be interacted with, updating any style appropriately.
func (w *DisableableWidget) Disable() {
	if w.Disabled() {
		return // Disabled already
	}

	w.disabled = true

	impl := w.super()
	if impl == nil {
		return
	}
	impl.Refresh()
}

// Disabled returns true if this widget is currently disabled or false if it can currently be interacted with.
func (w *DisableableWidget) Disabled() bool {
	return w.disabled
}

// NewSimpleRenderer creates a new SimpleRenderer to render a widget using a
// single fyne.CanvasObject.
//
// Since: 2.1
func NewSimpleRenderer(object fyne.CanvasObject) fyne.WidgetRenderer {
	return internalWidget.NewSimpleRenderer(object)
}

// Orientation controls the horizontal/vertical layout of a widget
type Orientation int

// Orientation constants to control widget layout
const (
	Horizontal Orientation = 0
	Vertical   Orientation = 1

	// Adaptive will switch between horizontal and vertical layouts according to device orientation.
	// This orientation is not always supported and interpretation can vary per-widget.
	//
	// Since: 2.5
	Adaptive Orientation = 2
)

type noCopy struct{}

func (*noCopy) Lock() {}

func (*noCopy) Unlock() {}

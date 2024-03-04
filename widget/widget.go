// Package widget defines the UI widgets within the Fyne toolkit.
package widget // import "fyne.io/fyne/v2/widget"

import (
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/internal/cache"
	internalWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

// BaseWidget provides a helper that handles basic widget behaviours.
type BaseWidget struct {
	size     async.Size
	position async.Position
	Hidden   bool

	impl         atomic.Pointer[fyne.Widget]
	propertyLock sync.RWMutex
	themeCache   fyne.Theme
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
func (w *BaseWidget) ExtendBaseWidget(wid fyne.Widget) {
	impl := w.super()
	if impl != nil {
		return
	}

	w.impl.Store(&wid)
}

// Size gets the current size of this widget.
func (w *BaseWidget) Size() fyne.Size {
	return w.size.Load()
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *BaseWidget) Resize(size fyne.Size) {
	if size == w.Size() {
		return
	}

	w.size.Store(size)

	impl := w.super()
	if impl == nil {
		return
	}
	cache.Renderer(impl).Layout(size)
}

// Position gets the current position of this widget, relative to its parent.
func (w *BaseWidget) Position() fyne.Position {
	return w.position.Load()
}

// Move the widget to a new position, relative to its parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *BaseWidget) Move(pos fyne.Position) {
	w.position.Store(pos)
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
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()

	return !w.Hidden
}

// Show this widget so it becomes visible
func (w *BaseWidget) Show() {
	if w.Visible() {
		return
	}

	w.SetFieldsAndRefresh(func() {
		w.Hidden = false
	})
}

// Hide this widget so it is no longer visible
func (w *BaseWidget) Hide() {
	if !w.Visible() {
		return
	}

	w.propertyLock.Lock()
	w.Hidden = true
	w.propertyLock.Unlock()

	impl := w.super()
	if impl == nil {
		return
	}
	canvas.Refresh(impl)
}

// Refresh causes this widget to be redrawn in it's current state
func (w *BaseWidget) Refresh() {
	impl := w.super()
	if impl == nil {
		return
	}

	w.propertyLock.Lock()
	w.themeCache = nil
	w.propertyLock.Unlock()

	render := cache.Renderer(impl)
	render.Refresh()
}

// Theme returns a cached Theme instance for this widget (or its extending widget).
// This will be the app theme in most cases, or a widget specific theme if it is inside a ThemeOverride container.
//
// Since: 2.5
func (w *BaseWidget) Theme() fyne.Theme {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	return w.themeWithLock()
}

func (w *BaseWidget) themeWithLock() fyne.Theme {
	cached := w.themeCache
	if cached == nil {
		cached = cache.WidgetTheme(w.super())
		// don't cache the default as it may change
		if cached == nil {
			return theme.Current()
		}

		w.themeCache = cached
	}

	return cached
}

// SetFieldsAndRefresh helps to make changes to a widget that should be followed by a refresh.
// This method is a guaranteed thread-safe way of directly manipulating widget fields.
// Widgets extending BaseWidget should use this in their setter functions.
//
// Since: 2.5
func (w *BaseWidget) SetFieldsAndRefresh(f func()) {
	w.propertyLock.Lock()
	f()
	w.propertyLock.Unlock()

	impl := w.super()
	if impl == nil {
		return
	}
	impl.Refresh()
}

// super will return the actual object that this represents.
// If extended then this is the extending widget, otherwise it is nil.
func (w *BaseWidget) super() fyne.Widget {
	impl := w.impl.Load()
	if impl == nil {
		return nil
	}

	return *impl
}

// DisableableWidget describes an extension to BaseWidget which can be disabled.
// Disabled widgets should have a visually distinct style when disabled, normally using theme.DisabledButtonColor.
type DisableableWidget struct {
	BaseWidget

	disabled atomic.Bool
}

// Enable this widget, updating any style or features appropriately.
func (w *DisableableWidget) Enable() {
	if !w.disabled.CompareAndSwap(true, false) {
		return // Enabled already
	}

	impl := w.super()
	if impl == nil {
		return
	}
	impl.Refresh()
}

// Disable this widget so that it cannot be interacted with, updating any style appropriately.
func (w *DisableableWidget) Disable() {
	if !w.disabled.CompareAndSwap(false, true) {
		return // Disabled already
	}

	impl := w.super()
	if impl == nil {
		return
	}
	impl.Refresh()
}

// Disabled returns true if this widget is currently disabled or false if it can currently be interacted with.
func (w *DisableableWidget) Disabled() bool {
	return w.disabled.Load()
}

// NewSimpleRenderer creates a new SimpleRenderer to render a widget using a
// single fyne.CanvasObject.
//
// Since: 2.1
func NewSimpleRenderer(object fyne.CanvasObject) fyne.WidgetRenderer {
	return internalWidget.NewSimpleRenderer(object)
}

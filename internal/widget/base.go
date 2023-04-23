package widget

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
)

// Base provides a helper that handles basic widget behaviours.
type Base struct {
	hidden   bool
	position fyne.Position
	size     fyne.Size

	impl         fyne.Widget
	propertyLock sync.RWMutex
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
func (w *Base) ExtendBaseWidget(wid fyne.Widget) {
	impl := w.super()
	if impl != nil {
		return
	}

	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	w.impl = wid
}

// Size gets the current size of this widget.
func (w *Base) Size() fyne.Size {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()

	return w.size
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *Base) Resize(size fyne.Size) {
	w.propertyLock.RLock()
	baseSize := w.size
	impl := w.impl
	w.propertyLock.RUnlock()
	if baseSize == size {
		return
	}

	w.propertyLock.Lock()
	w.size = size
	w.propertyLock.Unlock()

	if impl == nil {
		return
	}
	cache.Renderer(impl).Layout(size)
}

// Position gets the current position of this widget, relative to its parent.
func (w *Base) Position() fyne.Position {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()

	return w.position
}

// Move the widget to a new position, relative to its parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *Base) Move(pos fyne.Position) {
	w.propertyLock.Lock()
	w.position = pos
	w.propertyLock.Unlock()

	Repaint(w.super())
}

// MinSize for the widget - it should never be resized below this value.
func (w *Base) MinSize() fyne.Size {
	impl := w.super()

	r := cache.Renderer(impl)
	if r == nil {
		return fyne.NewSize(0, 0)
	}

	return r.MinSize()
}

// Visible returns whether or not this widget should be visible.
// Note that this may not mean it is currently visible if a parent has been hidden.
func (w *Base) Visible() bool {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()

	return !w.hidden
}

// Show this widget so it becomes visible
func (w *Base) Show() {
	if w.Visible() {
		return
	}

	w.setFieldsAndRefresh(func() {
		w.hidden = false
	})
}

// Hide this widget so it is no longer visible
func (w *Base) Hide() {
	if !w.Visible() {
		return
	}

	w.propertyLock.Lock()
	w.hidden = true
	impl := w.impl
	w.propertyLock.Unlock()

	if impl == nil {
		return
	}
	canvas.Refresh(impl)
}

// Refresh causes this widget to be redrawn in it's current state
func (w *Base) Refresh() {
	impl := w.super()
	if impl == nil {
		return
	}

	render := cache.Renderer(impl)
	render.Refresh()
}

// setFieldsAndRefresh helps to make changes to a widget that should be followed by a refresh.
// This method is a guaranteed thread-safe way of directly manipulating widget fields.
func (w *Base) setFieldsAndRefresh(f func()) {
	w.propertyLock.Lock()
	f()
	impl := w.impl
	w.propertyLock.Unlock()

	if impl == nil {
		return
	}
	impl.Refresh()
}

// super will return the actual object that this represents.
// If extended then this is the extending widget, otherwise it is nil.
func (w *Base) super() fyne.Widget {
	w.propertyLock.RLock()
	impl := w.impl
	w.propertyLock.RUnlock()
	return impl
}

// Repaint instructs the containing canvas to redraw, even if nothing changed.
// This method is a duplicate of what is in `canvas/canvas.go` to avoid a dependency loop or public API.
func Repaint(obj fyne.CanvasObject) {
	if fyne.CurrentApp() == nil || fyne.CurrentApp().Driver() == nil {
		return
	}

	c := fyne.CurrentApp().Driver().CanvasForObject(obj)
	if c != nil {
		if paint, ok := c.(interface{ SetDirty() }); ok {
			paint.SetDirty()
		}
	}
}

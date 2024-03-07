package widget

import (
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/internal/cache"
)

// Base provides a helper that handles basic widget behaviours.
type Base struct {
	hidden   atomic.Bool
	position async.Position
	size     async.Size
	impl     atomic.Pointer[fyne.Widget]
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
func (w *Base) ExtendBaseWidget(wid fyne.Widget) {
	impl := w.super()
	if impl != nil {
		return
	}

	w.impl.Store(&wid)
}

// Size gets the current size of this widget.
func (w *Base) Size() fyne.Size {
	return w.size.Load()
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *Base) Resize(size fyne.Size) {
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
func (w *Base) Position() fyne.Position {
	return w.position.Load()
}

// Move the widget to a new position, relative to its parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *Base) Move(pos fyne.Position) {
	w.position.Store(pos)

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
	return !w.hidden.Load()
}

// Show this widget so it becomes visible
func (w *Base) Show() {
	if !w.hidden.CompareAndSwap(true, false) {
		return // Visible already
	}

	impl := w.super()
	if impl == nil {
		return
	}
	impl.Refresh()
}

// Hide this widget so it is no longer visible
func (w *Base) Hide() {
	if !w.hidden.CompareAndSwap(false, true) {
		return // Hidden already
	}

	impl := w.super()
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

// super will return the actual object that this represents.
// If extended then this is the extending widget, otherwise it is nil.
func (w *Base) super() fyne.Widget {
	impl := w.impl.Load()
	if impl == nil {
		return nil
	}

	return *impl
}

// Repaint instructs the containing canvas to redraw, even if nothing changed.
// This method is a duplicate of what is in `canvas/canvas.go` to avoid a dependency loop or public API.
func Repaint(obj fyne.CanvasObject) {
	app := fyne.CurrentApp()
	if app == nil || app.Driver() == nil {
		return
	}

	c := app.Driver().CanvasForObject(obj)
	if c != nil {
		if paint, ok := c.(interface{ SetDirty() }); ok {
			paint.SetDirty()
		}
	}
}

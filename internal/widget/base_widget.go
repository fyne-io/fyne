package widget

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
	r := cache.Renderer(w.impl)
	if r == nil {
		return fyne.NewSize(0, 0)
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

	render := cache.Renderer(w.impl)
	render.Refresh()
}

// Implementor returns the implementing widget for a BaseWidget.
func Implementor(w *BaseWidget) fyne.Widget {
	if w.impl == nil {
		var x interface{}
		x = w
		return x.(fyne.Widget)
	}

	return w.impl
}

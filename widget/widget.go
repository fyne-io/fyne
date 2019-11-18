// Package widget defines the UI widgets within the Fyne toolkit
package widget // import "fyne.io/fyne/widget"

import (
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
)

var renderers sync.Map

// BaseWidget provides a helper that handles basic widget behaviours.
type BaseWidget struct {
	size     fyne.Size
	position fyne.Position
	Hidden   bool
	disabled bool

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
	w.size = size

	if w.impl == nil {
		return
	}
	Renderer(w.impl).Layout(size)
}

// Position gets the current position of this widget, relative to its parent.
func (w *BaseWidget) Position() fyne.Position {
	return w.position
}

// Move the widget to a new position, relative to its parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *BaseWidget) Move(pos fyne.Position) {
	w.position = pos

	if w.impl == nil {
		return
	}
	canvas.Refresh(w.impl)
}

// MinSize for the widget - it should never be resized below this value.
func (w *BaseWidget) MinSize() fyne.Size {
	if w.impl == nil || Renderer(w.impl) == nil {
		//	fyne.LogError("Renderer not configured, perhaps missing call to ExtendBaseWidget()?", nil)
		return fyne.NewSize(0, 0)
	}
	return Renderer(w.impl).MinSize()
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
	Refresh(w.impl)
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
	Refresh(w.impl)
}

func (w *BaseWidget) enable(parent fyne.Widget) {
	if !w.disabled {
		return
	}

	w.disabled = false
	Refresh(parent)
}

func (w *BaseWidget) disable(parent fyne.Widget) {
	if w.disabled {
		return
	}

	w.disabled = true
	Refresh(parent)
}

func (w *BaseWidget) super() fyne.Widget {
	return w.impl
}

type isBaseWidget interface {
	ExtendBaseWidget(fyne.Widget)
	super() fyne.Widget
}

// Renderer looks up the render implementation for a widget
func Renderer(wid fyne.Widget) fyne.WidgetRenderer {
	if wd, ok := wid.(isBaseWidget); ok {
		if wd.super() != nil {
			wid = wd.super()
		}
	}
	renderer, ok := renderers.Load(wid)
	if !ok {
		renderer = wid.CreateRenderer()
		renderers.Store(wid, renderer)
	}

	return renderer.(fyne.WidgetRenderer)
}

// DestroyRenderer frees a render implementation for a widget.
// This is typically for internal use only.
func DestroyRenderer(wid fyne.Widget) {
	Renderer(wid).Destroy()

	renderers.Delete(wid)
}

// Refresh instructs the containing canvas to refresh the specified widget.
func Refresh(wid fyne.Widget) {
	render := Renderer(wid)

	if render != nil {
		render.Refresh()
	}
}

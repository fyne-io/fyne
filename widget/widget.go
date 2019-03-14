// Package widget defines the UI widgets within the Fyne toolkit
package widget // import "fyne.io/fyne/widget"

import (
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
)

// A base widget class to define the standard widget behaviours.
type baseWidget struct {
	size     fyne.Size
	position fyne.Position
	Hidden   bool
}

// Get the current size of this widget.
func (w *baseWidget) Size() fyne.Size {
	return w.size
}

// Set a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *baseWidget) resize(size fyne.Size, parent fyne.Widget) {
	w.size = size

	Renderer(parent).Layout(size)
}

// Get the current position of this widget, relative to it's parent.
func (w *baseWidget) Position() fyne.Position {
	return w.position
}

// Move the widget to a new position, relative to it's parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *baseWidget) move(pos fyne.Position, parent fyne.Widget) {
	w.position = pos

	canvas.Refresh(parent)
}

func (w *baseWidget) minSize(parent fyne.Widget) fyne.Size {
	if Renderer(parent) == nil {
		return fyne.NewSize(0, 0)
	}
	return Renderer(parent).MinSize()
}

func (w *baseWidget) Visible() bool {
	return !w.Hidden
}

func (w *baseWidget) show(parent fyne.Widget) {
	w.Hidden = false
	for _, child := range Renderer(parent).Objects() {
		child.Show()
	}

	canvas.Refresh(parent)
}

func (w *baseWidget) hide(parent fyne.Widget) {
	w.Hidden = true
	for _, child := range Renderer(parent).Objects() {
		child.Hide()
	}

	canvas.Refresh(parent)
}

var renderers sync.Map

// Renderer looks up the render implementation for a widget
func Renderer(wid fyne.Widget) fyne.WidgetRenderer {
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

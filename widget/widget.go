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
	lck      sync.RWMutex
}

// Get the current size of this widget.
func (w *baseWidget) Size() fyne.Size {
	w.lck.RLock()
	defer w.lck.RUnlock()
	return w.size
}

// Set a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *baseWidget) resize(size fyne.Size, parent fyne.Widget) {
	w.lck.Lock()
	w.size = size
	w.lck.Unlock()

	Renderer(parent).Layout(size)
}

// Get the current position of this widget, relative to it's parent.
func (w *baseWidget) Position() fyne.Position {
	w.lck.RLock()
	defer w.lck.RUnlock()
	return w.position
}

// Move the widget to a new position, relative to it's parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *baseWidget) move(pos fyne.Position, parent fyne.Widget) {
	w.lck.Lock()
	w.position = pos
	w.lck.Unlock()

	canvas.Refresh(parent)
}

func (w *baseWidget) minSize(parent fyne.Widget) fyne.Size {
	if Renderer(parent) == nil {
		return fyne.NewSize(0, 0)
	}
	return Renderer(parent).MinSize()
}

func (w *baseWidget) Visible() bool {
	w.lck.RLock()
	defer w.lck.RUnlock()
	return !w.Hidden
}

func (w *baseWidget) show(parent fyne.Widget) {
	w.lck.Lock()
	w.Hidden = false
	w.lck.Unlock()
	for _, child := range Renderer(parent).Objects() {
		child.Show()
	}

	canvas.Refresh(parent)
}

func (w *baseWidget) hide(parent fyne.Widget) {
	w.lck.Lock()
	w.Hidden = true
	w.lck.Unlock()
	for _, child := range Renderer(parent).Objects() {
		child.Hide()
	}

	canvas.Refresh(parent)
}

var renderers sync.Map

var creatingRenderers = map[fyne.Widget]*sync.Mutex{}
var creatingRendererLck sync.Mutex

// Renderer looks up the render implementation for a widget
func Renderer(wid fyne.Widget) fyne.WidgetRenderer {
start:
	renderer, ok := renderers.Load(wid)
	if !ok {
		creatingRendererLck.Lock()
		if lck, prog := creatingRenderers[wid]; prog {
			creatingRendererLck.Unlock()
			lck.Lock()
			lck.Unlock()
			goto start
		}
		lck := new(sync.Mutex)
		lck.Lock()
		creatingRenderers[wid] = lck
		creatingRendererLck.Unlock()
		renderer = wid.CreateRenderer()
		renderers.Store(wid, renderer)
		creatingRendererLck.Lock()
		delete(creatingRenderers, wid)
		lck.Unlock()
		creatingRendererLck.Unlock()
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

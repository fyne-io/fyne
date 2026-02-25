package widget

import "fyne.io/fyne/v2"

// BaseRenderer is a renderer base that provides part of the widget.Renderer interface.
type BaseRenderer struct {
	objects []fyne.CanvasObject
}

// NewBaseRenderer creates a new BaseRenderer.
func NewBaseRenderer(objects []fyne.CanvasObject) BaseRenderer {
	return BaseRenderer{objects}
}

// Destroy does nothing in the base implementation.
func (r *BaseRenderer) Destroy() {
}

// Objects returns the objects that should be rendered.
func (r *BaseRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// SetObjects updates the objects of the renderer.
func (r *BaseRenderer) SetObjects(objects []fyne.CanvasObject) {
	r.objects = objects
}

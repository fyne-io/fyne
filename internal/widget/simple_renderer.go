package widget

import "fyne.io/fyne/v2"

var _ fyne.WidgetRenderer = (*SimpleRenderer)(nil)

// SimpleRenderer is a basic renderer that satisfies widget.Renderer interface by wrapping
// a single fyne.CanvasObject.
//
// Since: 2.1
type SimpleRenderer struct {
	object fyne.CanvasObject
}

// NewSimpleRenderer creates a new BaseRenderer.
//
// Since: 2.1
func NewSimpleRenderer(object fyne.CanvasObject) *SimpleRenderer {
	return &SimpleRenderer{object}
}

// Destroy does nothing in the base implementation.
//
// Implements: fyne.WidgetRenderer
//
// Since: 2.1
func (r *SimpleRenderer) Destroy() {
}

// Layout updates the contained object to be the requested size.
//
// Implements: fyne.WidgetRenderer
//
// Since: 2.1
func (r *SimpleRenderer) Layout(s fyne.Size) {
	r.object.Resize(s)
}

// MinSize returns the smallest size that this render can use, returned from the underlying object.
//
// Implements: fyne.WidgetRenderer
//
// Since: 2.1
func (r *SimpleRenderer) MinSize() fyne.Size {
	return r.object.MinSize()
}

// Objects returns the objects that should be rendered.
//
// Implements: fyne.WidgetRenderer
//
// Since: 2.1
func (r *SimpleRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.object}
}

// Refresh requests the underlying object to redraw.
//
// Implements: fyne.WidgetRenderer
//
// Since: 2.1
func (r *SimpleRenderer) Refresh() {
	r.object.Refresh()
}

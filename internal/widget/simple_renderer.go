package widget

import "fyne.io/fyne/v2"

var _ fyne.WidgetRenderer = (*SimpleRenderer)(nil)

// SimpleRenderer is a basic renderer that satisfies widget.Renderer interface by wrapping
// a single fyne.CanvasObject.
//
// Since: 2.1
type SimpleRenderer struct {
	objects []fyne.CanvasObject
}

// NewSimpleRenderer creates a new SimpleRenderer to render a widget using a
// single CanvasObject.
//
// Since: 2.1
func NewSimpleRenderer(object fyne.CanvasObject) *SimpleRenderer {
	return &SimpleRenderer{[]fyne.CanvasObject{object}}
}

// Destroy does nothing in this implementation.
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
	r.objects[0].Resize(s)
}

// MinSize returns the smallest size that this render can use, returned from the underlying object.
//
// Implements: fyne.WidgetRenderer
//
// Since: 2.1
func (r *SimpleRenderer) MinSize() fyne.Size {
	return r.objects[0].MinSize()
}

// Objects returns the objects that should be rendered.
//
// Implements: fyne.WidgetRenderer
//
// Since: 2.1
func (r *SimpleRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh requests the underlying object to redraw.
//
// Implements: fyne.WidgetRenderer
//
// Since: 2.1
func (r *SimpleRenderer) Refresh() {
	r.objects[0].Refresh()
}

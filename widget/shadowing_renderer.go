package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/internal/widget"
)

// When using the shadowingRenderer the embedding renderer should call
// layoutShadow(contentSize, contentPos) to lay out the shadow.
type shadowingRenderer struct {
	widget.BaseRenderer
	shadow fyne.CanvasObject
}

func newShadowingRenderer(objects []fyne.CanvasObject, level elevationLevel) *shadowingRenderer {
	var s fyne.CanvasObject
	if level > 0 {
		s = newShadow(shadowAround, level)
	}
	r := &shadowingRenderer{shadow: s}
	r.SetObjects(objects)
	return r
}

func (r *shadowingRenderer) layoutShadow(size fyne.Size, pos fyne.Position) {
	if r.shadow == nil {
		return
	}
	r.shadow.Resize(size)
	r.shadow.Move(pos)
}

// SetObjects updates the objects of the renderer.
func (r *shadowingRenderer) SetObjects(objects []fyne.CanvasObject) {
	if r.shadow != nil {
		objects = append([]fyne.CanvasObject{r.shadow}, objects...)
	}
	r.BaseRenderer.SetObjects(objects)
}

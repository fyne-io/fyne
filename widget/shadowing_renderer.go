package widget

import (
	"fyne.io/fyne"
)

// When using the shadowingRenderer the embedding renderer should call
// layoutShadow(contentSize, contentPos) to lay out the shadow.
type shadowingRenderer struct {
	baseRenderer
	shadow fyne.CanvasObject
}

func newShadowingRenderer(objects []fyne.CanvasObject, level elevationLevel) *shadowingRenderer {
	var s fyne.CanvasObject
	if level > 0 {
		s = newShadow(shadowAround, level)
	}
	r := &shadowingRenderer{shadow: s}
	r.setObjects(objects)
	return r
}

func (r *shadowingRenderer) layoutShadow(size fyne.Size, pos fyne.Position) {
	if r.shadow == nil {
		return
	}
	r.shadow.Resize(size)
	r.shadow.Move(pos)
}

func (r *shadowingRenderer) setObjects(objects []fyne.CanvasObject) {
	if r.shadow != nil {
		objects = append([]fyne.CanvasObject{r.shadow}, objects...)
	}
	r.baseRenderer.setObjects(objects)
}

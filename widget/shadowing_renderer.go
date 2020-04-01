package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

// When using the shadowingRenderer the embedding renderer should call
// layoutShadow(contentSize, contentPos) to lay out the shadow.
type shadowingRenderer struct {
	baseRenderer
	shadow fyne.CanvasObject
}

func newShadowingRenderer(objects []fyne.CanvasObject, level int) *shadowingRenderer {
	var shadow fyne.CanvasObject
	if level > 0 {
		shadow = newShadow(shadowAround, level*theme.Padding()/2)
		objects = append(objects, shadow)
	}
	return &shadowingRenderer{baseRenderer{objects}, shadow}
}

func (r *shadowingRenderer) layoutShadow(size fyne.Size, pos fyne.Position) {
	if r.shadow == nil {
		return
	}
	r.shadow.Resize(size)
	r.shadow.Move(pos)
}

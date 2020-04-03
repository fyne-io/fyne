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
	var shadow fyne.CanvasObject
	if level > 0 {
		shadow = newShadow(shadowAround, level)
		objects = append(objects, shadow)
	}
	return &shadowingRenderer{widget.NewBaseRenderer(objects), shadow}
}

func (r *shadowingRenderer) layoutShadow(size fyne.Size, pos fyne.Position) {
	if r.shadow == nil {
		return
	}
	r.shadow.Resize(size)
	r.shadow.Move(pos)
}

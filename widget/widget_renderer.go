package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

// When using the baseRenderer the embedding renderer should call
// layoutShadow(contentSize, contentPos) to lay out the shadow.
type baseRenderer struct {
	objects []fyne.CanvasObject
	sh      fyne.CanvasObject
}

func newBaseRenderer(objects []fyne.CanvasObject, level int) *baseRenderer {
	var shadow fyne.CanvasObject
	if level > 0 {
		shadow = newShadow(shadowAround, level*theme.Padding()/2)
		objects = append(objects, shadow)
	}
	return &baseRenderer{objects, shadow}
}

func (r *baseRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *baseRenderer) layoutShadow(size fyne.Size, pos fyne.Position) {
	if r.sh == nil {
		return
	}
	r.sh.Resize(size)
	r.sh.Move(pos)
}

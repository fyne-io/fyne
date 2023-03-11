package widget

import (
	"fyne.io/fyne/v2"
)

// ShadowingRenderer is a renderer that adds a shadow arount the rendered content.
// When using the ShadowingRenderer the embedding renderer should call
// LayoutShadow(contentSize, contentPos) to lay out the shadow.
type ShadowingRenderer struct {
	BaseRenderer
	shadow fyne.CanvasObject
}

// NewShadowingRenderer creates a ShadowingRenderer.
func NewShadowingRenderer(objects []fyne.CanvasObject, level ElevationLevel) *ShadowingRenderer {
	var s fyne.CanvasObject
	if level > 0 {
		s = NewShadow(ShadowAround, level)
	}
	r := &ShadowingRenderer{shadow: s}
	r.SetObjects(objects)
	return r
}

// LayoutShadow adjusts the size and position of the shadow if necessary.
func (r *ShadowingRenderer) LayoutShadow(size fyne.Size, pos fyne.Position) {
	if r.shadow == nil {
		return
	}
	r.shadow.Resize(size)
	r.shadow.Move(pos)
}

// SetObjects updates the renderer's objects including the shadow if necessary.
func (r *ShadowingRenderer) SetObjects(objects []fyne.CanvasObject) {
	if r.shadow != nil && len(objects) > 0 && r.shadow != objects[0] {
		objects = append([]fyne.CanvasObject{r.shadow}, objects...)
	}
	r.BaseRenderer.SetObjects(objects)
}

// RefreshShadow asks the shadow graphical element to update to current theme
func (r *ShadowingRenderer) RefreshShadow() {
	if r.shadow == nil {
		return
	}
	r.shadow.Refresh()
}

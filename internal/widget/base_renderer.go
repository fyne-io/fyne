package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

// NewBaseRenderer creates a new BaseRenderer.
func NewBaseRenderer(objects []fyne.CanvasObject) BaseRenderer {
	return BaseRenderer{objects}
}

// BaseRenderer is a renderer base providing the most common implementations of a part of the
// widget.Renderer interface.
type BaseRenderer struct {
	objects []fyne.CanvasObject
}

// BackgroundColor satisfies the fyne.WidgetRenderer interface.
func (r *BaseRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

// Destroy satisfies the fyne.WidgetRenderer interface.
func (r *BaseRenderer) Destroy() {
}

// Objects satisfies the fyne.WidgetRenderer interface.
func (r *BaseRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// SetObjects updates the objects of the renderer.
func (r *BaseRenderer) SetObjects(objects []fyne.CanvasObject) {
	r.objects = objects
}

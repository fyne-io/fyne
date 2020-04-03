package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

// BaseRenderer is a basic renderer handling objects and providing most common implementations.
type BaseRenderer struct {
	objects []fyne.CanvasObject
}

// NewBaseRenderer creates a BaseRenderer.
func NewBaseRenderer(objects []fyne.CanvasObject) BaseRenderer {
	return BaseRenderer{objects}
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

// SetObjects sets the objects of this renderer.
func (r *BaseRenderer) SetObjects(objects []fyne.CanvasObject) {
	r.objects = objects
}

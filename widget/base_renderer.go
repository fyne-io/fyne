package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

type baseRenderer struct {
	objects []fyne.CanvasObject
}

func (r *baseRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *baseRenderer) Destroy() {
}

func (r *baseRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *baseRenderer) SetObjects(objects []fyne.CanvasObject) {
	r.objects = objects
}

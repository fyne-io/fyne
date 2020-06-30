package painter

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
)

// VectorPad returns the number of additional points that should be added around a texture.
// This is to accommodate overflow caused by stroke and line endings etc.
// THe result is in fyne.Size type coordinates and should be scaled for output.
func VectorPad(obj fyne.CanvasObject) int {
	switch co := obj.(type) {
	case *canvas.Circle:
		if co.StrokeWidth > 0 && co.StrokeColor != nil {
			return int(co.StrokeWidth) + 2
		}
	case *canvas.Line:
		if co.StrokeWidth > 0 {
			return int(co.StrokeWidth) + 2
		}
	case *canvas.Rectangle:
		if co.StrokeWidth > 0 && co.StrokeColor != nil {
			return int(co.StrokeWidth) + 2
		}
	}

	return 0
}

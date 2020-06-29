package painter

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
)

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

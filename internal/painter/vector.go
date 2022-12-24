package painter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// VectorPad returns the number of additional points that should be added around a texture.
// This is to accommodate overflow caused by stroke and line endings etc.
// THe result is in fyne.Size type coordinates and should be scaled for output.
func VectorPad(obj fyne.CanvasObject) float32 {
	switch co := obj.(type) {
	case *canvas.Circle:
		if co.StrokeWidth > 0 && co.StrokeColor != nil {
			return co.StrokeWidth + 2
		}
		return 1 // anti-alias on circle fill
	case *canvas.Line:
		if co.StrokeWidth > 0 {
			return co.StrokeWidth + 2
		}
	case *canvas.Rectangle:
		if co.StrokeWidth > 0 && co.StrokeColor != nil {
			return co.StrokeWidth + 2
		}
	case *canvas.Text:
		if co.TextStyle.Italic {
			return co.TextSize / 5 // make sure that even a 20% lean does not overflow
		}
	}

	return 0
}

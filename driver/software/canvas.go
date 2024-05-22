package software

import (
	"fyne.io/fyne/v2/internal/driver/software"
	painter "fyne.io/fyne/v2/internal/painter/software"
)

// NewCanvas creates a new canvas in memory that can render without hardware support.
func NewCanvas() software.WindowlessCanvas {
	return software.NewCanvasWithPainter(painter.NewPainter())
}

// NewTransparentCanvas creates a new canvas in memory that can render without hardware support without a background color.
//
// Since: 2.2
func NewTransparentCanvas() software.WindowlessCanvas {
	return software.NewTransparentCanvasWithPainter(painter.NewPainter())
}

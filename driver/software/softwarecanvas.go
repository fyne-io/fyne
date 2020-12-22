package software

import (
	"fyne.io/fyne/internal/painter/software"
	"fyne.io/fyne/test"
)

// NewCanvas creates a new canvas in memory that can render without hardware support
func NewCanvas() test.WindowlessCanvas {
	return test.NewCanvasWithPainter(software.NewPainter())
}

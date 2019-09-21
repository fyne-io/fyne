package playground

import (
	"fyne.io/fyne/internal/painter/software"
	"fyne.io/fyne/test"
)

// NewSoftwareCanvas creates a new canvas in memory that can render without hardware support
func NewSoftwareCanvas() test.WindowlessCanvas {
	return test.NewCanvasWithPainter(software.NewPainter())
}

package playground

import (
	"fyne.io/fyne/v2/internal/driver/software"
	"fyne.io/fyne/v2/test"
)

// NewSoftwareCanvas creates a new canvas in memory that can render without hardware support
func NewSoftwareCanvas() test.WindowlessCanvas {
	return software.NewCanvas()
}

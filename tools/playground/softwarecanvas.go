package playground

import (
	"fyne.io/fyne/v2/internal/driver/software"
)

// NewSoftwareCanvas creates a new canvas in memory that can render without hardware support
func NewSoftwareCanvas() software.WindowlessCanvas {
	return software.NewCanvas()
}

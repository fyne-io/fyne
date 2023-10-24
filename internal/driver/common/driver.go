package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
)

const (
	// moved here from the platform driver packages to allow
	// referencing in widget package without an import cycle

	DesktopDoubleClickDelay = 300 // ms
	MobileDoubleClickDelay  = 500 // ms
)

// CanvasForObject returns the canvas for the specified object.
func CanvasForObject(obj fyne.CanvasObject) fyne.Canvas {
	return cache.GetCanvasForObject(obj)
}

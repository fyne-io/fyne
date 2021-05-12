package common

import (
	"sync"

	"fyne.io/fyne/v2"
)

var (
	canvasMutex sync.RWMutex
	canvases    = make(map[fyne.CanvasObject]fyne.Canvas)
)

// CanvasForObject returns the canvas for the specified object.
func CanvasForObject(obj fyne.CanvasObject) fyne.Canvas {
	canvasMutex.RLock()
	defer canvasMutex.RUnlock()
	return canvases[obj]
}

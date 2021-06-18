package cache

import (
	"sync"

	"fyne.io/fyne/v2"
)

var canvasesLock sync.RWMutex
var canvases = make(map[fyne.CanvasObject]*canvasInfo, 1024)

// GetCanvasForObject returns the canvas for the specified object.
func GetCanvasForObject(obj fyne.CanvasObject) fyne.Canvas {
	canvasesLock.RLock()
	cinfo, ok := canvases[obj]
	canvasesLock.RUnlock()
	if cinfo == nil || !ok {
		return nil
	}
	cinfo.setAlive()
	return cinfo.canvas
}

// SetCanvasForObject sets the canvas for the specified object.
func SetCanvasForObject(obj fyne.CanvasObject, canvas fyne.Canvas) {
	cinfo := &canvasInfo{canvas: canvas}
	cinfo.setAlive()
	canvasesLock.Lock()
	canvases[obj] = cinfo
	canvasesLock.Unlock()
}

type canvasInfo struct {
	expiringCache
	canvas fyne.Canvas
}

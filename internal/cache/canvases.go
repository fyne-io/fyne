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
// The passed function will be called if the item was not previously attached to this canvas
func SetCanvasForObject(obj fyne.CanvasObject, c fyne.Canvas, setup func()) {
	cinfo := &canvasInfo{canvas: c}
	cinfo.setAlive()
	canvasesLock.Lock()
	old, found := canvases[obj]
	canvases[obj] = cinfo
	canvasesLock.Unlock()

	if (!found || old.canvas != c) && setup != nil {
		setup()
	}
}

type canvasInfo struct {
	expiringCache
	canvas fyne.Canvas
}

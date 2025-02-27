package cache

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
)

var (
	canvases                 async.Map[fyne.CanvasObject, *canvasInfo]
	canvasCacheLastCleanSize int
	shouldCleanCanvases      bool
)

// GetCanvasForObject returns the canvas for the specified object.
func GetCanvasForObject(obj fyne.CanvasObject) fyne.Canvas {
	cinfo, ok := canvases.Load(obj)
	if cinfo == nil || !ok {
		return nil
	}
	cinfo.setAlive()
	return cinfo.canvas
}

// SetCanvasForObject sets the canvas for the specified object.
// The passed function will be called if the item was not previously attached to this canvas
func SetCanvasForObject(obj fyne.CanvasObject, c fyne.Canvas, setup func()) {
	old, found := canvases.Load(obj)
	needSetup := false
	if found {
		old.setAlive()
		if old.canvas != c {
			old.canvas = c
			needSetup = setup != nil
		}
	} else {
		needSetup = setup != nil
		cinfo := &canvasInfo{canvas: c}
		cinfo.setAlive()
		canvases.Store(obj, cinfo)
	}
	if needSetup {
		setup()
	}

	if canvases.Len() > 2*canvasCacheLastCleanSize {
		shouldCleanCanvases = true
	}
}

type canvasInfo struct {
	frameCounterCache
	canvas fyne.Canvas
}

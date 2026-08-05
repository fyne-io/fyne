package cache

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
)

var canvases async.Map[fyne.CanvasObject, *canvasInfo]

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
	// this runs for every visible object on every repaint, so avoid allocating
	// a canvasInfo we would only throw away - and keep the entry alive so a
	// constantly repainting window does not expire its own cache.
	if old, found := canvases.Load(obj); found && old.canvas == c {
		old.setAlive()
		return
	}

	// not cached, or moved to a different canvas - an object may only be on one
	// canvas at a time, so replace the entry and re-run setup for the new canvas.
	cinfo := &canvasInfo{canvas: c}
	cinfo.setAlive()
	canvases.Store(obj, cinfo)

	if setup != nil {
		setup()
	}
}

type canvasInfo struct {
	expiringCache
	canvas fyne.Canvas
}

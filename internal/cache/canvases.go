package cache

import (
	"sync/atomic"

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
	canvasEntry(obj, c, setup)
}

// CanvasRef holds an object's canvas registration between repaints.
type CanvasRef struct {
	info *canvasInfo
}

// SetCanvasForObjectRef sets the canvas for obj, skipping the cache lookup
// while ref still holds the current registration.
func SetCanvasForObjectRef(ref *CanvasRef, obj fyne.CanvasObject, c fyne.Canvas, setup func()) {
	if info := ref.info; info != nil && info.canvas == c && !info.replaced.Load() {
		info.setAlive()
		return
	}
	ref.info = canvasEntry(obj, c, setup)
}

func canvasEntry(obj fyne.CanvasObject, c fyne.Canvas, setup func()) *canvasInfo {
	// this runs for every visible object on every repaint, so avoid allocating
	// a canvasInfo we would only throw away - and keep the entry alive so a
	// constantly repainting window does not expire its own cache.
	old, found := canvases.Load(obj)
	if found && old.canvas == c {
		old.setAlive()
		return old
	}

	// not cached, or moved to a different canvas - an object may only be on one
	// canvas at a time, so replace the entry and re-run setup for the new canvas.
	if found {
		old.replaced.Store(true)
	}
	cinfo := &canvasInfo{canvas: c}
	cinfo.setAlive()
	canvases.Store(obj, cinfo)

	if setup != nil {
		setup()
	}
	return cinfo
}

type canvasInfo struct {
	expiringCache
	canvas   fyne.Canvas
	replaced atomic.Bool
}

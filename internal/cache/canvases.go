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

// CanvasRef is an object's registration with a canvas, kept by whoever
// registers the object on every repaint - a render cache node - so that
// SetCanvasForObjectRef skips the lookup while it is the registration the
// cache holds. The zero value holds none.
type CanvasRef struct {
	info *canvasInfo
}

// SetCanvasForObjectRef is SetCanvasForObject for a caller that keeps ref for
// obj between repaints: while ref holds the entry the cache has for obj, for
// this canvas, the entry is just kept alive; once it was replaced, by the
// object registering with another canvas, or cleaned, obj registers anew.
func SetCanvasForObjectRef(ref *CanvasRef, obj fyne.CanvasObject, c fyne.Canvas, setup func()) {
	if info := ref.info; info != nil && info.canvas == c && !info.replaced.Load() {
		info.setAlive()
		return
	}
	ref.info = canvasEntry(obj, c, setup)
}

// canvasEntry is SetCanvasForObject returning the entry obj has after.
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
	canvas fyne.Canvas
	// replaced is set once the entry is no longer the one the cache holds
	// for its object, so a CanvasRef still holding it registers anew.
	replaced atomic.Bool
}

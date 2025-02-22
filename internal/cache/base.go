package cache

import (
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/common"
)

var (
	ValidDuration     = 1 * time.Minute
	cleanTaskInterval = ValidDuration / 2

	lastClean                     time.Time
	skippedCleanWithCanvasRefresh = false

	// testing purpose only
	timeNow = time.Now
)

func init() {
	if t, err := time.ParseDuration(os.Getenv("FYNE_CACHE")); err == nil {
		ValidDuration = t
		cleanTaskInterval = ValidDuration / 2
	}
}

// ShouldClean returns whether the clean tasks (CleanTextTextures and Clean)
// should be invoked during the current iteration of the main loop.
func ShouldClean() bool {
	return shouldFullClean() || shouldCleanTextTextures
}

// CleanTextTextures runs the per-canvas cache clean text for the text texture cache.
// It should be run with a current GL context for each existing canvas
// during the main run loop if ShouldClean() returns true.
// Within the same iteration, Clean should be run once (not per-window)
// to clean the non-texture caches and mark the clean as complete.
func CleanTextTextures(canvas fyne.Canvas) {
	if shouldFullClean() || shouldCleanTextTextures {
		cleanCanvasTextTextureCache(canvas)
	}
}

// Clean runs the non-texture cache clean task, it should be run during the
// main loop, along with CleanTextTextures, if ShouldClean() returns true.
func Clean(canvasRefreshed bool) {
	now := time.Now()
	full := shouldFullClean()

	if full {
		destroyExpiredSvgs(now)
		destroyExpiredFontMetrics(now)
		if canvasRefreshed {
			// Destroy renderers on canvas refresh to avoid flickering screen.
			destroyExpiredRenderers(now)
			// canvases cache should be invalidated only on canvas refresh, otherwise there wouldn't
			// be a way to recover them later
			destroyExpiredCanvases(now)
		}
	}

	textTextureLastCleanSize = textTextures.Len()
	shouldCleanTextTextures = false
	if full {
		lastClean = timeNow()
	}
}

// CleanCanvas performs a complete remove of all the objects that belong to the specified
// canvas. Usually used to free all objects from a closing windows.
func CleanCanvas(canvas fyne.Canvas) {
	canvases.Range(func(obj fyne.CanvasObject, cinfo *canvasInfo) bool {
		if cinfo.canvas != canvas {
			return true
		}

		canvases.Delete(obj)

		wid, ok := obj.(fyne.Widget)
		if !ok {
			return true
		}
		rinfo, ok := renderers.LoadAndDelete(wid)
		if !ok {
			return true
		}
		rinfo.renderer.Destroy()
		overrides.Delete(wid)
		return true
	})
}

// ResetThemeCaches clears all the svg and text size cache maps
func ResetThemeCaches() {
	svgs.Clear()
	fontSizeCache.Clear()
}

// destroyExpiredCanvases deletes objects from the canvases cache.
func destroyExpiredCanvases(now time.Time) {
	canvases.Range(func(obj fyne.CanvasObject, cinfo *canvasInfo) bool {
		if cinfo.isExpired(now) {
			canvases.Delete(obj)
		}
		return true
	})
}

// destroyExpiredRenderers deletes the renderer from the cache and calls
// renderer.Destroy()
func destroyExpiredRenderers(now time.Time) {
	renderers.Range(func(wid fyne.Widget, rinfo *rendererInfo) bool {
		if rinfo.isExpired(now) {
			rinfo.renderer.Destroy()
			overrides.Delete(wid)
			renderers.Delete(wid)
		}
		return true
	})
}

func shouldFullClean() bool {
	return lastClean.Add(ValidDuration).Before(timeNow())
}

type expiringCache struct {
	expires time.Time
}

type frameCounterCache struct {
	lastAccessedFrame uint64
}

// setAlive updates expiration time.
func (c *frameCounterCache) setAlive() {
	c.lastAccessedFrame = common.CurrentFrameCounter()
}

// isExpired check if the cache data is expired.
func (c *frameCounterCache) isExpired(now time.Time) bool {
	return c.lastAccessedFrame < common.CurrentFrameCounter()-1
}

// isExpired check if the cache data is expired.
func (c *expiringCache) isExpired(now time.Time) bool {
	return c.expires.Before(now)
}

// setAlive updates expiration time.
func (c *expiringCache) setAlive() {
	c.expires = timeNow().Add(ValidDuration)
}

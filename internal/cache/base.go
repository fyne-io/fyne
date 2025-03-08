package cache

import (
	"os"
	"time"

	"fyne.io/fyne/v2"
)

var (
	ValidDuration = 2 * time.Minute

	lastClean                     time.Time
	skippedCleanWithCanvasRefresh = false

	framecounter uint64 = 1

	// testing purpose only
	timeNow = time.Now
)

func init() {
	if t, err := time.ParseDuration(os.Getenv("FYNE_CACHE")); err == nil {
		ValidDuration = t
	}
}

// IncrementFrameCounter increments the current frame counter
// that is used to track which cached objects were accessed in the last frame.
// It should be called at the end of each iteration of the main loop.
func IncrementFrameCounter() {
	framecounter += 1
}

// ShouldClean returns whether the clean tasks (CleanTextTextures and Clean)
// should be invoked during the current iteration of the main loop.
func ShouldClean() bool {
	return shouldCleanTextTextures ||
		shouldCleanObjectTextures ||
		shouldCleanRenderers ||
		shouldCleanCanvases ||
		shouldCleanFontSizeCache ||
		shouldFullClean()
}

// ShouldClean returns whether the clean tasks (CleanTextTextures and Clean)
// should be invoked during the current iteration of the main loop,
// AND the clean task will include a clean of the CanvasForObject map.
// If so, the driver should be sure to mark all objects -
// both visible and invisible, as alive before invoking the clean tasks.
func ShouldCleanCanvases() bool {
	return shouldCleanCanvases || shouldFullClean()
}

// CleanTextures runs the per-canvas cache clean text for the texture caches.
// It should be run with a current GL context for each existing canvas
// during the main run loop if ShouldClean() returns true.
// Within the same iteration, Clean should be run once (not per-window)
// after CleanTextures has been called for all canvases
// to clean the non-texture caches and mark the texture caches clean as complete.
func CleanTextures(canvas fyne.Canvas, texFree func(fyne.CanvasObject)) {
	full := shouldFullClean()
	if full || shouldCleanTextTextures {
		cleanTextTextureCache(canvas)
	}
	if (full || shouldCleanObjectTextures) && texFree != nil {
		rangeExpiredTexturesFor(canvas, texFree)
	}
}

// Clean runs the non-texture cache clean task, if ShouldClean() returns true,
// it should be run during the main loop, after having called
// CleanTextures for each canvas.
func Clean() {
	now := time.Now()
	full := shouldFullClean()

	if full {
		destroyExpiredSvgs(now)
		destroyExpiredFontMetrics(now)
	}
	if full || shouldCleanRenderers {
		destroyExpiredRenderers(now)
		rendererCacheLastCleanSize = renderers.Len()
		shouldCleanRenderers = false
	}
	if full || shouldCleanCanvases {
		destroyExpiredCanvases(now)
		canvasCacheLastCleanSize = canvases.Len()
		shouldCleanCanvases = false
	}
	if full || shouldCleanFontSizeCache {
		destroyExpiredFontMetrics(now)
		fontSizeCacheLastCleanSize = fontSizeCache.Len()
		shouldCleanFontSizeCache = false
	}

	// CleanTextures should have been called for each canvas
	// update the sizes of the texture caches
	if full || shouldCleanTextTextures {
		textTextureLastCleanSize = textTextures.Len()
		shouldCleanTextTextures = false
	}
	if full || shouldCleanObjectTextures {
		objectTexturesLastCleanSize = objectTextures.Len()
		shouldCleanObjectTextures = false
	}

	if full {
		lastClean = now
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
	c.lastAccessedFrame = framecounter
}

// isExpired check if the cache data is expired.
func (c *frameCounterCache) isExpired(now time.Time) bool {
	return c.lastAccessedFrame < framecounter
}

// isExpired check if the cache data is expired.
func (c *expiringCache) isExpired(now time.Time) bool {
	return c.expires.Before(now)
}

// setAlive updates expiration time.
func (c *expiringCache) setAlive() {
	c.expires = timeNow().Add(ValidDuration)
}

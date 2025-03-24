package cache

import (
	"os"
	"time"

	"fyne.io/fyne/v2"
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

// Clean run cache clean task, it should be called on paint events.
func Clean(canvasRefreshed bool) {
	now := timeNow()
	// do not run clean task too fast
	if now.Sub(lastClean) < 10*time.Second {
		if canvasRefreshed {
			skippedCleanWithCanvasRefresh = true
		}
		return
	}
	if skippedCleanWithCanvasRefresh {
		skippedCleanWithCanvasRefresh = false
		canvasRefreshed = true
	}
	if !canvasRefreshed && now.Sub(lastClean) < cleanTaskInterval {
		return
	}
	destroyExpiredSvgs(now)
	destroyExpiredFontMetrics(now)
	if canvasRefreshed {
		// Destroy renderers on canvas refresh to avoid flickering screen.
		destroyExpiredRenderers(now)
		// canvases cache should be invalidated only on canvas refresh, otherwise there wouldn't
		// be a way to recover them later
		destroyExpiredCanvases(now)
	}
	lastClean = timeNow()
}

// CleanCanvas performs a complete remove of all the objects that belong to the specified
// canvas. Usually used to free all objects from a closing windows.
func CleanCanvas(canvas fyne.Canvas) {
	canvasesLock.Lock()
	defer canvasesLock.Unlock()

	for obj, cinfo := range canvases {
		if cinfo.canvas != canvas {
			continue
		}

		delete(canvases, obj)

		wid, ok := obj.(fyne.Widget)
		if !ok {
			continue
		}

		renderersLock.Lock()
		rinfo, ok := renderers[wid]
		if !ok {
			renderersLock.Unlock()
			continue
		}

		delete(renderers, wid)
		renderersLock.Unlock()

		rinfo.renderer.Destroy()

		overridesLock.Lock()
		delete(overrides, wid)
		overridesLock.Unlock()
	}
}

// ResetThemeCaches clears all the svg and text size cache maps
func ResetThemeCaches() {
	canvasesLock.Lock()
	for k := range svgs {
		delete(svgs, k)
	}
	canvasesLock.Unlock()

	fontSizeLock.Lock()
	for k := range fontSizeCache {
		delete(fontSizeCache, k)
	}
	fontSizeLock.Unlock()

}

// destroyExpiredCanvases deletes objects from the canvases cache.
func destroyExpiredCanvases(now time.Time) {
	canvasesLock.Lock()
	defer canvasesLock.Unlock()

	for obj, cinfo := range canvases {
		if cinfo.isExpired(now) {
			delete(canvases, obj)
		}
	}
}

// destroyExpiredRenderers deletes the renderer from the cache and calls
// renderer.Destroy()
func destroyExpiredRenderers(now time.Time) {
	renderersLock.Lock()
	defer renderersLock.Unlock()

	for wid, rinfo := range renderers {
		if rinfo.isExpired(now) {
			rinfo.renderer.Destroy()

			overridesLock.Lock()
			delete(overrides, wid)
			overridesLock.Unlock()

			delete(renderers, wid)
		}
	}
}

type expiringCache struct {
	expires time.Time
}

// isExpired check if the cache data is expired.
func (c *expiringCache) isExpired(now time.Time) bool {
	return c.expires.Before(now)
}

// setAlive updates expiration time.
func (c *expiringCache) setAlive() {
	c.expires = timeNow().Add(ValidDuration)
}

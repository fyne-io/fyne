package cache

import (
	"os"
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

var (
	cacheDuration     = 1 * time.Minute
	cleanTaskInterval = cacheDuration / 2

	expiredObjects                = make([]fyne.CanvasObject, 0, 50)
	lastClean                     time.Time
	skippedCleanWithCanvasRefresh = false

	// testing purpose only
	timeNow func() time.Time = time.Now
)

func init() {
	if t, err := time.ParseDuration(os.Getenv("FYNE_CACHE")); err == nil {
		cacheDuration = t
		cleanTaskInterval = cacheDuration / 2
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
	deletingObjs := make([]fyne.CanvasObject, 0, 50)

	canvasesLock.RLock()
	for obj, cinfo := range canvases {
		if cinfo.canvas == canvas {
			deletingObjs = append(deletingObjs, obj)
		}
	}
	canvasesLock.RUnlock()
	if len(deletingObjs) == 0 {
		return
	}

	canvasesLock.Lock()
	for _, dobj := range deletingObjs {
		delete(canvases, dobj)
	}
	canvasesLock.Unlock()

	renderersLock.Lock()
	for _, dobj := range deletingObjs {
		wid, ok := dobj.(fyne.Widget)
		if !ok {
			continue
		}
		winfo, ok := renderers[wid]
		if !ok {
			continue
		}
		winfo.renderer.Destroy()
		delete(renderers, wid)
	}
	renderersLock.Unlock()
}

// ResetThemeCaches clears all the svg and text size cache maps
func ResetThemeCaches() {
	svgLock.Lock()
	svgs = make(map[string]*svgInfo)
	svgLock.Unlock()

	fontSizeLock.Lock()
	fontSizeCache = map[fontSizeEntry]fontMetric{}
	fontSizeLock.Unlock()
}

// destroyExpiredCanvases deletes objects from the canvases cache.
func destroyExpiredCanvases(now time.Time) {
	expiredObjects = expiredObjects[:0]
	canvasesLock.RLock()
	for obj, cinfo := range canvases {
		if cinfo.isExpired(now) {
			expiredObjects = append(expiredObjects, obj)
		}
	}
	canvasesLock.RUnlock()
	if len(expiredObjects) > 0 {
		canvasesLock.Lock()
		for i, exp := range expiredObjects {
			delete(canvases, exp)
			expiredObjects[i] = nil
		}
		canvasesLock.Unlock()
	}
}

// destroyExpiredRenderers deletes the renderer from the cache and calls
// renderer.Destroy()
func destroyExpiredRenderers(now time.Time) {
	expiredObjects = expiredObjects[:0]
	renderersLock.RLock()
	for wid, rinfo := range renderers {
		if rinfo.isExpired(now) {
			rinfo.renderer.Destroy()
			expiredObjects = append(expiredObjects, wid)
		}
	}
	renderersLock.RUnlock()
	if len(expiredObjects) > 0 {
		renderersLock.Lock()
		for i, exp := range expiredObjects {
			delete(renderers, exp.(fyne.Widget))
			expiredObjects[i] = nil
		}
		renderersLock.Unlock()
	}
}

type expiringCache struct {
	expireLock sync.RWMutex
	expires    time.Time
}

// isExpired check if the cache data is expired.
func (c *expiringCache) isExpired(now time.Time) bool {
	c.expireLock.RLock()
	defer c.expireLock.RUnlock()
	return c.expires.Before(now)
}

// setAlive updates expiration time.
func (c *expiringCache) setAlive() {
	t := timeNow().Add(cacheDuration)
	c.expireLock.Lock()
	c.expires = t
	c.expireLock.Unlock()
}

type expiringCacheNoLock struct {
	expires time.Time
}

// isExpired check if the cache data is expired.
func (c *expiringCacheNoLock) isExpired(now time.Time) bool {
	return c.expires.Before(now)
}

// setAlive updates expiration time.
func (c *expiringCacheNoLock) setAlive() {
	t := timeNow().Add(cacheDuration)
	c.expires = t
}

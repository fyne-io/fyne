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

	canvasRefreshCh = make(chan struct{}, 1)
	expiredObjects  = make([]fyne.CanvasObject, 0, 50)
	lastClean       time.Time
)

func init() {
	if t, err := time.ParseDuration(os.Getenv("FYNE_CACHE")); err == nil {
		cacheDuration = t
		cleanTaskInterval = cacheDuration / 2
	}
}

// CleanTask run cache clean task, it should be called on paint events.
func CleanTask() {
	now := time.Now()
	// do not run clean task so fast
	if now.Sub(lastClean) < 10*time.Second {
		return
	}
	canvasRefresh := false
	select {
	case <-canvasRefreshCh:
		canvasRefresh = true
	default:
		if now.Sub(lastClean) < cleanTaskInterval {
			return
		}
	}
	destroyExpiredSvgs(now)
	if canvasRefresh {
		// Destroy renderers on canvas refresh to avoid flickering screen.
		destroyExpiredRenderers(now)
		// canvases cache should be invalidated only on canvas refresh, otherwise there wouldn't
		// be a way to recover them later
		destroyExpiredCanvases(now)
	}
	lastClean = time.Now()
}

// ForceCleanFor forces a complete remove of all the objects that belong to the specified
// canvas. Usually used to free all objects from a closing windows.
func ForceCleanFor(canvas fyne.Canvas) {
	deletingObjs := make([]fyne.CanvasObject, 0, 50)

	// find all objects that belong to the specified canvas
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

	// remove them from canvas cache
	canvasesLock.Lock()
	for _, dobj := range deletingObjs {
		delete(canvases, dobj)
	}
	canvasesLock.Unlock()

	// destroy their renderers and delete them from the renderer
	// cache if they are widgets
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

// NotifyCanvasRefresh notifies to the caches that a canvas refresh was triggered.
func NotifyCanvasRefresh() {
	select {
	case canvasRefreshCh <- struct{}{}:
	default:
		return
	}
}

// ===============================================================
// Private functions
// ===============================================================

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

// destroyExpiredSvgs destroys expired svgs cache data.
func destroyExpiredSvgs(now time.Time) {
	expiredSvgs := make([]string, 0, 20)
	svgLock.RLock()
	for s, sinfo := range svgs {
		if sinfo.isExpired(now) {
			expiredSvgs = append(expiredSvgs, s)
		}
	}
	svgLock.RUnlock()
	if len(expiredSvgs) > 0 {
		svgLock.Lock()
		for _, exp := range expiredSvgs {
			delete(svgs, exp)
		}
		svgLock.Unlock()
	}
}

// ===============================================================
// Private types
// ===============================================================

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
	t := time.Now().Add(cacheDuration)
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
	t := time.Now().Add(cacheDuration)
	c.expires = t
}

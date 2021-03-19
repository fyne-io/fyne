package cache

import (
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

var cleanCh = make(chan struct{}, 1)
var cleanTaskOnce sync.Once

func init() {

	cleanTaskOnce.Do(func() {

		go func() {
			expired := make([]fyne.CanvasObject, 0, 50)

			cleanTask := func(excludeCanvases bool) {
				now := time.Now()
				// -- Renderers clean task
				expired = expired[:0]
				renderersLock.RLock()
				for wid, rinfo := range renderers {
					if rinfo.isExpired(now) {
						rinfo.renderer.Destroy()
						expired = append(expired, wid)
					}
				}
				renderersLock.RUnlock()
				if len(expired) > 0 {
					renderersLock.Lock()
					for i, exp := range expired {
						delete(renderers, exp.(fyne.Widget))
						expired[i] = nil
					}
					renderersLock.Unlock()
				}

				// -- Textures clean task
				// TODO find a way to clean textures here

				if excludeCanvases {
					return
				}

				// -- Canvases clean task
				expired = expired[:0]
				canvasesLock.RLock()
				for obj, cinfo := range canvases {
					if cinfo.isExpired(now) {
						expired = append(expired, obj)
					}
				}
				canvasesLock.RUnlock()
				if len(expired) > 0 {
					canvasesLock.Lock()
					for i, exp := range expired {
						delete(canvases, exp)
						expired[i] = nil
					}
					canvasesLock.Unlock()
				}
			}

			t := time.NewTicker(time.Minute)
			for {
				select {
				case <-cleanCh:
					cleanTask(false)
					// do not trigger another clean task so fast
					time.Sleep(10 * time.Second)
					t.Reset(time.Minute)
				case <-t.C:
					// canvases cache can't be invalidated using the ticker
					// because if we do it, there wouldn't be a way to recover them later
					cleanTask(true)
				}
			}
		}()
	})

}

// Clean triggers clean cache task.
func Clean() {
	select {
	case cleanCh <- struct{}{}:
	default:
		return
	}
}

// ===============================================================
// Privates
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
	t := time.Now().Add(1 * time.Minute)
	c.expireLock.Lock()
	c.expires = t
	c.expireLock.Unlock()
}

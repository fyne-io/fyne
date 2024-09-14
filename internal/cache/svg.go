package cache

import (
	"image"
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

var svgs = &sync.Map{} // make(map[string]*svgInfo)

// GetSvg gets svg image from cache if it exists.
func GetSvg(name string, o fyne.CanvasObject, w int, h int) *image.NRGBA {
	sinfo, ok := svgs.Load(overriddenName(name, o))
	if !ok || sinfo == nil {
		return nil
	}
	svginfo := sinfo.(*svgInfo)
	if svginfo.w != w || svginfo.h != h {
		return nil
	}

	svginfo.setAlive()
	return svginfo.pix
}

// SetSvg sets a svg into the cache map.
func SetSvg(name string, o fyne.CanvasObject, pix *image.NRGBA, w int, h int) {
	sinfo := &svgInfo{
		pix: pix,
		w:   w,
		h:   h,
	}
	sinfo.setAlive()
	svgs.Store(overriddenName(name, o), sinfo)
}

type svgInfo struct {
	// An svgInfo can be accessed from different goroutines, e.g., systray.
	// Use expiringCache instead of expiringCacheNoLock.
	expiringCache
	pix  *image.NRGBA
	w, h int
}

// destroyExpiredSvgs destroys expired svgs cache data.
func destroyExpiredSvgs(now time.Time) {
	svgs.Range(func(key, value any) bool {
		sinfo := value.(*svgInfo)
		if sinfo.isExpired(now) {
			svgs.Delete(key)
		}
		return true
	})
}

func overriddenName(name string, o fyne.CanvasObject) string {
	if o != nil { // for overridden themes get the cache key right
		if over, ok := overrides.Load(o); ok {
			return over.(*overrideScope).cacheID + name
		}
	}

	return name
}

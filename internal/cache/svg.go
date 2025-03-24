package cache

import (
	"image"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async/migration"
)

var (
	svgs     = make(map[string]*svgInfo)
	svgsLock migration.Mutex
)

// GetSvg gets svg image from cache if it exists.
func GetSvg(name string, o fyne.CanvasObject, w int, h int) *image.NRGBA {
	svgsLock.Lock()
	defer svgsLock.Unlock()

	svginfo, ok := svgs[overriddenName(name, o)]
	if !ok || svginfo == nil {
		return nil
	}

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

	svgsLock.Lock()
	defer svgsLock.Unlock()
	svgs[overriddenName(name, o)] = sinfo
}

type svgInfo struct {
	expiringCache
	pix  *image.NRGBA
	w, h int
}

// destroyExpiredSvgs destroys expired svgs cache data.
func destroyExpiredSvgs(now time.Time) {
	svgsLock.Lock()
	defer svgsLock.Unlock()

	for key, sinfo := range svgs {
		if sinfo.isExpired(now) {
			delete(svgs, key)
		}
	}
}

func overriddenName(name string, o fyne.CanvasObject) string {
	if o != nil { // for overridden themes get the cache key right
		overridesLock.Lock()
		defer overridesLock.Unlock()

		if over, ok := overrides[o]; ok {
			return over.cacheID + name
		}
	}

	return name
}

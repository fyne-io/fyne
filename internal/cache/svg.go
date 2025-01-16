package cache

import (
	"image"
	"time"

	"fyne.io/fyne/v2"
)

var svgs = make(map[string]*svgInfo)

// GetSvg gets svg image from cache if it exists.
func GetSvg(name string, o fyne.CanvasObject, w int, h int) *image.NRGBA {
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
	svgs[overriddenName(name, o)] = sinfo
}

type svgInfo struct {
	expiringCache
	pix  *image.NRGBA
	w, h int
}

// destroyExpiredSvgs destroys expired svgs cache data.
func destroyExpiredSvgs(now time.Time) {
	for key, sinfo := range svgs {
		if sinfo.isExpired(now) {
			delete(svgs, key)
		}
	}
}

func overriddenName(name string, o fyne.CanvasObject) string {
	if o != nil { // for overridden themes get the cache key right
		if over, ok := overrides[o]; ok {
			return over.cacheID + name
		}
	}

	return name
}

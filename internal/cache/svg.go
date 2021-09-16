package cache

import (
	"image"
	"sync"
	"time"
)

var svgLock sync.RWMutex
var svgs = make(map[string]*svgInfo)

// GetSvg gets svg image from cache if it exists.
func GetSvg(name string, w int, h int) *image.NRGBA {
	svgLock.RLock()
	sinfo, ok := svgs[name]
	svgLock.RUnlock()
	if !ok || sinfo == nil || sinfo.w != w || sinfo.h != h {
		return nil
	}
	sinfo.setAlive()
	return sinfo.pix
}

// SetSvg sets a svg into the cache map.
func SetSvg(name string, pix *image.NRGBA, w int, h int) {
	sinfo := &svgInfo{
		pix: pix,
		w:   w,
		h:   h,
	}
	sinfo.setAlive()
	svgLock.Lock()
	svgs[name] = sinfo
	svgLock.Unlock()
}

type svgInfo struct {
	expiringCacheNoLock
	pix  *image.NRGBA
	w, h int
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

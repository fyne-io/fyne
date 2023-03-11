package cache

import (
	"image"
	"sync"
	"time"
)

var svgs = &sync.Map{} // make(map[string]*svgInfo)

// GetSvg gets svg image from cache if it exists.
func GetSvg(name string, w int, h int) *image.NRGBA {
	sinfo, ok := svgs.Load(name)
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
func SetSvg(name string, pix *image.NRGBA, w int, h int) {
	sinfo := &svgInfo{
		pix: pix,
		w:   w,
		h:   h,
	}
	sinfo.setAlive()
	svgs.Store(name, sinfo)
}

type svgInfo struct {
	expiringCacheNoLock
	pix  *image.NRGBA
	w, h int
}

// destroyExpiredSvgs destroys expired svgs cache data.
func destroyExpiredSvgs(now time.Time) {
	expiredSvgs := make([]string, 0, 20)
	svgs.Range(func(key, value interface{}) bool {
		s, sinfo := key.(string), value.(*svgInfo)
		if sinfo.isExpired(now) {
			expiredSvgs = append(expiredSvgs, s)
		}
		return true
	})

	for _, exp := range expiredSvgs {
		svgs.Delete(exp)
	}
}

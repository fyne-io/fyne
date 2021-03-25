package cache

import (
	"image"
	"sync"
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

// ResetSvg clears all the svg cache map
func ResetSvg() {
	svgLock.Lock()
	svgs = make(map[string]*svgInfo)
	svgLock.Unlock()
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

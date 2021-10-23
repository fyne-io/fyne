package cache

import (
	"image"
	"sync"
	"time"
)

var imgLock sync.RWMutex
var imgs = make(map[string]*imgInfo)

// GetImg gets img image from cache if it exists.
func GetImg(name string, w int, h int) *image.NRGBA {
	imgLock.RLock()
	sinfo, ok := imgs[name]
	imgLock.RUnlock()
	if !ok || sinfo == nil || sinfo.w != w || sinfo.h != h {
		return nil
	}
	sinfo.setAlive()
	return sinfo.pix
}

// SetImg sets a img into the cache map.
func SetImg(name string, pix *image.NRGBA, w int, h int) {
	sinfo := &imgInfo{
		pix: pix,
		w:   w,
		h:   h,
	}
	sinfo.setAlive()
	imgLock.Lock()
	imgs[name] = sinfo
	imgLock.Unlock()
}

type imgInfo struct {
	expiringCacheNoLock
	pix  *image.NRGBA
	w, h int
}

// destroyExpiredImgs destroys expired imgs cache data.
func destroyExpiredImgs(now time.Time) {
	expiredImgs := make([]string, 0, 20)
	imgLock.RLock()
	for s, sinfo := range imgs {
		if sinfo.isExpired(now) {
			expiredImgs = append(expiredImgs, s)
		}
	}
	imgLock.RUnlock()
	if len(expiredImgs) > 0 {
		imgLock.Lock()
		for _, exp := range expiredImgs {
			delete(imgs, exp)
		}
		imgLock.Unlock()
	}
}

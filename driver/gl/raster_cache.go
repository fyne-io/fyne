package gl

import (
	"image"
	"os"
	"sync"
	"time"

	"fyne.io/fyne"
)

type rasterInfo struct {
	pix     *image.RGBA
	w, h    int
	expires time.Time
}

var cacheDuration = (time.Minute * 5)
var rasters = make(map[fyne.Resource]*rasterInfo)
var rasterMutex sync.RWMutex

func init() {
	if t, err := time.ParseDuration(os.Getenv("FYNE_CACHE")); err == nil {
		cacheDuration = t
	}

	go func() {
		delay := cacheDuration / 2
		if delay < time.Second {
			delay = time.Second
		}
		for {
			time.Sleep(delay)
			now := time.Now()
			rasterMutex.Lock()
			for k, v := range rasters {
				if v.expires.Before(now) {
					delete(rasters, k)
				}
			}
			rasterMutex.Unlock()
		}
	}()
}

func svgCacheGet(res fyne.Resource, w int, h int) *image.RGBA {
	rasterMutex.RLock()
	v, ok := rasters[res]
	if !ok || v == nil || v.w != w || v.h != h {
		rasterMutex.RUnlock()
		return nil
	}
	v.expires = time.Now().Add(cacheDuration)
	return v.pix
}

func svgCachePut(res fyne.Resource, pix *image.RGBA, w int, h int) {
	rasterMutex.Lock()
	rasters[res] = &rasterInfo{
		pix:     pix,
		w:       w,
		h:       h,
		expires: time.Now().Add(cacheDuration),
	}
	rasterMutex.Unlock()
}

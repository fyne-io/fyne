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

var cacheDuration = time.Minute * 5
var rasters = make(map[fyne.Resource]*rasterInfo)
var aspects = make(map[interface{}]float32, 16)
var rasterMutex sync.RWMutex
var janitorOnce sync.Once

func init() {
	if t, err := time.ParseDuration(os.Getenv("FYNE_CACHE")); err == nil {
		cacheDuration = t
	}

	svgCacheJanitor()
}

func svgCacheJanitor() {
	delay := cacheDuration / 2
	if delay < time.Second {
		delay = time.Second
	}

	go janitorOnce.Do(func() {
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
	})
}

func svgCacheGet(res fyne.Resource, w int, h int) *image.RGBA {
	rasterMutex.RLock()
	defer rasterMutex.RUnlock()
	v, ok := rasters[res]
	if !ok || v == nil || v.w != w || v.h != h {
		return nil
	}
	v.expires = time.Now().Add(cacheDuration)
	return v.pix
}

func svgCachePut(res fyne.Resource, pix *image.RGBA, w int, h int) {
	rasterMutex.Lock()
	defer rasterMutex.Unlock()
	defer func() {
		recover()
	}()

	rasters[res] = &rasterInfo{
		pix:     pix,
		w:       w,
		h:       h,
		expires: time.Now().Add(cacheDuration),
	}
}

func svgCacheReset() {
	rasterMutex.Lock()
	for k := range rasters {
		delete(rasters, k)
	}
	rasterMutex.Unlock()
}

func svgCacheMonitorTheme() {
	listener := make(chan fyne.Settings)
	fyne.CurrentApp().Settings().AddChangeListener(listener)
	go func() {
		for {
			<-listener
			svgCacheReset()
		}
	}()

}

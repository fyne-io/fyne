package painter

import (
	"image"
	"os"
	"sync"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
)

type rasterInfo struct {
	pix     *image.NRGBA
	w, h    int
	expires time.Time
}

var cacheDuration = time.Minute * 5
var rasters = make(map[string]*rasterInfo)
var aspects = make(map[interface{}]float32, 16)
var rasterMutex sync.RWMutex
var janitorOnce sync.Once

func init() {
	if t, err := time.ParseDuration(os.Getenv("FYNE_CACHE")); err == nil {
		cacheDuration = t
	}

	svgCacheJanitor()
}

// GetAspect looks up an aspect ratio of an image
func GetAspect(img *canvas.Image) float32 {
	aspect := aspects[img.Resource]
	if aspect == 0 {
		aspect = aspects[img]
	}

	return aspect
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

func svgCacheGet(name string, w int, h int) *image.NRGBA {
	rasterMutex.RLock()
	defer rasterMutex.RUnlock()
	v, ok := rasters[name]
	if !ok || v == nil || v.w != w || v.h != h {
		return nil
	}
	v.expires = time.Now().Add(cacheDuration)
	return v.pix
}

func svgCachePut(name string, pix *image.NRGBA, w int, h int) {
	rasterMutex.Lock()
	defer rasterMutex.Unlock()
	defer func() {
		recover()
	}()

	rasters[name] = &rasterInfo{
		pix:     pix,
		w:       w,
		h:       h,
		expires: time.Now().Add(cacheDuration),
	}
}

// SvgCacheReset clears the SVG cache.
func SvgCacheReset() {
	rasterMutex.Lock()
	rasters = make(map[string]*rasterInfo)
	rasterMutex.Unlock()
}

// SvgCacheMonitorTheme hooks up the SVG cache to listen for theme changes and resets the cache
func SvgCacheMonitorTheme() {
	listener := make(chan fyne.Settings)
	fyne.CurrentApp().Settings().AddChangeListener(listener)
	go func() {
		for {
			<-listener
			SvgCacheReset()
		}
	}()

}

package gl

import (
	"image"
	"sync"
	"time"

	"fyne.io/fyne"
	"github.com/steveoc64/memdebug"
)

type rasterInfo struct {
	pix     *image.RGBA
	w, h    int
	alpha   float64
	expires time.Time
}

var cacheDuration = (time.Minute * 5)
var rasters map[fyne.Resource]*rasterInfo
var rasterMutex sync.RWMutex

// CacheHoldDuration returns the cache hold duration (default 5 mins)
func CacheDuration() time.Duration {
	rasterMutex.RLock()
	defer rasterMutex.RUnlock()
	return cacheDuration
}

// SetCacheHoldDuration sets the new cache hold duration
func SetCacheDuration(t time.Duration) {
	rasterMutex.Lock()
	defer rasterMutex.Unlock()
	memdebug.Print(time.Now(), "cache duration set to", t)
	cacheDuration = t
}

func init() {
	memdebug.Print(time.Now(), "init rasters")
	rasters = make(map[fyne.Resource]*rasterInfo)

	janitor := func() {
		for {
			time.Sleep(time.Minute)
			now := time.Now()
			memdebug.Print(now, "check exp")
			rasterMutex.Lock()
			for k, v := range rasters {
				if v.expires.Before(now) {
					delete(rasters, k)
					memdebug.Print(time.Now(), "expires", v.expires)
				}
			}
			rasterMutex.Unlock()
		}
	}

	go janitor()

}

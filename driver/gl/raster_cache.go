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
	alpha   float64
	expires time.Time
}

var cacheDuration = (time.Minute * 5)
var rasters = make(map[fyne.Resource]*rasterInfo)
var rasterMutex sync.RWMutex

func init() {
	if t, err := time.ParseDuration(os.Getenv("FYNE_CACHE")); err == nil {
		cacheDuration = t
	}

	janitor := func() {
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
	}

	go janitor()

}

package gl

import (
	"image"
	"os"
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

func init() {
	memdebug.Print(time.Now(), "init rasters")
	rasters = make(map[fyne.Resource]*rasterInfo)

	if t, err := time.ParseDuration(os.Getenv("FYNE_CACHE")); err == nil {
		t1 := time.Now()
		memdebug.Print(t1, "parsed duration", t)
		cacheDuration = t
	}

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

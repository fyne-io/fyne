package cache

import (
	"image"
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
)

var (
	svgs async.Map[string, *svgInfo]

	// not an async.Map, ThemedResource.Content() can be called from any goroutine
	colorizedSvgsLock sync.Mutex
	colorizedSvgs     = map[colorizedSvgKey]*colorizedSvgInfo{}
)

// GetColorizedSvg gets the colorized version of an svg from cache if it exists.
func GetColorizedSvg(src []byte, clr color.NRGBA) ([]byte, bool) {
	colorizedSvgsLock.Lock()
	defer colorizedSvgsLock.Unlock()
	info, ok := colorizedSvgs[colorizedSvgKey{string(src), clr}]
	if !ok {
		return nil, false
	}

	info.setAlive()
	return info.content, true
}

// GetSvg gets svg image from cache if it exists.
func GetSvg(name string, o fyne.CanvasObject, w int, h int) *image.NRGBA {
	svginfo, ok := svgs.Load(overriddenName(name, o))
	if !ok || svginfo == nil {
		return nil
	}

	if svginfo.w != w || svginfo.h != h {
		return nil
	}

	svginfo.setAlive()
	return svginfo.pix
}

// SetColorizedSvg stores the colorized version of an svg in the cache map.
func SetColorizedSvg(src []byte, clr color.NRGBA, content []byte) {
	info := &colorizedSvgInfo{content: content}
	info.setAlive()
	colorizedSvgsLock.Lock()
	colorizedSvgs[colorizedSvgKey{string(src), clr}] = info
	colorizedSvgsLock.Unlock()
}

// SetSvg sets a svg into the cache map.
func SetSvg(name string, o fyne.CanvasObject, pix *image.NRGBA, w int, h int) {
	sinfo := &svgInfo{
		pix: pix,
		w:   w,
		h:   h,
	}
	sinfo.setAlive()
	svgs.Store(overriddenName(name, o), sinfo)
}

type colorizedSvgInfo struct {
	expiringCache
	content []byte
}

type colorizedSvgKey struct {
	src string
	clr color.NRGBA
}

type svgInfo struct {
	expiringCache
	pix  *image.NRGBA
	w, h int
}

// destroyExpiredSvgs destroys expired svgs cache data.
func destroyExpiredSvgs(now time.Time) {
	svgs.Range(func(key string, sinfo *svgInfo) bool {
		if sinfo.isExpired(now) {
			svgs.Delete(key)
		}
		return true
	})

	colorizedSvgsLock.Lock()
	for key, info := range colorizedSvgs {
		if info.isExpired(now) {
			delete(colorizedSvgs, key)
		}
	}
	colorizedSvgsLock.Unlock()
}

func overriddenName(name string, o fyne.CanvasObject) string {
	if o != nil { // for overridden themes get the cache key right
		if over, ok := overrides.Load(o); ok {
			return over.cacheID + name
		}
	}

	return name
}

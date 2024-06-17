package cache

import (
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

var (
	fontSizeCache = map[fontSizeEntry]*fontMetric{}
	fontSizeLock  = sync.RWMutex{}
)

type fontMetric struct {
	expiringCache
	size     fyne.Size
	baseLine float32
}

type fontSizeEntry struct {
	text   string
	size   float32
	style  fyne.TextStyle
	custom string
}

// GetFontMetrics looks up a calculated size and baseline required for the specified text parameters.
func GetFontMetrics(text string, fontSize float32, style fyne.TextStyle, source fyne.Resource) (size fyne.Size, base float32) {
	name := ""
	if source != nil {
		name = source.Name()
	}
	ent := fontSizeEntry{text, fontSize, style, name}
	fontSizeLock.RLock()
	ret, ok := fontSizeCache[ent]
	fontSizeLock.RUnlock()
	if !ok {
		return fyne.Size{Width: 0, Height: 0}, 0
	}
	ret.setAlive()
	return ret.size, ret.baseLine
}

// SetFontMetrics stores a calculated font size and baseline for parameters that were missing from the cache.
func SetFontMetrics(text string, fontSize float32, style fyne.TextStyle, source fyne.Resource, size fyne.Size, base float32) {
	name := ""
	if source != nil {
		name = source.Name()
	}
	ent := fontSizeEntry{text, fontSize, style, name}
	metric := &fontMetric{size: size, baseLine: base}
	metric.setAlive()
	fontSizeLock.Lock()
	fontSizeCache[ent] = metric
	fontSizeLock.Unlock()
}

// destroyExpiredFontMetrics destroys expired fontSizeCache entries
func destroyExpiredFontMetrics(now time.Time) {
	expiredObjs := make([]fontSizeEntry, 0, 50)
	fontSizeLock.RLock()
	for k, v := range fontSizeCache {
		if v.isExpired(now) {
			expiredObjs = append(expiredObjs, k)
		}
	}
	fontSizeLock.RUnlock()

	fontSizeLock.Lock()
	for _, k := range expiredObjs {
		delete(fontSizeCache, k)
	}
	fontSizeLock.Unlock()
}

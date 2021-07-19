package cache

import (
	"sync"

	"fyne.io/fyne/v2"
)

var (
	sizeCache = map[sizeEntry]fontMetric{}
	sizeLock  = sync.RWMutex{}
)

type fontMetric struct {
	bound    fyne.Size
	baseLine float32
}

type sizeEntry struct {
	text  string
	size  float32
	style fyne.TextStyle
}

// GetFontSize looks up a calculated size required for the specified text parameters.
func GetFontSize(text string, size float32, style fyne.TextStyle) (fyne.Size, float32) {
	ent := sizeEntry{text, size, style}
	sizeLock.RLock()
	ret, ok := sizeCache[ent]
	sizeLock.RUnlock()
	if !ok {
		return fyne.Size{Width: 0, Height: 0}, 0
	}
	return ret.bound, ret.baseLine
}

// SetFontSize stores a calculated font size for parameters that were missing from the cache.
func SetFontSize(text string, size float32, style fyne.TextStyle, bound fyne.Size, base float32) {
	ent := sizeEntry{text, size, style}
	sizeLock.Lock()
	sizeCache[ent] = fontMetric{bound: bound, baseLine: base}
	sizeLock.Unlock()
}

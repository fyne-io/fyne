package cache

import (
	"sync"

	"fyne.io/fyne/v2"
)

var (
	fontSizeCache = map[fontSizeEntry]fontMetric{}
	fontSizeLock  = sync.RWMutex{}
)

type fontMetric struct {
	size     fyne.Size
	baseLine float32
}

type fontSizeEntry struct {
	text  string
	size  float32
	style fyne.TextStyle
}

// GetFontMetrics looks up a calculated size and baseline required for the specified text parameters.
func GetFontMetrics(text string, fontSize float32, style fyne.TextStyle) (size fyne.Size, base float32) {
	ent := fontSizeEntry{text, fontSize, style}
	fontSizeLock.RLock()
	ret, ok := fontSizeCache[ent]
	fontSizeLock.RUnlock()
	if !ok {
		return fyne.Size{Width: 0, Height: 0}, 0
	}
	return ret.size, ret.baseLine
}

// SetFontMetrics stores a calculated font size and baseline for parameters that were missing from the cache.
func SetFontMetrics(text string, fontSize float32, style fyne.TextStyle, size fyne.Size, base float32) {
	ent := fontSizeEntry{text, fontSize, style}
	fontSizeLock.Lock()
	fontSizeCache[ent] = fontMetric{size: size, baseLine: base}
	fontSizeLock.Unlock()
}

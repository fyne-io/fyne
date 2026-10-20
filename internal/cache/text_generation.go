package cache

import "sync/atomic"

// fontMetricsGen counts the times the font metrics were cleared, which a
// font or theme change does.
var fontMetricsGen atomic.Uint64

// FontMetricsGeneration changes whenever the font metrics are cleared. A
// measurement kept outside the cache, a canvas.Text's MinSize, is good while
// it stays the same.
func FontMetricsGeneration() uint64 {
	return fontMetricsGen.Load()
}

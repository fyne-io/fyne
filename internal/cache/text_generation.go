package cache

import "sync/atomic"

var fontMetricsGen atomic.Uint64

// FontMetricsGeneration changes whenever the font metrics are cleared.
func FontMetricsGeneration() uint64 {
	return fontMetricsGen.Load()
}

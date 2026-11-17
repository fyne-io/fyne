package cache

import "sync/atomic"

var (
	minSizeGen    atomic.Uint64
	minSizePasses atomic.Uint64
)

// BeginMinSizePass starts a pass in which a widget computes its MinSize once.
func BeginMinSizePass() {
	minSizeGen.Store(minSizePasses.Add(1))
}

// EndMinSizePass ends the current min size pass.
func EndMinSizePass() {
	minSizeGen.Store(0)
}

// MinSizeGeneration returns the current min size pass, or 0 outside a pass.
func MinSizeGeneration() uint64 {
	return minSizeGen.Load()
}

package cache

import "sync/atomic"

// minSizeGen is the min size pass a canvas is running, zero between passes.
// minSizePasses counts them, so no pass gets the generation of one before.
var (
	minSizeGen    atomic.Uint64
	minSizePasses atomic.Uint64
)

// BeginMinSizePass starts a pass over a canvas's objects in which widgets keep
// the MinSize they computed (see MinSizeGeneration). The pass asks every
// object for its min size, and a container's layout asks its children again
// for theirs, at every level, so without it a widget's min size was computed
// once per ancestor. EndMinSizePass ends it.
func BeginMinSizePass() {
	minSizeGen.Store(minSizePasses.Add(1))
}

// EndMinSizePass ends the min size pass: widgets compute their MinSize anew
// again, every time, as they do outside a pass.
func EndMinSizePass() {
	minSizeGen.Store(0)
}

// MinSizeGeneration identifies the min size pass running, or is zero between
// passes. A widget whose MinSize was computed in the pass running returns it
// again until the pass ends, or until it is resized or refreshed.
func MinSizeGeneration() uint64 {
	return minSizeGen.Load()
}

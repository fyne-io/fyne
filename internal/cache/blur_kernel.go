package cache

import (
	"time"

	"fyne.io/fyne/v2/internal/async"
)

var blurKernels async.Map[float32, *blurKernelInfo]

type blurKernelInfo struct {
	expiringCache
	values []float32
}

// GetBlurKernel returns a cached Gaussian kernel for the given pixel-space radius.
func GetBlurKernel(radius float32) ([]float32, bool) {
	info, ok := blurKernels.Load(radius)
	if info == nil || !ok {
		return nil, false
	}
	info.setAlive()
	return info.values, true
}

// SetBlurKernel stores a Gaussian kernel for the given pixel-space radius.
func SetBlurKernel(radius float32, values []float32) {
	info := &blurKernelInfo{values: values}
	info.setAlive()
	blurKernels.Store(radius, info)
}

// destroyExpiredBlurKernels removes blur kernel cache entries that have not been
// used recently.
func destroyExpiredBlurKernels(now time.Time) {
	blurKernels.Range(func(radius float32, info *blurKernelInfo) bool {
		if info.isExpired(now) {
			blurKernels.Delete(radius)
		}
		return true
	})
}

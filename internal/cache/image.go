package cache

import (
	"image"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
)

var images async.Map[string, *imageInfo]

// GetImage gets a decoded raster image from the cache if it exists.
func GetImage(key string, o fyne.CanvasObject) image.Image {
	info, ok := images.Load(overriddenName(key, o))
	if !ok || info == nil {
		return nil
	}

	info.setAlive()
	return info.img
}

// SetImage sets a decoded raster image into the cache.
func SetImage(key string, o fyne.CanvasObject, img image.Image) {
	if img == nil {
		return
	}

	info := &imageInfo{img: img}
	info.setAlive()
	images.Store(overriddenName(key, o), info)
}

type imageInfo struct {
	expiringCache
	img image.Image
}

// destroyExpiredImages destroys expired decoded raster image cache data.
func destroyExpiredImages(now time.Time) {
	images.Range(func(key string, info *imageInfo) bool {
		if info.isExpired(now) {
			images.Delete(key)
		}
		return true
	})
}

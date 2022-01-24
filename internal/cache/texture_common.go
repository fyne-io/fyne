package cache

import (
	"sync"

	"fyne.io/fyne/v2"
)

var textures = sync.Map{} // map[fyne.CanvasObject]*textureInfo

// DeleteTexture deletes the texture from the cache map.
func DeleteTexture(obj fyne.CanvasObject) {
	textures.Delete(obj)
}

// GetTexture gets cached texture.
func GetTexture(obj fyne.CanvasObject) (TextureType, bool) {
	t, ok := textures.Load(obj)
	if t == nil || !ok {
		return NoTexture, false
	}
	texInfo := t.(*textureInfo)
	texInfo.setAlive()
	return texInfo.texture, true
}

// RangeExpiredTexturesFor range over the expired textures for the specified canvas.
//
// Note: If this is used to free textures, then it should be called inside a current
// gl context to ensure textures are deleted from gl.
func RangeExpiredTexturesFor(canvas fyne.Canvas, f func(fyne.CanvasObject)) {
	now := timeNow()
	textures.Range(func(key, value interface{}) bool {
		obj, tinfo := key.(fyne.CanvasObject), value.(*textureInfo)
		if tinfo.isExpired(now) && tinfo.canvas == canvas {
			f(obj)
		}
		return true
	})
}

// RangeTexturesFor range over the textures for the specified canvas.
//
// Note: If this is used to free textures, then it should be called inside a current
// gl context to ensure textures are deleted from gl.
func RangeTexturesFor(canvas fyne.Canvas, f func(fyne.CanvasObject)) {
	textures.Range(func(key, value interface{}) bool {
		obj, tinfo := key.(fyne.CanvasObject), value.(*textureInfo)
		if tinfo.canvas == canvas {
			f(obj)
		}
		return true
	})
}

// SetTexture sets cached texture.
func SetTexture(obj fyne.CanvasObject, texture TextureType, canvas fyne.Canvas) {
	texInfo := &textureInfo{texture: texture}
	texInfo.canvas = canvas
	texInfo.setAlive()
	textures.Store(obj, texInfo)
}

// textureCacheBase defines base texture cache object.
type textureCacheBase struct {
	expiringCacheNoLock
	canvas fyne.Canvas
}

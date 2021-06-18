package cache

import (
	"fyne.io/fyne/v2"
)

// NOTE: Texture cache functions should always be called in
// the same goroutine.

var textures = make(map[fyne.CanvasObject]*textureInfo, 1024)

// DeleteTexture deletes the texture from the cache map.
func DeleteTexture(obj fyne.CanvasObject) {
	delete(textures, obj)
}

// GetTexture gets cached texture.
func GetTexture(obj fyne.CanvasObject) (TextureType, bool) {
	texInfo, ok := textures[obj]
	if texInfo == nil || !ok {
		return noTexture, false
	}
	texInfo.setAlive()
	return texInfo.texture, true
}

// RangeExpiredTexturesFor range over the expired textures for the specified canvas.
//
// Note: If this is used to free textures, then it should be called inside a current
// gl context to ensure textures are deleted from gl.
func RangeExpiredTexturesFor(canvas fyne.Canvas, f func(fyne.CanvasObject)) {
	now := timeNow()
	for obj, tinfo := range textures {
		if tinfo.isExpired(now) && tinfo.canvas == canvas {
			f(obj)
		}
	}
}

// RangeTexturesFor range over the textures for the specified canvas.
//
// Note: If this is used to free textures, then it should be called inside a current
// gl context to ensure textures are deleted from gl.
func RangeTexturesFor(canvas fyne.Canvas, f func(fyne.CanvasObject)) {
	for obj, tinfo := range textures {
		if tinfo.canvas == canvas {
			f(obj)
		}
	}
}

// SetTexture sets cached texture.
func SetTexture(obj fyne.CanvasObject, texture TextureType, canvas fyne.Canvas) {
	texInfo := &textureInfo{texture: texture}
	texInfo.canvas = canvas
	texInfo.setAlive()
	textures[obj] = texInfo
}

// textureCacheBase defines base texture cache object.
type textureCacheBase struct {
	expiringCacheNoLock
	canvas fyne.Canvas
}

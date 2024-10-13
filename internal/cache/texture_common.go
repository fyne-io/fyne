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

// GetTextTexture gets cached texture for a text run.
func GetTextTexture(ent FontCacheEntry) (TextureType, bool) {
	return load(ent)
}

// GetTexture gets cached texture.
func GetTexture(obj fyne.CanvasObject) (TextureType, bool) {
	return load(obj)
}

func load(obj any) (TextureType, bool) {
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
	textures.Range(func(key, value any) bool {
		if _, ok := key.(FontCacheEntry); ok {
			tinfo := value.(*textureInfo)

			// just free text directly when that string/style combo is done
			if tinfo.isExpired(now) && tinfo.canvas == canvas {
				textures.Delete(key)
				tinfo.textFree()
			}

			return true
		}
		obj, tinfo := key.(fyne.CanvasObject), value.(*textureInfo)
		if tinfo.isExpired(now) && tinfo.canvas == canvas {
			f(obj)
		}
		return true
	})
}

// RangeTexturesFor range over the textures for the specified canvas.
// It will not return the texture for a `canvas.Text` as their render lifecycle is handled separately.
//
// Note: If this is used to free textures, then it should be called inside a current
// gl context to ensure textures are deleted from gl.
func RangeTexturesFor(canvas fyne.Canvas, f func(fyne.CanvasObject)) {
	textures.Range(func(key, value any) bool {
		if _, ok := key.(FontCacheEntry); ok {
			return true // do nothing, text cache lives outside the scope of an object
		}

		obj, tinfo := key.(fyne.CanvasObject), value.(*textureInfo)
		if tinfo.canvas == canvas {
			f(obj)
		}
		return true
	})
}

// SetTextTexture sets cached texture for a text run.
func SetTextTexture(ent FontCacheEntry, texture TextureType, canvas fyne.Canvas, free func()) {
	store(ent, texture, canvas, free)
}

// SetTexture sets cached texture.
func SetTexture(obj fyne.CanvasObject, texture TextureType, canvas fyne.Canvas) {
	store(obj, texture, canvas, nil)
}

func store(obj any, texture TextureType, canvas fyne.Canvas, free func()) {
	texInfo := &textureInfo{texture: texture}
	if free != nil {
		texInfo.textFree = free
	}
	texInfo.canvas = canvas
	texInfo.setAlive()
	textures.Store(obj, texInfo)
}

// textureCacheBase defines base texture cache object.
type textureCacheBase struct {
	expiringCacheNoLock
	canvas fyne.Canvas
}

package cache

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
)

var (
	textTextures             async.Map[FontCacheEntry, *textureInfo]
	textTextureLastCleanSize int
	shouldCleanTextTextures  bool

	objectTextures              async.Map[fyne.CanvasObject, *textureInfo]
	objectTexturesLastCleanSize int
	shouldCleanObjectTextures   bool
)

// DeleteTexture deletes the texture from the cache map.
func DeleteTexture(obj fyne.CanvasObject) {
	objectTextures.Delete(obj)
}

// GetTextTexture gets cached texture for a text run.
func GetTextTexture(ent FontCacheEntry) (TextureType, bool) {
	texInfo, ok := textTextures.Load(ent)
	if texInfo == nil || !ok {
		return NoTexture, false
	}
	texInfo.setAlive()
	return texInfo.texture, true
}

// GetTexture gets cached texture.
func GetTexture(obj fyne.CanvasObject) (TextureType, bool) {
	texInfo, ok := objectTextures.Load(obj)
	if texInfo == nil || !ok {
		return NoTexture, false
	}
	texInfo.setAlive()
	return texInfo.texture, true
}

// rangeExpiredTexturesFor range over the expired textures for the specified canvas.
// Note that this function *does not* range over text textures, which are cleaned
// on a different schedule.
//
// Note: If this is used to free textures, then it should be called inside a current
// gl context to ensure textures are deleted from gl.
func rangeExpiredTexturesFor(canvas fyne.Canvas, f func(fyne.CanvasObject)) {
	now := timeNow()

	objectTextures.Range(func(obj fyne.CanvasObject, tinfo *textureInfo) bool {
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
	// Do nothing for texture cache, it lives outside the scope of an object.
	objectTextures.Range(func(obj fyne.CanvasObject, tinfo *textureInfo) bool {
		if tinfo.canvas == canvas {
			f(obj)
		}
		return true
	})
}

// DeleteTextTexturesFor deletes all text textures for the given canvas.
func DeleteTextTexturesFor(canvas fyne.Canvas) {
	textTextures.Range(func(key FontCacheEntry, tinfo *textureInfo) bool {
		if tinfo.canvas == canvas {
			textTextures.Delete(key)
			tinfo.textFree()
		}
		return true
	})
}

// SetTextTexture sets cached texture for a text run.
func SetTextTexture(ent FontCacheEntry, texture TextureType, canvas fyne.Canvas, free func()) {
	tinfo := prepareTexture(texture, canvas, free)
	textTextures.Store(ent, tinfo)
	if textTextures.Len() > 2*textTextureLastCleanSize {
		shouldCleanTextTextures = true
	}
}

// SetTexture sets cached texture.
func SetTexture(obj fyne.CanvasObject, texture TextureType, canvas fyne.Canvas) {
	tinfo := prepareTexture(texture, canvas, nil)
	objectTextures.Store(obj, tinfo)
	if objectTextures.Len() > 2*objectTexturesLastCleanSize {
		shouldCleanObjectTextures = true
	}
}

func prepareTexture(texture TextureType, canvas fyne.Canvas, free func()) *textureInfo {
	tinfo := &textureInfo{texture: texture, textFree: free}
	tinfo.canvas = canvas
	tinfo.setAlive()
	return tinfo
}

// textureCacheBase defines base texture cache object.
type textureCacheBase struct {
	frameCounterCache
	canvas fyne.Canvas
}

func cleanTextTextureCache(forCanvas fyne.Canvas) {
	textTextures.Range(func(key FontCacheEntry, tinfo *textureInfo) bool {
		// Just free text directly when that string/style combo is done.
		if tinfo.isExpired(time.Time{}) && (forCanvas == nil || tinfo.canvas == forCanvas) {
			textTextures.Delete(key)
			tinfo.textFree()
		}
		return true
	})
}
